package jobs

import (
	"context"
	"log"
	"time"

	"worldcup-mate/internal/services"
)

var matchSyncerStop chan struct{}

// StartMatchSyncer launches the World Cup data syncer. It observes the given
// parent context so it can cancel the running sync round and stop on shutdown
// (OBS-04).
func StartMatchSyncer(ctx context.Context) {
	if !services.IsMatchSyncEnabled() {
		log.Println("Match data sync disabled")
		return
	}
	matchSyncerStop = make(chan struct{})

	go func() {
		log.Println("Match data syncer started")
		if err := services.MarkInterruptedMatchSync(); err != nil {
			log.Printf("mark interrupted match sync failed: %v", err)
		}
		runMatchSync(ctx, "startup")

		for {
			interval := services.NextMatchSyncInterval()
			select {
			case <-time.After(interval):
				runMatchSync(ctx, "scheduled")
			case <-ctx.Done():
				log.Println("Match data syncer stopping")
				return
			case <-matchSyncerStop:
				log.Println("Match data syncer stopped")
				return
			}
		}
	}()
}

// StopMatchSyncer asks the match syncer loop to stop.
func StopMatchSyncer() {
	if matchSyncerStop != nil {
		select {
		case matchSyncerStop <- struct{}{}:
		default:
		}
	}
}

func runMatchSync(parent context.Context, reason string) {
	ctx, cancel := context.WithTimeout(parent, 30*time.Second)
	defer cancel()

	if _, err := services.SyncMatchesWithDefault(ctx, reason); err != nil {
		log.Printf("match data sync failed: %v", err)
	}
}
