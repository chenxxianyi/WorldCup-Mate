package models

import "time"

type UserFavoriteMatch struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex:idx_user_match;not null" json:"user_id"`
	MatchID   uint      `gorm:"uniqueIndex:idx_user_match;not null" json:"match_id"`
	Match     Match     `gorm:"foreignKey:MatchID" json:"match,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
