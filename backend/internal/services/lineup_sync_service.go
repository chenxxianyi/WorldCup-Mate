package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"worldcup-mate/internal/models"
	"worldcup-mate/internal/providers/apifootball"
	"worldcup-mate/internal/providers/footballdata"
	"worldcup-mate/internal/repositories"

	"gorm.io/gorm"
)

const syncResourceLineups = "lineups"

type LineupSyncConfig struct {
	Enabled             bool
	PrimaryProvider     string
	EnhancedProvider    string
	FootballDataAPIKey  string
	FootballDataBaseURL string
	APIFootballKey      string
	APIFootballBaseURL  string
	PregameWindow       time.Duration
	LiveInterval        time.Duration
}

type LineupSyncResult struct {
	Provider    string    `json:"provider"`
	Resource    string    `json:"resource"`
	Reason      string    `json:"reason"`
	Matches     int       `json:"matches"`
	Available   int       `json:"available"`
	Partial     int       `json:"partial"`
	Unavailable int       `json:"unavailable"`
	Updated     int       `json:"updated"`
	Skipped     int       `json:"skipped"`
	Failed      int       `json:"failed"`
	StartedAt   time.Time `json:"started_at"`
	FinishedAt  time.Time `json:"finished_at"`
}

type LineupPlayerDTO struct {
	PlayerID      *uint  `json:"player_id"`
	Name          string `json:"name"`
	NameEn        string `json:"name_en"`
	ShirtNumber   int    `json:"shirt_number"`
	Position      string `json:"position"`
	PositionLabel string `json:"position_label"`
	Role          string `json:"role"`
	Grid          string `json:"grid"`
	PhotoURL      string `json:"photo_url"`
}

type TeamLineupDTO struct {
	TeamID      uint              `json:"team_id"`
	TeamName    string            `json:"team_name"`
	TeamCode    string            `json:"team_code"`
	Side        string            `json:"side"`
	Formation   string            `json:"formation"`
	CoachName   string            `json:"coach_name"`
	Status      string            `json:"status"`
	StartXI     []LineupPlayerDTO `json:"start_xi"`
	Substitutes []LineupPlayerDTO `json:"substitutes"`
}

type MatchLineupsDTO struct {
	MatchID uint           `json:"match_id"`
	Status  string         `json:"status"`
	Source  string         `json:"source"`
	Message string         `json:"message"`
	Home    *TeamLineupDTO `json:"home"`
	Away    *TeamLineupDTO `json:"away"`
}

type ExternalMatchMappingInput struct {
	Provider         string `json:"provider"`
	ExternalMatchID  string `json:"external_match_id" binding:"required"`
	ExternalHomeID   string `json:"external_home_id"`
	ExternalAwayID   string `json:"external_away_id"`
	ExternalHomeName string `json:"external_home_name"`
	ExternalAwayName string `json:"external_away_name"`
}

var (
	lineupSyncConfig LineupSyncConfig
	lineupSyncMu     sync.RWMutex
	lineupSyncRunMu  sync.Mutex
)

func ConfigureLineupSync(cfg LineupSyncConfig) {
	if strings.TrimSpace(cfg.PrimaryProvider) == "" {
		cfg.PrimaryProvider = syncProviderFootballData
	}
	if strings.TrimSpace(cfg.EnhancedProvider) == "" {
		cfg.EnhancedProvider = syncProviderAPIFootball
	}
	if strings.TrimSpace(cfg.FootballDataBaseURL) == "" {
		cfg.FootballDataBaseURL = "https://api.football-data.org/v4"
	}
	if strings.TrimSpace(cfg.APIFootballBaseURL) == "" {
		cfg.APIFootballBaseURL = "https://v3.football.api-sports.io"
	}
	if cfg.PregameWindow <= 0 {
		cfg.PregameWindow = 90 * time.Minute
	}
	if cfg.LiveInterval <= 0 {
		cfg.LiveInterval = 15 * time.Minute
	}
	lineupSyncMu.Lock()
	lineupSyncConfig = cfg
	lineupSyncMu.Unlock()
}

