package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/providers/footballdata"
	"worldcup-mate/internal/repositories"

	"gorm.io/gorm"
)

// League sync is a brand-new, additive pipeline for the top-5 leagues.
// It intentionally does NOT touch the existing World Cup sync path
// (sync_service.go): different config, different state rows
// (resource = "matches:<CODE>"), and its own upsert helpers.

const (
	leagueSyncProvider = "football-data"
	leagueStage        = "regular_season"
)

type LeagueSyncTarget struct {
	Code   string
	Season int
}

type LeagueSyncConfig struct {
	Targets  []LeagueSyncTarget
	Interval time.Duration
	BaseURL  string
	APIKey   string
}

var (
	leagueSyncCfg   LeagueSyncConfig
	leagueSyncMu    sync.RWMutex
	leagueSyncRun   sync.Mutex
	requestThrottle = 9 * time.Second // 10 req/min free tier: 10 reqs over ~81s ≈ 7.4/min, leaving headroom for the WC pipeline
)

// ConfigureLeagueSync parses SYNC_COMPETITIONS like "PL:2025,PD:2025,BL1:2025,SA:2025,FL1:2025".
// Empty targets disable the league sync entirely (default = legacy behavior).
func ConfigureLeagueSync(raw string, interval time.Duration, baseURL, apiKey string) {
	cfg := LeagueSyncConfig{
		Interval: interval,
		BaseURL:  baseURL,
		APIKey:   apiKey,
	}
	if cfg.Interval <= 0 {
		cfg.Interval = 30 * time.Minute
	}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		fields := strings.Split(part, ":")
		code := strings.ToUpper(strings.TrimSpace(fields[0]))
		if code == "" {
			continue
		}
		target := LeagueSyncTarget{Code: code}
		if len(fields) > 1 {
			if season, err := strconv.Atoi(strings.TrimSpace(fields[1])); err == nil && season > 0 {
				target.Season = season
			}
		}
		cfg.Targets = append(cfg.Targets, target)
	}
	leagueSyncMu.Lock()
	leagueSyncCfg = cfg
	leagueSyncMu.Unlock()
}

func GetLeagueSyncConfig() LeagueSyncConfig {
	leagueSyncMu.RLock()
	defer leagueSyncMu.RUnlock()
	return leagueSyncCfg
}

func IsLeagueSyncEnabled() bool {
	cfg := GetLeagueSyncConfig()
	return len(cfg.Targets) > 0 && cfg.APIKey != ""
}

func NextLeagueSyncInterval() time.Duration {
	return GetLeagueSyncConfig().Interval
}

