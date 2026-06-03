package handlers

import (
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
	utils.Success(c, gin.H{"token": token, "user": user})
}

func Login(c *gin.Context) {
	var input services.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	token, user, err := services.Login(input)
	if err != nil {
		utils.Error(c, 401, err.Error())
		return
	}
	utils.Success(c, gin.H{"token": token, "user": user})
}

func Logout(c *gin.Context) {
	utils.Success(c, nil)
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
