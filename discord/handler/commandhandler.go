package handler

import (
	"context"
	"strings"
	"sync"
	"time"

	"yey007.github.io/software/pingmc/discord/commands"
	"yey007.github.io/software/pingmc/discord/utils"

	"github.com/andersfylling/disgord"
)

type CommandMap = map[string]commands.Command

type CommandHandler struct {
	commandMap CommandMap
	lock       sync.RWMutex
}

//New returns a new CommandHandler with an empty command set
func New() CommandHandler {
	return CommandHandler{commandMap: CommandMap{}}
}

//HandleEventAsync handles a message event in another goroutine
func (h *CommandHandler) HandleEventAsync(session disgord.Session, evt *disgord.MessageCreate) {
	go h.handleEvent(session, evt)
}

//handleEvent handles a message send event from discord
func (h *CommandHandler) handleEvent(session disgord.Session, evt *disgord.MessageCreate) {
	msg := evt.Message
	args := strings.Split(msg.Content, " ")

	h.lock.RLock()
	cmd, ok := h.commandMap[args[0]]
	h.lock.RUnlock()
	if !ok {
		session.SendMsg(evt.Ctx, msg.ChannelID, utils.CreateError("That command does not exist."))
		return
	}
	if len(args)-1 == len(cmd.Help().Arguments) {
		ctx, cancel := context.WithTimeout(evt.Ctx, 3*time.Second)
		defer cancel()

		embed := cmd.Run(ctx, msg, args[1:])
		session.SendMsg(ctx, msg.ChannelID, embed)
	} else {
		session.SendMsg(evt.Ctx, msg.ChannelID, utils.CreateError("Wrong number of arguments."))
	}
}

func (h *CommandHandler) Handle(cmd commands.Command) {
	h.lock.Lock()
	h.commandMap[cmd.Name()] = cmd
	h.lock.Unlock()
}

func (h *CommandHandler) Commands() CommandMap {
	return h.commandMap
}
