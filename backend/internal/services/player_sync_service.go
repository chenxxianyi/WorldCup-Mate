package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"worldcup-mate/internal/models"
	"worldcup-mate/internal/providers/apifootball"
	"worldcup-mate/internal/repositories"

	"gorm.io/gorm"
)

const (
	syncProviderAPIFootball = "api-football"
	syncResourcePlayers     = "players"
)

type PlayerSyncConfig struct {
	Enabled            bool
	Provider           string
	APIFootballKey     string
	APIFootballBaseURL string
	Interval           time.Duration
	SyncOnStartup      bool
}

type PlayerSyncResult struct {
	Provider    string    `json:"provider"`
	Resource    string    `json:"resource"`
	Reason      string    `json:"reason"`
	Teams       int       `json:"teams"`
	Total       int       `json:"total"`
	Created     int       `json:"created"`
	Updated     int       `json:"updated"`
	Deactivated int64     `json:"deactivated"`
	Skipped     int       `json:"skipped"`
	StartedAt   time.Time `json:"started_at"`
	FinishedAt  time.Time `json:"finished_at"`
}

type ExternalTeamMappingInput struct {
	ExternalTeamID   string `json:"external_team_id" binding:"required"`
	ExternalTeamName string `json:"external_team_name"`
	Provider         string `json:"provider"`
}

var (
	playerSyncConfig PlayerSyncConfig
	playerSyncMu     sync.RWMutex
	playerSyncRunMu  sync.Mutex
)

func ConfigurePlayerSync(cfg PlayerSyncConfig) {
	if strings.TrimSpace(cfg.Provider) == "" {
		cfg.Provider = syncProviderAPIFootball
	}
	if strings.TrimSpace(cfg.APIFootballBaseURL) == "" {
		cfg.APIFootballBaseURL = "https://v3.football.api-sports.io"
	}
	if cfg.Interval <= 0 {
		cfg.Interval = 168 * time.Hour
	}
	playerSyncMu.Lock()
	playerSyncConfig = cfg
	playerSyncMu.Unlock()
}

func GetPlayerSyncConfig() PlayerSyncConfig {
	playerSyncMu.RLock()
	defer playerSyncMu.RUnlock()
	return playerSyncConfig
}

func IsPlayerSyncEnabled() bool {
	cfg := GetPlayerSyncConfig()
	return cfg.Enabled && strings.TrimSpace(cfg.APIFootballKey) != ""
}

func NextPlayerSyncInterval() time.Duration {
	return GetPlayerSyncConfig().Interval
}

func SyncTeamPlayersWithDefault(ctx context.Context, teamID uint, reason string) (*PlayerSyncResult, error) {
	cfg := GetPlayerSyncConfig()
	return SyncTeamPlayers(ctx, cfg, teamID, reason)
}

func SyncAllMappedTeamPlayersWithDefault(ctx context.Context, reason string) (*PlayerSyncResult, error) {
	cfg := GetPlayerSyncConfig()
	return SyncAllMappedTeamPlayers(ctx, cfg, reason)
}

func SyncTeamPlayers(ctx context.Context, cfg PlayerSyncConfig, teamID uint, reason string) (*PlayerSyncResult, error) {
	playerSyncRunMu.Lock()
	defer playerSyncRunMu.Unlock()
	return syncTeamPlayersUnlocked(ctx, cfg, teamID, reason)
}

