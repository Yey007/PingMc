package help

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/andersfylling/disgord"
	"yey007.github.io/software/pingmc/discord/utils"
)

var helpmap map[string]command
var doOnce sync.Once

//OnHelpRequest gives help about a command
func OnHelpRequest(info utils.SessionInfo, msg *disgord.Message, args []string) {

	if len(args) == 3 {

		doOnce.Do(func() {
			helpdata, err := ioutil.ReadFile("data/helpmap.json")
			if err != nil {
				utils.ShowError(info, msg.ChannelID, "Something went wrong reading helpmap.json. Please report this to the author (Yey007#3321).")
				return
			}

			err = json.Unmarshal(helpdata, &helpmap)
			if err != nil {
				utils.ShowError(info, msg.ChannelID, "Something went wrong reading helpmap.json. Please report this to the author (Yey007#3321).")
				return
			}
		})

		cmd, ok := helpmap[args[2]]

		if ok {

			//build and format help string
			help := "**" + cmd.Cmd + "** - " + cmd.Desc + "\n\n"
			for _, a := range cmd.Args {
				help += a.Arg + " - *" + a.Desc + "*\n"
			}

			info.Session.SendMsg(info.Con, msg.ChannelID, disgord.Embed{Title: "Help", Color: 0x1a5fba,
				Description: help})
		} else {
			utils.ShowError(info, msg.ChannelID, "Command not found")
		}
	} else {
		utils.ShowError(info, msg.ChannelID, "Wrong number of arguments")
	}
}
