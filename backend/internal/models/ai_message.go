package models

import "time"

type AIMessage struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	ConversationID   uint      `gorm:"not null;index" json:"conversation_id"`
	UserID           uint      `gorm:"not null;index" json:"user_id"`
	Role             string    `gorm:"size:20;not null;index" json:"role"`
	Content          string    `gorm:"type:text;not null" json:"content"`
	Provider         string    `gorm:"size:40" json:"provider"`
	Model            string    `gorm:"size:80" json:"model"`
	PromptTokens     int       `gorm:"default:0" json:"prompt_tokens"`
	CompletionTokens int       `gorm:"default:0" json:"completion_tokens"`
	TotalTokens      int       `gorm:"default:0" json:"total_tokens"`
	CreatedAt        time.Time `json:"created_at"`
}
