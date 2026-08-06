package jobs

import (
	"context"
	"log"
	"time"

	"worldcup-mate/internal/services"
)

var leagueSyncerStop chan struct{}

// StartLeagueSyncer schedules league data sync rounds. It is a brand-new
// worker independent from the World Cup match syncer (match_syncer.go).
// It observes the parent context so it can stop cleanly on shutdown (OBS-04).
func StartLeagueSyncer(ctx context.Context) {
	if !services.IsLeagueSyncEnabled() {
		log.Println("League data sync disabled")
		return
	}
	leagueSyncerStop = make(chan struct{})

	go func() {
		log.Println("League data syncer started")
		runLeagueSync(ctx, "startup")

		for {
			interval := services.NextLeagueSyncInterval()
			select {
			case <-time.After(interval):
				runLeagueSync(ctx, "scheduled")
			case <-ctx.Done():
				log.Println("League data syncer stopping")
				return
			case <-leagueSyncerStop:
				log.Println("League data syncer stopped")
				return
			}
		}
	}()
}

// StopLeagueSyncer asks the league syncer loop to stop.
func StopLeagueSyncer() {
	if leagueSyncerStop != nil {
		select {
		case leagueSyncerStop <- struct{}{}:
		default:
		}
	}
}

func runLeagueSync(parent context.Context, reason string) {
	ctx, cancel := context.WithTimeout(parent, 10*time.Minute)
	defer cancel()

	results := services.SyncAllLeagues(ctx, reason)
	for _, r := range results {
		log.Printf("league sync %s (%s): total=%d created=%d updated=%d skipped=%d",
			r.Resource, r.Reason, r.Total, r.Created, r.Updated, r.Skipped)
	}
}
