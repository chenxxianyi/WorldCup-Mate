package models

import "time"

type SyncState struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Provider     string     `gorm:"size:30;index:idx_sync_state,unique;not null" json:"provider"`
	Resource     string     `gorm:"size:50;index:idx_sync_state,unique;not null" json:"resource"`
	Status       string     `gorm:"size:20;default:idle" json:"status"`
	LastSyncedAt *time.Time `json:"last_synced_at"`
	NextSyncAt   *time.Time `json:"next_sync_at"`
	LastError    string     `gorm:"type:text" json:"last_error"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
