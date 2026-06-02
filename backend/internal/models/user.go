package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Avatar       string         `gorm:"size:255" json:"avatar"`
	Timezone     string         `gorm:"size:50;default:Asia/Shanghai" json:"timezone"`
	Language     string         `gorm:"size:10;default:zh-CN" json:"language"`
	Role         string         `gorm:"size:20;default:user" json:"role"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
