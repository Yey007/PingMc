package discord

import (
	"context"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
)

type sessionInfo struct {
	con     context.Context
	session disgord.Session
}

func showUsage(info sessionInfo, channelID snowflake.Snowflake) (*disgord.Message, error) {
	msg, err := info.session.SendMsg(info.con, channelID, disgord.Embed{Title: "Commands", Color: 0x1a5fba,
		Description: `.pingmc help - displays this message
		.pingmc ping ip serverType channelName
			*ip - the IP to ping. If the port isn't provided, 25565 is assumed*
			*serverType - the type of server to ping (forge or other)*
			*channelName - the channel to post updates in*`})
	return msg, err
}

func showError(info sessionInfo, channelID snowflake.Snowflake, message string) (*disgord.Message, error) {
	msg, err := info.session.SendMsg(info.con, channelID, disgord.Embed{Title: "Error", Color: 0xb00000,
		Description: message + "\n **`.pingmc help` for help**"})
	return msg, err
}

func showSuccess(info sessionInfo, channelID snowflake.Snowflake, message string) (*disgord.Message, error) {
	msg, err := info.session.SendMsg(info.con, channelID, disgord.Embed{Title: "Success!", Color: 0x00ad37,
		Description: message})
	return msg, err
}