func GetLineupSyncConfig() LineupSyncConfig {
	lineupSyncMu.RLock()
	defer lineupSyncMu.RUnlock()
	return lineupSyncConfig
}

func IsLineupSyncEnabled() bool {
	cfg := GetLineupSyncConfig()
	return cfg.Enabled && (strings.TrimSpace(cfg.FootballDataAPIKey) != "" || strings.TrimSpace(cfg.APIFootballKey) != "")
}

func NextLineupSyncInterval() time.Duration {
	return GetLineupSyncConfig().LiveInterval
}

func GetMatchLineups(matchID uint) (*MatchLineupsDTO, error) {
	if _, err := repositories.GetMatchByID(matchID); err != nil {
		return nil, err
	}

	lineups, err := repositories.GetLineupsByMatch(matchID)
	if err != nil {
		return nil, err
	}
	res := &MatchLineupsDTO{
		MatchID: matchID,
		Status:  "unavailable",
		Message: "首发阵容暂未公布",
	}
	if len(lineups) == 0 {
		return res, nil
	}

	statuses := make([]string, 0, len(lineups))
	for _, lineup := range lineups {
		team := buildTeamLineupDTO(lineup)
		statuses = append(statuses, lineup.Status)
		if res.Source == "" || lineup.Source == syncProviderAPIFootball {
			res.Source = lineup.Source
		}
		switch lineup.Side {
		case "home":
			res.Home = team
		case "away":
			res.Away = team
		}
	}

	res.Status = combinedLineupStatus(statuses)
	if res.Status != "unavailable" {
		res.Message = ""
	}
	return res, nil
}

func SyncMatchLineups(ctx context.Context, matchID uint, reason string) (*LineupSyncResult, error) {
	lineupSyncRunMu.Lock()
	defer lineupSyncRunMu.Unlock()

	cfg := GetLineupSyncConfig()
	if cfg.EnhancedProvider == syncProviderAPIFootball && strings.TrimSpace(cfg.APIFootballKey) != "" {
		result, err := syncAPIFootballMatchLineups(ctx, cfg, matchID, reason)
		if err == nil && result.Updated > 0 {
			return result, nil
		}
	}
	return syncFootballDataMatchLineups(ctx, cfg, matchID, reason)
}

func SyncAPIFootballMatchLineups(ctx context.Context, matchID uint, reason string) (*LineupSyncResult, error) {
	lineupSyncRunMu.Lock()
	defer lineupSyncRunMu.Unlock()
	cfg := GetLineupSyncConfig()
	return syncAPIFootballMatchLineups(ctx, cfg, matchID, reason)
}

func SyncLiveWindowLineups(ctx context.Context, reason string) (*LineupSyncResult, error) {
	lineupSyncRunMu.Lock()
	defer lineupSyncRunMu.Unlock()

	cfg := GetLineupSyncConfig()
	if !cfg.Enabled {
		return nil, errors.New("lineup sync is disabled")
	}

	startedAt := time.Now().UTC()
	result := &LineupSyncResult{
		Provider:  cfg.PrimaryProvider,
		Resource:  syncResourceLineups,
		Reason:    reason,
		StartedAt: startedAt,
	}
	_ = saveLineupSyncState(cfg.PrimaryProvider, "running", "", startedAt, nextLineupSyncAt(cfg))

	matches, err := repositories.ListMatchesNeedingLineupSync(time.Now().UTC(), cfg.PregameWindow)
	if err != nil {
		_ = saveLineupSyncState(cfg.PrimaryProvider, "failed", err.Error(), startedAt, nextLineupSyncAt(cfg))
		return nil, err
	}
	result.Matches = len(matches)

	var firstErr error
	for i, match := range matches {
		matchResult, err := syncMatchLineupsUnlocked(ctx, cfg, match.ID, reason)
		if err != nil {
			result.Failed++
			if firstErr == nil {
				firstErr = err
			}
		} else {
			result.Available += matchResult.Available
			result.Partial += matchResult.Partial
			result.Unavailable += matchResult.Unavailable
			result.Updated += matchResult.Updated
			result.Skipped += matchResult.Skipped
		}
		if i < len(matches)-1 {
			select {
			case <-ctx.Done():
				err := ctx.Err()
				_ = saveLineupSyncState(cfg.PrimaryProvider, "failed", err.Error(), startedAt, nextLineupSyncAt(cfg))
				return nil, err
			case <-time.After(800 * time.Millisecond):
			}
		}
	}

	result.FinishedAt = time.Now().UTC()
	if result.Failed > 0 && result.Updated == 0 {
		_ = saveLineupSyncState(cfg.PrimaryProvider, "failed", firstErr.Error(), result.FinishedAt, nextLineupSyncAt(cfg))
		return result, firstErr
	}
	status := "success"
	errText := ""
	if result.Failed > 0 {
		status = "partial"
		errText = firstErr.Error()
	}
	_ = saveLineupSyncState(cfg.PrimaryProvider, status, errText, result.FinishedAt, nextLineupSyncAt(cfg))
	return result, nil
}

