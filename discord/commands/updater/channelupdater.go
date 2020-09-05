package updater

import (
	"errors"
	"net"
	"strconv"
	"strings"
	"time"

	"yey007.github.io/software/pingmc/discord/utils/colors"
	"yey007.github.io/software/pingmc/networking/fml1"
	"yey007.github.io/software/pingmc/networking/vanilla"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
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
		return
	}

	server, err := createServer(info, serverType)
	cancellationChan := addNewPing(msg.ChannelID)
	if cancellationChan == nil {
		utils.ShowError(info, channelID, "A ping is already running in this channel.")
		return
	}

	datastring := "(0/0)"
	var pingTime time.Duration = 5
	wasOffline := true

	online := func() {
		embed := disgord.Embed{Title: "Online " + datastring, Color: colors.Green,
			Description: "The Minecraft server came online"}
		info.Session.SendMsg(info.Con, pingChannel.ID, embed)
	}

	join := func() {
		embed := disgord.Embed{Title: "Join " + datastring, Color: colors.Blue,
			Description: "Someone entered the Minecraft server"}
		info.Session.SendMsg(info.Con, pingChannel.ID, embed)
	}

	leave := func() {
		embed := disgord.Embed{Title: "Leave " + datastring, Color: colors.Orange,
			Description: "Someone left the Minecraft server"}
		info.Session.SendMsg(info.Con, pingChannel.ID, embed)
	}

	offline := func() {
		pingTime = 5
		if !wasOffline {
			embed := disgord.Embed{Title: "Offline", Color: colors.Red,
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

					//check if ping has been cancelled
					select {
					case _, ok := <-(*cancellationChan):
						if ok {
							//there was something on the channel. WTF?
						} else {
							//the channel was closed, which means we have to stop pinging.
							return
						}
					default:
						//nothing on the channel, keep going
					}

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

func validateInputs(info utils.SessionInfo, channelID snowflake.Snowflake, args []string) bool {

	for i := range args {
		args[i] = strings.Trim(args[i], "\r\n ")
	}

	if len(args) != 4 {
		utils.ShowError(info, channelID, "Wrong number of arguments")
		return false
	}

	result := true
	ip := args[2]
	serverType := args[3]

	ip = strings.Split(ip, ":")[0]

	if net.ParseIP(ip) == nil {
		result = result && false
		utils.ShowError(info, channelID, "Invalid IP address")
	}

	if serverType != "vanilla" && serverType != "forge1" {
		result = result && false
		utils.ShowError(info, channelID, "Invalid server type")
	}
	return result
}

func createServer(info utils.SessionInfo, serverType string) (networking.McServer, error) {
	if serverType == "vanilla" {
		return new(vanilla.Server), nil
	} else if serverType == "forge1" {
		return new(fml1.Server), nil
	}
	return nil, errors.New("Incorrect server type. How did this happen?")
}
