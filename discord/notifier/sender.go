package notifier

import (
	"fmt"

	"github.com/andersfylling/disgord"
	"yey007.github.io/software/pingmc/discord/utils"
	"yey007.github.io/software/pingmc/networking"
)

func createServerOnline(pingData networking.PingData, serverName string) *disgord.Embed {
	return &disgord.Embed{
		Title:       fmt.Sprintf("Online (%v/%v)", pingData.Players.Online, pingData.Players.Max),
		Color:       utils.GREEN,
		Description: fmt.Sprintf("The server %v may have come online", utils.Block(serverName)),
	}
}

func createServerOffline(serverName string) *disgord.Embed {
	return &disgord.Embed{
		Title:       "Offline",
		Color:       utils.RED,
		Description: fmt.Sprintf("The server %v is probably offline", utils.Block(serverName)),
	}
}

func createPlayerJoin(pingData networking.PingData, serverName string) *disgord.Embed {
	return &disgord.Embed{
		Title:       fmt.Sprintf("Join (%v/%v)", pingData.Players.Online, pingData.Players.Max),
		Color:       utils.BLUE,
		Description: fmt.Sprintf("A player may have come online on the server %v", utils.Block(serverName)),
	}
}

func createPlayerLeave(pingData networking.PingData, serverName string) *disgord.Embed {
	return &disgord.Embed{
		Title:       fmt.Sprintf("Leave (%v/%v)", pingData.Players.Online, pingData.Players.Max),
		Color:       utils.ORANGE,
		Description: fmt.Sprintf("A player may have went offline on the server %v", utils.Block(serverName)),
	}
}
