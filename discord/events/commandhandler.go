package events

import (
	"context"
	"strings"

	"github.com/andersfylling/disgord"
	"yey007.github.io/software/pingmc/discord/commands/help"
	"yey007.github.io/software/pingmc/discord/commands/updater"
	"yey007.github.io/software/pingmc/discord/utils"
)

var commandMap = map[string]func(utils.SessionInfo, *disgord.Message, []string){
	"ping": updater.OnPingRequest,
	"help": help.OnHelpRequest,
}

//OnMessageCreate handles a message send event from discord
func OnMessageCreate(session disgord.Session, evt *disgord.MessageCreate) {

	con := context.Background()
	msg := evt.Message
	args := strings.Split(msg.Content, " ")
	info := utils.SessionInfo{Con: con, Session: session}

	if len(args) >= 2 && args[0] == ".pingmc" {

		if command, ok := commandMap[args[1]]; ok {
			go command(info, msg, args)
		} else {
			utils.ShowError(info, msg.ChannelID, "Command not found")
		}
	}
}
