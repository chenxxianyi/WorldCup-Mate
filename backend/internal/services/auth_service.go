package services

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif" // register gif decoder (first frame)
	"image/jpeg"
	_ "image/png" // register png decoder
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/image/webp"

	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"
	"worldcup-mate/internal/utils"

	"gorm.io/gorm"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileInput struct {
	Timezone          string `json:"timezone"`
	Language          string `json:"language"`
	NotificationEmail string `json:"notification_email" binding:"omitempty,email"`
	// Avatar is intentionally NOT accepted here (SEC-06): it can only be
	// changed through the avatar upload endpoint.
}

func Register(input RegisterInput) (*models.User, error) {
	_, err := repositories.GetUserByEmail(input.Email)
	if err == nil {
		return nil, errors.New("email already registered")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err := utils.ValidatePassword(input.Password); err != nil {
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

	// ADM-06: disabled accounts cannot sign in.
	if user.Status == "disabled" {
		return "", nil, errors.New("account disabled")
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
	// SEC-06: avatar is managed exclusively by UploadAvatar (random filename,
	// re-encode, old-file cleanup); this DTO must not bypass that flow.
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
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

func ChangePassword(userID uint, input ChangePasswordInput) error {
	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return err
	}
	if !utils.CheckPassword(input.OldPassword, user.PasswordHash) {
		return errors.New("旧密码错误")
	}
	if err := utils.ValidatePassword(input.NewPassword); err != nil {
		return err
	}
	hash, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = hash
	// SEC-04: mark the change so pre-change access tokens die immediately.
	now := time.Now()
	user.PasswordChangedAt = &now
	return repositories.UpdateUser(user)
}

var allowedExts = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}

// uploadRoot is set from config at startup (SEC-06E); defaults to "uploads".
var uploadRoot = "uploads"

// SetUploadRoot configures the upload directory (called from main with
// cfg.UploadDir). SEC-06E: upload paths are configuration-driven.
func SetUploadRoot(dir string) {
	if strings.TrimSpace(dir) != "" {
		uploadRoot = dir
	}
}

func UploadAvatar(userID uint, file *multipart.FileHeader) (string, error) {
	if file.Size > 5*1024*1024 {
		return "", errors.New("文件大小不能超过 5MB")
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// SEC-06C: real pixel decode + JPEG re-encode. Fake images (polyglot
	// files with a spoofed extension/MIME) fail at decode time, and all
	// EXIF/metadata is stripped by the re-encode.
	clean, err := reencodeAvatar(src)
	if err != nil {
		return "", err
	}

	dir := filepath.Join(uploadRoot, "avatars")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	// SEC-06B: cryptographically random filename — no user ID, no timestamp
	// that could correlate accounts.
	filename := randomHex(16) + ".jpg"
	dst := filepath.Join(dir, filename)
	if err := os.WriteFile(dst, clean, 0644); err != nil {
		return "", err
	}

	avatarURL := "/uploads/avatars/" + filename

	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return "", err
	}
	// SEC-06D: replace the previous avatar file so old files do not
	// accumulate on every change.
	oldAvatar := user.Avatar
	user.Avatar = avatarURL
	if err := repositories.UpdateUser(user); err != nil {
		_ = os.Remove(dst) // roll back the new file if the DB update fails
		return "", err
	}
	if oldAvatar != "" && strings.HasPrefix(oldAvatar, "/uploads/avatars/") {
		_ = os.Remove(filepath.Join(dir, filepath.Base(oldAvatar)))
	}
	return avatarURL, nil
}

// randomHex returns n cryptographically random bytes as hex.
func randomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d_%d", time.Now().UnixNano(), n)
	}
	return hex.EncodeToString(b)
}

// maxAvatarPixels caps decoded dimensions (decompression-bomb guard).
// Checked via DecodeConfig BEFORE full decode allocates memory. 16.7M
// (4096×4096) is far beyond any avatar while bounding worst-case memory.
const maxAvatarPixels = 4096 * 4096

func validateAvatarDimensions(dx, dy int) bool {
	return dx > 0 && dy > 0 && int64(dx)*int64(dy) <= maxAvatarPixels
}

// webpDimensions extracts the canvas size from a WebP VP8X header.
// Lossy VP8/VP8L frames carry no canvas header, so those report ok=false
// and fall through to the post-decode check (still fail-closed).
func webpDimensions(data []byte) (w, h int, ok bool) {
	if len(data) < 30 || string(data[0:4]) != "RIFF" || string(data[8:12]) != "WEBP" {
		return 0, 0, false
	}
	if string(data[12:16]) != "VP8X" {
		return 0, 0, false
	}
	w = int(data[24]) | int(data[25])<<8 | int(data[26])<<16
	h = int(data[27]) | int(data[28])<<8 | int(data[29])<<16
	return w, h, true
}

// reencodeAvatar decodes the image (jpg/png/gif/webp) and re-encodes it as
// JPEG, dropping all metadata. Fake images fail at decode time. Transparent
// images are composited onto white (JPEG has no alpha). Dimensions are
// pre-checked via DecodeConfig so oversized images never allocate pixels.
func reencodeAvatar(src io.Reader) ([]byte, error) {
	data, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("空文件")
	}
	// SEC-06C: dimension pre-check via headers only — a 5MB highly-compressed
	// image can carry tens of millions of pixels which would otherwise
	// allocate hundreds of MB before the decode-time check could reject it.
	// stdlib formats go through DecodeConfig; webp VP8X is parsed manually.
	if cfg, _, err := image.DecodeConfig(bytes.NewReader(data)); err == nil {
		if !validateAvatarDimensions(cfg.Width, cfg.Height) {
			return nil, errors.New("图片尺寸过大")
		}
	} else if w, h, ok := webpDimensions(data); ok && !validateAvatarDimensions(w, h) {
		return nil, errors.New("图片尺寸过大")
	}
	img, err := decodeImage(data)
	if err != nil {
		return nil, errors.New("无效的图片文件，仅支持 jpg/png/gif/webp")
	}
	b := img.Bounds()
	if !validateAvatarDimensions(b.Dx(), b.Dy()) {
		return nil, errors.New("图片尺寸过大")
	}
	// Composite onto white so transparent PNG/GIF do not come out with a
	// black background after the JPEG re-encode. Build the canvas at (0,0)
	// and align the source at b.Min — do not assume decoders return Min==0.
	white := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(white, white.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	draw.Draw(white, white.Bounds(), img, b.Min, draw.Over)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, white, &jpeg.Options{Quality: 88}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// decodeImage tries the stdlib registry (jpg/png/gif) first, then the
// explicit webp decoder. Note: golang.org/x/image/webp supports lossy VP8
// only — lossless/animated WebP are rejected (fail-closed, acceptable for
// avatars); the frontend accept list stays aligned with this.
func decodeImage(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err == nil {
		return img, nil
	}
	return webp.Decode(bytes.NewReader(data))
}
