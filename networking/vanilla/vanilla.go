package vanilla

import (
	"net"

	"yey007.github.io/software/pingmc/networking"
)

//Server represents a vanilla server
type Server struct {
}

//Ping pings a vanilla server for player data
func (s *Server) Ping(conn net.Conn) (networking.PingData, error) {

	err := sendHandshake(conn)
	err = networking.SendRequest(conn)
	response, err := recieveResponse(conn)
	return response, err
}
