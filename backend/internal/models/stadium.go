package models

import (
	"time"

	"gorm.io/gorm"
)

type Stadium struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	NameEn      string         `gorm:"size:100;index" json:"name_en"`
	CityID      uint           `json:"city_id"`
	City        City           `gorm:"foreignKey:CityID" json:"city,omitempty"`
	Capacity    int            `json:"capacity"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
