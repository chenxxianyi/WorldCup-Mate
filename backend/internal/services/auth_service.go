package services

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	Avatar            string `json:"avatar"`
	Timezone          string `json:"timezone"`
	Language          string `json:"language"`
	NotificationEmail string `json:"notification_email" binding:"omitempty,email"`
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
	if notificationEmail := strings.TrimSpace(input.NotificationEmail); notificationEmail != "" {
		user.NotificationEmail = notificationEmail
	}
	if err := repositories.UpdateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

type ChangePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func ChangePassword(userID uint, input ChangePasswordInput) error {
	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return err
	}
	if !utils.CheckPassword(input.OldPassword, user.PasswordHash) {
		return errors.New("旧密码错误")
	}
	hash, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = hash
	return repositories.UpdateUser(user)
}

var allowedExts = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
var allowedAvatarMIMEs = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

func UploadAvatar(userID uint, file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExts[ext] {
		return "", errors.New("不支持的文件格式，仅支持 jpg/png/gif/webp")
	}
	if file.Size > 5*1024*1024 {
		return "", errors.New("文件大小不能超过 5MB")
	}

	dir := "uploads/avatars"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	filename := fmt.Sprintf("%d_%d%s", userID, time.Now().UnixMilli(), ext)
	dst := filepath.Join(dir, filename)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	header := make([]byte, 512)
	n, err := src.Read(header)
	if err != nil && err != io.EOF {
		return "", err
	}
	mimeType := http.DetectContentType(header[:n])
	if !allowedAvatarMIMEs[mimeType] {
		return "", errors.New("不支持的文件内容类型，仅支持 jpg/png/gif/webp")
	}
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()
	if _, err = io.Copy(out, src); err != nil {
		return "", err
	}

	avatarURL := "/uploads/avatars/" + filename

	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return "", err
	}
	user.Avatar = avatarURL
	if err := repositories.UpdateUser(user); err != nil {
		return "", err
	}
	return avatarURL, nil
}
