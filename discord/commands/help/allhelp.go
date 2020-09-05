package help

import (
	"github.com/andersfylling/disgord"
	"yey007.github.io/software/pingmc/discord/utils"
	"yey007.github.io/software/pingmc/discord/utils/colors"
)

//OnCommandsRequest displays the list of all commands
func OnCommandsRequest(info utils.SessionInfo, msg *disgord.Message, args []string) {
	helpmap, err := readJSON()
	if err != nil {
		utils.ShowError(info, msg.ChannelID, "Error reading helpmap.json. Please report this to the author, Yey007#3321")
	}

	var help string
	for _, v := range helpmap {
		help += utils.Bold(utils.Block(v.Cmd)) + " - " + v.Desc + "\n"
	}
	info.Session.SendMsg(info.Con, msg.ChannelID, disgord.Embed{Title: "Commands", Color: colors.Blue,
		Description: help})
}
