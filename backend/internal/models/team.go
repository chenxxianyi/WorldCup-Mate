package models

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:50;not null" json:"name"`
	NameEn      string         `gorm:"size:50" json:"name_en"`
	FIFACode    string         `gorm:"size:10;uniqueIndex" json:"fifa_code"`
	FlagURL     string         `gorm:"size:255" json:"flag_url"`
	Continent   string         `gorm:"size:20" json:"continent"`
	GroupID     uint           `json:"group_id"`
	Group       Group          `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Description string         `gorm:"type:text" json:"description"`
	Coach       string         `gorm:"size:100" json:"coach"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