func EnsureAPIFootballMatchMapping(ctx context.Context, matchID uint) (*models.ExternalMatchMapping, error) {
	cfg := GetLineupSyncConfig()
	return ensureAPIFootballMatchMapping(ctx, cfg, matchID)
}

func GetMatchExternalMapping(matchID uint, provider string) (*models.ExternalMatchMapping, error) {
	if strings.TrimSpace(provider) == "" {
		provider = syncProviderAPIFootball
	}
	return repositories.GetExternalMatchMapping(matchID, provider)
}

func UpsertMatchExternalMapping(matchID uint, input ExternalMatchMappingInput) (*models.ExternalMatchMapping, error) {
	provider := strings.TrimSpace(input.Provider)
	if provider == "" {
		provider = syncProviderAPIFootball
	}
	mapping := &models.ExternalMatchMapping{
		MatchID:          matchID,
		Provider:         provider,
		ExternalMatchID:  strings.TrimSpace(input.ExternalMatchID),
		ExternalHomeID:   strings.TrimSpace(input.ExternalHomeID),
		ExternalAwayID:   strings.TrimSpace(input.ExternalAwayID),
		ExternalHomeName: strings.TrimSpace(input.ExternalHomeName),
		ExternalAwayName: strings.TrimSpace(input.ExternalAwayName),
		MatchedBy:        "manual",
	}
	if mapping.ExternalMatchID == "" {
		return nil, errors.New("external_match_id is required")
	}
	if err := repositories.UpsertExternalMatchMapping(mapping); err != nil {
		return nil, err
	}
	return repositories.GetExternalMatchMapping(matchID, provider)
}

func MarkInterruptedLineupSync() error {
	return repositories.MarkInterruptedSyncState(
		GetLineupSyncConfig().PrimaryProvider,
		syncResourceLineups,
		"previous sync was interrupted before completion",
		time.Now().UTC(),
	)
}

func syncMatchLineupsUnlocked(ctx context.Context, cfg LineupSyncConfig, matchID uint, reason string) (*LineupSyncResult, error) {
	if cfg.EnhancedProvider == syncProviderAPIFootball && strings.TrimSpace(cfg.APIFootballKey) != "" {
		if result, err := syncAPIFootballMatchLineups(ctx, cfg, matchID, reason); err == nil && result.Updated > 0 {
			return result, nil
		}
	}
	return syncFootballDataMatchLineups(ctx, cfg, matchID, reason)
}

