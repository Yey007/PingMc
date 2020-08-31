package bukkit

import (
	"net"

	"yey007.github.io/software/pingmc/networking"
)

//Server represents a bukkit server
type Server struct {
}

//Ping pings a bukkit server for player data
func (b *Server) Ping(conn net.Conn) (networking.PingData, error) {

	err := sendHandshake(conn)
	err = sendRequest(conn)
	response, err := recieveResponse(conn)
	return response, err
}
