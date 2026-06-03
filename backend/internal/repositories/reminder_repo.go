package repositories

import (
	"time"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
)

func CreateReminder(reminder *models.Reminder) error {
	return database.DB.Create(reminder).Error
}

func GetRemindersByUserID(userID uint) ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := database.DB.Preload("Match").Preload("Match.HomeTeam").Preload("Match.AwayTeam").
		Where("user_id = ?", userID).Order("remind_at ASC").Find(&reminders).Error
	return reminders, err
}

func GetReminderByID(id uint) (*models.Reminder, error) {
	var reminder models.Reminder
	err := database.DB.First(&reminder, id).Error
	return &reminder, err
}

func UpdateReminder(reminder *models.Reminder) error {
	return database.DB.Save(reminder).Error
}

func DeleteReminder(id uint) error {
	return database.DB.Delete(&models.Reminder{}, id).Error
}

func GetPendingReminders(now time.Time) ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := database.DB.Preload("Match").Preload("Match.HomeTeam").Preload("Match.AwayTeam").Preload("User").
		Where("status = ? AND remind_at <= ?", "pending", now).
		Find(&reminders).Error
	return reminders, err
}
