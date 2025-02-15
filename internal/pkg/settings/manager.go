package settings

import (
	"fmt"

	"github.com/disgoorg/snowflake/v2"
	"github.com/makeitchaccha/loggings/internal/pkg/model"
	"gorm.io/gorm"
)

type Manager struct {
	db     *gorm.DB
	guilds map[snowflake.ID]Guild // guildID -> Guild
}

func NewManager(db *gorm.DB) *Manager {
	mgr := &Manager{
		db:     db,
		guilds: make(map[snowflake.ID]Guild),
	}

	// load guilds
	var guilds []model.GuildSettings

	if err := db.Find(&guilds).Error; err != nil {
		panic("failed to load guilds")
	}

	for _, guild := range guilds {
		mgr.guilds[snowflake.ID(guild.GuildID)] = UnmarshalModel(&guild)
	}

	return mgr
}

func (m *Manager) Guild(guildID snowflake.ID) (Guild, bool) {
	guild, ok := m.guilds[guildID]
	return guild, ok
}

func (m *Manager) SaveGuild(guild Guild) error {
	err := m.db.Save(guild.MarshalModel()).Error
	if err != nil {
		return fmt.Errorf("failed to save guild: %w", err)
	}

	m.guilds[guild.ID] = guild
	return nil
}
