package services

import (
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"
)

func GetNotifications(userID uint, page, pageSize int) ([]models.Notification, int64, error) {
	return repositories.GetNotificationsByUserID(userID, page, pageSize)
}

func CountUnreadNotifications(userID uint) (int64, error) {
	return repositories.CountUnreadNotifications(userID)
}

func MarkNotificationRead(userID, id uint) error {
	return repositories.MarkNotificationRead(userID, id)
}

func MarkAllNotificationsRead(userID uint) error {
	return repositories.MarkAllNotificationsRead(userID)
}
