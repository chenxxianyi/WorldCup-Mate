package jobs

import (
	"log"
	"time"

	"worldcup-mate/internal/services"
)

func StartReminderScanner() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		log.Println("Reminder scanner started")
		for range ticker.C {
			services.ScanAndSendReminders()
		}
	}()
}
