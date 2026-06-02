package services

import (
	"fmt"
	"time"

	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"
)

type CreateReminderInput struct {
	MatchID             uint   `json:"matchId" binding:"required"`
	RemindBeforeMinutes int    `json:"remindBeforeMinutes"`
	Channel             string `json:"channel"`
}

func CreateReminder(userID uint, input CreateReminderInput) (*models.Reminder, error) {
	match, err := repositories.GetMatchByID(input.MatchID)
	if err != nil {
		return nil, fmt.Errorf("match not found")
	}

	minutes := input.RemindBeforeMinutes
	if minutes <= 0 {
		minutes = 30
	}
	channel := input.Channel
	if channel == "" {
		channel = "site"
	}

	reminder := &models.Reminder{
		UserID:              userID,
		MatchID:             input.MatchID,
		RemindBeforeMinutes: minutes,
		RemindAt:            match.KickoffTimeUTC.Add(-time.Duration(minutes) * time.Minute),
		Channel:             channel,
		Status:              "pending",
	}
	if err := repositories.CreateReminder(reminder); err != nil {
		return nil, err
	}
	return reminder, nil
}

func GetReminders(userID uint) ([]models.Reminder, error) {
	return repositories.GetRemindersByUserID(userID)
}

func UpdateReminder(id uint, userID uint, minutes int, channel string) (*models.Reminder, error) {
	reminder, err := repositories.GetReminderByID(id)
	if err != nil {
		return nil, fmt.Errorf("reminder not found")
	}
	if reminder.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	if minutes > 0 {
		reminder.RemindBeforeMinutes = minutes
		match, err := repositories.GetMatchByID(reminder.MatchID)
		if err == nil {
			reminder.RemindAt = match.KickoffTimeUTC.Add(-time.Duration(minutes) * time.Minute)
		}
	}
	if channel != "" {
		reminder.Channel = channel
	}
	if err := repositories.UpdateReminder(reminder); err != nil {
		return nil, err
	}
	return reminder, nil
}

func DeleteReminder(id uint, userID uint) error {
	reminder, err := repositories.GetReminderByID(id)
	if err != nil {
		return fmt.Errorf("reminder not found")
	}
	if reminder.UserID != userID {
		return fmt.Errorf("unauthorized")
	}
	return repositories.DeleteReminder(id)
}

func ScanAndSendReminders() {
	now := time.Now().UTC()
	reminders, err := repositories.GetPendingReminders(now)
	if err != nil {
		return
	}
	for _, r := range reminders {
		match, err := repositories.GetMatchByID(r.MatchID)
		if err != nil {
			continue
		}
		title := fmt.Sprintf("比赛提醒：%s vs %s", match.HomeTeam.Name, match.AwayTeam.Name)
		content := fmt.Sprintf("比赛即将在 %d 分钟后开始", r.RemindBeforeMinutes)

		notification := &models.Notification{
			UserID:  r.UserID,
			Title:   title,
			Content: content,
			Type:    "reminder",
		}
		_ = repositories.CreateNotification(notification)
		r.Status = "sent"
		_ = repositories.UpdateReminder(&r)
	}
}
