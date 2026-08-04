package models

import "time"

// Competition represents a football league (e.g. Premier League) or a
// tournament (e.g. World Cup). It is the new top-level dimension added for
// multi-competition support. Existing tables are untouched; matches/teams
// reference a competition through newly added nullable columns.
type Competition struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"size:20;uniqueIndex;not null" json:"code"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	NameEn    string    `gorm:"size:100" json:"name_en"`
	Country   string    `gorm:"size:50" json:"country"`
	LogoURL   string    `gorm:"size:255" json:"logo_url"`
	Format    string    `gorm:"size:20;default:league" json:"format"` // league | cup
	Season    int       `gorm:"index" json:"season"`                  // starting year, e.g. 2025 for 2025-26 season
	Status    string    `gorm:"size:20;default:active" json:"status"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
