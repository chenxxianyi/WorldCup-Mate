package repositories

import (
	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
)

func ListCities() ([]models.City, error) {
	var cities []models.City
	err := database.DB.Find(&cities).Error
	return cities, err
}

func GetCityByID(id uint) (*models.City, error) {
	var city models.City
	err := database.DB.First(&city, id).Error
	return &city, err
}

func GetCityByName(name string) (*models.City, error) {
	var city models.City
	err := database.DB.Where("name = ?", name).First(&city).Error
	return &city, err
}

func CreateCity(city *models.City) error {
	return database.DB.Create(city).Error
}

func UpdateCity(city *models.City) error {
	return database.DB.Save(city).Error
}

func DeleteCity(id uint) error {
	return database.DB.Delete(&models.City{}, id).Error
}
