package networking

//PingData represents the most basic data returned from a ping
type PingData struct {
	Players Players `json:"players"`
}

//Players represents the player counts on the server
type Players struct {
	Max    int `json:"max"`
	Online int `json:"online"`
}
