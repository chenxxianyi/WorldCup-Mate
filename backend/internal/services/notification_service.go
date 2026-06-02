package services

import (
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"
)

func GetNotifications(userID uint, page, pageSize int) ([]models.Notification, int64, error) {
	return repositories.GetNotificationsByUserID(userID, page, pageSize)
}

func MarkNotificationRead(id uint) error {
	return repositories.MarkNotificationRead(id)
}

func MarkAllNotificationsRead(userID uint) error {
	return repositories.MarkAllNotificationsRead(userID)
}
