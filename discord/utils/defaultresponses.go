package utils

import (
	"fmt"

	"github.com/andersfylling/disgord"
)

//CreateError displays an error
func CreateError(message string) disgord.Embed {
	return disgord.Embed{
		Title: "Error",
		Color: RED,
		Description: fmt.Sprintf(
			"%v \n\n %v for help\n %v for command list",
			message,
			Block(".pingmc help <command>"),
			Block(".pingmc commands"),
		),
	}
}

//CreateSuccess display a success
func CreateSuccess(message string) disgord.Embed {
	return disgord.Embed{
		Title:       "Success!",
		Color:       GREEN,
		Description: message,
	}
}
