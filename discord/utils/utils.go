package utils

import (
	"context"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
	"yey007.github.io/software/pingmc/discord/utils/colors"
)

//SessionInfo represent common information about a session
type SessionInfo struct {
	Con     context.Context
	Session disgord.Session
}

//ShowError displays an error
func ShowError(info SessionInfo, channelID snowflake.Snowflake, message string) (*disgord.Message, error) {
	msg, err := info.Session.SendMsg(info.Con, channelID, disgord.Embed{Title: "Error", Color: colors.Red,
		Description: message + "\n\n" + Block(Bold(".pingmc help <command>")) + "for help\n" + Block(Bold(".pingmc commands")) + " for command list"})
	return msg, err
}

//ShowSuccess display a success
func ShowSuccess(info SessionInfo, channelID snowflake.Snowflake, message string) (*disgord.Message, error) {
	msg, err := info.Session.SendMsg(info.Con, channelID, disgord.Embed{Title: "Success!", Color: colors.Green,
		Description: message})
	return msg, err
}

//Block puts the given text in a single-line code block
func Block(s string) string {
	return "`" + s + "`"
}

//Bold bolds the given text
func Bold(s string) string {
	return "**" + s + "**"
}

//Italics italicizes the given text
func Italics(s string) string {
	return "*" + s + "*"
}
