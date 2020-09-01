package discord

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
	"yey007.github.io/software/pingmc/networking"
	"yey007.github.io/software/pingmc/networking/fml1"
	"yey007.github.io/software/pingmc/networking/vanilla"
)

func onPingRequest(info sessionInfo, guildID snowflake.Snowflake, channelID snowflake.Snowflake, args []string) {

	ip := args[2]
	serverType := args[3]
	channelName := args[4]

	if validateInputs(info, channelID, args) == false {
		return
	}

	pingChannel, err := createPingChannel(info, guildID, channelName)
	if err != nil {
		showError(info, channelID, "Unable to create channel. Reason is unknown.")
	}

	server, err := createServer(info, serverType)

	datastring := "(0/0)"
	var pingTime time.Duration = 5
	wasOffline := true

	online := func() error {
		embed := disgord.Embed{Title: "Online " + datastring, Color: 0x00ad37,
			Description: "The Minecraft server came online"}
		_, err := info.session.SendMsg(info.con, pingChannel.ID, embed)
		return err
	}

	join := func() error {
		embed := disgord.Embed{Title: "Join " + datastring, Color: 0x1a5fba,
			Description: "Someone entered the Minecraft server"}
		_, err := info.session.SendMsg(info.con, pingChannel.ID, embed)
		return err
	}

	leave := func() error {
		embed := disgord.Embed{Title: "Leave " + datastring, Color: 0xe8bb35,
			Description: "Someone left the Minecraft server"}
		_, err := info.session.SendMsg(info.con, pingChannel.ID, embed)
		return err
	}

	offline := func() error {
		pingTime = 5
		if !wasOffline {
			embed := disgord.Embed{Title: "Offline", Color: 0xb00000,
				Description: "The Minecraft server may be offline"}
			wasOffline = true
			_, err := info.session.SendMsg(info.con, pingChannel.ID, embed)
			return err
		}
		return nil
	}

	if err == nil {

		var lastData networking.PingData

		showSuccess(info, channelID, "Started listening on `"+ip+"`")

		for {
			conn, err := net.Dial("tcp", ip)

			if err == nil {
				data, err := server.Ping(conn)

				if err == nil {

					datastring = "(" + strconv.Itoa(data.Players.Online) + "/" + strconv.Itoa(data.Players.Max) + ")"
					if wasOffline {
						if online() != nil {
							return
						}
					}

					wasOffline = false
					pingTime = 5

					if data.Players.Online > lastData.Players.Online {
						if join() != nil {
							return
						}
					} else if data.Players.Online < lastData.Players.Online {
						if leave() != nil {
							return
						}
					}
					lastData = data

				} else {
					fmt.Println(err)
					if offline() != nil {
						return
					}
				}
			} else {
				fmt.Println(err)
				if offline() != nil {
					return
				}
			}
			fmt.Println("ping")
			time.Sleep(pingTime * time.Second)
		}
	}
}

func validateInputs(info sessionInfo, channelID snowflake.Snowflake, args []string) bool {

	for i := range args {
		args[i] = strings.Trim(args[i], "\r\n ")
	}

	result := true
	ip := args[2]
	serverType := args[3]

	ip = strings.Split(ip, ":")[0]

	if net.ParseIP(ip) == nil {
		result = result && false
		showError(info, channelID, "Invalid IP address")
	}

	if serverType != "vanilla" && serverType != "forge1" {
		fmt.Println(serverType)
		result = result && false
		showError(info, channelID, "Invalid server type")
	}
	return result
}

func createPingChannel(info sessionInfo, guildID snowflake.Snowflake, channelName string) (*disgord.Channel, error) {
	guild, err := info.session.GetGuild(info.con, guildID)
	if err != nil {
		return &disgord.Channel{}, err
	}

	roles, err := guild.RoleByName("@everyone")
	if err != nil {
		return &disgord.Channel{}, err
	}

	permissions := disgord.PermissionOverwrite{ID: roles[0].ID, Type: "role",
		Deny: disgord.PermissionSendMessages}

	params := disgord.CreateGuildChannelParams{Type: disgord.ChannelTypeGuildText,
		PermissionOverwrites: []disgord.PermissionOverwrite{permissions}}

	channel, err := info.session.CreateGuildChannel(info.con, guildID, channelName, &params)
	if err != nil {
		return &disgord.Channel{}, err
	}
	return channel, nil
}

func createServer(info sessionInfo, serverType string) (networking.McServer, error) {
	if serverType == "vanilla" {
		return new(vanilla.Server), nil
	} else if serverType == "forge1" {
		return new(fml1.Server), nil
	}
	return nil, errors.New("Incorrect server type")
}
