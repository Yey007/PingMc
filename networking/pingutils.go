package networking

import (
	"encoding/binary"
	"fmt"
	"net"
)

//ReadVarInt reads a variable length integer from the given connection
func ReadVarInt(conn net.Conn) (int, error) {

	var err error

	numRead := 0
	result := 0
	readBuf := make([]byte, 1)

	cycle := func() {
		readCount, err := conn.Read(readBuf)
		if err == nil && readCount == 1 {
			var value int = int(readBuf[0] & 0b01111111)
			result |= value << (7 * numRead)
			numRead++

			if numRead > 5 {
				fmt.Println("VarInt from server is too large!")
			}

		}
	}

	cycle()
	for (readBuf[0] & 0b10000000) != 0 {
		cycle()
	}

	return result, err
}

//ReadString reads a variable length string from the given connection
func ReadString(conn net.Conn) (string, error) {
	length, err := ReadVarInt(conn)

	var readBuf []byte
	var readCount int

	if err == nil {
		readBuf = make([]byte, length)
		readCount, err = conn.Read(readBuf)
	}

	if readCount == length && err == nil {
		return string(readBuf), err
	}

	return "", err
}

//WriteVarInt encodes an integer into a byte array which can be sent over a connection
func WriteVarInt(value int, data []byte) []byte {

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

//WriteString encodes a string into a byte array which can be sent over a connection
func WriteString(str string, data []byte) []byte {
	temp := []byte(str)
	data = WriteVarInt(len(temp), data)
	data = append(data, temp...)
	return data
}

//WriteShort encodes a short into a byte array which can be sent over a connection
func WriteShort(short uint16, data []byte) []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, short)
	data = append(data, buf...)
	return data
}

//sendRequest sends an Server List Ping request packet to a server
func sendRequest(conn net.Conn) error {
	data := make([]byte, 0)
	packetID := 0x00
	data = WriteVarInt(packetID, data)

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
