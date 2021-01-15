package allhelp

import (
	"context"

	"yey007.github.io/software/pingmc/discord/handler"

	"yey007.github.io/software/pingmc/discord/commands"

	"github.com/andersfylling/disgord"
	"yey007.github.io/software/pingmc/discord/utils"
)

type allHelpCommand struct {
	name       string
	commandMap handler.CommandMap
	help       commands.CommandHelp
}

//New returns a new commands.Command object that has access to the given session
func New(commandMap handler.CommandMap) commands.Command {
	command := allHelpCommand{
		name:       "commands",
		commandMap: commandMap,
		help: commands.CommandHelp{
			Description: "Returns a list of all available commands",
			Arguments:   []commands.CommandArgument{},
		},
	}
	command.help.Usage = commands.BuildUsage(command.name, command.help.Arguments)
	return command
}

func (a allHelpCommand) Name() string {
	return a.name
}

func (a allHelpCommand) Help() commands.CommandHelp {
	return a.help
}

//Run displays the list of all commands to the user
func (a allHelpCommand) Run(ctx context.Context, msg *disgord.Message, args []string) disgord.Embed {

	var help string
	for _, v := range a.commandMap {
		commandHelp := v.Help()
		help += utils.Bold(utils.Block(commandHelp.Usage)) + " - " + commandHelp.Description + "\n"
	}
	return disgord.Embed{
		Title:       "Commands",
		Color:       utils.BLUE,
		Description: help,
	}
}
