package cancelping

import (
	"context"
	"fmt"

	"yey007.github.io/software/pingmc/discord/utils"

	"github.com/andersfylling/disgord"
	"yey007.github.io/software/pingmc/data"
	"yey007.github.io/software/pingmc/discord/commands"
)

type cancelPingCommand struct {
	name string
	help commands.CommandHelp
	repo *data.PingRepo
}

func New(repo *data.PingRepo) commands.Command {
	command := cancelPingCommand{
		name: "cancelping",
		help: commands.CommandHelp{
			Description: "remove/cancel a ping",
			Arguments: []commands.CommandArgument{
				{
					Name:        "pingname",
					Description: "The name you gave this ping",
				},
			},
		},
		repo: repo,
	}
	command.help.Usage = commands.BuildUsage(command.name, command.help.Arguments)
	return command
}

func (c cancelPingCommand) Name() string {
	return c.name
}

func (c cancelPingCommand) Help() commands.CommandHelp {
	return c.help
}

func (c cancelPingCommand) Run(ctx context.Context, msg *disgord.Message, args []string) disgord.Embed {
	toBeDeleted, err := c.repo.GetByNameGuildID(ctx, args[0], msg.GuildID.String())
	if err != nil {
		return utils.CreateError("That ping doesn't appear to exist.")
	}
	err = c.repo.Remove(ctx, toBeDeleted.ID)
	if err != nil {
		return utils.CreateError("Something went wrong deleting the ping. Please report this issue.")
	}
	return utils.CreateSuccess(fmt.Sprintf("Deleted ping with name %v", utils.Block(args[0])))
}
