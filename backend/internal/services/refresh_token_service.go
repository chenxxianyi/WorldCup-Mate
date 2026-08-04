package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"
	"worldcup-mate/internal/utils"

	"gorm.io/gorm"
)

const refreshTokenTTL = 14 * 24 * time.Hour // refresh token lifetime

// Errors returned by the refresh-token flow (SEC-04).
var (
	ErrRefreshTokenInvalid  = errors.New("invalid refresh token")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrRefreshTokenRevoked  = errors.New("refresh token revoked")
	ErrRefreshTokenReplayed = errors.New("refresh token reused")
)

func hashRefreshToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// IssueRefreshToken creates a new opaque refresh token for the user and
// persists only its hash. The plaintext is returned exactly once.
func IssueRefreshToken(userID uint) (string, error) {
	// GC: drop already-revoked/expired records of this user so the table
	// does not grow unboundedly with each rotation (SEC-04).
	_ = database.DB.Where("user_id = ? AND (revoked_at IS NOT NULL OR expires_at < ?)",
		userID, time.Now()).Delete(&models.RefreshToken{}).Error

	rawBytes := make([]byte, 32)
	if _, err := rand.Read(rawBytes); err != nil {
		return "", err
	}
	raw := base64.RawURLEncoding.EncodeToString(rawBytes)
	token := models.RefreshToken{
		UserID:    userID,
		TokenHash: hashRefreshToken(raw),
		ExpiresAt: time.Now().Add(refreshTokenTTL),
	}
	if err := database.DB.Create(&token).Error; err != nil {
		return "", err
	}
	return raw, nil
}

// ValidateAndRotate verifies a refresh token, revokes it (rotation), and
// issues a fresh access+refresh pair. Reuse of an already-revoked token is
// treated as replay: every session of that user is revoked (SEC-04C).
func ValidateAndRotate(raw string) (userID uint, accessToken string, newRefresh string, err error) {
	if raw == "" {
		return 0, "", "", ErrRefreshTokenInvalid
	}
	var record models.RefreshToken
	result := database.DB.Where("token_hash = ?", hashRefreshToken(raw)).First(&record)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, "", "", ErrRefreshTokenInvalid
	}
	if result.Error != nil {
		return 0, "", "", result.Error
	}

	if record.RevokedAt != nil {
		// Replay of an already-rotated token → revoke the whole account.
		_ = database.DB.Model(&models.RefreshToken{}).
			Where("user_id = ? AND revoked_at IS NULL", record.UserID).
			Update("revoked_at", time.Now())
		return 0, "", "", ErrRefreshTokenReplayed
	}
	if time.Now().After(record.ExpiresAt) {
		return 0, "", "", ErrRefreshTokenExpired
	}

	// Atomic takeover: only the first concurrent request can flip
	// revoked_at (RowsAffected == 1); losers are treated as replay
	// (TOCTOU-safe rotation, SEC-04).
	res := database.DB.Model(&models.RefreshToken{}).
		Where("id = ? AND revoked_at IS NULL", record.ID).
		Update("revoked_at", time.Now())
	if res.Error != nil {
		return 0, "", "", res.Error
	}
	if res.RowsAffected == 0 {
		_ = database.DB.Model(&models.RefreshToken{}).
			Where("user_id = ? AND revoked_at IS NULL", record.UserID).
			Update("revoked_at", time.Now())
		return 0, "", "", ErrRefreshTokenReplayed
	}

	user, err := repositories.GetUserByID(record.UserID)
	if err != nil {
		return 0, "", "", err
	}
	// Disabled accounts must not be able to keep refreshing sessions.
	if user.Status == "disabled" {
		_ = RevokeAllRefreshTokens(user.ID)
		return 0, "", "", ErrRefreshTokenInvalid
	}
	accessToken, err = utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return 0, "", "", err
	}
	newRefresh, err = IssueRefreshToken(user.ID)
	if err != nil {
		return 0, "", "", err
	}
	return user.ID, accessToken, newRefresh, nil
}

// RevokeAllRefreshTokens invalidates every session of the user.
// Used on logout, password change and replay detection (SEC-04C/D).
func RevokeAllRefreshTokens(userID uint) error {
	return database.DB.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", time.Now()).Error
}
