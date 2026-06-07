package repositories

import (
	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
)

func CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func ListUsersByUsername(username string) ([]models.User, error) {
	var users []models.User
	err := database.DB.Where("username = ?", username).Order("id ASC").Find(&users).Error
	return users, err
}

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := database.DB.First(&user, id).Error
	return &user, err
}

func UpdateUser(user *models.User) error {
	return database.DB.Save(user).Error
}

func ListUsers(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	database.DB.Model(&models.User{}).Count(&total)
	err := database.DB.Offset((page-1)*pageSize).Limit(pageSize).Find(&users).Error
	return users, total, err
}
