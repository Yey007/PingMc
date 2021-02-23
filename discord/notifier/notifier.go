package notifier

import (
	"context"
	"sync"
	"time"

	"github.com/andersfylling/snowflake/v4"

	"yey007.github.io/software/pingmc/data"

	"github.com/andersfylling/disgord"
)

const fastPingInterval = 5
const slowPingInterval = 10

func StartNotifier(ctx context.Context, session disgord.Session, repo *data.PingRepo) {
	go notifyQuickly(ctx, session, repo)
	go notifySlowly(ctx, session, repo)
}

func notifyQuickly(ctx context.Context, session disgord.Session, repo *data.PingRepo) {
	for range time.Tick(fastPingInterval * time.Second) {
		ch := make(chan data.PingPair)
		go repo.GetAll(ctx, ch)

		wg := sync.WaitGroup{}
		for pair := range ch {
			wg.Add(1)
			go func(pair data.PingPair) {
				//Only ping if there was previous data (was online)
				defer wg.Done()
				if pair.Temp.HasPrevious {
					embed, shouldSend := createNotification(pair.Ping, pair.Temp, repo)
					if shouldSend {
						_, _ = session.SendMsg(ctx, snowflake.ParseSnowflakeString(pair.Ping.ChannelID), embed)
					}
				}
			}(pair)
		}
		wg.Wait()
	}
}

func notifySlowly(ctx context.Context, session disgord.Session, repo *data.PingRepo) {
	for range time.Tick(slowPingInterval * time.Second) {
		ch := make(chan data.PingPair)
		go repo.GetAll(ctx, ch)

		wg := sync.WaitGroup{}
		for pair := range ch {
			wg.Add(1)
			go func(pair data.PingPair) {
				//Only ping if there wasn't previous data (was offline)
				defer wg.Done()
				if !pair.Temp.HasPrevious {
					embed, shouldSend := createNotification(pair.Ping, pair.Temp, repo)
					if shouldSend {
						_, _ = session.SendMsg(ctx, snowflake.ParseSnowflakeString(pair.Ping.ChannelID), embed)
					}
				}
			}(pair)
		}
		wg.Wait()
	}
}
