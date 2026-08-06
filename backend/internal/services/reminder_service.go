package services

import (
	"crypto/rand"
	"encoding/hex"
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

// claimDueReminders atomically marks due pending reminders as "sending"
// using a claim_token. Each scanner generates a unique worker ID and
// updates only pending reminders whose next_retry_at has passed. The
// UPDATE is conditional on status=pending, so concurrent scanners
// cannot overlap (REL-08). A reminder left in "sending" for longer
// than the claim timeout is reclaimed automatically on the next scan.
func claimDueReminders(now time.Time, limit int) ([]models.Reminder, error) {
	workerID := generateWorkerID()
	claimTimeout := 5 * time.Minute

	// Reclaim stale claims from crashed workers first.
	database.DB.Model(&models.Reminder{}).
		Where("status = ? AND claimed_at IS NOT NULL AND claimed_at < ?", "sending", now.Add(-claimTimeout)).
		Updates(map[string]interface{}{"status": "pending", "claim_token": "", "claimed_at": nil, "worker_id": ""})

	// Claim pending reminders atomically. The conditional WHERE clause
	// ensures only one scanner can claim each row (REL-08).
	var due []models.Reminder
	claimQuery := database.DB.Model(&models.Reminder{}).
		Where("status = ? AND (next_retry_at IS NULL OR next_retry_at <= ?) AND remind_at <= ?", "pending", now, now).
		Order("remind_at ASC").Limit(limit)

	token := generateClaimToken()
	err := claimQuery.Updates(map[string]interface{}{
		"status":      "sending",
		"claim_token": token,
		"claimed_at":  now,
		"worker_id":   workerID,
	}).Error
	if err != nil {
		return nil, err
	}

	// Fetch only the reminders claimed by this worker (token match).
	if err := database.DB.Preload("User").
		Where("claim_token = ? AND status = ?", token, "sending").
		Order("remind_at ASC").Find(&due).Error; err != nil {
		return nil, err
	}
	return due, nil
}

// generateClaimToken returns a random hex token used to scope a
// scanner's claim so it never processes another worker's reminders.
func generateClaimToken() string {
	b := make([]byte, 18)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// generateWorkerID returns a short random identifier for the current
// scanner process (used for observability in audit logs).
func generateWorkerID() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
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
		// Use NotificationID as idempotency key (REL-08).
		if r.RetryCount == 0 && (r.NotificationID == nil || *r.NotificationID == 0) {
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
			r.NotificationID = &notification.ID
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
// attempts only. On retry, apply exponential backoff via NextRetryAt so
// failed reminders are not hot-retried and appear in the admin view
// (ADM-13 / REL-08).
func requeueReminder(r *models.Reminder, cause error) {
	r.RetryCount++
	if r.RetryCount >= reminderMaxEmailRetries {
		r.Status = "failed"
		r.LastError = cause.Error()
		log.Printf("Reminder %d failed permanently after %d attempts: %v", r.ID, r.RetryCount, cause)
	} else {
		r.Status = "pending"
		// Exponential backoff: 1min, 2min, 4min.
		backoff := time.Duration(1<<uint(r.RetryCount-1)) * time.Minute
		now := time.Now().UTC()
		next := now.Add(backoff)
		r.NextRetryAt = &next
		r.LastError = cause.Error()
		log.Printf("Reminder %d requeued for retry %d with backoff %v: %v", r.ID, r.RetryCount, backoff, cause)
	}
	_ = repositories.UpdateReminder(r)
}

// requeueTransient moves a reminder back to pending WITHOUT touching
// RetryCount: transient failures (match lookup, notification insert) must
// be retried from scratch, including re-creating the notification.
func requeueTransient(r *models.Reminder, cause error) {
	log.Printf("Reminder %d requeued (transient): %v", r.ID, cause)
	r.Status = "pending"
	r.LastError = cause.Error()
	now := time.Now().UTC()
	r.NextRetryAt = &now
	_ = repositories.UpdateReminder(r)
}
