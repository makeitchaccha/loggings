package bot

import (
	"context"
	"fmt"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/makeitchaccha/loggings/internal/pkg/command"
	"github.com/makeitchaccha/loggings/internal/pkg/logging"
	"github.com/makeitchaccha/loggings/internal/pkg/model"
	"github.com/makeitchaccha/loggings/internal/pkg/settings"
	"gorm.io/gorm"
)

type Bot struct {
	client bot.Client

	manager *settings.Manager

	commands map[string]command.Command
}

func New(token string, db *gorm.DB) (*Bot, error) {
	b := &Bot{
		commands: make(map[string]command.Command),
	}

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildPresences, gateway.IntentGuildMembers)),
		bot.WithEventListenerFunc(b.onGuildsReady),
		bot.WithEventListenerFunc(b.onGuildJoin),
		bot.WithEventListenerFunc(b.onGuildMemberJoin),
		bot.WithEventListenerFunc(b.onApplicationCommandInteractionCreate),
		bot.WithEventListenerFunc(b.onAutoCompleteInteractionCreate),
	)

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.GuildSettings{})

	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	b.manager = settings.NewManager(db)

	b.client = client

	b.registerCommands()

	return b, nil
}

func (b *Bot) Open(ctx context.Context) error {
	return b.client.OpenGateway(ctx)
}

func (b *Bot) Close(ctx context.Context) {
	b.client.Close(ctx)
}

func (b *Bot) registerCommands() {
	commands := []command.Command{
		&command.SettingsCommand{
			Events:  logging.Events,
			Manager: b.manager,
		},
	}

	for _, c := range commands {
		b.commands[c.Name()] = c
	}
}

func (b *Bot) onGuildsReady(event *events.GuildReady) {
	b.initGuild(event.Guild.ID)
}

func (b *Bot) onGuildJoin(event *events.GuildJoin) {
	b.initGuild(event.Guild.ID)
}

func (b *Bot) initGuild(guildId snowflake.ID) {
	_, ok := b.manager.Guild(guildId)
	if !ok {
		guild := settings.NewGuild(guildId)
		b.manager.SaveGuild(guild)
	}
}

func (b *Bot) onGuildMemberJoin(event *events.GuildMemberJoin) {
	guild, _ := b.manager.Guild(event.GuildID)

	if !guild.Loggings.MemberJoin.Enabled {
		return
	}

	embed, err := guild.Loggings.MemberJoin.Embed(event)

	if err != nil {
		// fallback to default embed
		b.client.Rest().CreateMessage(guild.Loggings.MemberJoin.Channel, discord.MessageCreate{
			Content: "-# メンバー参加メッセージの作成に失敗しました。設定からカスタム埋め込みを修正してください",
			Embeds: []discord.Embed{
				guild.Loggings.MemberJoin.DefaultEmbed(event),
			},
		})
		return
	}

	b.client.Rest().CreateMessage(guild.Loggings.MemberJoin.Channel, discord.MessageCreate{
		Embeds: []discord.Embed{
			embed,
		},
	})
}

func (b *Bot) onApplicationCommandInteractionCreate(event *events.ApplicationCommandInteractionCreate) {
	cmd, ok := b.commands[event.SlashCommandInteractionData().CommandName()]
	if !ok {
		return
	}

	err := cmd.Execute(event)

	if err != nil {
		event.CreateMessage(discord.MessageCreate{
			Content: fmt.Sprintf("error: %v", err),
		})
	}
}

func (b *Bot) onAutoCompleteInteractionCreate(event *events.AutocompleteInteractionCreate) {
	cmd, ok := b.commands[event.Data.CommandName]
	if !ok {
		return
	}

	if cmd, ok := cmd.(command.Autocomplete); ok {
		cmd.Autocomplete(event)
	}
}
