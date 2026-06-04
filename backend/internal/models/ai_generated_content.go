package models

import "time"

type AIGeneratedContent struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          *uint     `gorm:"index" json:"user_id"`
	Type            string    `gorm:"size:40;not null;index" json:"type"`
	TargetType      string    `gorm:"size:40;not null;index" json:"target_type"`
	TargetID        uint      `gorm:"default:0;index" json:"target_id"`
	ContentJSON     string    `gorm:"type:longtext" json:"content_json"`
	ContentMarkdown string    `gorm:"type:longtext" json:"content_markdown"`
	Provider        string    `gorm:"size:40" json:"provider"`
	Model           string    `gorm:"size:80" json:"model"`
	CacheKey        string    `gorm:"size:180;uniqueIndex" json:"cache_key"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
