package repositories

import (
	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"

	"gorm.io/gorm"
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

func MarkNotificationRead(userID, id uint) error {
	result := database.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func MarkAllNotificationsRead(userID uint) error {
	return database.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}
