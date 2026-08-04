package handlers

import (
	"errors"

	"worldcup-mate/internal/ratelimit"
	"worldcup-mate/internal/services"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var input services.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	// Note: registration intentionally does NOT consult the login-failure
	// freeze, otherwise attackers could block sign-ups for a target email
	// by failing to log in repeatedly (pre-registration DoS).
	user, err := services.Register(input)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		utils.Error(c, 500, "failed to generate token")
		return
	}
	refreshToken, err := services.IssueRefreshToken(user.ID)
	if err != nil {
		utils.Error(c, 500, "failed to issue refresh token")
		return
	}
	utils.Success(c, gin.H{"token": token, "refresh_token": refreshToken, "user": user})
}

func Login(c *gin.Context) {
	var input services.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	// SEC-02: account-level freeze after consecutive failures.
	if locked, err := ratelimit.LoginLocked(input.Email); err == nil && locked {
		utils.Error(c, 429, "too many failed attempts, try again later")
		return
	}
	token, user, err := services.Login(input)
	if err != nil {
		ratelimit.RecordLoginFailure(input.Email)
		utils.Error(c, 401, err.Error())
		return
	}
	ratelimit.ClearLoginFailures(input.Email)
	refreshToken, err := services.IssueRefreshToken(user.ID)
	if err != nil {
		utils.Error(c, 500, "failed to issue refresh token")
		return
	}
	utils.Success(c, gin.H{"token": token, "refresh_token": refreshToken, "user": user})
}

// Refresh exchanges a valid refresh token for a fresh access+refresh pair
// (rotation). Replayed tokens revoke the whole account (SEC-04C).
func Refresh(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	_, accessToken, newRefresh, err := services.ValidateAndRotate(input.RefreshToken)
	if err != nil {
		switch err {
		case services.ErrRefreshTokenReplayed:
			utils.Error(c, 401, "会话已失效，请重新登录")
		case services.ErrRefreshTokenExpired:
			utils.Error(c, 401, "登录已过期，请重新登录")
		default:
			utils.Error(c, 401, "登录已失效，请重新登录")
		}
		return
	}
	utils.Success(c, gin.H{"token": accessToken, "refresh_token": newRefresh})
}

// Logout revokes every refresh-token session of the caller. Idempotent:
// without a valid access token it still returns success (SEC-04C).
func Logout(c *gin.Context) {
	userID, err := userIDFromBearer(c)
	if err == nil && userID > 0 {
		_ = services.RevokeAllRefreshTokens(userID)
	}
	utils.Success(c, nil)
}

// userIDFromBearer extracts the user id from the Authorization header
// without requiring the JWT middleware (logout stays unauthenticated).
func userIDFromBearer(c *gin.Context) (uint, error) {
	header := c.GetHeader("Authorization")
	if len(header) < 8 || header[:7] != "Bearer " {
		return 0, errors.New("missing bearer token")
	}
	claims, err := utils.ParseToken(header[7:])
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

func GetProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	user, err := services.GetProfile(userID)
	if err != nil {
		utils.Error(c, 404, "user not found")
		return
	}
	utils.Success(c, user)
}

func UpdateProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var input services.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	user, err := services.UpdateProfile(userID, input)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, user)
}

func ChangePassword(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var input services.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	if err := services.ChangePassword(userID, input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	// SEC-04D: password change revokes every existing session.
	_ = services.RevokeAllRefreshTokens(userID)
	utils.Success(c, nil)
}

func UploadAvatar(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	file, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, 400, "请选择要上传的文件")
		return
	}
	avatarURL, err := services.UploadAvatar(userID, file)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, gin.H{"avatar": avatarURL})
}
