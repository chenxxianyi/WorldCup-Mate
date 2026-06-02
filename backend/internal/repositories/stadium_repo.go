package repositories

import (
	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
)

func ListStadiums() ([]models.Stadium, error) {
	var stadiums []models.Stadium
	err := database.DB.Preload("City").Find(&stadiums).Error
	return stadiums, err
}

func GetStadiumByID(id uint) (*models.Stadium, error) {
	var stadium models.Stadium
	err := database.DB.Preload("City").First(&stadium, id).Error
	return &stadium, err
}

func GetStadiumByName(name string) (*models.Stadium, error) {
	var stadium models.Stadium
	err := database.DB.Where("name = ?", name).First(&stadium).Error
	return &stadium, err
}

func CreateStadium(stadium *models.Stadium) error {
	return database.DB.Create(stadium).Error
}

func UpdateStadium(stadium *models.Stadium) error {
	return database.DB.Save(stadium).Error
}

func DeleteStadium(id uint) error {
	return database.DB.Delete(&models.Stadium{}, id).Error
}
