package models

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:50;not null;index" json:"name"`
	NameEn       string         `gorm:"size:50;index" json:"name_en"`
	FIFACode     *string        `gorm:"size:10;uniqueIndex" json:"fifa_code"`            // national teams only; NULL for clubs (column stays nullable, unique index untouched)
	ExternalCode *string        `gorm:"size:20;index" json:"external_code"`              // provider code (e.g. football-data tla) for clubs
	TeamType     string         `gorm:"size:20;default:national;index" json:"team_type"` // national | club
	FlagURL      string         `gorm:"size:255" json:"flag_url"`
	Continent    string         `gorm:"size:20;index" json:"continent"`
	Country      string         `gorm:"size:50" json:"country"` // club country
	Venue        string         `gorm:"size:100" json:"venue"`  // club home stadium
	GroupID      *uint          `gorm:"index" json:"group_id"`  // national teams only; NULL for clubs
	Group        Group          `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Description  string         `gorm:"type:text" json:"description"`
	Coach        string         `gorm:"size:100" json:"coach"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