func syncFootballDataMatchLineups(ctx context.Context, cfg LineupSyncConfig, matchID uint, reason string) (*LineupSyncResult, error) {
	if !cfg.Enabled {
		return nil, errors.New("lineup sync is disabled")
	}
	if cfg.PrimaryProvider != syncProviderFootballData {
		return nil, fmt.Errorf("unsupported lineup provider: %s", cfg.PrimaryProvider)
	}
	if strings.TrimSpace(cfg.FootballDataAPIKey) == "" {
		return nil, errors.New("FOOTBALL_DATA_API_KEY is empty")
	}

	startedAt := time.Now().UTC()
	result := &LineupSyncResult{
		Provider:  syncProviderFootballData,
		Resource:  syncResourceLineups,
		Reason:    reason,
		Matches:   1,
		StartedAt: startedAt,
	}
	_ = saveLineupSyncState(syncProviderFootballData, "running", "", startedAt, nextLineupSyncAt(cfg))

	match, err := repositories.GetMatchByID(matchID)
	if err != nil {
		_ = saveLineupSyncState(syncProviderFootballData, "failed", err.Error(), startedAt, nextLineupSyncAt(cfg))
		return nil, err
	}
	if match.ExternalProvider != syncProviderFootballData || strings.TrimSpace(match.ExternalID) == "" {
		err := errors.New("match has no football-data external id")
		_ = saveLineupSyncState(syncProviderFootballData, "failed", err.Error(), startedAt, nextLineupSyncAt(cfg))
		return nil, err
	}

	client := footballdata.NewClient(cfg.FootballDataBaseURL, cfg.FootballDataAPIKey)
	detail, err := client.Match(ctx, match.ExternalID)
	if err != nil {
		_ = saveLineupSyncState(syncProviderFootballData, "failed", err.Error(), startedAt, nextLineupSyncAt(cfg))
		return nil, err
	}

	now := time.Now().UTC()
	homeStatus, homeUpdated, err := saveFootballDataTeamLineup(match, "home", match.HomeTeamID, detail.HomeTeam, &now)
	if err != nil {
		result.Failed++
	}
	awayStatus, awayUpdated, err := saveFootballDataTeamLineup(match, "away", match.AwayTeamID, detail.AwayTeam, &now)
	if err != nil {
		result.Failed++
	}
	result.Updated = homeUpdated + awayUpdated
	countLineupStatus(result, homeStatus)
	countLineupStatus(result, awayStatus)
	result.FinishedAt = time.Now().UTC()

	if result.Failed > 0 && result.Updated == 0 {
		err := errors.New("failed to save football-data lineups")
		_ = saveLineupSyncState(syncProviderFootballData, "failed", err.Error(), result.FinishedAt, nextLineupSyncAt(cfg))
		return result, err
	}
	_ = saveLineupSyncState(syncProviderFootballData, "success", "", result.FinishedAt, nextLineupSyncAt(cfg))
	return result, nil
}

func syncAPIFootballMatchLineups(ctx context.Context, cfg LineupSyncConfig, matchID uint, reason string) (*LineupSyncResult, error) {
	if !cfg.Enabled {
		return nil, errors.New("lineup sync is disabled")
	}
	if strings.TrimSpace(cfg.APIFootballKey) == "" {
		return nil, errors.New("API_FOOTBALL_KEY is empty")
	}

	startedAt := time.Now().UTC()
	result := &LineupSyncResult{
		Provider:  syncProviderAPIFootball,
		Resource:  syncResourceLineups,
		Reason:    reason,
		Matches:   1,
		StartedAt: startedAt,
	}
	_ = saveLineupSyncState(syncProviderAPIFootball, "running", "", startedAt, nextLineupSyncAt(cfg))

	match, err := repositories.GetMatchByID(matchID)
	if err != nil {
		_ = saveLineupSyncState(syncProviderAPIFootball, "failed", err.Error(), startedAt, nextLineupSyncAt(cfg))
		return nil, err
	}
	mapping, err := ensureAPIFootballMatchMapping(ctx, cfg, matchID)
	if err != nil {
		_ = saveLineupSyncState(syncProviderAPIFootball, "failed", err.Error(), startedAt, nextLineupSyncAt(cfg))
		return nil, err
	}

	client := apifootball.NewClient(cfg.APIFootballBaseURL, cfg.APIFootballKey)
	data, err := client.FixtureLineups(ctx, mapping.ExternalMatchID)
	if err != nil {
		_ = saveLineupSyncState(syncProviderAPIFootball, "failed", err.Error(), startedAt, nextLineupSyncAt(cfg))
		return nil, err
	}
	if len(data.Response) == 0 {
		result.Unavailable = 2
		result.Skipped = 1
		result.FinishedAt = time.Now().UTC()
		_ = saveLineupSyncState(syncProviderAPIFootball, "success", "", result.FinishedAt, nextLineupSyncAt(cfg))
		return result, nil
	}

	now := time.Now().UTC()
	for _, entry := range data.Response {
		side, teamID := apiFootballLineupSide(match, mapping, entry.Team.ID)
		if side == "" || teamID == 0 {
			result.Skipped++
			continue
		}
		status, updated, err := saveAPIFootballTeamLineup(match, mapping, side, teamID, entry, &now)
		if err != nil {
			result.Failed++
			continue
		}
		result.Updated += updated
		countLineupStatus(result, status)
	}
	result.FinishedAt = time.Now().UTC()
	if result.Failed > 0 && result.Updated == 0 {
		err := errors.New("failed to save api-football lineups")
		_ = saveLineupSyncState(syncProviderAPIFootball, "failed", err.Error(), result.FinishedAt, nextLineupSyncAt(cfg))
		return result, err
	}
	_ = saveLineupSyncState(syncProviderAPIFootball, "success", "", result.FinishedAt, nextLineupSyncAt(cfg))
	return result, nil
}

