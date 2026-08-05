package handlers_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/testutil"
)

// TestCompetitionAdminConfig covers the admin-configurable competition
// switcher: create/enable/disable and the public API honoring the status.
func TestCompetitionAdminConfig(t *testing.T) {
	testutil.Setup(t)
	defer testutil.ResetDB(t)
	admin := testutil.CreateUser(t, "compadmin@test.dev", "admin", "active")
	tok := testutil.TokenFor(t, admin)

	// Seed a known active competition (like seedCompetitions does).
	if err := database.DB.Create(&models.Competition{
		Code: "PL", Name: "英超", NameEn: "Premier League", Country: "England",
		Format: "league", Season: 2026, Status: "active", SortOrder: 1,
	}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	// Create a new league through the admin API.
	w := perform(t, http.MethodPost, "/api/admin/competitions", tok, map[string]any{
		"code": "epl2", "name": "测试联赛", "name_en": "Test League",
		"country": "Testland", "format": "league", "season": 2026,
		"status": "active", "sort_order": 9,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("create: %d, body = %s", w.Code, w.Body.String())
	}
	_, _, data := read(t, w)
	compID := uint(data["id"].(float64))

	// Duplicate code -> 409 (client error, not 500).
	w = perform(t, http.MethodPost, "/api/admin/competitions", tok, map[string]any{
		"code": "epl2", "name": "Duplicate",
	})
	if w.Code != http.StatusConflict {
		t.Fatalf("duplicate code: %d, want 409", w.Code)
	}

	// Public list contains the new league.
	w = perform(t, http.MethodGet, "/api/competitions", "", nil)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"code":"EPL2"`) {
		t.Fatalf("public list missing new league: %s", w.Body.String())
	}

	// Disable it via the admin API.
	w = perform(t, http.MethodPut, "/api/admin/competitions/"+strconv.FormatUint(uint64(compID), 10), tok,
		map[string]any{"status": "inactive"})
	if w.Code != http.StatusOK {
		t.Fatalf("disable: %d, body = %s", w.Code, w.Body.String())
	}

	// Public list no longer contains it; its standings endpoint is 404.
	w = perform(t, http.MethodGet, "/api/competitions", "", nil)
	if strings.Contains(w.Body.String(), `"code":"EPL2"`) {
		t.Fatalf("disabled league still in public list: %s", w.Body.String())
	}
	w = perform(t, http.MethodGet, "/api/competitions/EPL2/standings", "", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("disabled league standings: %d, want 404", w.Code)
	}

	// Invalid status value is rejected (400) instead of silently hiding it.
	w = perform(t, http.MethodPut, "/api/admin/competitions/"+strconv.FormatUint(uint64(compID), 10), tok,
		map[string]any{"status": "enabled"})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid status: %d, want 400", w.Code)
	}

	// Re-enable restores visibility.
	w = perform(t, http.MethodPut, "/api/admin/competitions/"+strconv.FormatUint(uint64(compID), 10), tok,
		map[string]any{"status": "active", "sort_order": 2})
	if w.Code != http.StatusOK {
		t.Fatalf("enable: %d", w.Code)
	}
	w = perform(t, http.MethodGet, "/api/competitions", "", nil)
	if !strings.Contains(w.Body.String(), `"code":"EPL2"`) {
		t.Fatalf("re-enabled league missing: %s", w.Body.String())
	}
}
