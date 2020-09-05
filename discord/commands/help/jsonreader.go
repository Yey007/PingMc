package help

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type command struct {
	Cmd  string     `json:"cmd"`
	Desc string     `json:"desc"`
	Args []argument `json:"args"`
}

type argument struct {
	Arg  string `json:"arg"`
	Desc string `json:"desc"`
}

var doOnce sync.Once
var constmap map[string]command

func readJSON() (map[string]command, error) {

	var err error

	doOnce.Do(func() {
		var helpdata []byte
		helpdata, err = ioutil.ReadFile("data/helpmap.json")
		err = json.Unmarshal(helpdata, &constmap)
	})

	return constmap, err
}
