package models

import (
	"time"

	"gorm.io/gorm"
)

type MatchLineupPlayer struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	MatchLineupID  uint           `gorm:"index;not null" json:"match_lineup_id"`
	MatchLineup    MatchLineup    `gorm:"foreignKey:MatchLineupID" json:"lineup,omitempty"`
	MatchID        uint           `gorm:"index;not null" json:"match_id"`
	TeamID         uint           `gorm:"index;not null" json:"team_id"`
	PlayerID       *uint          `gorm:"index" json:"player_id"`
	Player         *Player        `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
	Source         string         `gorm:"size:50;index;not null" json:"source"`
	SourcePlayerID string         `gorm:"size:100;index" json:"source_player_id"`
	Name           string         `gorm:"size:100;not null" json:"name"`
	NameEn         string         `gorm:"size:100" json:"name_en"`
	ShirtNumber    int            `gorm:"index" json:"shirt_number"`
	Position       string         `gorm:"size:20;index" json:"position"`
	PositionLabel  string         `gorm:"size:20" json:"position_label"`
	Role           string         `gorm:"size:20;index;not null" json:"role"`
	Grid           string         `gorm:"size:20" json:"grid"`
	SortOrder      int            `gorm:"index" json:"sort_order"`
	PhotoURL       string         `gorm:"size:500" json:"photo_url"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
