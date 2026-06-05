package repositories

import (
	"errors"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"

	"gorm.io/gorm"
)

func GetExternalMatchMapping(matchID uint, provider string) (*models.ExternalMatchMapping, error) {
	var mapping models.ExternalMatchMapping
	err := database.DB.
		Preload("Match").
		Where("match_id = ? AND provider = ?", matchID, provider).
		First(&mapping).Error
	return &mapping, err
}

func UpsertExternalMatchMapping(mapping *models.ExternalMatchMapping) error {
	var existing models.ExternalMatchMapping
	err := database.DB.
		Where("match_id = ? AND provider = ?", mapping.MatchID, mapping.Provider).
		First(&existing).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return database.DB.Create(mapping).Error
	}

	existing.ExternalMatchID = mapping.ExternalMatchID
	existing.ExternalHomeID = mapping.ExternalHomeID
	existing.ExternalAwayID = mapping.ExternalAwayID
	existing.ExternalHomeName = mapping.ExternalHomeName
	existing.ExternalAwayName = mapping.ExternalAwayName
	existing.MatchedBy = mapping.MatchedBy
	return database.DB.Save(&existing).Error
}
