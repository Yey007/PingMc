package networking

import (
	"encoding/json"
)

//PingData represents the most basic data returned from a ping
type PingData struct {
	Desc    json.RawMessage `json:"description"`
	Players players         `json:"players"`
	Version version         `json:"version"`
	ModInfo modinfo         `json:"modinfo"`
}

//forge
type modinfo struct {
	Type    string `json:"type"`
	ModList []mod  `json:"modList"`
}

type players struct {
	Max    int      `json:"max"`
	Online int      `json:"online"`
	Sample []player `json:"sample"` //vanilla
}

//Description represents a vanilla and fml2 description
type Description struct {
	Text string `json:"text"`
}

type version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type player struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

//forge
type mod struct {
	ModID   string `json:"modid"`
	Version string `json:"version"`
}