func saveFootballDataTeamLineup(match *models.Match, side string, teamID uint, team footballdata.MatchLineupTeam, syncedAt *time.Time) (string, int, error) {
	status := lineupStatus(len(team.Lineup), len(team.Bench))
	lineup := &models.MatchLineup{
		MatchID:        match.ID,
		TeamID:         teamID,
		Side:           side,
		Formation:      strings.TrimSpace(team.Formation),
		Source:         syncProviderFootballData,
		SourceMatchID:  match.ExternalID,
		ExternalTeamID: strconv.FormatInt(team.ID, 10),
		Status:         status,
		LastSyncedAt:   syncedAt,
	}
	saved, err := repositories.UpsertMatchLineup(lineup)
	if err != nil {
		return status, 0, err
	}
	players := footballDataLineupPlayers(match.ID, teamID, team)
	if len(players) == 0 {
		return status, 0, nil
	}
	if err := repositories.ReplaceLineupPlayers(saved.ID, players); err != nil {
		return status, 0, err
	}
	return status, 1, nil
}

func saveAPIFootballTeamLineup(match *models.Match, mapping *models.ExternalMatchMapping, side string, teamID uint, entry apifootball.FixtureLineupEntry, syncedAt *time.Time) (string, int, error) {
	status := lineupStatus(len(entry.StartXI), len(entry.Substitutes))
	lineup := &models.MatchLineup{
		MatchID:        match.ID,
		TeamID:         teamID,
		Side:           side,
		Formation:      strings.TrimSpace(entry.Formation),
		CoachName:      strings.TrimSpace(entry.Coach.Name),
		CoachNameEn:    strings.TrimSpace(entry.Coach.Name),
		Source:         syncProviderAPIFootball,
		SourceMatchID:  mapping.ExternalMatchID,
		ExternalTeamID: strconv.FormatInt(entry.Team.ID, 10),
		Status:         status,
		LastSyncedAt:   syncedAt,
	}
	saved, err := repositories.UpsertMatchLineup(lineup)
	if err != nil {
		return status, 0, err
	}
	players := apiFootballLineupPlayers(match.ID, teamID, entry)
	if len(players) == 0 {
		return status, 0, nil
	}
	if err := repositories.ReplaceLineupPlayers(saved.ID, players); err != nil {
		return status, 0, err
	}
	return status, 1, nil
}

