package forge

import (
	"net"

	"yey007.github.io/software/pingmc/networking"
)

//Server represents a forge server
type Server struct {
}

//Ping pings a forge server for player data
func (b *Server) Ping(conn net.Conn) (networking.PingData, error) {
	return networking.PingData{}, nil
}
