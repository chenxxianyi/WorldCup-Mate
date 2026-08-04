package services

import (
	"encoding/json"
	"log"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
)

// RecordAudit appends an admin write operation to admin_audit_logs.
// before/after are JSON-serialized; serialization failures degrade to "".
// The audit trail is best-effort: a failure is logged, never fails the
// business operation.
func RecordAudit(adminID uint, adminEmail, object, objectID, action string, before, after interface{}) {
	entry := models.AdminAuditLog{
		AdminID:    adminID,
		AdminEmail: adminEmail,
		Object:     object,
		ObjectID:   objectID,
		Action:     action,
	}
	if before != nil {
		if b, err := json.Marshal(before); err == nil {
			entry.Before = string(b)
		}
	}
	if after != nil {
		if a, err := json.Marshal(after); err == nil {
			entry.After = string(a)
		}
	}
	if err := database.DB.Create(&entry).Error; err != nil {
		log.Printf("[audit] failed to record audit: %v", err)
	}
}
