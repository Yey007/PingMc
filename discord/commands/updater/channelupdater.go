package updater

import (
	"net"
	"strconv"
	"time"

	"github.com/andersfylling/disgord"
	"yey007.github.io/software/pingmc/discord/utils"
	"yey007.github.io/software/pingmc/networking"
)

//OnPingRequest starts pinging a server
func OnPingRequest(info utils.SessionInfo, msg *disgord.Message, args []string) {

	channelID := msg.ChannelID
	ip := args[2]
	serverType := args[3]

	if validateInputs(info, channelID, args) == false {
		return
	}

	pingChannel, err := info.Session.GetChannel(info.Con, channelID)
	if err != nil {
		utils.ShowError(info, channelID, "Unable to create channel. Reason is unknown.")
	}

	server, err := createServer(info, serverType)

	datastring := "(0/0)"
	var pingTime time.Duration = 5
	wasOffline := true

	online := func() {
		embed := disgord.Embed{Title: "Online " + datastring, Color: 0x00ad37,
			Description: "The Minecraft server came online"}
		info.Session.SendMsg(info.Con, pingChannel.ID, embed)
	}

	join := func() {
		embed := disgord.Embed{Title: "Join " + datastring, Color: 0x1a5fba,
			Description: "Someone entered the Minecraft server"}
		info.Session.SendMsg(info.Con, pingChannel.ID, embed)
	}

	leave := func() {
		embed := disgord.Embed{Title: "Leave " + datastring, Color: 0xe8bb35,
			Description: "Someone left the Minecraft server"}
		info.Session.SendMsg(info.Con, pingChannel.ID, embed)
	}

	offline := func() {
		pingTime = 5
		if !wasOffline {
			embed := disgord.Embed{Title: "Offline", Color: 0xb00000,
				Description: "The Minecraft server may be offline"}
			wasOffline = true
			info.Session.SendMsg(info.Con, pingChannel.ID, embed)
		}
	}

	if err == nil {

		var lastData networking.PingData

		utils.ShowSuccess(info, channelID, "Started listening on `"+ip+"`")

		for {
			conn, err := net.Dial("tcp", ip)

			if err == nil {
				data, err := server.Ping(conn)

				if err == nil {

					//check if channel is alive
					_, err := info.Session.GetChannel(info.Con, pingChannel.ID)
					if err != nil {
						return
					}

					datastring = "(" + strconv.Itoa(data.Players.Online) + "/" + strconv.Itoa(data.Players.Max) + ")"
					if wasOffline {
						online()
						wasOffline = false
						pingTime = 5
						lastData = data
						continue
					}

					if data.Players.Online > lastData.Players.Online {
						join()
					} else if data.Players.Online < lastData.Players.Online {
						leave()
					}
					lastData = data

				} else {
					offline()
				}
			} else {
				offline()
			}
			time.Sleep(pingTime * time.Second)
		}
	}
}
