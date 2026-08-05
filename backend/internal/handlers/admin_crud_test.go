package handlers_test

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/testutil"
)

// adminToken creates an admin user and returns its token.
func adminToken(t *testing.T) string {
	t.Helper()
	admin := testutil.CreateUser(t, "crudadmin@test.dev", "admin", "active")
	return testutil.TokenFor(t, admin)
}

// TestAdminTeamCRUD covers create/409/update/delete/404 for teams and the
// "cannot delete a referenced team" rule (ADM-03).
func TestAdminTeamCRUD(t *testing.T) {
	testutil.Setup(t)
	defer testutil.ResetDB(t)
	tok := adminToken(t)

	// Create (national team with FIFA code).
	fifa := "TST"
	w := perform(t, http.MethodPost, "/api/admin/teams", tok, map[string]any{
		"name": "Testland", "name_en": "Testland", "fifa_code": fifa,
		"team_type": "national", "continent": "Europe", "country": "Testland",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("create team: status = %d, body = %s", w.Code, w.Body.String())
	}
	code, _, data := read(t, w)
	if code != 0 {
		t.Fatalf("create team business code = %d", code)
	}
	teamID := uint(data["id"].(float64))

	// Duplicate FIFA code -> 409.
	w = perform(t, http.MethodPost, "/api/admin/teams", tok, map[string]any{
		"name": "Otherland", "name_en": "Otherland", "fifa_code": fifa, "team_type": "national",
	})
	if w.Code != http.StatusConflict {
		t.Fatalf("duplicate fifa_code: status = %d, want 409", w.Code)
	}

	// Update colliding with another team's code -> 409.
	other := seedTeam(t, "FreeTeam")
	freeCode := "FR2"
	other.FIFACode = &freeCode
	if err := database.DB.Save(other).Error; err != nil {
		t.Fatalf("seed other team: %v", err)
	}
	w = perform(t, http.MethodPut, "/api/admin/teams/"+strconv.FormatUint(uint64(teamID), 10), tok,
		map[string]any{"fifa_code": freeCode})
	if w.Code != http.StatusConflict {
		t.Fatalf("update to taken fifa_code: status = %d, want 409", w.Code)
	}

	// Update name.
	w = perform(t, http.MethodPut, "/api/admin/teams/"+strconv.FormatUint(uint64(teamID), 10), tok,
		map[string]any{"name": "Testland FC"})
	if w.Code != http.StatusOK {
		t.Fatalf("update team: status = %d", w.Code)
	}

	// A referenced team cannot be deleted (create a match using it).
	away := seedTeam(t, "Awayland")
	match := seedMatch(t, teamID, away.ID)
	match.HomeTeamID = teamID
	if err := database.DB.Save(match).Error; err != nil {
		t.Fatalf("link match: %v", err)
	}
	w = perform(t, http.MethodDelete, "/api/admin/teams/"+strconv.FormatUint(uint64(teamID), 10), tok, nil)
	if w.Code != http.StatusConflict {
		t.Fatalf("delete referenced team: status = %d, want 409", w.Code)
	}

	// Unlink and delete succeeds.
	if err := database.DB.Delete(match).Error; err != nil {
		t.Fatalf("delete match: %v", err)
	}
	w = perform(t, http.MethodDelete, "/api/admin/teams/"+strconv.FormatUint(uint64(teamID), 10), tok, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("delete team: status = %d, want 200", w.Code)
	}
	// Second delete -> 404.
	w = perform(t, http.MethodDelete, "/api/admin/teams/"+strconv.FormatUint(uint64(teamID), 10), tok, nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("delete missing team: status = %d, want 404", w.Code)
	}
}

// TestAdminGroupCityStadiumCRUD covers the simple admin resources.
func TestAdminGroupCityStadiumCRUD(t *testing.T) {
	testutil.Setup(t)
	defer testutil.ResetDB(t)
	tok := adminToken(t)

	// Group.
	w := perform(t, http.MethodPost, "/api/admin/groups", tok, map[string]any{"name": "Group A", "stage": "group"})
	if w.Code != http.StatusOK {
		t.Fatalf("create group: %d", w.Code)
	}
	_, _, gdata := read(t, w)
	groupID := uint(gdata["id"].(float64))
	w = perform(t, http.MethodPut, "/api/admin/groups/"+strconv.FormatUint(uint64(groupID), 10), tok,
		map[string]any{"name": "Group B"})
	if w.Code != http.StatusOK {
		t.Fatalf("update group: %d", w.Code)
	}
	w = perform(t, http.MethodDelete, "/api/admin/groups/"+strconv.FormatUint(uint64(groupID), 10), tok, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("delete group: %d", w.Code)
	}

	// City.
	w = perform(t, http.MethodPost, "/api/admin/cities", tok, map[string]any{
		"name": "London", "country": "England", "timezone": "Europe/London",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("create city: %d", w.Code)
	}
	_, _, cdata := read(t, w)
	cityID := uint(cdata["id"].(float64))
	w = perform(t, http.MethodPut, "/api/admin/cities/"+strconv.FormatUint(uint64(cityID), 10), tok,
		map[string]any{"name": "Greater London"})
	if w.Code != http.StatusOK {
		t.Fatalf("update city: %d", w.Code)
	}

	// Stadium requires an existing city.
	w = perform(t, http.MethodPost, "/api/admin/stadiums", tok, map[string]any{
		"name": "Wembley", "city_id": cityID, "capacity": 90000,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("create stadium: %d, body = %s", w.Code, w.Body.String())
	}
	_, _, sdata := read(t, w)
	stadiumID := uint(sdata["id"].(float64))
	w = perform(t, http.MethodDelete, "/api/admin/stadiums/"+strconv.FormatUint(uint64(stadiumID), 10), tok, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("delete stadium: %d", w.Code)
	}
}

// TestAdminMatchScoreAndStandings covers match create/score and the legacy
// standings recalculation path.
func TestAdminMatchScoreAndStandings(t *testing.T) {
	testutil.Setup(t)
	defer testutil.ResetDB(t)
	tok := adminToken(t)

	// Teams in a group, group, city, stadium (legacy standings are computed
	// per group).
	w := perform(t, http.MethodPost, "/api/admin/groups", tok, map[string]any{"name": "League Table", "stage": "league"})
	_, _, grp := read(t, w)
	groupID := uint(grp["id"].(float64))

	w = perform(t, http.MethodPost, "/api/admin/teams", tok, map[string]any{
		"name": "Home United", "team_type": "club", "country": "England", "group_id": groupID,
	})
	_, _, home := read(t, w)
	homeID := uint(home["id"].(float64))
	w = perform(t, http.MethodPost, "/api/admin/teams", tok, map[string]any{
		"name": "Away City", "team_type": "club", "country": "England", "group_id": groupID,
	})
	_, _, away := read(t, w)
	awayID := uint(away["id"].(float64))

	w = perform(t, http.MethodPost, "/api/admin/cities", tok, map[string]any{
		"name": "Manchester", "country": "England", "timezone": "Europe/London",
	})
	_, _, city := read(t, w)
	cityID := uint(city["id"].(float64))
	w = perform(t, http.MethodPost, "/api/admin/stadiums", tok, map[string]any{
		"name": "Etihad", "city_id": cityID,
	})
	_, _, stadium := read(t, w)
	stadiumID := uint(stadium["id"].(float64))

	// Create match.
	w = perform(t, http.MethodPost, "/api/admin/matches", tok, map[string]any{
		"match_no": 1, "home_team_id": homeID, "away_team_id": awayID,
		"stage": "league", "group_id": groupID, "stadium_id": stadiumID, "city_id": cityID,
		"kickoff_time_utc": "2026-09-01T19:00:00Z", "importance_level": 2,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("create match: %d, body = %s", w.Code, w.Body.String())
	}
	_, _, mdata := read(t, w)
	matchID := uint(mdata["id"].(float64))

	// Finish the match: score first (the status machine refuses to finish
	// a match without a score), then the finished transition auto-recalculates
	// group standings (ADM-05).
	w = perform(t, http.MethodPut, "/api/admin/matches/"+strconv.FormatUint(uint64(matchID), 10)+"/score", tok,
		map[string]any{"home_score": 2, "away_score": 1})
	if w.Code != http.StatusOK {
		t.Fatalf("set score: %d, body = %s", w.Code, w.Body.String())
	}
	w = perform(t, http.MethodPut, "/api/admin/matches/"+strconv.FormatUint(uint64(matchID), 10)+"/status", tok,
		map[string]any{"status": "finished"})
	if w.Code != http.StatusOK {
		t.Fatalf("set status: %d, body = %s", w.Code, w.Body.String())
	}

	// Recalculate legacy standings and verify rows are produced (the
	// standings list endpoint returns a bare array in `data`).
	w = perform(t, http.MethodPost, "/api/admin/standings/recalculate", tok, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("recalculate: %d, body = %s", w.Code, w.Body.String())
	}
	w = perform(t, http.MethodGet, "/api/admin/standings?page=1&page_size=50", tok, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("list standings: %d", w.Code)
	}
	var env struct {
		Code int              `json:"code"`
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("decode standings: %v", err)
	}
	if env.Code != 0 || len(env.Data) == 0 {
		t.Fatalf("recalculate produced no standings rows: %s", w.Body.String())
	}
	// Verify the 2:1 score: winner 3 points, loser 0, played=1 each.
	byName := map[string]map[string]any{}
	for _, row := range env.Data {
		team, _ := row["team"].(map[string]any)
		name, _ := team["name"].(string)
		byName[name] = row
	}
	homeRow, ok := byName["Home United"]
	if !ok || int(homeRow["points"].(float64)) != 3 || int(homeRow["played"].(float64)) != 1 {
		t.Fatalf("home standing wrong: %v", homeRow)
	}
	awayRow, ok := byName["Away City"]
	if !ok || int(awayRow["points"].(float64)) != 0 || int(awayRow["lost"].(float64)) != 1 {
		t.Fatalf("away standing wrong: %v", awayRow)
	}
}

// TestAdminUserStatusDisable uses the admin endpoint to disable a user and
// verifies their existing token dies immediately (ADM-06).
func TestAdminUserStatusDisable(t *testing.T) {
	testutil.Setup(t)
	defer testutil.ResetDB(t)
	tok := adminToken(t)
	victim := testutil.CreateUser(t, "victim@test.dev", "user", "active")
	victimTok := testutil.TokenFor(t, victim)

	w := perform(t, http.MethodPut, "/api/admin/users/"+strconv.FormatUint(uint64(victim.ID), 10)+"/status", tok,
		map[string]any{"status": "disabled"})
	if w.Code != http.StatusOK {
		t.Fatalf("disable user: %d, body = %s", w.Code, w.Body.String())
	}

	w = perform(t, http.MethodGet, "/api/user/profile", victimTok, nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("disabled user token: status = %d, want 401", w.Code)
	}
}
