package jobs

import (
	"context"
	"log"
	"time"

	"worldcup-mate/internal/services"
)

func StartLineupSyncer() {
	if !services.IsLineupSyncEnabled() {
		log.Println("Lineup sync disabled")
		return
	}

	go func() {
		log.Println("Lineup syncer started")
		if err := services.MarkInterruptedLineupSync(); err != nil {
			log.Printf("mark interrupted lineup sync failed: %v", err)
		}

		for {
			runLineupSync("scheduled")
			time.Sleep(services.NextLineupSyncInterval())
		}
	}()
}

func runLineupSync(reason string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if _, err := services.SyncLiveWindowLineups(ctx, reason); err != nil {
		log.Printf("lineup sync failed: %v", err)
	}
}
