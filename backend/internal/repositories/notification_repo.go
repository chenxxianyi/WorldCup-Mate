package repositories

import (
	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
)

func CreateNotification(n *models.Notification) error {
	return database.DB.Create(n).Error
}

func GetNotificationsByUserID(userID uint, page, pageSize int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64
	database.DB.Model(&models.Notification{}).Where("user_id = ?", userID).Count(&total)
	err := database.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&notifications).Error
	return notifications, total, err
}

func CountUnreadNotifications(userID uint) (int64, error) {
	var total int64
	err := database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&total).Error
	return total, err
}

func MarkNotificationRead(id uint) error {
	return database.DB.Model(&models.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

func MarkAllNotificationsRead(userID uint) error {
	return database.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}
