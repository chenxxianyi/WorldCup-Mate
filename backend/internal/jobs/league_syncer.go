package jobs

import (
	"context"
	"log"
	"time"

	"worldcup-mate/internal/services"
)

// StartLeagueSyncer schedules league data sync rounds. It is a brand-new
// worker independent from the World Cup match syncer (match_syncer.go).
func StartLeagueSyncer() {
	if !services.IsLeagueSyncEnabled() {
		log.Println("League data sync disabled")
		return
	}

	go func() {
		log.Println("League data syncer started")
		runLeagueSync("startup")

		for {
			interval := services.NextLeagueSyncInterval()
			time.Sleep(interval)
			runLeagueSync("scheduled")
		}
	}()
}

func runLeagueSync(reason string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	results := services.SyncAllLeagues(ctx, reason)
	for _, r := range results {
		log.Printf("league sync %s (%s): total=%d created=%d updated=%d skipped=%d",
			r.Resource, r.Reason, r.Total, r.Created, r.Updated, r.Skipped)
	}
}
