package models

import (
	"time"

	"gorm.io/gorm"
)

type Reminder struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	UserID              uint           `gorm:"not null;index" json:"user_id"`
	User                User           `gorm:"foreignKey:UserID" json:"-"`
	MatchID             uint           `gorm:"not null;index" json:"match_id"`
	Match               Match          `gorm:"foreignKey:MatchID" json:"match,omitempty"`
	RemindBeforeMinutes int            `gorm:"default:30" json:"remind_before_minutes"`
	RemindAt            time.Time      `gorm:"not null;index" json:"remind_at"`
	Channel             string         `gorm:"size:20;default:site" json:"channel"`
	Status              string         `gorm:"size:20;default:pending;index" json:"status"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}
