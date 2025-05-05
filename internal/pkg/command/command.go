package command

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/makeitchaccha/loggings/internal/pkg/logging"
	"github.com/makeitchaccha/loggings/internal/pkg/settings"
	"github.com/samber/lo"
)

type Command interface {
	Name() string
	Create() discord.ApplicationCommandCreate
	Execute(event *events.ApplicationCommandInteractionCreate) error
}

type Autocomplete interface {
	Autocomplete(event *events.AutocompleteInteractionCreate) error
}

type SettingsCommand struct {
	Events  []logging.Event
	Manager *settings.Manager
}

var _ Command = (*SettingsCommand)(nil)

func (s *SettingsCommand) Name() string {
	return "logging"
}

func (s *SettingsCommand) Create() discord.ApplicationCommandCreate {
	return discord.SlashCommandCreate{
		Name:        s.Name(),
		Description: "log settings",
		DescriptionLocalizations: map[discord.Locale]string{
			discord.LocaleJapanese: "ログ設定",
		},
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        "set-channel",
				Description: "set log channel",
				DescriptionLocalizations: map[discord.Locale]string{
					discord.LocaleJapanese: "ログチャンネルを設定します",
				},
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:        "type",
						Description: "event type e.g. *member join, delete message*",
						DescriptionLocalizations: map[discord.Locale]string{
							discord.LocaleJapanese: "イベントの種類 例: *メンバーの参加, メッセージの削除*",
						},
						Required:     true,
						Autocomplete: true,
					},
					discord.ApplicationCommandOptionChannel{
						Name:        "channel",
						Description: "channel to set",
						DescriptionLocalizations: map[discord.Locale]string{
							discord.LocaleJapanese: "設定するチャンネル",
						},
						Required: false,
					},
				},
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        "set-format",
				Description: "set log format",
				DescriptionLocalizations: map[discord.Locale]string{
					discord.LocaleJapanese: "ログのフォーマットを設定します",
				},
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:        "type",
						Description: "event type e.g. *member join, delete message*",
						DescriptionLocalizations: map[discord.Locale]string{
							discord.LocaleJapanese: "イベントの種類 例: *メンバーの参加, メッセージの削除*",
						},
						Required:     true,
						Autocomplete: true,
					},
					discord.ApplicationCommandOptionString{
						Name:        "format",
						Description: "format to set in embed json",
						DescriptionLocalizations: map[discord.Locale]string{
							discord.LocaleJapanese: "設定するフォーマット(埋め込みjson)",
						},
						Required: false,
					},
				},
			},
		},
	}
}

func (s *SettingsCommand) Execute(event *events.ApplicationCommandInteractionCreate) error {
	data := event.SlashCommandInteractionData()

	if data.SubCommandName == nil {
		event.CreateMessage(discord.MessageCreate{
			Content: "error: subcommand not found",
		})
	}

	if event.GuildID() == nil {
		event.CreateMessage(discord.MessageCreate{
			Content: "サーバー内でのみ使用できます",
		})
		return nil
	}

	switch *data.SubCommandName {
	case "set-channel":
		return s.executeSetChannel(event)
	case "set-format":
		return s.executeSetFormat(event)
	}
	return nil
}

func (s *SettingsCommand) executeSetChannel(event *events.ApplicationCommandInteractionCreate) error {
	data := event.SlashCommandInteractionData()

	t, ok := logging.ParseEvent(data.String("type"))
	if !ok {
		event.CreateMessage(discord.MessageCreate{
			Content: "エラー: 有効なイベントタイプではありません",
		})
		return nil
	}

	// set channel
	channelId := event.Channel().ID()

	if channel, ok := data.OptChannel("channel"); ok {
		channelId = channel.ID
	}

	// set channel
	guild, _ := s.Manager.Guild(*event.GuildID())

	switch t {
	case logging.EventMemberJoin:
		memberJoin := guild.Loggings.MemberJoin
		memberJoin.Enabled = true
		memberJoin.Channel = channelId
		guild.Loggings.MemberJoin = memberJoin
	}

	if err := s.Manager.SaveGuild(guild); err != nil {
		event.CreateMessage(discord.MessageCreate{
			Content: "エラー: チャンネルの設定に失敗しました",
		})
	}

	event.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("ログチャンネルを <#%s> に設定しました", channelId),
	})

	return nil
}

func (s *SettingsCommand) executeSetFormat(event *events.ApplicationCommandInteractionCreate) error {
	data := event.SlashCommandInteractionData()

	t, ok := logging.ParseEvent(data.String("type"))

	if !ok {
		event.CreateMessage(discord.MessageCreate{
			Content: "エラー: 有効なイベントタイプではありません",
		})
		return nil
	}

	format := data.String("format")

	var embed discord.Embed

	// validate format
	if format != "" {
		format := settings.SanitizeFormat(format)

		if err := json.Unmarshal([]byte(format), &embed); err != nil {
			event.CreateMessage(discord.MessageCreate{
				Content: "エラー: フォーマットが無効です",
			})
			return nil
		}
	}

	guild, _ := s.Manager.Guild(*event.GuildID())

	switch t {
	case logging.EventMemberJoin:
		memberJoin := guild.Loggings.MemberJoin
		memberJoin.CustomEmbedJson = format
		guild.Loggings.MemberJoin = memberJoin
	}

	if err := s.Manager.SaveGuild(guild); err != nil {
		event.CreateMessage(discord.MessageCreate{
			Content: "エラー: フォーマットの設定に失敗しました",
		})
	}

	embeds := []discord.Embed{}
	if format != "" {
		embeds = append(embeds, embed)
	}

	event.CreateMessage(discord.MessageCreate{
		Content: "ログフォーマットを設定しました",
		Embeds:  embeds,
	})
	return nil
}

func (s *SettingsCommand) Autocomplete(event *events.AutocompleteInteractionCreate) error {
	if event.Data.SubCommandName == nil {
		return nil
	}

	switch *event.Data.SubCommandName {
	case "set-channel":
		return s.autocompleteType(event)
	case "set-format":
		return s.autocompleteType(event)
	}
	return nil
}

func (s *SettingsCommand) autocompleteType(event *events.AutocompleteInteractionCreate) error {
	t := event.Data.String("type")

	prefixeds := lo.Filter(s.Events, func(e logging.Event, _ int) bool {
		return strings.HasPrefix(e.String(), t)
	})

	choices := lo.Map(prefixeds, func(e logging.Event, _ int) discord.AutocompleteChoice {
		return discord.AutocompleteChoiceString{
			Name:  e.String(),
			Value: e.String(),
		}
	})

	event.AutocompleteResult(choices)

	return nil
}
