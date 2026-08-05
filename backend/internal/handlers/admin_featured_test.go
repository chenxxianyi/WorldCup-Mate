package handlers_test

import (
	"net/http"
	"strings"
	"testing"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/testutil"
)

// TestFeaturedConfigFlow covers the admin-configurable homepage hero:
// upsert copy + pinned focus match, public exposure, disable, validation.
func TestFeaturedConfigFlow(t *testing.T) {
	testutil.Setup(t)
	defer testutil.ResetDB(t)
	admin := testutil.CreateUser(t, "featadmin@test.dev", "admin", "active")
	tok := testutil.TokenFor(t, admin)

	if err := database.DB.Create(&models.Competition{
		Code: "PL", Name: "英超", NameEn: "Premier League", Country: "England",
		Format: "league", Season: 2026, Status: "active", SortOrder: 1,
	}).Error; err != nil {
		t.Fatalf("seed competition: %v", err)
	}
	compID := uint(1)
	home := seedTeam(t, "FeatHome")
	away := seedTeam(t, "FeatAway")
	match := seedMatch(t, home.ID, away.ID)
	season := 2026
	match.CompetitionID = &compID
	match.Season = &season
	if err := database.DB.Save(match).Error; err != nil {
		t.Fatalf("link match: %v", err)
	}

	// Upsert hero config with a pinned focus match.
	w := perform(t, http.MethodPut, "/api/admin/featured/PL", tok, map[string]any{
		"match_id": match.ID, "tagline": "每一轮，都有新的主角",
		"description": "悬念持续到最后一分钟", "stage_label": "MATCHWEEK 12", "enabled": true,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("upsert: %d, body = %s", w.Code, w.Body.String())
	}

	// Public endpoint exposes it under the competition code.
	w = perform(t, http.MethodGet, "/api/featured", "", nil)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"PL"`) ||
		!strings.Contains(w.Body.String(), "每一轮，都有新的主角") {
		t.Fatalf("public featured missing config: %s", w.Body.String())
	}

	// Focus-match picker lists the competition's matches.
	w = perform(t, http.MethodGet, "/api/admin/featured/PL/matches", tok, nil)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "FeatHome") {
		t.Fatalf("picker: %d, body = %s", w.Code, w.Body.String())
	}

	// A match from another competition is rejected (400).
	other := seedTeam(t, "OtherHome")
	otherMatch := seedMatch(t, other.ID, seedTeam(t, "OtherAway").ID)
	otherMatch.Season = &season
	if err := database.DB.Save(otherMatch).Error; err != nil {
		t.Fatalf("save other match: %v", err)
	}
	w = perform(t, http.MethodPut, "/api/admin/featured/PL", tok, map[string]any{
		"match_id": otherMatch.ID, "tagline": "x", "enabled": true,
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("foreign match: %d, want 400", w.Code)
	}

	// A real cross-competition match (owned by another competition) is
	// also rejected.
	if err := database.DB.Create(&models.Competition{
		Code: "PD", Name: "西甲", NameEn: "La Liga", Country: "Spain",
		Format: "league", Season: 2026, Status: "active", SortOrder: 2,
	}).Error; err != nil {
		t.Fatalf("seed PD: %v", err)
	}
	pdCompID := uint(2)
	laHome := seedTeam(t, "LaHome")
	laMatch := seedMatch(t, laHome.ID, seedTeam(t, "LaAway").ID)
	laMatch.CompetitionID = &pdCompID
	laMatch.Season = &season
	if err := database.DB.Save(laMatch).Error; err != nil {
		t.Fatalf("save la match: %v", err)
	}
	w = perform(t, http.MethodPut, "/api/admin/featured/PL", tok, map[string]any{
		"match_id": laMatch.ID, "tagline": "x", "enabled": true,
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("cross-competition match: %d, want 400", w.Code)
	}

	// Disabling removes it from the public endpoint.
	w = perform(t, http.MethodPut, "/api/admin/featured/PL", tok, map[string]any{
		"tagline": "每一轮，都有新的主角", "enabled": false,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("disable: %d", w.Code)
	}
	w = perform(t, http.MethodGet, "/api/featured", "", nil)
	if strings.Contains(w.Body.String(), `"PL"`) {
		t.Fatalf("disabled config still public: %s", w.Body.String())
	}
}
