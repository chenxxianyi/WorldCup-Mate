package models

import "time"

// AdminAuditLog records every admin write operation (ADM-01):
// who (admin), what (object + id), which action, before/after state.
type AdminAuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	AdminID    uint      `gorm:"not null;index" json:"admin_id"`
	AdminEmail string    `gorm:"size:100" json:"admin_email"`
	Object     string    `gorm:"size:50;not null;index" json:"object"` // team / group / city / stadium / match / standing / user / competition
	ObjectID   string    `gorm:"size:50;index" json:"object_id"`
	Action     string    `gorm:"size:30;not null;index" json:"action"` // create / update / delete / recalculate / status
	Before     string    `gorm:"type:text" json:"before"`              // JSON snapshot
	After      string    `gorm:"type:text" json:"after"`               // JSON snapshot
	CreatedAt  time.Time `json:"created_at"`
}