type LeagueSyncResult struct {
	Provider   string    `json:"provider"`
	Resource   string    `json:"resource"`
	Reason     string    `json:"reason"`
	Total      int       `json:"total"`
	Created    int       `json:"created"`
	Updated    int       `json:"updated"`
	Skipped    int       `json:"skipped"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
}

// SyncAllLeagues runs a sync round for every configured league.
func SyncAllLeagues(ctx context.Context, reason string) []*LeagueSyncResult {
	cfg := GetLeagueSyncConfig()
	if !IsLeagueSyncEnabled() {
		return nil
	}
	var results []*LeagueSyncResult
	for _, target := range cfg.Targets {
		leagueSyncRun.Lock()
		result, _ := syncLeague(ctx, cfg, target, reason)
		leagueSyncRun.Unlock()
		if result != nil {
			results = append(results, result)
		}
		// Free tier allows 10 req/min; leave headroom before the next league.
		time.Sleep(requestThrottle)
	}
	return results
}

// SyncLeague triggers a sync for a single league (used by the admin endpoint).
func SyncLeague(ctx context.Context, code string, season int, reason string) (*LeagueSyncResult, error) {
	cfg := GetLeagueSyncConfig()
	if !IsLeagueSyncEnabled() {
		return nil, errors.New("league sync is disabled")
	}
	leagueSyncRun.Lock()
	defer leagueSyncRun.Unlock()
	target := LeagueSyncTarget{Code: strings.ToUpper(code), Season: season}
	return syncLeague(ctx, cfg, target, reason)
}

func syncLeague(ctx context.Context, cfg LeagueSyncConfig, target LeagueSyncTarget, reason string) (*LeagueSyncResult, error) {
	competition, err := repositories.GetCompetitionByCode(target.Code)
	if err != nil {
		_ = saveLeagueSyncState(target.Code, "failed", "unknown competition: "+err.Error())
		return nil, err
	}
	season := target.Season
	if season <= 0 {
		season = competition.Season
	}
	resource := "matches:" + target.Code

	// REL-06: cross-instance distributed lock (manual + scheduled syncs
	// share the same key). Redis outage degrades to the process-local
	// mutex only — sync is not on a critical path.
	lockKey := syncLockKey(leagueSyncProvider, resource, target.Code, season)
	lockOwner := randomLockOwner()
	acquired, lockErr := tryAcquireSyncLock(ctx, lockKey, lockOwner, syncLockTTL)
	if lockErr != nil {
		log.Printf("[sync-lock] redis unavailable, running with process-local lock only: %v", lockErr)
	} else if !acquired {
		return nil, fmt.Errorf("%w: %s season %d", ErrSyncAlreadyRunning, target.Code, season)
	} else {
		defer releaseSyncLock(context.Background(), lockKey, lockOwner)
	}

	startedAt := time.Now().UTC()
	next := startedAt.Add(cfg.Interval)
	_ = saveLeagueSyncStateDetail(resource, "running", "", startedAt, &next)

	result := &LeagueSyncResult{
		Provider:  leagueSyncProvider,
		Resource:  resource,
		Reason:    reason,
		StartedAt: startedAt,
	}

	client := footballdata.NewClient(cfg.BaseURL, cfg.APIKey)
	data, err := client.CompetitionMatches(ctx, target.Code, season)
	if err != nil {
		next := startedAt.Add(cfg.Interval)
		_ = saveLeagueSyncStateDetail(resource, "failed", err.Error(), startedAt, &next)
		return nil, err
	}

	result.Total = len(data.Matches)
	for _, externalMatch := range data.Matches {
		action, err := upsertLeagueMatch(competition, season, externalMatch)
		if err != nil {
			result.Skipped++
			continue
		}
		switch action {
		case "created":
			result.Created++
		case "updated":
			result.Updated++
		default:
			result.Skipped++
		}
	}
	if result.Created+result.Updated > 0 {
		_ = repositories.DeleteSeedDemoMatches()
	}

	// Official standings (TOTAL / HOME / AWAY) into league_standings.
	time.Sleep(requestThrottle)
	if err := SyncLeagueStandings(ctx, target.Code, season); err != nil {
		result.Skipped++ // standings failure is not fatal for the match sync
	}

	result.FinishedAt = time.Now().UTC()
	next = result.FinishedAt.Add(cfg.Interval)
	_ = saveLeagueSyncStateDetail(resource, "success", "", result.FinishedAt, &next)
	return result, nil
}

// upsertLeagueMatch creates or updates a league match using the new nullable
// columns on matches (competition_id / season / matchday). Clubs are upserted
// by external_code (tla); group/stage/venue semantics are league-flavored.
func upsertLeagueMatch(competition *models.Competition, season int, externalMatch footballdata.Match) (string, error) {
	externalID := strconv.FormatInt(externalMatch.ID, 10)
	if externalID == "0" {
		return "skipped", errors.New("missing external match id")
	}

	homeTeam, err := findOrCreateClub(externalMatch.HomeTeam, competition.Country)
	if err != nil {
		return "skipped", err
	}
	awayTeam, err := findOrCreateClub(externalMatch.AwayTeam, competition.Country)
	if err != nil {
		return "skipped", err
	}

	var stadiumID, cityID uint
	if strings.TrimSpace(externalMatch.Venue) != "" {
		stadiumID, cityID, err = findOrCreateLeagueStadium(externalMatch.Venue)
	} else {
		// matches.city_id / stadium_id have NOT NULL FK constraints
		// (fk_matches_city / fk_matches_stadium); point them at the
		// fallback TBD venue when the provider omits the venue.
		cityID, stadiumID, err = ensureFallbackVenue()
	}
	if err != nil {
		return "skipped", err
	}

	now := time.Now().UTC()
	status := normalizeExternalStatus(externalMatch.Status)
	homeScore, awayScore := scoreForStatus(externalMatch)
	winnerTeamID := winnerID(externalMatch.Score.Winner, homeTeam.ID, awayTeam.ID)

	competitionID := competition.ID
	matchday := externalMatch.Matchday
	existing, err := repositories.GetMatchByExternal(leagueSyncProvider, externalID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "skipped", err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		match := models.Match{
			MatchNo:          int(externalMatch.ID),
			HomeTeamID:       homeTeam.ID,
			AwayTeamID:       awayTeam.ID,
			Stage:            leagueStage,
			CompetitionID:    &competitionID,
			Season:           &season,
			Matchday:         &matchday,
			StadiumID:        stadiumID,
			CityID:           cityID,
			KickoffTimeUTC:   externalMatch.UTCDate.UTC(),
			HomeScore:        homeScore,
			AwayScore:        awayScore,
			Status:           status,
			StatusDetail:     externalMatch.Status,
			WinnerTeamID:     winnerTeamID,
			ExternalProvider: leagueSyncProvider,
			ExternalID:       externalID,
			LastSyncedAt:     &now,
		}
		if err := repositories.CreateMatch(&match); err != nil {
			return "skipped", err
		}
		return "created", nil
	}

	existing.HomeTeamID = homeTeam.ID
	existing.AwayTeamID = awayTeam.ID
	existing.Stage = leagueStage
	existing.CompetitionID = &competitionID
	existing.Season = &season
	existing.Matchday = &matchday
	existing.StadiumID = stadiumID
	existing.CityID = cityID
	existing.KickoffTimeUTC = externalMatch.UTCDate.UTC()
	existing.HomeScore = homeScore
	existing.AwayScore = awayScore
	existing.Status = status
	existing.StatusDetail = externalMatch.Status
	existing.WinnerTeamID = winnerTeamID
	existing.LastSyncedAt = &now
	if err := repositories.UpdateMatch(existing); err != nil {
		return "skipped", err
	}
	return "updated", nil
}

// findOrCreateClub upserts a club by the provider's immutable team ID.
// TLA values are display codes only: they are not globally unique (for
// example Barcelona and Bayern can both use FCB).
func findOrCreateClub(team footballdata.Team, country string) (*models.Team, error) {
	code := strings.ToUpper(strings.TrimSpace(team.TLA))
	if code == "" {
		code = fmt.Sprintf("TBD%d", team.ID)
	}
	externalID := strconv.FormatInt(team.ID, 10)
	name := strings.TrimSpace(team.ShortName)
	if name == "" {
		name = strings.TrimSpace(team.Name)
	}
	if name == "" {
		name = "TBD"
	}
	name = localizeFootballDataClub(externalID, name)
	normalizedCountry := strings.TrimSpace(country)

	var existing models.Team
	result := database.DB.Where("external_provider = ? AND external_id = ?", leagueSyncProvider, externalID).Limit(1).Find(&existing)
	if result.Error != nil {
		return nil, result.Error
	}
	// One-time migration path for clubs imported before provider IDs were
	// stored. Country scopes the legacy TLA lookup and prevents collisions.
	if result.RowsAffected == 0 {
		result = database.DB.Where("external_code = ? AND country = ?", code, normalizedCountry).Limit(1).Find(&existing)
		if result.Error != nil {
			return nil, result.Error
		}
	}
	if result.RowsAffected > 0 {
		existing.Name = name
		existing.NameEn = team.Name
		existing.TeamType = "club"
		existing.Country = normalizedCountry
		existing.ExternalProvider = leagueSyncProvider
		existing.ExternalID = &externalID
		existing.ExternalCode = &code
		// Provider crest URLs are authoritative. Always refresh them so a club
		// that once collided by TLA does not retain another club's badge.
		if team.Crest != "" {
			existing.FlagURL = team.Crest
		}
		return &existing, database.DB.Save(&existing).Error
	}

	created := models.Team{
		Name:             name,
		NameEn:           team.Name,
		ExternalCode:     &code,
		ExternalProvider: leagueSyncProvider,
		ExternalID:       &externalID,
		TeamType:         "club",
		FlagURL:          team.Crest,
		Country:          normalizedCountry,
	}
	if err := repositories.CreateTeam(&created); err != nil {
		return nil, err
	}
	return &created, nil
}

// findOrCreateLeagueStadium resolves a venue name into a stadium row and
// returns (stadiumID, cityID). The stadium is attached to the shared "TBD"
// city (the matches.city_id FK requires a valid city row).
func findOrCreateLeagueStadium(venue string) (uint, uint, error) {
	existing, err := repositories.GetStadiumByName(venue)
	if err == nil {
		return existing.ID, existing.CityID, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, err
	}

	cityID, _, err := ensureFallbackVenue()
	if err != nil {
		return 0, 0, err
	}
	created := models.Stadium{
		Name:   venue,
		NameEn: venue,
		CityID: cityID,
	}
	if err := repositories.CreateStadium(&created); err != nil {
		return 0, 0, err
	}
	return created.ID, cityID, nil
}

func saveLeagueSyncState(code, status, lastError string) error {
	return saveLeagueSyncStateDetail("matches:"+code, status, lastError, time.Now().UTC(), nil)
}

func saveLeagueSyncStateDetail(resource, status, lastError string, lastSyncedAt time.Time, next *time.Time) error {
	state := &models.SyncState{
		Provider:     leagueSyncProvider,
		Resource:     resource,
		Status:       status,
		LastSyncedAt: &lastSyncedAt,
		NextSyncAt:   next,
		LastError:    lastError,
	}
	return repositories.UpsertSyncState(state)
}