func footballDataLineupPlayers(matchID, teamID uint, team footballdata.MatchLineupTeam) []models.MatchLineupPlayer {
	players := make([]models.MatchLineupPlayer, 0, len(team.Lineup)+len(team.Bench))
	appendPlayer := func(role string, offset int, player footballdata.LineupPlayer) {
		nameEn := strings.TrimSpace(player.Name)
		if nameEn == "" {
			return
		}
		position, positionLabel := normalizeLineupPosition(player.Position)
		sourcePlayerID := ""
		if player.ID > 0 {
			sourcePlayerID = strconv.FormatInt(player.ID, 10)
		}
		localPlayer := findLineupPlayer(teamID, "", nameEn, player.ShirtNumber, position)
		players = append(players, models.MatchLineupPlayer{
			MatchID:        matchID,
			TeamID:         teamID,
			PlayerID:       localPlayerID(localPlayer),
			Source:         syncProviderFootballData,
			SourcePlayerID: sourcePlayerID,
			Name:           displayLineupName(localPlayer, nameEn),
			NameEn:         displayLineupNameEn(localPlayer, nameEn),
			ShirtNumber:    displayShirtNumber(localPlayer, player.ShirtNumber),
			Position:       displayPosition(localPlayer, position),
			PositionLabel:  displayPositionLabel(localPlayer, positionLabel),
			Role:           role,
			SortOrder:      offset,
			PhotoURL:       displayPhotoURL(localPlayer, ""),
		})
	}
	for i, player := range team.Lineup {
		appendPlayer("start_xi", i+1, player)
	}
	for i, player := range team.Bench {
		appendPlayer("substitute", i+1, player)
	}
	return players
}

func apiFootballLineupPlayers(matchID, teamID uint, entry apifootball.FixtureLineupEntry) []models.MatchLineupPlayer {
	players := make([]models.MatchLineupPlayer, 0, len(entry.StartXI)+len(entry.Substitutes))
	appendPlayer := func(role string, offset int, row apifootball.FixtureLineupPlayerRow) {
		player := row.Player
		nameEn := strings.TrimSpace(player.Name)
		if nameEn == "" {
			return
		}
		position, positionLabel := normalizeLineupPosition(player.Pos)
		sourcePlayerID := ""
		if player.ID > 0 {
			sourcePlayerID = strconv.FormatInt(player.ID, 10)
		}
		localPlayer := findLineupPlayer(teamID, sourcePlayerID, nameEn, player.Number, position)
		players = append(players, models.MatchLineupPlayer{
			MatchID:        matchID,
			TeamID:         teamID,
			PlayerID:       localPlayerID(localPlayer),
			Source:         syncProviderAPIFootball,
			SourcePlayerID: sourcePlayerID,
			Name:           displayLineupName(localPlayer, nameEn),
			NameEn:         displayLineupNameEn(localPlayer, nameEn),
			ShirtNumber:    displayShirtNumber(localPlayer, player.Number),
			Position:       displayPosition(localPlayer, position),
			PositionLabel:  displayPositionLabel(localPlayer, positionLabel),
			Role:           role,
			Grid:           strings.TrimSpace(player.Grid),
			SortOrder:      offset,
			PhotoURL:       displayPhotoURL(localPlayer, ""),
		})
	}
	for i, row := range entry.StartXI {
		appendPlayer("start_xi", i+1, row)
	}
	for i, row := range entry.Substitutes {
		appendPlayer("substitute", i+1, row)
	}
	return players
}

