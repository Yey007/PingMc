package updater

import (
	"errors"
	"net"
	"strings"

	"github.com/andersfylling/snowflake/v4"
	"yey007.github.io/software/pingmc/discord/utils"
	"yey007.github.io/software/pingmc/networking"
	"yey007.github.io/software/pingmc/networking/fml1"
	"yey007.github.io/software/pingmc/networking/vanilla"
)

func validateInputs(info utils.SessionInfo, channelID snowflake.Snowflake, args []string) bool {

	for i := range args {
		args[i] = strings.Trim(args[i], "\r\n ")
	}

	if len(args) != 4 {
		utils.ShowError(info, channelID, "Wrong number of arguments")
		return false
	}

	result := true
	ip := args[2]
	serverType := args[3]

	ip = strings.Split(ip, ":")[0]

	if net.ParseIP(ip) == nil {
		result = result && false
		utils.ShowError(info, channelID, "Invalid IP address")
	}

	if serverType != "vanilla" && serverType != "forge1" {
		result = result && false
		utils.ShowError(info, channelID, "Invalid server type")
	}
	return result
}

func createServer(info utils.SessionInfo, serverType string) (networking.McServer, error) {
	if serverType == "vanilla" {
		return new(vanilla.Server), nil
	} else if serverType == "forge1" {
		return new(fml1.Server), nil
	}
	return nil, errors.New("Incorrect server type. How did this happen?")
}
