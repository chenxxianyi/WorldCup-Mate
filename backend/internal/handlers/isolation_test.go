package handlers_test

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/services"
	"worldcup-mate/internal/testutil"
)

// seedTeam inserts a minimal team row and returns it.
func seedTeam(t *testing.T, name string) *models.Team {
	t.Helper()
	team := models.Team{Name: name, NameEn: name, TeamType: "club", ExternalProvider: "football-data"}
	if err := database.DB.Create(&team).Error; err != nil {
		t.Fatalf("seed team: %v", err)
	}
	return &team
}

// seedMatch inserts a minimal scheduled match row.
func seedMatch(t *testing.T, home, away uint) *models.Match {
	t.Helper()
	m := models.Match{
		MatchNo:         int(time.Now().UnixNano() % 1e9),
		Stage:           "Group Stage",
		Status:          "scheduled",
		ImportanceLevel: 2,
		HomeTeamID:      home,
		AwayTeamID:      away,
		KickoffTimeUTC:  time.Now().UTC().Add(24 * time.Hour),
	}
	if err := database.DB.Create(&m).Error; err != nil {
		t.Fatalf("seed match: %v", err)
	}
	return &m
}

// TestUserDataIsolation verifies one user can never see another user's
// favorites, reminders, or notifications (SEC-01).
func TestUserDataIsolation(t *testing.T) {
	testutil.Setup(t)
	defer testutil.ResetDB(t)
	alice := testutil.CreateUser(t, "alice@test.dev", "user", "active")
	bob := testutil.CreateUser(t, "bob@test.dev", "user", "active")
	aliceTok, bobTok := testutil.TokenFor(t, alice), testutil.TokenFor(t, bob)

	// Alice favorites a team; Bob's list must stay empty.
	team := seedTeam(t, "FCTest")
	w := perform(t, http.MethodPost, "/api/favorites/teams/"+strconv.Itoa(int(team.ID)), aliceTok, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("favorite status = %d", w.Code)
	}
	// ListFavoriteTeams returns a bare array in `data`; assert emptiness
	// directly on the body to avoid envelope-decoding complexity.
	bobList := perform(t, http.MethodGet, "/api/favorites/teams", bobTok, nil)
	if bobList.Code != http.StatusOK || !strings.Contains(bobList.Body.String(), `"data":[]`) {
		t.Fatalf("bob must not see alice's favorite teams: %s", bobList.Body.String())
	}
	// Alice sets a reminder; Bob's list must stay empty.
	match := seedMatch(t, team.ID, seedTeam(t, "FCVisit").ID)
	_ = perform(t, http.MethodPost, "/api/reminders", aliceTok, map[string]any{
		"matchId": match.ID, "remindBeforeMinutes": 30,
	})
	bobReminders := perform(t, http.MethodGet, "/api/reminders", bobTok, nil)
	if bobReminders.Code != http.StatusOK || !strings.Contains(bobReminders.Body.String(), `"data":[]`) {
		t.Fatalf("bob must not see alice's reminders: %s", bobReminders.Body.String())
	}
	// Alice gets a notification (inserted by the system); Bob's list empty.
	if err := database.DB.Create(&models.Notification{
		UserID: alice.ID, Title: "match start", Content: "x", Type: "reminder",
	}).Error; err != nil {
		t.Fatalf("seed notification: %v", err)
	}
	code, _, data := read(t, perform(t, http.MethodGet, "/api/notifications?page=1&page_size=20", bobTok, nil))
	list, ok := data["list"].([]any)
	if code != 0 || !ok || len(list) != 0 {
		t.Fatalf("bob sees alice's notifications: %v", data)
	}
}

// TestRoleEnforcement verifies admin endpoints reject normal users with 403
// and unauthenticated requests with 401.
func TestRoleEnforcement(t *testing.T) {
	testutil.Setup(t)
	defer testutil.ResetDB(t)
	user := testutil.CreateUser(t, "plain@test.dev", "user", "active")
	admin := testutil.CreateUser(t, "admin@test.dev", "admin", "active")

	// Unauthenticated.
	if w := perform(t, http.MethodGet, "/api/admin/teams?page=1&page_size=10", "", nil); w.Code != http.StatusUnauthorized {
		t.Fatalf("no token: status = %d, want 401", w.Code)
	}
	// Normal user.
	for _, path := range []string{
		"/api/admin/teams?page=1&page_size=10",
		"/api/admin/dashboard",
		"/api/admin/users?page=1&page_size=10",
	} {
		w := perform(t, http.MethodGet, path, testutil.TokenFor(t, user), nil)
		if w.Code != http.StatusForbidden {
			t.Fatalf("user GET %s: status = %d, want 403", path, w.Code)
		}
	}
	// Admin.
	w := perform(t, http.MethodGet, "/api/admin/teams?page=1&page_size=10", testutil.TokenFor(t, admin), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("admin GET teams: status = %d, want 200", w.Code)
	}
}

// TestSyncLockConflict returns 409 when another worker holds the league
// sync lock (REL-06 follow-up test).
func TestSyncLockConflict(t *testing.T) {
	testutil.Setup(t)
	defer testutil.ResetDB(t)
	admin := testutil.CreateUser(t, "syncadmin@test.dev", "admin", "active")
	token := testutil.TokenFor(t, admin)

	// Enable the league sync pipeline (env config is absent in tests).
	services.ConfigureLeagueSync("PL:2025", 30*time.Minute, "http://localhost:8080", "test-key")

	// The league must exist before the sync reaches the lock check.
	if err := database.DB.Create(&models.Competition{
		Code: "PL", Name: "Premier League", NameEn: "Premier League",
		Country: "England", Format: "league", Season: 2025, Status: "active",
	}).Error; err != nil {
		t.Fatalf("seed competition: %v", err)
	}

	// A concurrent worker holds the lock for PL/2025 (same key format as
	// syncLockKey("football-data", "matches:PL", "PL", 2025)).
	ctx := context.Background()
	if err := database.RDB.Set(ctx, "lock:sync:football-data:matches:PL:PL:2025", "other-worker", 10*time.Minute).Err(); err != nil {
		t.Fatalf("seed lock: %v", err)
	}

	w := perform(t, http.MethodPost, "/api/admin/sync/matches?code=PL&season=2025", token, nil)
	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409 (lock held), body = %s", w.Code, w.Body.String())
	}
}
