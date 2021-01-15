package ping

import (
	"context"
	"fmt"
	"net"

	"github.com/asaskevich/govalidator"

	"yey007.github.io/software/pingmc/discord/commands"

	"yey007.github.io/software/pingmc/networking"

	"yey007.github.io/software/pingmc/discord/utils"

	"yey007.github.io/software/pingmc/data"

	"github.com/andersfylling/disgord"
)

type pingCommand struct {
	name string
	help commands.CommandHelp
	repo *data.PingRepo
}

//TODO: Storage object
//New returns a new commands.Command object that has access to the given storage
func New(repo *data.PingRepo) commands.Command {
	command := pingCommand{
		name: "ping",
		help: commands.CommandHelp{
			Description: "Starts pinging a minecraft server",
			Arguments: []commands.CommandArgument{
				{
					Name:        "name",
					Description: "the name of this ping",
				},
				{
					Name:        "address",
					Description: "the address of the minecraft server. If a port is not specified 25565 is assumed",
				},
				{
					Name: "serverType",
					Description: fmt.Sprintf(
						"the type of the server to be pinged. Can be %v or %v",
						utils.Block("vanilla"),
						utils.Block("forge"),
					),
				},
			},
		},
		repo: repo,
	}
	command.help.Usage = commands.BuildUsage(command.name, command.help.Arguments)
	return command
}

func (p pingCommand) Name() string {
	return p.name
}

func (p pingCommand) Help() commands.CommandHelp {
	return p.help
}

//Run adds a new ping to the database
func (p pingCommand) Run(ctx context.Context, msg *disgord.Message, args []string) disgord.Embed {

	//region Validate inputs
	trueAddress := args[1]
	errorStr := ""
	ok := true
	if len(args[0]) > 20 {
		errorStr += "Name too long. Must be less than 20 characters."
		ok = false
	}
	if !govalidator.IsDialString(args[1]) {
		if !govalidator.IsDialString(net.JoinHostPort(args[1], "25565")) {
			errorStr += "\nInvalid host"
			ok = false
		} else {
			trueAddress = net.JoinHostPort(args[1], "25565")
		}
	}
	if args[2] != "vanilla" && args[2] != "forge" {
		errorStr += "\nInvalid server type"
		ok = false
	}
	if !ok {
		return utils.CreateError(errorStr)
	}
	//endregion

	var serverType uint8
	if args[2] == "vanilla" {
		serverType = networking.ServerTypeVanilla
	} else {
		serverType = networking.ServerTypeForge
	}

	ping := data.RecurringPing{
		ChannelID: msg.ChannelID.String(),
		Server: networking.McServer{
			Name:    args[0],
			Address: trueAddress,
			Type:    serverType,
		},
	}

	err := p.repo.Create(ctx, ping)
	if err != nil {
		return utils.CreateError("Failed to create ping! Something seems to have gone wrong with our database.")
	}

	return utils.CreateSuccess(fmt.Sprintf(
		"Added a ping to address %v with id %v",
		utils.Block(ping.Server.Address),
		utils.Block(ping.ID),
	))
}
