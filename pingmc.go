package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"yey007.github.io/software/pingmc/data"

	"yey007.github.io/software/pingmc/discord"
)

func main() {
	if os.Getenv("CONTAINER") != "TRUE" {
		err := godotenv.Load("bot.env", "db.env")
		if err != nil {
			log.Fatal("Error loading env files")
		}
	}

	r, err := data.NewPingRepo(data.Config{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DATABASE"),
		Port:     os.Getenv("DB_PORT"),
	})

	if err != nil {
		panic(err)
	}

	discord.Start(os.Getenv("TOKEN"), r)
}
