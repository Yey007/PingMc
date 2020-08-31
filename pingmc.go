package main

import (
	"flag"

	"yey007.github.io/software/pingmc/discord"
)

var (
	token string
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	discord.Init(token)
}
