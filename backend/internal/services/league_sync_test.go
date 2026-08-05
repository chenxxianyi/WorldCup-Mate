package services_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/services"
	"worldcup-mate/internal/testutil"
)

// sampleMatchesJSON is a minimal valid football-data matches response.
const sampleMatchesJSON = `{
  "count": 1,
  "matches": [{
    "id": 424242,
    "utcDate": "2026-09-01T19:00:00Z",
    "status": "SCHEDULED",
    "matchday": 1,
    "stage": "REGULAR_SEASON",
    "group": "",
    "venue": "Emirates Stadium",
    "homeTeam": {"id": 57, "name": "Arsenal FC", "shortName": "Arsenal", "tla": "ARS", "crest": ""},
    "awayTeam": {"id": 61, "name": "Chelsea FC", "shortName": "Chelsea", "tla": "CHE", "crest": ""},
    "score": {
      "winner": null, "duration": "REGULAR",
      "fullTime": {"home": null, "away": null},
      "regularTime": {"home": null, "away": null},
      "halfTime": {"home": null, "away": null}
    }
  }]
}`

// TestLeagueSyncIsIdempotent: syncing the same league twice upserts instead
// of duplicating teams/matches (REL-04 + QA-03A).
func TestLeagueSyncIsIdempotent(t *testing.T) {
	testutil.SetupServices(t)

	var gotAuthHeader string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuthHeader = r.Header.Get("X-Auth-Token")
		if r.URL.Path != "/competitions/PL/matches" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, sampleMatchesJSON)
	}))
	defer srv.Close()

	services.ConfigureLeagueSync("PL:2025", 30*60*1e9, srv.URL, "qa3-test-key")
	if err := database.DB.Create(&models.Competition{
		Code: "PL", Name: "Premier League", NameEn: "Premier League",
		Country: "England", Format: "league", Season: 2025, Status: "active",
	}).Error; err != nil {
		t.Fatalf("seed competition: %v", err)
	}

	ctx := context.Background()
	first, err := services.SyncLeague(ctx, "PL", 2025, "manual")
	if err != nil {
		t.Fatalf("first sync: %v", err)
	}
	if first.Created != 1 {
		t.Fatalf("first sync created = %d, want 1 (teams+match upserted)", first.Created)
	}
	if gotAuthHeader != "qa3-test-key" {
		t.Fatalf("auth header = %q, want the configured API key", gotAuthHeader)
	}

	var teams, matches int64
	database.DB.Model(&models.Team{}).Where("external_provider = ?", "football-data").Count(&teams)
	database.DB.Model(&models.Match{}).Where("external_provider = ?", "football-data").Count(&matches)
	if teams != 2 || matches != 1 {
		t.Fatalf("after first sync: teams=%d matches=%d, want 2/1", teams, matches)
	}

	// Second sync: nothing new created, one match updated.
	second, err := services.SyncLeague(ctx, "PL", 2025, "manual")
	if err != nil {
		t.Fatalf("second sync: %v", err)
	}
	if second.Created != 0 {
		t.Fatalf("second sync created = %d, want 0", second.Created)
	}
	if second.Updated != 1 {
		t.Fatalf("second sync updated = %d, want 1", second.Updated)
	}
	database.DB.Model(&models.Team{}).Where("external_provider = ?", "football-data").Count(&teams)
	database.DB.Model(&models.Match{}).Where("external_provider = ?", "football-data").Count(&matches)
	if teams != 2 || matches != 1 {
		t.Fatalf("after second sync: teams=%d matches=%d, want 2/1 (idempotent)", teams, matches)
	}
}

// TestLeagueSyncTransientFailureRetries: a 503 then a 200 must succeed via
// the retry/backoff layer (REL-05 + QA-03A).
func TestLeagueSyncTransientFailureRetries(t *testing.T) {
	testutil.SetupServices(t)

	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, sampleMatchesJSON)
	}))
	defer srv.Close()

	services.ConfigureLeagueSync("PL:2025", 30*60*1e9, srv.URL, "qa3-test-key")
	if err := database.DB.Create(&models.Competition{
		Code: "PL", Name: "Premier League", NameEn: "Premier League",
		Country: "England", Format: "league", Season: 2025, Status: "active",
	}).Error; err != nil {
		t.Fatalf("seed competition: %v", err)
	}

	_, err := services.SyncLeague(context.Background(), "PL", 2025, "manual")
	if err != nil {
		t.Fatalf("sync after transient failure: %v", err)
	}
	if calls < 2 {
		t.Fatalf("provider calls = %d, want >= 2 (retried after 503)", calls)
	}
	if !strings.Contains(errOrEmpty(""), "") {
		t.Fatal("unreachable")
	}
}

// errOrEmpty is a tiny helper kept for symmetry in assertions above.
func errOrEmpty(s string) string { return s }
