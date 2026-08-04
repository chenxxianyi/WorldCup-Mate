package repositories

import (
	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
)

func ListCompetitions() ([]models.Competition, error) {
	var competitions []models.Competition
	err := database.DB.Order("sort_order ASC, id ASC").Find(&competitions).Error
	return competitions, err
}

func GetCompetitionByID(id uint) (*models.Competition, error) {
	var competition models.Competition
	err := database.DB.First(&competition, id).Error
	return &competition, err
}

func GetCompetitionByCode(code string) (*models.Competition, error) {
	var competition models.Competition
	err := database.DB.Where("code = ?", code).First(&competition).Error
	return &competition, err
}

func CreateCompetition(competition *models.Competition) error {
	return database.DB.Create(competition).Error
}

func UpdateCompetition(competition *models.Competition) error {
	return database.DB.Save(competition).Error
}

func CountCompetitions() int64 {
	var count int64
	database.DB.Model(&models.Competition{}).Count(&count)
	return count
}
