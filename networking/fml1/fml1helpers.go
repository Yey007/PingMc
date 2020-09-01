package fml1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"yey007.github.io/software/pingmc/networking"
)

func recieveResponse(conn net.Conn) (networking.PingData, error) {
	networking.ReadVarInt(conn)
	networking.ReadVarInt(conn)
	length, err := networking.ReadVarInt(conn)
	if err != nil {
		return networking.PingData{}, err
	}

	readBuf := make([]byte, length)
	readCount := 0

	//For some reason, it doesn't work correctly unless I do this garbadge
	for readCount < length {
		small := make([]byte, 1)
		conn.Read(small)
		if !(small[0] == 0) {
			readBuf = append(readBuf, small...)
			readCount++
		}
	}

	readBuf = bytes.Trim(readBuf, "\x00")

	var data networking.PingData
	err = json.Unmarshal(readBuf, &data)
	if err != nil {
		return networking.PingData{}, err
	}

	fmt.Println(data)

	if readCount == length && err == nil {
		return data, err
	}
	return networking.PingData{}, errors.New("Something went wrong")
}

func remove(s []byte, i int) []byte {
	return append(s[:i], s[i+1:]...)
}