func ensureAPIFootballMatchMapping(ctx context.Context, cfg LineupSyncConfig, matchID uint) (*models.ExternalMatchMapping, error) {
	mapping, err := repositories.GetExternalMatchMapping(matchID, syncProviderAPIFootball)
	if err == nil {
		return mapping, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if strings.TrimSpace(cfg.APIFootballKey) == "" {
		return nil, errors.New("API_FOOTBALL_KEY is empty")
	}

	match, err := repositories.GetMatchByID(matchID)
	if err != nil {
		return nil, err
	}
	homeMapping, err := repositories.GetExternalTeamMapping(match.HomeTeamID, syncProviderAPIFootball)
	if err != nil {
		return nil, fmt.Errorf("home external team mapping not found: %w", err)
	}
	awayMapping, err := repositories.GetExternalTeamMapping(match.AwayTeamID, syncProviderAPIFootball)
	if err != nil {
		return nil, fmt.Errorf("away external team mapping not found: %w", err)
	}

	client := apifootball.NewClient(cfg.APIFootballBaseURL, cfg.APIFootballKey)
	fixtures, err := client.SearchFixtures(ctx, match.KickoffTimeUTC.UTC().Format("2006-01-02"), match.KickoffTimeUTC.UTC().Year(), homeMapping.ExternalTeamID)
	if err != nil {
		return nil, err
	}

	var best *apifootball.FixtureSearchEntry
	bestDelta := math.MaxFloat64
	for i := range fixtures.Response {
		candidate := &fixtures.Response[i]
		if strconv.FormatInt(candidate.Teams.Home.ID, 10) != homeMapping.ExternalTeamID {
			continue
		}
		if strconv.FormatInt(candidate.Teams.Away.ID, 10) != awayMapping.ExternalTeamID {
			continue
		}
		date, err := time.Parse(time.RFC3339, candidate.Fixture.Date)
		if err != nil {
			continue
		}
		delta := math.Abs(date.UTC().Sub(match.KickoffTimeUTC.UTC()).Seconds())
		if delta <= (2*time.Hour).Seconds() && delta < bestDelta {
			best = candidate
			bestDelta = delta
		}
	}
	if best == nil {
		return nil, errors.New("api-football fixture mapping candidate not found")
	}

	created := &models.ExternalMatchMapping{
		MatchID:          match.ID,
		Provider:         syncProviderAPIFootball,
		ExternalMatchID:  strconv.FormatInt(best.Fixture.ID, 10),
		ExternalHomeID:   strconv.FormatInt(best.Teams.Home.ID, 10),
		ExternalAwayID:   strconv.FormatInt(best.Teams.Away.ID, 10),
		ExternalHomeName: best.Teams.Home.Name,
		ExternalAwayName: best.Teams.Away.Name,
		MatchedBy:        "auto_date_team_time",
	}
	if err := repositories.UpsertExternalMatchMapping(created); err != nil {
		return nil, err
	}
	return repositories.GetExternalMatchMapping(match.ID, syncProviderAPIFootball)
}

func buildTeamLineupDTO(lineup models.MatchLineup) *TeamLineupDTO {
	dto := &TeamLineupDTO{
		TeamID:      lineup.TeamID,
		TeamName:    lineup.Team.Name,
		TeamCode:    lineup.Team.FIFACode,
		Side:        lineup.Side,
		Formation:   lineup.Formation,
		CoachName:   lineup.CoachName,
		Status:      lineup.Status,
		StartXI:     []LineupPlayerDTO{},
		Substitutes: []LineupPlayerDTO{},
	}
	for _, player := range lineup.Players {
		value := buildLineupPlayerDTO(player)
		if player.Role == "substitute" {
			dto.Substitutes = append(dto.Substitutes, value)
		} else {
			dto.StartXI = append(dto.StartXI, value)
		}
	}
	return dto
}

func buildLineupPlayerDTO(value models.MatchLineupPlayer) LineupPlayerDTO {
	name := value.Name
	nameEn := value.NameEn
	photoURL := value.PhotoURL
	position := value.Position
	positionLabel := value.PositionLabel
	shirtNumber := value.ShirtNumber
	playerID := value.PlayerID
	if value.Player != nil {
		if value.Player.Name != "" {
			name = value.Player.Name
		}
		if value.Player.NameEn != "" {
			nameEn = value.Player.NameEn
		}
		if value.Player.PhotoURL != "" {
			photoURL = value.Player.PhotoURL
		}
		if value.Player.Position != "" {
			position = value.Player.Position
		}
		if value.Player.PositionLabel != "" {
			positionLabel = value.Player.PositionLabel
		}
		if value.Player.ShirtNumber > 0 {
			shirtNumber = value.Player.ShirtNumber
		}
	}
	return LineupPlayerDTO{
		PlayerID:      playerID,
		Name:          name,
		NameEn:        nameEn,
		ShirtNumber:   shirtNumber,
		Position:      position,
		PositionLabel: positionLabel,
		Role:          value.Role,
		Grid:          value.Grid,
		PhotoURL:      photoURL,
	}
}

func combinedLineupStatus(statuses []string) string {
	if len(statuses) == 0 {
		return "unavailable"
	}
	available := 0
	partial := 0
	for _, status := range statuses {
		switch status {
		case "available":
			available++
		case "partial":
			partial++
		}
	}
	if available == len(statuses) {
		return "available"
	}
	if available > 0 || partial > 0 {
		return "partial"
	}
	return "unavailable"
}

func lineupStatus(starters, substitutes int) string {
	if starters >= 11 {
		return "available"
	}
	if starters > 0 || substitutes > 0 {
		return "partial"
	}
	return "unavailable"
}

func countLineupStatus(result *LineupSyncResult, status string) {
	switch status {
	case "available":
		result.Available++
	case "partial":
		result.Partial++
	default:
		result.Unavailable++
	}
}

func normalizeLineupPosition(value string) (string, string) {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "G", "GK", "GOALKEEPER":
		return "GK", "门将"
	case "D", "DF", "DEFENDER":
		return "DF", "后卫"
	case "M", "MF", "MIDFIELDER":
		return "MF", "中场"
	case "F", "FW", "ATTACKER", "FORWARD":
		return "FW", "前锋"
	default:
		return strings.TrimSpace(value), ""
	}
}

