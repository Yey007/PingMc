package networking

import (
	"errors"
	"net"

	"gorm.io/gorm"
)

const (
	ServerTypeVanilla = iota
	ServerTypeForge
)

var (
	ErrUnknownServerType = errors.New("unknown server type identifier")
)

//McServer represents a remote minecraft server which can be pinged
type McServer struct {
	gorm.Model
	Address         string
	Type            uint8
	RecurringPingID uint
}

func (m *McServer) Ping() (*PingData, error) {
	conn, err := net.Dial("tcp", m.Address)
	if err != nil {
		return nil, err
	}
	if m.Type == ServerTypeForge {
		return pingForge(conn)
	} else if m.Type == ServerTypeVanilla {
		return pingVanilla(conn)
	} else {
		return nil, ErrUnknownServerType
	}
}

func pingForge(conn net.Conn) (*PingData, error) {
	defer conn.Close()
	err := forgeSendHandshake(conn)
	err = sendRequest(conn)
	response, err := forgeReceiveResponse(conn)
	return response, err
}

func pingVanilla(conn net.Conn) (*PingData, error) {
	defer conn.Close()
	err := vanillaSendHandshake(conn)
	err = sendRequest(conn)
	response, err := vanillaReceiveResponse(conn)
	return response, err
}
