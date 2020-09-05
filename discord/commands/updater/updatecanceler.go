package updater

import (
	"sync"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
	"yey007.github.io/software/pingmc/discord/utils"
)

var runningPings = map[snowflake.Snowflake]*chan bool{}
var lock = sync.RWMutex{}

//OnPingCancelRequest stops a running ping in the current channel
func OnPingCancelRequest(info utils.SessionInfo, msg *disgord.Message, args []string) {

	lock.RLock()
	defer lock.RUnlock()

	ch, exists := runningPings[msg.ChannelID]

	if exists {
		close(*ch)
		delete(runningPings, msg.ChannelID)
		utils.ShowSuccess(info, msg.ChannelID, "Ping cancelled")
	} else {
		utils.ShowError(info, msg.ChannelID, "No running ping found.")
	}
}

func addNewPing(channelID snowflake.Snowflake) *chan bool {
	ch := make(chan bool, 1)
	lock.Lock()
	defer lock.Unlock()

	_, exists := runningPings[channelID]

	if exists {
		return nil
	}
	runningPings[channelID] = &ch

	return &ch
}
