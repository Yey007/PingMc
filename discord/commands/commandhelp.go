package commands

import (
	"fmt"

	"yey007.github.io/software/pingmc/discord/utils"
)

type CommandHelp struct {
	Description string
	Arguments   []CommandArgument
	Usage       string
}

type CommandArgument struct {
	Name        string
	Description string
}

//BuildUsage builds the usage string for a command give a name and arguments
func BuildUsage(name string, arguments []CommandArgument) string {
	usage := utils.PREFIX + " " + name
	for _, a := range arguments {
		usage += fmt.Sprintf(" <%v>", a.Name)
	}
	return usage
}
