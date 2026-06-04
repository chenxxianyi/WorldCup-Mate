package models

import (
	"time"

	"gorm.io/gorm"
)

type City struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:50;not null" json:"name"`
	NameEn      string         `gorm:"size:50;index" json:"name_en"`
	Country     string         `gorm:"size:50" json:"country"`
	Timezone    string         `gorm:"size:50" json:"timezone"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
