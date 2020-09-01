package vanilla

import (
	"encoding/json"
	"errors"
	"net"
	"strconv"
	"strings"

	"yey007.github.io/software/pingmc/networking"
)

func sendHandshake(conn net.Conn) error {

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

	data = networking.WriteVarInt(packetID, data)
	data = networking.WriteVarInt(protocolVersion, data)
	data = networking.WriteString(serverAddress, data)
	data = networking.WriteShort(serverPort, data)
	data = networking.WriteVarInt(nextState, data)

	length := len(data)
	lengthData := make([]byte, 0)
	lengthData = networking.WriteVarInt(length, lengthData)

	_, err := conn.Write(lengthData)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}

func sendRequest(conn net.Conn) error {
	data := make([]byte, 0)
	packetID := 0x00
	data = networking.WriteVarInt(packetID, data)

	length := len(data)
	lengthData := make([]byte, 0)
	lengthData = networking.WriteVarInt(length, lengthData)

	_, err := conn.Write(lengthData)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}

func recieveResponse(conn net.Conn) (networking.PingData, error) {
	networking.ReadVarInt(conn)
	networking.ReadVarInt(conn)
	length, err := networking.ReadVarInt(conn)
	if err != nil {
		return networking.PingData{}, err
	}

	readBuf := make([]byte, length)
	readCount, err := conn.Read(readBuf)
	if err != nil {
		return networking.PingData{}, err
	}

	var data networking.PingData
	err = json.Unmarshal(readBuf, &data)
	if err != nil {
		return networking.PingData{}, err
	}

	if readCount == length && err == nil {
		return data, err
	}
	return networking.PingData{}, errors.New("Something went wrong")
}
