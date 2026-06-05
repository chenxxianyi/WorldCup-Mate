package jobs

import (
	"context"
	"log"
	"time"

	"worldcup-mate/internal/services"
)

func StartPlayerSyncer() {
	if !services.IsPlayerSyncEnabled() {
		log.Println("Player data sync disabled")
		return
	}

	go func() {
		log.Println("Player data syncer started")
		if err := services.MarkInterruptedPlayerSync(); err != nil {
			log.Printf("mark interrupted player sync failed: %v", err)
		}

		if services.GetPlayerSyncConfig().SyncOnStartup {
			time.Sleep(2 * time.Minute)
			runPlayerSync("startup")
		}

		for {
			time.Sleep(services.NextPlayerSyncInterval())
			runPlayerSync("scheduled")
		}
	}()
}

func runPlayerSync(reason string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if _, err := services.SyncAllMappedTeamPlayersWithDefault(ctx, reason); err != nil {
		log.Printf("player data sync failed: %v", err)
	}
}
