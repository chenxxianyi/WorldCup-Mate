package services

import (
	"fmt"
	"log"
	"time"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"
	"worldcup-mate/internal/utils"
)

type CreateReminderInput struct {
	MatchID             uint   `json:"matchId" binding:"required"`
	RemindBeforeMinutes int    `json:"remindBeforeMinutes"`
	Channel             string `json:"channel"`
}

type CreateReminderBatchInput struct {
	MatchID uint   `json:"match_id" binding:"required"`
	Minutes []int  `json:"minutes" binding:"required"`
	Channel string `json:"channel"`
}

func CreateReminder(userID uint, input CreateReminderInput) (*models.Reminder, error) {
	match, err := repositories.GetMatchByID(input.MatchID)
	if err != nil {
		return nil, fmt.Errorf("match not found")
	}

	reminder, err := buildReminder(userID, match, input.RemindBeforeMinutes, input.Channel)
	if err != nil {
		return nil, err
	}
	if err := repositories.CreateReminder(reminder); err != nil {
		return nil, err
	}
	return reminder, nil
}

func CreateReminderBatch(userID uint, input CreateReminderBatchInput) ([]models.Reminder, error) {
	match, err := repositories.GetMatchByID(input.MatchID)
	if err != nil {
		return nil, fmt.Errorf("match not found")
	}

	existing, err := repositories.GetRemindersByUserID(userID)
	if err != nil {
		return nil, err
	}
	existingKeys := map[string]bool{}
	for _, item := range existing {
		existingKeys[reminderKey(item.MatchID, item.RemindBeforeMinutes, item.Channel)] = true
	}

	seenMinutes := map[int]bool{}
	var created []models.Reminder
	for _, minutes := range input.Minutes {
		if seenMinutes[minutes] {
			continue
		}
		seenMinutes[minutes] = true

		reminder, err := buildReminder(userID, match, minutes, input.Channel)
		if err != nil {
			return nil, err
		}
		key := reminderKey(reminder.MatchID, reminder.RemindBeforeMinutes, reminder.Channel)
		if existingKeys[key] {
			continue
		}
		if err := repositories.CreateReminder(reminder); err != nil {
			return nil, err
		}
		created = append(created, *reminder)
	}

	if len(created) == 0 {
		return nil, fmt.Errorf("reminder already exists")
	}
	return created, nil
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

func buildReminder(userID uint, match *models.Match, minutes int, channel string) (*models.Reminder, error) {
	if match.Status == "live" || match.Status == "finished" || match.Status == "cancelled" {
		return nil, fmt.Errorf("match already started or ended")
	}
	if minutes <= 0 {
		minutes = 30
	}
	if channel == "" {
		channel = "site"
	}
	if channel != "site" && channel != "email" {
		return nil, fmt.Errorf("unsupported reminder channel")
	}

	remindAt := match.KickoffTimeUTC.Add(-time.Duration(minutes) * time.Minute)
	if !remindAt.After(time.Now().UTC()) {
		return nil, fmt.Errorf("reminder time has passed")
	}

	return &models.Reminder{
		UserID:              userID,
		MatchID:             match.ID,
		RemindBeforeMinutes: minutes,
		RemindAt:            remindAt,
		Channel:             channel,
		Status:              "pending",
	}, nil
}

func reminderKey(matchID uint, minutes int, channel string) string {
	return fmt.Sprintf("%d:%d:%s", matchID, minutes, channel)
}

// reminderMaxEmailRetries: after this many failed email attempts the
// reminder is marked failed and abandoned (QA-03).
const reminderMaxEmailRetries = 3

// sendEmail is injected so tests can simulate SMTP failures.
var sendEmail = func(to, subject, htmlBody string) error {
	return utils.SendEmail(to, subject, htmlBody)
}

// claimDueReminders atomically marks due pending reminders as "sending".
// Concurrent scanners (multiple instances) race on the conditional UPDATE:
// only the winner's RowsAffected counts, so each reminder is processed
// exactly once (QA-03B).
func claimDueReminders(now time.Time, limit int) ([]models.Reminder, error) {
	var due []models.Reminder
	if err := database.DB.Where("status = ? AND remind_at <= ?", "pending", now).
		Order("remind_at ASC").Limit(limit).Find(&due).Error; err != nil {
		return nil, err
	}
	if len(due) == 0 {
		return nil, nil
	}
	ids := make([]uint, 0, len(due))
	for _, r := range due {
		ids = append(ids, r.ID)
	}
	res := database.DB.Model(&models.Reminder{}).
		Where("id IN ? AND status = ?", ids, "pending").
		Update("status", "sending")
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil // another scanner claimed them first
	}
	var claimed []models.Reminder
	if err := database.DB.Preload("User").Where("id IN ? AND status = ?", ids, "sending").Find(&claimed).Error; err != nil {
		return nil, err
	}
	return claimed, nil
}

