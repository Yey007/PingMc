package pinglist

import (
	"context"
	"fmt"

	"yey007.github.io/software/pingmc/networking"

	"yey007.github.io/software/pingmc/discord/utils"

	"github.com/andersfylling/disgord"
	"yey007.github.io/software/pingmc/data"
	"yey007.github.io/software/pingmc/discord/commands"
)

type pingListCommand struct {
	name string
	help commands.CommandHelp
	repo *data.PingRepo
}

func New(repo *data.PingRepo) commands.Command {
	command := pingListCommand{
		name: "pinglist",
		help: commands.CommandHelp{
			Description: "returns a list of pings for this server",
			Arguments:   []commands.CommandArgument{},
		},
		repo: repo,
	}
	command.help.Usage = commands.BuildUsage(command.name, command.help.Arguments)
	return command
}

func (p pingListCommand) Name() string {
	return p.name
}

func (p pingListCommand) Help() commands.CommandHelp {
	return p.help
}

func (p pingListCommand) Run(ctx context.Context, msg *disgord.Message, args []string) disgord.Embed {
	pings, err := p.repo.GetByGuildID(ctx, msg.GuildID.String())
	if err != nil {
		return disgord.Embed{Title: "Pings", Description: "Nothing to see here!", Color: utils.BLUE}
	}

	// Build description
	desc := ""
	for _, ping := range pings {

		typeString := ""
		if ping.Server.Type == networking.ServerTypeVanilla {
			typeString = "vanilla"
		} else {
			typeString = "forge"
		}

		desc += fmt.Sprintf(
			"%v to %v with server type %v\n",
			utils.Block(ping.Name),
			utils.Block(ping.Server.Address),
			utils.Block(typeString),
		)
	}

	//TODO: Weird formatting
	return disgord.Embed{Title: "Pings", Description: desc, Color: utils.BLUE}
}
