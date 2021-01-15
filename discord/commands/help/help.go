package help

import (
	"context"

	"yey007.github.io/software/pingmc/discord/commands"
	"yey007.github.io/software/pingmc/discord/handler"

	"github.com/andersfylling/disgord"
	"yey007.github.io/software/pingmc/discord/utils"
)

type helpCommand struct {
	name       string
	commandMap handler.CommandMap
	help       commands.CommandHelp
}

//New returns a new commands.Command object that uses the given handler.CommandMap
func New(commandMap handler.CommandMap) commands.Command {
	command := helpCommand{
		name:       "help",
		commandMap: commandMap,
		help: commands.CommandHelp{
			Description: "Gives help about a command",
			Arguments: []commands.CommandArgument{
				{
					Name:        "command",
					Description: "the command to provide help about",
				},
			},
		},
	}
	command.help.Usage = commands.BuildUsage(command.name, command.help.Arguments)
	return command
}

func (h helpCommand) Name() string {
	return h.name
}

func (h helpCommand) Help() commands.CommandHelp {
	return h.help
}

//Run gives help about a helpCommand
func (h helpCommand) Run(ctx context.Context, msg *disgord.Message, args []string) disgord.Embed {
	cmd, ok := h.commandMap[args[0]]

	if !ok {
		return utils.CreateError("Command not found")
	}

	cmdHelp := cmd.Help()

	help := utils.Bold(utils.Block(cmdHelp.Usage)) + " - " + cmdHelp.Description + "\n\n"
	for _, a := range cmdHelp.Arguments {
		help += a.Name + " - " + a.Description + "\n"
	}

	return disgord.Embed{
		Title:       "Help",
		Color:       utils.BLUE,
		Description: help,
	}
}