// ClaimDueRemindersForTest exposes the claim step for service tests.
func ClaimDueRemindersForTest(now time.Time, limit int) ([]models.Reminder, error) {
	return claimDueReminders(now, limit)
}

// SetSendEmailForTest replaces the mail sender and returns the previous one.
func SetSendEmailForTest(fn func(to, subject, htmlBody string) error) func(to, subject, htmlBody string) error {
	prev := sendEmail
	sendEmail = fn
	return prev
}

func ScanAndSendReminders() {
	now := time.Now().UTC()
	reminders, err := claimDueReminders(now, 100)
	if err != nil || len(reminders) == 0 {
		return
	}
	for _, r := range reminders {
		match, err := repositories.GetMatchByID(r.MatchID)
		if err != nil {
			// Transient lookup failure: requeue (unclaimed) instead of
			// leaving the reminder stuck in "sending" forever (QA-03 review).
			requeueTransient(&r, err)
			continue
		}
		title := fmt.Sprintf("比赛提醒：%s vs %s", match.HomeTeam.Name, match.AwayTeam.Name)
		content := fmt.Sprintf("比赛即将在 %d 分钟后开始", r.RemindBeforeMinutes)

		// In-app notification is created only on the FIRST attempt;
		// retries after a failed email must not duplicate it.
		if r.RetryCount == 0 {
			notification := &models.Notification{
				UserID:  r.UserID,
				Title:   title,
				Content: content,
				Type:    "reminder",
			}
			if err := repositories.CreateNotification(notification); err != nil {
				// Notification was NOT created: requeue WITHOUT counting a
				// retry, so the next attempt re-creates it (RetryCount only
				// tracks email attempts — QA-03 review).
				log.Printf("Notification create failed for reminder %d: %v", r.ID, err)
				requeueTransient(&r, err)
				continue
			}
		}

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
</div>`, title, content, match.HomeTeam.Name+" vs "+match.AwayTeam.Name, match.Stadium.Name)
			if err := sendEmail(toEmail, subject, htmlBody); err != nil {
				// QA-03B: failed email -> requeue for retry; give up after
				// the cap with a terminal "failed" status.
				log.Printf("Email send failed for reminder %d: %v", r.ID, err)
				requeueReminder(&r, err)
				continue
			}
		}

		r.Status = "sent"
		_ = repositories.UpdateReminder(&r)
	}
}

// requeueReminder moves a failed email reminder back to pending (or to
// failed at the retry cap), logging the reason. RetryCount counts EMAIL
// attempts only.
func requeueReminder(r *models.Reminder, cause error) {
	r.RetryCount++
	if r.RetryCount >= reminderMaxEmailRetries {
		r.Status = "failed"
		log.Printf("Reminder %d failed permanently after %d attempts: %v", r.ID, r.RetryCount, cause)
	} else {
		r.Status = "pending"
	}
	_ = repositories.UpdateReminder(r)
}

// requeueTransient moves a reminder back to pending WITHOUT touching
// RetryCount: transient failures (match lookup, notification insert) must
// be retried from scratch, including re-creating the notification.
func requeueTransient(r *models.Reminder, cause error) {
	log.Printf("Reminder %d requeued (transient): %v", r.ID, cause)
	r.Status = "pending"
	_ = repositories.UpdateReminder(r)
}
