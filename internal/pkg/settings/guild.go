package settings

import (
	"github.com/disgoorg/snowflake/v2"
	"github.com/makeitchaccha/loggings/internal/pkg/model"
)

type Guild struct {
	ID       snowflake.ID
	Loggings Loggings
}

type Loggings struct {
	MemberJoin MemberJoin
}

func NewGuild(id snowflake.ID) Guild {
	return Guild{ID: id, Loggings: Loggings{
		MemberJoin: MemberJoin{
			Enabled:         false,
			Channel:         0,
			CustomEmbedJson: "",
		},
	}}
}

func (g *Guild) MarshalModel() *model.GuildSettings {
	return &model.GuildSettings{
		GuildID: uint64(g.ID),

		// MemberJoin
		MemberJoinEnabled: g.Loggings.MemberJoin.Enabled,
		MemberJoinChannel: uint64(g.Loggings.MemberJoin.Channel),
		MemberJoinFormat:  g.Loggings.MemberJoin.CustomEmbedJson,
	}
}

func UnmarshalModel(guild *model.GuildSettings) Guild {
	return Guild{
		ID: snowflake.ID(guild.GuildID),
		Loggings: Loggings{
			MemberJoin: MemberJoin{
				Enabled:         guild.MemberJoinEnabled,
				Channel:         snowflake.ID(guild.MemberJoinChannel),
				CustomEmbedJson: guild.MemberJoinFormat,
			},
		},
	}
}
