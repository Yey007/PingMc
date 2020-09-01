package fml1

import (
	"net"

	"yey007.github.io/software/pingmc/networking"
)

//Server represents a forge server
type Server struct {
}

//Ping pings a forge server for player data
func (s *Server) Ping(conn net.Conn) (networking.PingData, error) {
	err := networking.SendHandshake(conn)
	err = networking.SendRequest(conn)
	response, err := recieveResponse(conn)
	return response, err
}
