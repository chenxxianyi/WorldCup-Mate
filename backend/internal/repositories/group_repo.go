package repositories

import (
	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
)

func ListGroups() ([]models.Group, error) {
	var groups []models.Group
	err := database.DB.Find(&groups).Error
	return groups, err
}

func GetGroupByID(id uint) (*models.Group, error) {
	var group models.Group
	err := database.DB.First(&group, id).Error
	return &group, err
}

func GetGroupByName(name string) (*models.Group, error) {
	var group models.Group
	err := database.DB.Where("name = ?", name).First(&group).Error
	return &group, err
}

func CreateGroup(group *models.Group) error {
	return database.DB.Create(group).Error
}

func UpdateGroup(group *models.Group) error {
	return database.DB.Save(group).Error
}

func CountGroups() int64 {
	var count int64
	database.DB.Model(&models.Group{}).Count(&count)
	return count
}
