package models

import "time"

// FeaturedConfig is the per-competition homepage hero configuration
// (admin-configurable focus match + promo copy).
type FeaturedConfig struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	CompetitionID uint         `gorm:"uniqueIndex;not null" json:"competition_id"`
	Competition   Competition  `gorm:"foreignKey:CompetitionID" json:"-"`
	MatchID       *uint        `gorm:"index" json:"match_id"` // pinned focus match; nil = automatic (first upcoming)
	Tagline       string       `gorm:"size:200" json:"tagline"`
	Description   string       `gorm:"size:500" json:"description"`
	StageLabel    string       `gorm:"size:50" json:"stage_label"`
	Enabled       bool         `gorm:"default:true" json:"enabled"`
	UpdatedAt     time.Time    `json:"updated_at"`
}
