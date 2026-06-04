package models

import "time"

type AIUsageLog struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           *uint     `gorm:"index" json:"user_id"`
	IP               string    `gorm:"size:64;index" json:"ip"`
	Endpoint         string    `gorm:"size:80;index" json:"endpoint"`
	Provider         string    `gorm:"size:40" json:"provider"`
	Model            string    `gorm:"size:80" json:"model"`
	PromptTokens     int       `gorm:"default:0" json:"prompt_tokens"`
	CompletionTokens int       `gorm:"default:0" json:"completion_tokens"`
	TotalTokens      int       `gorm:"default:0" json:"total_tokens"`
	Status           string    `gorm:"size:20;index" json:"status"`
	ErrorMessage     string    `gorm:"type:text" json:"error_message"`
	LatencyMS        int64     `gorm:"default:0" json:"latency_ms"`
	CreatedAt        time.Time `gorm:"index" json:"created_at"`
}
