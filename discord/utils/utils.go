package utils

import (
	"context"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
)

//SessionInfo represent common information about a session
type SessionInfo struct {
	Con     context.Context
	Session disgord.Session
}

//ShowError displays an error
func ShowError(info SessionInfo, channelID snowflake.Snowflake, message string) (*disgord.Message, error) {
	msg, err := info.Session.SendMsg(info.Con, channelID, disgord.Embed{Title: "Error", Color: 0xb00000,
		Description: message + "\n **`.pingmc help <command>` for help** \n **`.pingmc commands` for command list**"})
	return msg, err
}

//ShowSuccess display a success
func ShowSuccess(info SessionInfo, channelID snowflake.Snowflake, message string) (*disgord.Message, error) {
	msg, err := info.Session.SendMsg(info.Con, channelID, disgord.Embed{Title: "Success!", Color: 0x00ad37,
		Description: message})
	return msg, err
}
