package networking

import "net"

//McServer represents a remote minecraft server which can be pinged
type McServer interface {
	Ping(conn net.Conn) (PingData, error)
}
