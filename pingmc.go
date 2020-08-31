package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
)

var (
	token string
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

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
		if len(args) == 4 && args[1] == "ping" {
			go onPingRequest(con, session, msg, args)
		} else {
			go showUsage(con, session, msg)
		}
	}
}

func onPingRequest(con context.Context, session disgord.Session, msg *disgord.Message, args []string) {
	/*
		guild, _ := session.GetGuild(con, msg.GuildID)
		roles, _ := guild.RoleByName("@everyone")
		permissions := disgord.PermissionOverwrite{ID: roles[0].ID, Type: "role", Deny: disgord.PermissionVoiceConnect}
		params := disgord.CreateGuildChannelParams{Type: disgord.ChannelTypeGuildVoice, PermissionOverwrites: []disgord.PermissionOverwrite{permissions}}
		channel, err := session.CreateGuildChannel(con, msg.GuildID, args[2], &params)
	*/
	guild, _ := session.GetGuild(con, msg.GuildID)
	roles, _ := guild.RoleByName("@everyone")
	permissions := disgord.PermissionOverwrite{ID: roles[0].ID, Type: "role", Deny: disgord.PermissionSendMessages}
	params := disgord.CreateGuildChannelParams{Type: disgord.ChannelTypeGuildText, PermissionOverwrites: []disgord.PermissionOverwrite{permissions}}
	channel, err := session.CreateGuildChannel(con, msg.GuildID, args[2], &params)

	if err == nil {
		msg.Reply(con, session, "Started listening on `"+args[3]+"`")
		var lastData pingData
		for {
			time.Sleep(5 * time.Second)
			conn, err := net.Dial("tcp", args[3])
			if err == nil {
				data := ping(conn)
				datastring := "(" + strconv.Itoa(data.Play.Online) + "/" + strconv.Itoa(data.Play.Max) + ")"
				if data.Play.Online > lastData.Play.Online {
					embed := disgord.Embed{Title: "Join " + datastring, Color: 0x00ad37, Description: "Someone entered the Minecraft server"}
					session.SendMsg(con, channel.ID, embed)
				} else if data.Play.Online < lastData.Play.Online {
					embed := disgord.Embed{Title: "Join " + datastring, Color: 0xb00000, Description: "Someone left the Minecraft server"}
					session.SendMsg(con, channel.ID, embed)
				}

				lastData = data
			}
		}
	}
}

func showUsage(con context.Context, session disgord.Session, msg *disgord.Message) {
	msg.Reply(con, session, "```Usage: \n"+
		".pingmc help - displays this message \n"+
		".pingmc ping channelName IP```")
}

type pingData struct {
	Des  description `json:"description"`
	Play players     `json:"players"`
	Ver  version     `json:"version"`
}

type description struct {
	Text string `json:"text"`
}

type players struct {
	Max    int `json:"max"`
	Online int `json:"online"`
}

type version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

func ping(conn net.Conn) pingData {

	sendHandshake(conn)
	sendRequest(conn)
	response := recieveResponse(conn)
	return response
}

func sendHandshake(conn net.Conn) {

	data := make([]byte, 0)
	packetID := 0x00
	protocolVersion := 736

	temp := strings.Split(conn.RemoteAddr().String(), ":")

	serverAddress := temp[0]
	var serverPort uint16

	if len(temp) == 2 {
		var err error
		var i int
		i, err = strconv.Atoi(temp[1])
		if err != nil {
			serverPort = 25565
		} else {
			serverPort = uint16(i)
		}
	} else {
		serverPort = 25565
	}

	nextState := 1

	data = writeVarInt(packetID, data)
	data = writeVarInt(protocolVersion, data)
	data = writeString(serverAddress, data)
	data = writeShort(serverPort, data)
	data = writeVarInt(nextState, data)

	length := len(data)
	lengthData := make([]byte, 0)
	lengthData = writeVarInt(length, lengthData)

	conn.Write(lengthData)
	conn.Write(data)
}

func sendRequest(conn net.Conn) {
	data := make([]byte, 0)
	packetID := 0x00
	data = writeVarInt(packetID, data)

	length := len(data)
	lengthData := make([]byte, 0)
	lengthData = writeVarInt(length, lengthData)

	conn.Write(lengthData)
	conn.Write(data)
}

func recieveResponse(conn net.Conn) pingData {
	readVarInt(conn)
	readVarInt(conn)
	length := readVarInt(conn)
	readBuf := make([]byte, length)
	readCount, err := conn.Read(readBuf)
	var data pingData
	json.Unmarshal(readBuf, &data)

	if readCount == length && err == nil {
		return data
	}
	fmt.Println("Error recieving response from server.")
	fmt.Println(err)
	return data
}

func readVarInt(conn net.Conn) int {
	numRead := 0
	result := 0
	readBuf := make([]byte, 1)

	cycle := func() {
		readCount, err := conn.Read(readBuf)
		if err == nil && readCount == 1 {
			var value int = int(readBuf[0] & 0b01111111)
			result |= (value << (7 * numRead))
			numRead++

			if numRead > 5 {
				fmt.Println("VarInt from server is too large!")
			}

		} else {
			fmt.Println("Error receiving VarInt from server.")
			fmt.Println(err)
		}
	}

	cycle()
	for (readBuf[0] & 0b10000000) != 0 {
		cycle()
	}

	return result
}

func readString(conn net.Conn) string {
	length := readVarInt(conn)
	readBuf := make([]byte, length)
	readCount, err := conn.Read(readBuf)

	if readCount == length && err == nil {
		return string(readBuf)
	}

	fmt.Println("Error recieving string from server")
	fmt.Println(err)
	return ""
}

func writeVarInt(value int, data []byte) []byte {

	cycle := func() {
		temp := byte(value & 0b01111111)
		value = value >> 7
		if value != 0 {
			temp = (temp | 0b10000000)
		}
		data = append(data, temp)
	}

	cycle()
	for value != 0 {
		cycle()
	}
	return data
}

func writeString(str string, data []byte) []byte {
	temp := []byte(str)
	data = writeVarInt(len(temp), data)
	data = append(data, temp...)
	return data
}

func writeShort(short uint16, data []byte) []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, short)
	data = append(data, buf...)
	return data
}