func SyncAllMappedTeamPlayers(ctx context.Context, cfg PlayerSyncConfig, reason string) (*PlayerSyncResult, error) {
	playerSyncRunMu.Lock()
	defer playerSyncRunMu.Unlock()

	if err := validatePlayerSyncConfig(cfg); err != nil {
		return nil, err
	}

	startedAt := time.Now().UTC()
	_ = savePlayerSyncState(cfg, "running", "", startedAt, nextPlayerSyncAt(cfg))

	result := &PlayerSyncResult{
		Provider:  cfg.Provider,
		Resource:  syncResourcePlayers,
		Reason:    reason,
		StartedAt: startedAt,
	}

	mappings, err := repositories.ListExternalTeamMappings(cfg.Provider)
	if err != nil {
		_ = savePlayerSyncState(cfg, "failed", err.Error(), startedAt, nextPlayerSyncAt(cfg))
		return nil, err
	}
	if len(mappings) == 0 {
		err := errors.New("no external team mappings configured")
		_ = savePlayerSyncState(cfg, "failed", err.Error(), startedAt, nextPlayerSyncAt(cfg))
		return nil, err
	}

	failed := 0
	var firstErr error
	for i, mapping := range mappings {
		teamResult, err := syncMappedTeamPlayers(ctx, cfg, mapping.TeamID, &mapping)
		if err != nil {
			failed++
			if firstErr == nil {
				firstErr = err
			}
			result.Skipped++
			continue
		}
		result.Teams++
		result.Total += teamResult.Total
		result.Created += teamResult.Created
		result.Updated += teamResult.Updated
		result.Deactivated += teamResult.Deactivated
		result.Skipped += teamResult.Skipped

		if i < len(mappings)-1 {
			select {
			case <-ctx.Done():
				err := ctx.Err()
				_ = savePlayerSyncState(cfg, "failed", err.Error(), startedAt, nextPlayerSyncAt(cfg))
				return nil, err
			case <-time.After(500 * time.Millisecond):
			}
		}
	}

	result.FinishedAt = time.Now().UTC()
	if failed > 0 && result.Teams == 0 {
		_ = savePlayerSyncState(cfg, "failed", firstErr.Error(), result.FinishedAt, nextPlayerSyncAt(cfg))
		return result, firstErr
	}
	if failed > 0 {
		_ = savePlayerSyncState(cfg, "partial", firstErr.Error(), result.FinishedAt, nextPlayerSyncAt(cfg))
		return result, nil
	}
	_ = savePlayerSyncState(cfg, "success", "", result.FinishedAt, nextPlayerSyncAt(cfg))
	return result, nil
}

func syncTeamPlayersUnlocked(ctx context.Context, cfg PlayerSyncConfig, teamID uint, reason string) (*PlayerSyncResult, error) {
	if err := validatePlayerSyncConfig(cfg); err != nil {
		return nil, err
	}

	startedAt := time.Now().UTC()
	_ = savePlayerSyncState(cfg, "running", "", startedAt, nextPlayerSyncAt(cfg))

	mapping, err := repositories.GetExternalTeamMapping(teamID, cfg.Provider)
	if err != nil {
		message := mappingErrorMessage(err, teamID, cfg.Provider)
		_ = savePlayerSyncState(cfg, "failed", message, startedAt, nextPlayerSyncAt(cfg))
		return nil, errors.New(message)
	}

	result, err := syncMappedTeamPlayers(ctx, cfg, teamID, mapping)
	if err != nil {
		_ = savePlayerSyncState(cfg, "failed", err.Error(), startedAt, nextPlayerSyncAt(cfg))
		return nil, err
	}
	result.Reason = reason
	result.StartedAt = startedAt
	result.FinishedAt = time.Now().UTC()
	_ = savePlayerSyncState(cfg, "success", "", result.FinishedAt, nextPlayerSyncAt(cfg))
	return result, nil
}

func syncMappedTeamPlayers(ctx context.Context, cfg PlayerSyncConfig, teamID uint, mapping *models.ExternalTeamMapping) (*PlayerSyncResult, error) {
	client := apifootball.NewClient(cfg.APIFootballBaseURL, cfg.APIFootballKey)
	data, err := client.TeamSquad(ctx, mapping.ExternalTeamID)
	if err != nil {
		return nil, err
	}

	result := &PlayerSyncResult{
		Provider: cfg.Provider,
		Resource: syncResourcePlayers,
		Teams:    1,
	}

	now := time.Now().UTC()
	activeIDs := []string{}
	for _, entry := range data.Response {
		if mapping.ExternalTeamName == "" && entry.Team.Name != "" {
			mapping.ExternalTeamName = entry.Team.Name
			_ = repositories.UpsertExternalTeamMapping(mapping)
		}

		for _, externalPlayer := range entry.Players {
			sourcePlayerID := strconv.FormatInt(externalPlayer.ID, 10)
			if externalPlayer.ID <= 0 || strings.TrimSpace(externalPlayer.Name) == "" {
				result.Skipped++
				continue
			}

			nameEn := strings.TrimSpace(externalPlayer.Name)
			position, positionLabel := normalizeAPIPlayerPosition(externalPlayer.Position)
			player := &models.Player{
				TeamID:         teamID,
				Name:           localizeAPIPlayerName(sourcePlayerID, nameEn),
				NameEn:         nameEn,
				ShirtNumber:    externalPlayer.Number,
				Position:       position,
				PositionLabel:  positionLabel,
				PhotoURL:       strings.TrimSpace(externalPlayer.Photo),
				Source:         cfg.Provider,
				SourcePlayerID: sourcePlayerID,
				ExternalTeamID: mapping.ExternalTeamID,
				IsActive:       true,
				LastSyncedAt:   &now,
			}
			action, err := repositories.UpsertPlayer(player)
			if err != nil {
				result.Skipped++
				continue
			}
			activeIDs = append(activeIDs, sourcePlayerID)
			result.Total++
			switch action {
			case "created":
				result.Created++
			case "updated":
				result.Updated++
			default:
				result.Skipped++
			}
		}
	}

	if len(activeIDs) > 0 {
		deactivated, err := repositories.MarkMissingPlayersInactive(teamID, cfg.Provider, activeIDs, now)
		if err != nil {
			return nil, err
		}
		result.Deactivated = deactivated
	}
	return result, nil
}

