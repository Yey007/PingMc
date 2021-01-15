package commands

import (
	"context"

	"github.com/andersfylling/disgord"
)

type Command interface {
	Name() string
	Help() CommandHelp
	Run(context.Context, *disgord.Message, []string) disgord.Embed
}
