package models

import (
	"time"

	"gorm.io/gorm"
)

type Reminder struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	UserID              uint           `gorm:"not null;index" json:"user_id"`
	User                User           `gorm:"foreignKey:UserID" json:"-"`
	MatchID             uint           `gorm:"not null;index" json:"match_id"`
	Match               Match          `gorm:"foreignKey:MatchID" json:"match,omitempty"`
	RemindBeforeMinutes int            `gorm:"default:30" json:"remind_before_minutes"`
	RemindAt            time.Time      `gorm:"not null;index" json:"remind_at"`
	Channel             string         `gorm:"size:20;default:site" json:"channel"`
	Status              string         `gorm:"size:20;default:pending;index" json:"status"`
	RetryCount          int            `gorm:"default:0" json:"retry_count"`
	// ClaimToken is set atomically when a scanner claims this reminder in the
	// multi-instance model. After a claim, the scanner only queries records
	// carrying its own token, so overlapping claim windows never observe
	// another worker's reminders (REL-08).
	ClaimToken string `gorm:"size:36;default:'';index" json:"claim_token"`
	// ClaimedAt records when the reminder was claimed, enabling timed
	// recovery of claims abandoned by crashed workers (sending timeout).
	ClaimedAt *time.Time `gorm:"index" json:"claimed_at"`
	// NextRetryAt schedules the next processing attempt (exponential backoff).
	// A reminder is only picked up by the scanner when status=pending AND
	// remind_at/next_retry_at has passed, so failures are not retried hot.
	NextRetryAt *time.Time `gorm:"index" json:"next_retry_at"`
	// LastError records the terminal or most recent failure reason (ADM-13).
	LastError string `gorm:"size:500" json:"last_error"`
	// NotificationID is the idempotency key: the in-app notification created
	// on the first successful delivery attempt is stored here so retries and
	// re-claims never duplicate it (REL-08).
	NotificationID *uint  `gorm:"index" json:"notification_id"`
	WorkerID       string `gorm:"size:36;default:''" json:"worker_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
