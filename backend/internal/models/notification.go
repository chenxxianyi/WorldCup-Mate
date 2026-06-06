package models

import "time"

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index;index:idx_notifications_user_read,priority:1" json:"user_id"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	Type      string    `gorm:"size:30;default:reminder" json:"type"`
	TargetType string   `gorm:"size:40;index" json:"target_type"`
	TargetID   uint     `gorm:"default:0;index" json:"target_id"`
	IsRead    bool      `gorm:"default:false;index:idx_notifications_user_read,priority:2" json:"is_read"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}
