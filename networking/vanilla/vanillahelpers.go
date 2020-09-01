package vanilla

import (
	"encoding/json"
	"errors"
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
