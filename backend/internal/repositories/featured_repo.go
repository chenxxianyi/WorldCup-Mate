package repositories

import (
	"errors"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"

	"gorm.io/gorm"
)

// GetFeaturedConfigs returns all per-competition hero configurations,
// preloaded with competition codes.
func GetFeaturedConfigs() ([]models.FeaturedConfig, error) {
	var configs []models.FeaturedConfig
	err := database.DB.Preload("Competition").Find(&configs).Error
	return configs, err
}

// GetFeaturedConfigByCompetition returns one competition's hero config.
func GetFeaturedConfigByCompetition(competitionID uint) (*models.FeaturedConfig, error) {
	var cfg models.FeaturedConfig
	err := database.DB.Where("competition_id = ?", competitionID).First(&cfg).Error
	return &cfg, err
}

// UpsertFeaturedConfig creates or updates the hero config of a competition.
func UpsertFeaturedConfig(competitionID uint, cfg *models.FeaturedConfig) error {
	cfg.CompetitionID = competitionID
	var existing models.FeaturedConfig
	err := database.DB.Where("competition_id = ?", competitionID).First(&existing).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err == nil {
		cfg.ID = existing.ID
		return database.DB.Save(cfg).Error
	}
	return database.DB.Create(cfg).Error
}
