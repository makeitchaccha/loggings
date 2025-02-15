package settings

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

type MemberJoin struct {
	Enabled         bool
	Channel         snowflake.ID // only valid if enabled
	CustomEmbedJson string       // only valid if enabled
}

func (m MemberJoin) Embed(event *events.GuildMemberJoin) (discord.Embed, error) {
	if m.CustomEmbedJson != "" {
		return m.customEmbed(event)
	}
	return m.DefaultEmbed(event), nil
}

func (m MemberJoin) customEmbed(event *events.GuildMemberJoin) (discord.Embed, error) {

	replaces := map[string]string{
		"{display name}": event.Member.EffectiveName(),
		"{username}":     event.Member.User.Username,
		"{mention}":      event.Member.User.Mention(),
		"{avatar url}":   event.Member.User.EffectiveAvatarURL(),
		"{timestamp}":    time.Now().UTC().Format(time.RFC3339),
	}

	customEmbedJson := m.CustomEmbedJson
	for k, v := range replaces {
		customEmbedJson = strings.ReplaceAll(customEmbedJson, k, v)
	}

	var embed discord.Embed
	// unmarshal m.CustomEmbedJson to embed
	if err := json.Unmarshal([]byte(customEmbedJson), &embed); err != nil {
		return discord.Embed{}, err
	}

	return embed, nil
}

func (m MemberJoin) DefaultEmbed(event *events.GuildMemberJoin) discord.Embed {
	return discord.Embed{
		Title:       "メンバーが参加しました",
		Description: event.Member.User.Mention(),
		Color:       0x00ff00,
	}
}