func GetTeamPlayerMapping(teamID uint, provider string) (*models.ExternalTeamMapping, error) {
	if strings.TrimSpace(provider) == "" {
		provider = GetPlayerSyncConfig().Provider
	}
	return repositories.GetExternalTeamMapping(teamID, provider)
}

func UpsertTeamPlayerMapping(teamID uint, input ExternalTeamMappingInput) (*models.ExternalTeamMapping, error) {
	provider := strings.TrimSpace(input.Provider)
	if provider == "" {
		provider = GetPlayerSyncConfig().Provider
	}
	if provider == "" {
		provider = syncProviderAPIFootball
	}

	mapping := &models.ExternalTeamMapping{
		TeamID:           teamID,
		Provider:         provider,
		ExternalTeamID:   strings.TrimSpace(input.ExternalTeamID),
		ExternalTeamName: strings.TrimSpace(input.ExternalTeamName),
	}
	if mapping.ExternalTeamID == "" {
		return nil, errors.New("external_team_id is required")
	}
	if err := repositories.UpsertExternalTeamMapping(mapping); err != nil {
		return nil, err
	}
	return repositories.GetExternalTeamMapping(teamID, provider)
}

func MarkInterruptedPlayerSync() error {
	return repositories.MarkInterruptedSyncState(
		GetPlayerSyncConfig().Provider,
		syncResourcePlayers,
		"previous sync was interrupted before completion",
		time.Now().UTC(),
	)
}

func validatePlayerSyncConfig(cfg PlayerSyncConfig) error {
	if !cfg.Enabled {
		return errors.New("player data sync is disabled")
	}
	if cfg.Provider != syncProviderAPIFootball {
		return fmt.Errorf("unsupported player data provider: %s", cfg.Provider)
	}
	if strings.TrimSpace(cfg.APIFootballKey) == "" {
		return errors.New("API_FOOTBALL_KEY is empty")
	}
	return nil
}

func normalizeAPIPlayerPosition(value string) (string, string) {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "GOALKEEPER":
		return "GK", "门将"
	case "DEFENDER":
		return "DF", "后卫"
	case "MIDFIELDER":
		return "MF", "中场"
	case "ATTACKER", "FORWARD":
		return "FW", "前锋"
	default:
		return strings.TrimSpace(value), ""
	}
}

func savePlayerSyncState(cfg PlayerSyncConfig, status, lastError string, lastSyncedAt time.Time, next *time.Time) error {
	state := &models.SyncState{
		Provider:     cfg.Provider,
		Resource:     syncResourcePlayers,
		Status:       status,
		LastSyncedAt: &lastSyncedAt,
		NextSyncAt:   next,
		LastError:    lastError,
	}
	return repositories.UpsertSyncState(state)
}

func nextPlayerSyncAt(cfg PlayerSyncConfig) *time.Time {
	next := time.Now().UTC().Add(cfg.Interval)
	return &next
}

func mappingErrorMessage(err error, teamID uint, provider string) string {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Sprintf("external team mapping not found for team %d and provider %s", teamID, provider)
	}
	return err.Error()
}