func findLineupPlayer(teamID uint, sourcePlayerID, nameEn string, shirtNumber int, position string) *models.Player {
	if sourcePlayerID != "" {
		if player, err := repositories.GetPlayerBySource(teamID, syncProviderAPIFootball, sourcePlayerID); err == nil {
			return player
		}
	}
	if player, err := repositories.FindPlayerForLineup(teamID, nameEn, shirtNumber, position); err == nil {
		return player
	}
	return nil
}

func localPlayerID(player *models.Player) *uint {
	if player == nil {
		return nil
	}
	return &player.ID
}

func displayLineupName(player *models.Player, fallback string) string {
	if player != nil && player.Name != "" {
		return player.Name
	}
	return strings.TrimSpace(fallback)
}

func displayLineupNameEn(player *models.Player, fallback string) string {
	if player != nil && player.NameEn != "" {
		return player.NameEn
	}
	return strings.TrimSpace(fallback)
}

func displayShirtNumber(player *models.Player, fallback int) int {
	if player != nil && player.ShirtNumber > 0 {
		return player.ShirtNumber
	}
	return fallback
}

func displayPosition(player *models.Player, fallback string) string {
	if player != nil && player.Position != "" {
		return player.Position
	}
	return fallback
}

func displayPositionLabel(player *models.Player, fallback string) string {
	if player != nil && player.PositionLabel != "" {
		return player.PositionLabel
	}
	return fallback
}

func displayPhotoURL(player *models.Player, fallback string) string {
	if player != nil && player.PhotoURL != "" {
		return player.PhotoURL
	}
	return strings.TrimSpace(fallback)
}

func apiFootballLineupSide(match *models.Match, mapping *models.ExternalMatchMapping, externalTeamID int64) (string, uint) {
	value := strconv.FormatInt(externalTeamID, 10)
	if value == mapping.ExternalHomeID {
		return "home", match.HomeTeamID
	}
	if value == mapping.ExternalAwayID {
		return "away", match.AwayTeamID
	}
	return "", 0
}

func saveLineupSyncState(provider, status, lastError string, lastSyncedAt time.Time, next *time.Time) error {
	state := &models.SyncState{
		Provider:     provider,
		Resource:     syncResourceLineups,
		Status:       status,
		LastSyncedAt: &lastSyncedAt,
		NextSyncAt:   next,
		LastError:    lastError,
	}
	return repositories.UpsertSyncState(state)
}

func nextLineupSyncAt(cfg LineupSyncConfig) *time.Time {
	next := time.Now().UTC().Add(cfg.LiveInterval)
	return &next
}
