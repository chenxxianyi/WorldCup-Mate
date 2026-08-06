package jobs

import (
	"context"
	"log"
	"time"

	"worldcup-mate/internal/services"
)

var reminderScannerStop chan struct{}

// StartReminderScanner starts the reminder scanning loop. It captures the
// given process-lifetime context (or creates an uncancelled one when nil) so
// worker goroutines can be stopped cleanly on shutdown (OBS-04).
func StartReminderScanner(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	reminderScannerStop = make(chan struct{})

	go func() {
		log.Println("Reminder scanner started")
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				services.ScanAndSendReminders()
			case <-ctx.Done():
				log.Println("Reminder scanner stopping")
				return
			case <-reminderScannerStop:
				log.Println("Reminder scanner stopped")
				return
			}
		}
	}()
}

// StopReminderScanner asks the reminder scanner loop to stop.
func StopReminderScanner() {
	if reminderScannerStop != nil {
		select {
		case reminderScannerStop <- struct{}{}:
		default:
		}
	}
}
