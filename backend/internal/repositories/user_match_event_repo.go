package repositories

import (
	"errors"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"

	"gorm.io/gorm"
)

func HasUserMatchEvent(userID, matchID uint, eventType string) (bool, error) {
	var count int64
	err := database.DB.Model(&models.UserMatchEventLog{}).
		Where("user_id = ? AND match_id = ? AND event_type = ?", userID, matchID, eventType).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CreateUserMatchEvent(userID, matchID uint, eventType string) error {
	event := models.UserMatchEventLog{
		UserID:    userID,
		MatchID:   matchID,
		EventType: eventType,
	}
	err := database.DB.Create(&event).Error
	if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		return err
	}
	return nil
}
