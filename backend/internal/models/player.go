package models

import (
	"time"

	"gorm.io/gorm"
)

type Player struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	TeamID         uint           `gorm:"index;not null" json:"team_id"`
	Team           Team           `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Name           string         `gorm:"size:100;not null;index" json:"name"`
	NameEn         string         `gorm:"size:100;index" json:"name_en"`
	ShirtNumber    int            `gorm:"index" json:"shirt_number"`
	Position       string         `gorm:"size:20;index" json:"position"`
	PositionLabel  string         `gorm:"size:20" json:"position_label"`
	PhotoURL       string         `gorm:"size:500" json:"photo_url"`
	Club           string         `gorm:"size:100" json:"club"`
	Source         string         `gorm:"size:50;index" json:"source"`
	SourcePlayerID string         `gorm:"size:100;index" json:"source_player_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
