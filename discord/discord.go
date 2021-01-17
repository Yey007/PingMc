package discord

import (
	"context"
	"fmt"

	"yey007.github.io/software/pingmc/discord/commands/pinglist"

	"yey007.github.io/software/pingmc/discord/commands/cancelping"

	"yey007.github.io/software/pingmc/data"

	"yey007.github.io/software/pingmc/discord/utils"

	"yey007.github.io/software/pingmc/discord/commands/ping"

	"yey007.github.io/software/pingmc/discord/commands/allhelp"
	"yey007.github.io/software/pingmc/discord/commands/help"
	"yey007.github.io/software/pingmc/discord/handler"

	"yey007.github.io/software/pingmc/discord/notifier"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
)

//Start initializes the bot
func Start(token string) {
	client := disgord.New(disgord.Config{
		BotToken: token,
	})

	fmt.Println("PingMC v0.2 running", client)

	h := handler.New()
	r, err := data.NewPingRepo()
	if err != nil {
		panic(err)
	}

	h.Handle(help.New(h.Commands()))
	h.Handle(allhelp.New(h.Commands()))
	h.Handle(ping.New(r))
	h.Handle(cancelping.New(r))
	h.Handle(pinglist.New(r))

	filter, _ := std.NewMsgFilter(context.Background(), client)
	filter.SetPrefix(utils.PREFIX + " ")

	client.On(disgord.EvtMessageCreate,
		filter.NotByBot,
		filter.HasPrefix,
		std.CopyMsgEvt,
		filter.StripPrefix,
		h.HandleEventAsync)

	client.On(disgord.EvtReady, func(session disgord.Session, event *disgord.Ready) {
		notifier.StartNotifier(event.Ctx, session, r)
	})

	defer client.StayConnectedUntilInterrupted(context.Background())

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
}
