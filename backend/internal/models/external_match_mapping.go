package models

import "time"

type ExternalMatchMapping struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	MatchID          uint      `gorm:"uniqueIndex:idx_provider_match;not null" json:"match_id"`
	Match            Match     `gorm:"foreignKey:MatchID" json:"match,omitempty"`
	Provider         string    `gorm:"uniqueIndex:idx_provider_match;size:50;not null" json:"provider"`
	ExternalMatchID  string    `gorm:"size:100;index;not null" json:"external_match_id"`
	ExternalHomeID   string    `gorm:"size:100;index" json:"external_home_id"`
	ExternalAwayID   string    `gorm:"size:100;index" json:"external_away_id"`
	ExternalHomeName string    `gorm:"size:100" json:"external_home_name"`
	ExternalAwayName string    `gorm:"size:100" json:"external_away_name"`
	MatchedBy        string    `gorm:"size:50" json:"matched_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
