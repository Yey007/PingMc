package networking

import (
	"encoding/json"
	"errors"
	"net"
	"strconv"
	"strings"
)

//vanillaReceiveResponse receives a response from a vanilla server
func vanillaReceiveResponse(conn net.Conn) (*PingData, error) {
	// There is some data we don't need so we discard it (packet length)
	ReadVarInt(conn)
	ReadVarInt(conn)
	length, err := ReadVarInt(conn)
	if err != nil {
		return nil, err
	}

	readBuf := make([]byte, length)
	readCount, err := conn.Read(readBuf)
	if err != nil {
		return nil, err
	}

	var data PingData
	err = json.Unmarshal(readBuf, &data)
	if err != nil {
		return nil, err
	}

	//Error is nil by now
	if readCount == length {
		return &data, err
	}
	return nil, errors.New("something went wrong")
}

//vanillaSendHandshake sends a vanilla handshake packet to a vanilla server
func vanillaSendHandshake(conn net.Conn) error {

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

	data = WriteVarInt(packetID, data)
	data = WriteVarInt(protocolVersion, data)
	data = WriteString(serverAddress, data)
	data = WriteShort(serverPort, data)
	data = WriteVarInt(nextState, data)

	length := len(data)
	lengthData := make([]byte, 0)
	lengthData = WriteVarInt(length, lengthData)

	_, err := conn.Write(lengthData)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}
