package repositories

import (
	"errors"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"

	"gorm.io/gorm"
)

func GetExternalTeamMapping(teamID uint, provider string) (*models.ExternalTeamMapping, error) {
	var mapping models.ExternalTeamMapping
	err := database.DB.
		Preload("Team").
		Where("team_id = ? AND provider = ?", teamID, provider).
		First(&mapping).Error
	return &mapping, err
}

func ListExternalTeamMappings(provider string) ([]models.ExternalTeamMapping, error) {
	var mappings []models.ExternalTeamMapping
	q := database.DB.Preload("Team").Order("team_id ASC")
	if provider != "" {
		q = q.Where("provider = ?", provider)
	}
	err := q.Find(&mappings).Error
	return mappings, err
}

func UpsertExternalTeamMapping(mapping *models.ExternalTeamMapping) error {
	var existing models.ExternalTeamMapping
	err := database.DB.
		Where("team_id = ? AND provider = ?", mapping.TeamID, mapping.Provider).
		First(&existing).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return database.DB.Create(mapping).Error
	}

	existing.ExternalTeamID = mapping.ExternalTeamID
	existing.ExternalTeamName = mapping.ExternalTeamName
	return database.DB.Save(&existing).Error
}
