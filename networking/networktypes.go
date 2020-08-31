package networking

//PingData represents the most basic data returned from a ping
type PingData struct {
	Des  description `json:"description"`
	Play players     `json:"players"`
	Ver  version     `json:"version"`
}

type description struct {
	Text string `json:"text"`
}

type players struct {
	Max    int `json:"max"`
	Online int `json:"online"`
}

type version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}
