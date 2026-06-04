package models

import (
	"time"

	"gorm.io/gorm"
)

type AIConversation struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	Title       string         `gorm:"size:120;not null" json:"title"`
	ContextType string         `gorm:"size:30;default:general;index" json:"context_type"`
	ContextID   uint           `gorm:"default:0;index" json:"context_id"`
	LastMessage string         `gorm:"type:text" json:"last_message"`
	Messages    []AIMessage    `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
