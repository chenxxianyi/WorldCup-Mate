package models

import "time"

type UserMatchEventLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_match_event" json:"user_id"`
	MatchID   uint      `gorm:"not null;uniqueIndex:idx_user_match_event" json:"match_id"`
	EventType string    `gorm:"size:40;not null;uniqueIndex:idx_user_match_event" json:"event_type"`
	CreatedAt time.Time `json:"created_at"`
}
