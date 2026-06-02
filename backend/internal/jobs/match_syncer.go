package jobs

import (
	"context"
	"log"
	"time"

	"worldcup-mate/internal/services"
)

func StartMatchSyncer() {
	if !services.IsMatchSyncEnabled() {
		log.Println("Match data sync disabled")
		return
	}

	go func() {
		log.Println("Match data syncer started")
		if err := services.MarkInterruptedMatchSync(); err != nil {
			log.Printf("mark interrupted match sync failed: %v", err)
		}
		runMatchSync("startup")

		for {
			interval := services.NextMatchSyncInterval()
			time.Sleep(interval)
			runMatchSync("scheduled")
		}
	}()
}

func runMatchSync(reason string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err := services.SyncMatchesWithDefault(ctx, reason); err != nil {
		log.Printf("match data sync failed: %v", err)
	}
}
