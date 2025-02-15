package model

import "gorm.io/gorm"

type GuildSettings struct {
	gorm.Model
	GuildID uint64 `gorm:"primaryKey"`

	// Loggings channels
	MemberJoinEnabled bool
	MemberJoinChannel uint64
	MemberJoinFormat  string
}
