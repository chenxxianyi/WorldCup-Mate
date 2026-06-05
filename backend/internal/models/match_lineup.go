package models

import (
	"time"

	"gorm.io/gorm"
)

type MatchLineup struct {
	ID             uint                `gorm:"primaryKey" json:"id"`
	MatchID        uint                `gorm:"uniqueIndex:idx_match_team_lineup;not null" json:"match_id"`
	Match          Match               `gorm:"foreignKey:MatchID" json:"match,omitempty"`
	TeamID         uint                `gorm:"uniqueIndex:idx_match_team_lineup;index;not null" json:"team_id"`
	Team           Team                `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Players        []MatchLineupPlayer `gorm:"foreignKey:MatchLineupID" json:"players,omitempty"`
	Side           string              `gorm:"size:10;index;not null" json:"side"`
	Formation      string              `gorm:"size:20" json:"formation"`
	CoachName      string              `gorm:"size:100" json:"coach_name"`
	CoachNameEn    string              `gorm:"size:100" json:"coach_name_en"`
	Source         string              `gorm:"size:50;index;not null" json:"source"`
	SourceMatchID  string              `gorm:"size:100;index" json:"source_match_id"`
	ExternalTeamID string              `gorm:"size:100;index" json:"external_team_id"`
	Status         string              `gorm:"size:20;default:unavailable;index" json:"status"`
	LastSyncedAt   *time.Time          `json:"last_synced_at"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	DeletedAt      gorm.DeletedAt      `gorm:"index" json:"-"`
}
