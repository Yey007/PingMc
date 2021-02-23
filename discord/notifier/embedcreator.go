package notifier

import (
	"fmt"
	"yey007.github.io/software/pingmc/data"

	"github.com/andersfylling/disgord"
	"yey007.github.io/software/pingmc/discord/utils"
	"yey007.github.io/software/pingmc/networking"
)

func createNotification(
	p data.RecurringPing,
	temp data.TempPingData,
	repo *data.PingRepo,
) (disgord.Embed, bool) {
	pingdata, err := p.Server.Ping()

	if temp.HasPrevious {
		if err == nil {
			defer repo.CreateTemp(data.TempPingData{
				RecurringPingID: p.ID,
				PreviousData:    pingdata,
				HasPrevious:     true,
			})
			//Server was online and is still online
			if pingdata.Players.Online > temp.PreviousData.Players.Online {
				//Number of players more
				return createPlayerJoin(pingdata, p.Name), true
			} else if pingdata.Players.Online < temp.PreviousData.Players.Online {
				//Number of players less
				return createPlayerLeave(pingdata, p.Name), true
			} else if pingdata.Players.Max != temp.PreviousData.Players.Max {
				//Max players different
				return createServerOnline(pingdata, p.Name), true
			}
		} else {
			//Server was online and is now offline
			repo.CreateTemp(data.TempPingData{
				RecurringPingID: p.ID,
				PreviousData:    pingdata,
				HasPrevious:     false,
			})
			return createServerOffline(p.Name), true
		}
	} else {
		if err == nil {
			//Server was offline and is now online
			repo.CreateTemp(data.TempPingData{
				RecurringPingID: p.ID,
				PreviousData:    pingdata,
				HasPrevious:     true,
			})
			return createServerOnline(pingdata, p.Name), true
		}
		//Server was offline and is still offline, do nothing
	}
	return disgord.Embed{}, false
}

func createServerOnline(pingData networking.PingData, serverName string) disgord.Embed {
	return disgord.Embed{
		Title:       fmt.Sprintf("Online (%v/%v)", pingData.Players.Online, pingData.Players.Max),
		Color:       utils.GREEN,
		Description: fmt.Sprintf("The server %v may have come online", utils.Block(serverName)),
	}
}

func createServerOffline(serverName string) disgord.Embed {
	return disgord.Embed{
		Title:       "Offline",
		Color:       utils.RED,
		Description: fmt.Sprintf("The server %v is probably offline", utils.Block(serverName)),
	}
}

func createPlayerJoin(pingData networking.PingData, serverName string) disgord.Embed {
	return disgord.Embed{
		Title:       fmt.Sprintf("Join (%v/%v)", pingData.Players.Online, pingData.Players.Max),
		Color:       utils.BLUE,
		Description: fmt.Sprintf("A player may have come online on the server %v", utils.Block(serverName)),
	}
}

func createPlayerLeave(pingData networking.PingData, serverName string) disgord.Embed {
	return disgord.Embed{
		Title:       fmt.Sprintf("Leave (%v/%v)", pingData.Players.Online, pingData.Players.Max),
		Color:       utils.ORANGE,
		Description: fmt.Sprintf("A player may have went offline on the server %v", utils.Block(serverName)),
	}
}
