package handlers

import (
	"errors"
	"strconv"

	"worldcup-mate/internal/services"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListNotifications(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)
	notifications, total, err := services.GetNotifications(userID, page, pageSize)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Paginated(c, notifications, total, page, pageSize)
}

func CountUnreadNotifications(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	total, err := services.CountUnreadNotifications(userID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, gin.H{"count": total})
}

func MarkNotificationRead(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid notification id")
		return
	}
	userID := c.MustGet("user_id").(uint)
	if err := services.MarkNotificationRead(uint(id), userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, 404, "notification not found")
			return
		}
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, gin.H{"read": true})
}

func MarkAllNotificationsRead(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	if err := services.MarkAllNotificationsRead(userID); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, gin.H{"read_all": true})
}
