package services

import (
	"errors"

	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"
	"worldcup-mate/internal/utils"

	"gorm.io/gorm"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileInput struct {
	Avatar   string `json:"avatar"`
	Timezone string `json:"timezone"`
	Language string `json:"language"`
}

func Register(input RegisterInput) (*models.User, error) {
	_, err := repositories.GetUserByEmail(input.Email)
	if err == nil {
		return nil, errors.New("email already registered")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hash,
		Role:         "user",
		Timezone:     "Asia/Shanghai",
		Language:     "zh-CN",
	}
	if err := repositories.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func Login(input LoginInput) (string, *models.User, error) {
	user, err := repositories.GetUserByEmail(input.Email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	if !utils.CheckPassword(input.Password, user.PasswordHash) {
		return "", nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func GetProfile(userID uint) (*models.User, error) {
	return repositories.GetUserByID(userID)
}

func UpdateProfile(userID uint, input UpdateProfileInput) (*models.User, error) {
	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if input.Avatar != "" {
		user.Avatar = input.Avatar
	}
	if input.Timezone != "" {
		user.Timezone = input.Timezone
	}
	if input.Language != "" {
		user.Language = input.Language
	}
	if err := repositories.UpdateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}
