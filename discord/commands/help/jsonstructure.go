package help

type command struct {
	Cmd  string     `json:"cmd"`
	Desc string     `json:"desc"`
	Args []argument `json:"args"`
}

type argument struct {
	Arg  string `json:"arg"`
	Desc string `json:"desc"`
}
