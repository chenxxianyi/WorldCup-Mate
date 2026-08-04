package models

import "time"

// RefreshToken is the persisted refresh-token record (SEC-04).
// Only a SHA-256 hash of the token is stored — the plaintext is handed to
// the client exactly once at issue time and can never be recovered.
// Revocation is soft (RevokedAt set) so replay detection can distinguish
// "unknown token" from "already used token" (rotation).
type RefreshToken struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"index" json:"user_id"`
	TokenHash string     `gorm:"size:64;uniqueIndex" json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `gorm:"index" json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
