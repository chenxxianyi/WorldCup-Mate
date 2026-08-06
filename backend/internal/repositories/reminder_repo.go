package repositories

import (
	"errors"
	"time"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"

	"gorm.io/gorm"
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
	err := database.DB.Preload("User").Preload("Match").Preload("Match.HomeTeam").Preload("Match.AwayTeam").First(&reminder, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
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
	err := database.DB.Preload("Match").Preload("Match.HomeTeam").Preload("Match.AwayTeam").Preload("Match.Stadium").Preload("User").
		Where("status = ? AND remind_at <= ?", "pending", now).
		Find(&reminders).Error
	return reminders, err
}

// CountRemindersByStatus returns the number of reminders with the given status (ADM-13).
func CountRemindersByStatus(status string) int64 {
	var n int64
	_ = database.DB.Model(&models.Reminder{}).Where("status = ?", status).Count(&n).Error
	return n
}

// ListRemindersByMatchID lists reminders for a given match with user/preload info (ADM-13).
func ListRemindersByMatchID(matchID uint) ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := database.DB.Where("match_id = ?", matchID).
		Order("remind_at ASC").Find(&reminders).Error
	return reminders, err
}

// ListRemindersWithStats returns a paginated list of reminders with user and match
// preloads for admin views (ADM-13).
func ListRemindersWithStats(page, pageSize int) ([]models.Reminder, int64, error) {
	var total int64
	q := database.DB.Model(&models.Reminder{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	var reminders []models.Reminder
	err := database.DB.Preload("User").Preload("Match").
		Order("remind_at ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&reminders).Error
	return reminders, total, err
}
