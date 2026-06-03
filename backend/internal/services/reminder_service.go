package services

import (
	"fmt"
	"log"
	"time"

	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"
	"worldcup-mate/internal/utils"
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

		// Always create in-app notification
		notification := &models.Notification{
			UserID:  r.UserID,
			Title:   title,
			Content: content,
			Type:    "reminder",
		}
		_ = repositories.CreateNotification(notification)

		// Send email if channel is "email"
		// Use NotificationEmail if set, fallback to registration Email
		toEmail := r.User.NotificationEmail
		if toEmail == "" {
			toEmail = r.User.Email
		}
		if r.Channel == "email" && toEmail != "" {
			subject := fmt.Sprintf("⚽ %s", title)
			htmlBody := fmt.Sprintf(`<div style="font-family:sans-serif;max-width:480px;margin:0 auto;padding:20px;">
<h2 style="color:#1a1a2e;">⚽ %s</h2>
<p style="font-size:16px;color:#333;">%s</p>
<div style="background:#f0f4ff;padding:16px;border-radius:12px;margin:16px 0;">
<p style="margin:0;font-size:14px;color:#555;">🏟️ %s</p>
<p style="margin:8px 0 0;font-size:14px;color:#555;">📍 %s</p>
</div>
<p style="font-size:12px;color:#999;">来自 WorldCup Mate 比赛提醒</p>
</div>`, title, content, match.HomeTeam.Name+" vs "+match.AwayTeam.Name, match.Stadium)
			if err := utils.SendEmail(toEmail, subject, htmlBody); err != nil {
				log.Printf("Email send failed for reminder %d: %v", r.ID, err)
			}
		}

		r.Status = "sent"
		_ = repositories.UpdateReminder(&r)
	}
}
