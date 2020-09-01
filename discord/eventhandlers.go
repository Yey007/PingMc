package discord

import (
	"context"
	"strings"

	"github.com/andersfylling/disgord"
)

func onMessageCreate(session disgord.Session, evt *disgord.MessageCreate) {
	con := context.Background()
	msg := evt.Message
	args := strings.Split(msg.Content, " ")
	info := sessionInfo{con, session}

	if len(args) >= 2 && args[0] == ".pingmc" {
		if len(args) == 5 && args[1] == "ping" {
			go onPingRequest(info, msg.GuildID, msg.ChannelID, args)
		} else {
			go showUsage(info, msg.ChannelID)
		}
	}
}
