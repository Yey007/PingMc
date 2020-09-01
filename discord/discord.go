package discord

import (
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
)

//Init inititalizes the bot
func Init(token string) {
	client := disgord.New(disgord.Config{
		BotToken: token,
	})

	if client == nil {
		fmt.Println("Unable to create bot with the given token.")
		return
	}
	fmt.Println("PingMC v0.1 running", client)

	defer client.StayConnectedUntilInterrupted(context.Background())

	client.On(disgord.EvtMessageCreate, onMessageCreate)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
}
