package data

import (
	"gorm.io/gorm"
	"yey007.github.io/software/pingmc/networking"
)

type RecurringPing struct {
	gorm.Model
	Name      string
	GuildID   string
	ChannelID string
	Server    networking.McServer
}

type TempPingData struct {
	RecurringPingID uint
	PreviousData    *networking.PingData
	HasPrevious     bool
}
