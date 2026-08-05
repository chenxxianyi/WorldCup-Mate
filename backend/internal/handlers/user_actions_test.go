package handlers_test

import (
	"net/http"
	"strconv"
	"testing"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/testutil"
)

// TestNotificationReadFlows covers mark-one-read, mark-all-read and the
// unread counter (scoped to the caller's own notifications).
func TestNotificationReadFlows(t *testing.T) {
	testutil.Setup(t)
	defer testutil.ResetDB(t)
	user := testutil.CreateUser(t, "notif@test.dev", "user", "active")
	tok := testutil.TokenFor(t, user)

	for i := 0; i < 3; i++ {
		if err := database.DB.Create(&models.Notification{
			UserID: user.ID, Title: "n", Content: "x", Type: "reminder",
		}).Error; err != nil {
			t.Fatalf("seed notification: %v", err)
		}
	}
	// Unread counter: 3.
	w := perform(t, http.MethodGet, "/api/notifications/unread-count", tok, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("unread count: %d", w.Code)
	}
	_, _, data := read(t, w)
	if int(data["count"].(float64)) != 3 {
		t.Fatalf("unread count = %v, want 3", data["count"])
	}

	// Mark one read -> 2.
	var n models.Notification
	if err := database.DB.Where("user_id = ?", user.ID).First(&n).Error; err != nil {
		t.Fatalf("fetch notification: %v", err)
	}
	firstID := n.ID
	w = perform(t, http.MethodPut, "/api/notifications/"+strconv.FormatUint(uint64(firstID), 10)+"/read", tok, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("mark read: %d", w.Code)
	}
	w = perform(t, http.MethodGet, "/api/notifications/unread-count", tok, nil)
	_, _, data = read(t, w)
	if int(data["count"].(float64)) != 2 {
		t.Fatalf("unread after one read = %v, want 2", data["count"])
	}

	// Mark all read -> 0.
	w = perform(t, http.MethodPut, "/api/notifications/read-all", tok, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("read all: %d", w.Code)
	}
	w = perform(t, http.MethodGet, "/api/notifications/unread-count", tok, nil)
	_, _, data = read(t, w)
	if int(data["count"].(float64)) != 0 {
		t.Fatalf("unread after read-all = %v, want 0", data["count"])
	}

	// Another user cannot read or touch someone else's notification.
	other := testutil.CreateUser(t, "notif2@test.dev", "user", "active")
	otherTok := testutil.TokenFor(t, other)
	w = perform(t, http.MethodPut, "/api/notifications/"+strconv.FormatUint(uint64(firstID), 10)+"/read", otherTok, nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("cross-user mark-read: status = %d, want 404", w.Code)
	}
}

// TestReminderBatchAndScoping covers the batch reminder endpoint and the
// "cannot create a reminder for a match that is already live/finished" rule.
func TestReminderBatchAndScoping(t *testing.T) {
	testutil.Setup(t)
	defer testutil.ResetDB(t)
	user := testutil.CreateUser(t, "rmd@test.dev", "user", "active")
	tok := testutil.TokenFor(t, user)

	home := seedTeam(t, "BatchHome")
	away := seedTeam(t, "BatchAway")
	match := seedMatch(t, home.ID, away.ID)
	match.Status = "scheduled"
	if err := database.DB.Save(match).Error; err != nil {
		t.Fatalf("save match: %v", err)
	}

	// Batch of three lead times.
	w := perform(t, http.MethodPost, "/api/reminders/batch", tok, map[string]any{
		"match_id": match.ID, "minutes": []int{15, 30, 60},
	})
	if w.Code != http.StatusOK {
		t.Fatalf("batch create: %d, body = %s", w.Code, w.Body.String())
	}
	var count int64
	database.DB.Model(&models.Reminder{}).Where("user_id = ? AND match_id = ?", user.ID, match.ID).Count(&count)
	if count != 3 {
		t.Fatalf("batch reminders = %d, want 3", count)
	}

	// A live match must reject new reminders (single and batch).
	match.Status = "live"
	if err := database.DB.Save(match).Error; err != nil {
		t.Fatalf("save match: %v", err)
	}
	w = perform(t, http.MethodPost, "/api/reminders", tok, map[string]any{
		"matchId": match.ID, "remindBeforeMinutes": 30,
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("reminder on live match: status = %d, want 400, body = %s", w.Code, w.Body.String())
	}
	w = perform(t, http.MethodPost, "/api/reminders/batch", tok, map[string]any{
		"match_id": match.ID, "minutes": []int{15},
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("batch reminder on live match: status = %d, want 400, body = %s", w.Code, w.Body.String())
	}
}
