package help

import (
	"github.com/andersfylling/disgord"
	"yey007.github.io/software/pingmc/discord/utils"
	"yey007.github.io/software/pingmc/discord/utils/colors"
)

//OnHelpRequest gives help about a command
func OnHelpRequest(info utils.SessionInfo, msg *disgord.Message, args []string) {

	if len(args) == 3 {

		helpmap, err := readJSON()
		if err != nil {
			utils.ShowError(info, msg.ChannelID, "Error reading helpmap.json. Please report this to the author, Yey007#3321")
		}

		cmd, ok := helpmap[args[2]]

		if ok {

			//build and format help string
			help := utils.Bold(utils.Block(cmd.Cmd)) + " - " + cmd.Desc + "\n\n"
			for _, a := range cmd.Args {
				help += a.Arg + " - " + utils.Italics(a.Desc) + "\n"
			}

			info.Session.SendMsg(info.Con, msg.ChannelID, disgord.Embed{Title: "Help", Color: colors.Blue,
				Description: help})
		} else {
			utils.ShowError(info, msg.ChannelID, "Command not found")
		}
	} else {
		utils.ShowError(info, msg.ChannelID, "Wrong number of arguments")
	}
}
