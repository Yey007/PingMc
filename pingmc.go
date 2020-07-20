package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	token string
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Unable to create bot with the given token.")
		fmt.Println(err)
		return
	}

	discord.AddHandler(messageCreate)

	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	err = discord.Open()
	if err != nil {
		fmt.Println("Unable to open a connection to discord")
		fmt.Println(err)
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	args := strings.Split(m.Content, " ")
	length := len(args)
	usage := "```Usage:\n" +
		".pingmc help - displays this message \n" +
		".pingmc start channelName IP secondsBetweenPings - starts pinging an ip for player counts\n" +
		"There is no stop command, just delete the channel.```"

	//base command
	if args[0] == ".pingmc" {
		if length >= 2 {
			if args[1] == "start" {
				if length == 5 {

					ch, err := s.GuildChannelCreate(m.GuildID, args[2]+" 0/0", discordgo.ChannelTypeGuildVoice)
					temp := strings.Split(args[3], ":")
					pingString := "Started pinging `" + temp[0] + "`"
					if len(temp) == 2 {
						pingString += " on `" + temp[1] + "`"
					}
					s.ChannelMessageSend(m.ChannelID, pingString)
					if err == nil {
						for true {
							err := update(ch.ID, args[3], args[2], s)
							if err != nil {
								fmt.Println(err)
								stopString := "Stopped pinging `" + temp[0] + "`"
								if len(temp) == 2 {
									stopString += " on `" + temp[1] + "`"
								}
								s.ChannelMessageSend(m.ChannelID, stopString)
								return
							}
							sleepTime, err := strconv.Atoi(args[4])
							if err != nil {
								sleepTime = 20
							}
							time.Sleep(time.Duration(sleepTime) * time.Second)
						}
					}

				} else {
					//wrong number of arguments - display usage
					s.ChannelMessageSend(m.ChannelID, "Wrong number of arguments!")
					s.ChannelMessageSend(m.ChannelID, usage)
				}
			} else if args[1] == "help" {
				//display usage
				s.ChannelMessageSend(m.ChannelID, usage)
			} else {
				//non existent sub command - display usage
				s.ChannelMessageSend(m.ChannelID, "That command does not exist!")
				s.ChannelMessageSend(m.ChannelID, usage)
			}
		} else {
			//not enough arguments
			s.ChannelMessageSend(m.ChannelID, "Wrong number of arguments!")
			s.ChannelMessageSend(m.ChannelID, usage)
		}
	}
}

func update(channelID, ip, channelName string, s *discordgo.Session) error {
	conn, err := net.Dial("tcp", ip)
	if err == nil {
		data := ping(conn)
		_, err = s.ChannelEdit(channelID, (channelName + " " + strconv.Itoa(data.Play.Online) + "/" + strconv.Itoa(data.Play.Max)))
		if err != nil {
			conn.Close()
		}
		return err //we must notify the calling function as well
	}
	if conn != nil {
		conn.Close()
	}
	return err
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
