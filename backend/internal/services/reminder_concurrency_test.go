package services_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/services"
	"worldcup-mate/internal/testutil"
)

// seedReminder inserts a due reminder owned by user with the given channel.
func seedReminder(t *testing.T, userID, matchID uint, channel string, remindAt time.Time) *models.Reminder {
	t.Helper()
	r := models.Reminder{
		UserID: userID, MatchID: matchID, Channel: channel,
		RemindBeforeMinutes: 30, RemindAt: remindAt, Status: "pending",
	}
	if err := database.DB.Create(&r).Error; err != nil {
		t.Fatalf("seed reminder: %v", err)
	}
	return &r
}

func reminderCount(t *testing.T, status string) int64 {
	t.Helper()
	var n int64
	if err := database.DB.Model(&models.Reminder{}).Where("status = ?", status).Count(&n).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	return n
}

func notificationCount(t *testing.T) int64 {
	t.Helper()
	var n int64
	if err := database.DB.Model(&models.Notification{}).Count(&n).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	return n
}

// TestClaimDueRemindersSingleWinner: the conditional UPDATE yields exactly
// one winner; the loser claims zero rows.
func TestClaimDueRemindersSingleWinner(t *testing.T) {
	testutil.SetupServices(t)
	user := testutil.CreateUser(t, "claim@test.dev", "user", "active")
	home := testutil.CreateTeam(t, "ClaimHome")
	away := testutil.CreateTeam(t, "ClaimAway")
	match := testutil.CreateMatch(t, home.ID, away.ID)

	due := time.Now().UTC().Add(-time.Minute)
	seedReminder(t, user.ID, match.ID, "site", due)

	claimed, err := services.ClaimDueRemindersForTest(due.Add(time.Second), 100)
	if err != nil {
		t.Fatalf("first claim: %v", err)
	}
	if len(claimed) != 1 {
		t.Fatalf("first claim won %d reminders, want 1", len(claimed))
	}
	// Second scanner running at the same time wins nothing.
	again, err := services.ClaimDueRemindersForTest(due.Add(time.Second), 100)
	if err != nil {
		t.Fatalf("second claim: %v", err)
	}
	if len(again) != 0 {
		t.Fatalf("second claim won %d reminders, want 0", len(again))
	}
}

// TestScanAndSendConcurrentScanners: two scanner goroutines run together;
// every due reminder is processed exactly once (one notification each).
func TestScanAndSendConcurrentScanners(t *testing.T) {
	testutil.SetupServices(t)
	user := testutil.CreateUser(t, "conc@test.dev", "user", "active")
	home := testutil.CreateTeam(t, "ConcHome")
	away := testutil.CreateTeam(t, "ConcAway")
	match := testutil.CreateMatch(t, home.ID, away.ID)

	due := time.Now().UTC().Add(-time.Minute)
	for i := 0; i < 3; i++ {
		seedReminder(t, user.ID, match.ID, "site", due)
	}

	// Real concurrency: both scanners race the claim UPDATE. The in-memory
	// SQLite serializes writes, so the race manifests as RowsAffected
	// winner-take-all — exactly the multi-instance semantics.
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			services.ScanAndSendReminders()
		}()
	}
	wg.Wait()

	if got := reminderCount(t, "sent"); got != 3 {
		t.Fatalf("sent reminders = %d, want 3 (each exactly once)", got)
	}
	if got := notificationCount(t); got != 3 {
		t.Fatalf("notifications = %d, want 3 (no duplicates)", got)
	}
}

// TestScanWithMissingMatchRequeues: a reminder whose match disappeared is
// requeued to pending instead of being stuck in "sending" (QA-03 review).
func TestScanWithMissingMatchRequeues(t *testing.T) {
	testutil.SetupServices(t)
	user := testutil.CreateUser(t, "missing@test.dev", "user", "active")
	home := testutil.CreateTeam(t, "GoneHome")
	away := testutil.CreateTeam(t, "GoneAway")
	match := testutil.CreateMatch(t, home.ID, away.ID)
	due := time.Now().UTC().Add(-time.Minute)
	seedReminder(t, user.ID, match.ID, "site", due)

	if err := database.DB.Delete(match).Error; err != nil {
		t.Fatalf("delete match: %v", err)
	}
	services.ScanAndSendReminders()

	var r models.Reminder
	if err := database.DB.First(&r).Error; err != nil {
		t.Fatalf("fetch reminder: %v", err)
	}
	if r.Status != "pending" {
		t.Fatalf("status = %q, want pending (not stuck in sending)", r.Status)
	}
	if r.RetryCount != 0 {
		t.Fatalf("retry_count = %d, want 0 (transient requeue must not count)", r.RetryCount)
	}
}
func TestEmailFailureRetriesThenFails(t *testing.T) {
	testutil.SetupServices(t)
	user := testutil.CreateUser(t, "mail@test.dev", "user", "active")
	user.NotificationEmail = "mail@test.dev"
	if err := database.DB.Save(user).Error; err != nil {
		t.Fatalf("save user: %v", err)
	}
	home := testutil.CreateTeam(t, "MailHome")
	away := testutil.CreateTeam(t, "MailAway")
	match := testutil.CreateMatch(t, home.ID, away.ID)
	due := time.Now().UTC().Add(-time.Minute)
	seedReminder(t, user.ID, match.ID, "email", due)

	orig := services.SetSendEmailForTest(func(to, subject, html string) error {
		return errors.New("smtp down")
	})
	defer func() { services.SetSendEmailForTest(orig) }()

	// Scan 3 times: each attempt fails, retry count climbs to the cap.
	for i := 0; i < 3; i++ {
		services.ScanAndSendReminders()
	}
	var r models.Reminder
	if err := database.DB.First(&r).Error; err != nil {
		t.Fatalf("fetch reminder: %v", err)
	}
	if r.Status != "failed" {
		t.Fatalf("status = %q, want failed", r.Status)
	}
	if r.RetryCount != 3 {
		t.Fatalf("retry_count = %d, want 3", r.RetryCount)
	}
	// The in-app notification was created once on the first attempt only.
	if got := notificationCount(t); got != 1 {
		t.Fatalf("notifications = %d, want 1 (no duplicates on retry)", got)
	}

	// A subsequent scan leaves the failed reminder untouched.
	services.ScanAndSendReminders()
	if got := reminderCount(t, "failed"); got != 1 {
		t.Fatalf("failed reminders = %d, want 1", got)
	}
}
