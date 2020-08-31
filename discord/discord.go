package discord

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	"yey007.github.io/software/pingmc/networking"
	"yey007.github.io/software/pingmc/networking/bukkit"
	"yey007.github.io/software/pingmc/networking/forge"
)

//Init inititalizes the bot
func Init(token string) {
	client := disgord.New(disgord.Config{
		BotToken: token,
	})

	if client == nil {
		fmt.Println("Unable to create bot with the given token.")
		return
	}
	fmt.Println("PingMC v0.1 running", client)

	defer client.StayConnectedUntilInterrupted(context.Background())

	client.On(disgord.EvtMessageCreate, onMessageCreate)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
}

func onMessageCreate(session disgord.Session, evt *disgord.MessageCreate) {
	con := context.Background()
	msg := evt.Message
	args := strings.Split(msg.Content, " ")

	if len(args) >= 2 && args[0] == ".pingmc" {
		if len(args) == 5 && args[1] == "ping" {
			go onPingRequest(con, session, msg, args)
		} else {
			go showUsage(con, session, msg)
		}
	}
}

func onPingRequest(con context.Context, session disgord.Session, msg *disgord.Message, args []string) {

	ip := args[2]
	serverType := args[3]
	channelName := args[4]

	guild, err := session.GetGuild(con, msg.GuildID)
	if err != nil {
		return
	}

	roles, err := guild.RoleByName("@everyone")
	if err != nil {
		return
	}

	permissions := disgord.PermissionOverwrite{ID: roles[0].ID, Type: "role",
		Deny: disgord.PermissionSendMessages}

	params := disgord.CreateGuildChannelParams{Type: disgord.ChannelTypeGuildText,
		PermissionOverwrites: []disgord.PermissionOverwrite{permissions}}

	channel, err := session.CreateGuildChannel(con, msg.GuildID, channelName, &params)
	if err != nil {
		return
	}

	var server networking.McServer
	if serverType == "bukkit" {
		server = new(bukkit.Server)
	} else if serverType == "forge" {
		server = new(forge.Server)
	} else {
		msg.Reply(con, session, "Incorrect server type")
	}

	if err == nil {

		var lastData networking.PingData
		var pingTime time.Duration = 5
		msg.Reply(con, session, "Started listening on `"+ip+"`")

		for {
			conn, err := net.Dial("tcp", ip)

			if err == nil {
				data, err := server.Ping(conn)

				if err == nil {

					datastring := "(" + strconv.Itoa(data.Play.Online) + "/" + strconv.Itoa(data.Play.Max) + ")"

					if pingTime == 40 {
						embed := disgord.Embed{Title: "Online " + datastring, Color: 0x00ad37,
							Description: "The Minecraft server came online"}
						session.SendMsg(con, channel.ID, embed)
					}
					pingTime = 5

					if data.Play.Online > lastData.Play.Online {
						embed := disgord.Embed{Title: "Join " + datastring, Color: 0x1a5fba,
							Description: "Someone entered the Minecraft server"}
						session.SendMsg(con, channel.ID, embed)
					} else if data.Play.Online < lastData.Play.Online {

						embed := disgord.Embed{Title: "Leave " + datastring, Color: 0xe8bb35,
							Description: "Someone left the Minecraft server"}
						session.SendMsg(con, channel.ID, embed)
					}
					lastData = data

				} else {
					pingTime = 40
					embed := disgord.Embed{Title: "Offline", Color: 0xb00000,
						Description: "The Minecraft server may be offline"}
					session.SendMsg(con, channel.ID, embed)
				}
			} else {
				pingTime = 40
				embed := disgord.Embed{Title: "Offline", Color: 0xb00000,
					Description: "The Minecraft server may be offline"}
				session.SendMsg(con, channel.ID, embed)
			}
			fmt.Println("Ping!")
			time.Sleep(pingTime * time.Second)
		}
	}
}

func showUsage(con context.Context, session disgord.Session, msg *disgord.Message) {
	msg.Reply(con, session, "```Usage: \n"+
		".pingmc help - displays this message \n"+
		".pingmc ping channelName IP```")
}
