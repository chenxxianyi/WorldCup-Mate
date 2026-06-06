package jobs

import (
	"context"
	"log"
	"time"

	"worldcup-mate/internal/services"
)

func StartLineupAlertScanner() {
	if !services.IsLineupAlertEnabled() {
		log.Println("Lineup alert scanner disabled")
		return
	}

	go func() {
		interval := services.NextLineupAlertInterval()
		log.Printf("Lineup alert scanner started (interval: %v)", interval)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		// Run once immediately
		runLineupAlertScan()

		for range ticker.C {
			runLineupAlertScan()
		}
	}()
}

func runLineupAlertScan() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := services.ScanAndSendLineupAlerts(ctx); err != nil {
		log.Printf("Lineup alert scan error: %v", err)
	}
}
