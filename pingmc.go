package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"yey007.github.io/software/pingmc/discord"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	discord.Start(os.Getenv("TOKEN"))
}
