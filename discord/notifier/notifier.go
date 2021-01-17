package notifier

import (
	"context"
	"time"

	"github.com/andersfylling/snowflake/v4"

	"yey007.github.io/software/pingmc/data"

	"github.com/andersfylling/disgord"
)

func StartNotifier(ctx context.Context, session disgord.Session, repo *data.PingRepo) {
	go notifyQuickly(ctx, session, repo)
	go notifySlowly(ctx, session, repo)
}

func notifyQuickly(ctx context.Context, session disgord.Session, repo *data.PingRepo) {
	for range time.Tick(5 * time.Second) {
		repo.GetAllInBatches(ctx, func(ping data.RecurringPing, tempData data.TempPingData) {
			//Only ping if there was previous data (was online)
			if tempData.HasPrevious {
				embed := createNotification(ping, tempData, repo)
				if embed != nil {
					session.SendMsg(ctx, snowflake.ParseSnowflakeString(ping.ChannelID), *embed)
				}
			}
		})
	}
}

func notifySlowly(ctx context.Context, session disgord.Session, repo *data.PingRepo) {
	for range time.Tick(10 * time.Second) {
		repo.GetAllInBatches(ctx, func(ping data.RecurringPing, tempData data.TempPingData) {
			//Only ping if there wasn't previous data (was offline)
			if !tempData.HasPrevious {
				embed := createNotification(ping, tempData, repo)
				if embed != nil {
					session.SendMsg(ctx, snowflake.ParseSnowflakeString(ping.ChannelID), *embed)
				}
			}
		})
	}
}

func createNotification(
	p data.RecurringPing,
	temp data.TempPingData,
	repo *data.PingRepo,
) *disgord.Embed {
	//TODO: How many times does a ping have to fail for us to just stop pinging it?
	pingdata, err := p.Server.Ping()

	if temp.HasPrevious {
		if err == nil {
			defer repo.CreateTemp(data.TempPingData{
				RecurringPingID: p.ID,
				PreviousData:    pingdata,
				HasPrevious:     true,
			})
			//Server was online and is still online
			if pingdata.Players.Online > temp.PreviousData.Players.Online {
				//Number of players more
				return createPlayerJoin(*pingdata, p.Name)
			} else if pingdata.Players.Online < temp.PreviousData.Players.Online {
				//Number of players less
				return createPlayerLeave(*pingdata, p.Name)
			} else if pingdata.Players.Max != temp.PreviousData.Players.Max {
				//Max players different
				return createServerOnline(*pingdata, p.Name)
			}
		} else {
			//Server was online and is now offline
			defer repo.CreateTemp(data.TempPingData{
				RecurringPingID: p.ID,
				PreviousData:    nil,
				HasPrevious:     false,
			})
			return createServerOffline(p.Name)
		}
	} else {
		if err == nil {
			//Server was offline and is now online
			defer repo.CreateTemp(data.TempPingData{
				RecurringPingID: p.ID,
				PreviousData:    pingdata,
				HasPrevious:     true,
			})
			return createServerOnline(*pingdata, p.Name)
		}
		//Server was offline and is still offline, do nothing
	}
	return nil
}
