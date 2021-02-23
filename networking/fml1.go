package networking

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"strconv"
	"strings"
)

//forgeReceiveResponse receives a response from a forge server
func forgeReceiveResponse(conn net.Conn) (PingData, error) {
	ReadVarInt(conn)
	ReadVarInt(conn)
	length, err := ReadVarInt(conn)
	if err != nil {
		return PingData{}, err
	}

	readBuf := make([]byte, length)
	readCount := 0

	//For some reason, it doesn't work correctly unless I do this garbage
	for readCount < length {
		small := make([]byte, 1)
		conn.Read(small)
		if !(small[0] == 0) {
			readBuf = append(readBuf, small...)
			readCount++
		}
	}

	readBuf = bytes.Trim(readBuf, "\x00")

	var data PingData
	err = json.Unmarshal(readBuf, &data)
	if err != nil {
		return PingData{}, err
	}

	if readCount == length {
		return data, err
	}
	return PingData{}, errors.New("something went wrong")
}

//forgeSendHandshake sends a forge handshake packet to a forge server
func forgeSendHandshake(conn net.Conn) error {

	data := make([]byte, 0)
	packetID := 0x00
	protocolVersion := 736

	temp := strings.Split(conn.RemoteAddr().String(), ":")

	serverAddress := temp[0] + "\000FML\000"
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
