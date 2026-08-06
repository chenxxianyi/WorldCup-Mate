package repositories

import (
	"time"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
)

func GetSyncStates() ([]models.SyncState, error) {
	var states []models.SyncState
	err := database.DB.Order("provider ASC, resource ASC").Find(&states).Error
	return states, err
}

func GetSyncState(provider, resource string) (*models.SyncState, error) {
	var state models.SyncState
	err := database.DB.Where("provider = ? AND resource = ?", provider, resource).First(&state).Error
	return &state, err
}

func UpsertSyncState(state *models.SyncState) error {
	var existing models.SyncState
	err := database.DB.Where("provider = ? AND resource = ?", state.Provider, state.Resource).First(&existing).Error
	if err != nil {
		return database.DB.Create(state).Error
	}
	existing.Status = state.Status
	existing.LastSyncedAt = state.LastSyncedAt
	existing.NextSyncAt = state.NextSyncAt
	existing.LastError = state.LastError
	return database.DB.Save(&existing).Error
}

func MarkInterruptedSyncState(provider, resource, message string, now time.Time) error {
	return database.DB.Model(&models.SyncState{}).
		Where("provider = ? AND resource = ? AND status = ?", provider, resource, "running").
		Updates(map[string]interface{}{
			"status":         "failed",
			"last_error":     message,
			"last_synced_at": now,
			"next_sync_at":   nil,
		}).Error
}

// ListSyncHistory returns the last N sync states ordered by id descending.
// It is used by ADM-13 / sync history operations.
func ListSyncHistory(limit int) ([]models.SyncState, error) {
	if limit <= 0 {
		limit = 50
	}
	var states []models.SyncState
	err := database.DB.Order("id DESC").Limit(limit).Find(&states).Error
	return states, err
}
