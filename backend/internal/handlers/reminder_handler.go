package handlers

import (
	"strconv"

	"worldcup-mate/internal/services"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
)

func CreateReminder(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var input services.CreateReminderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	reminder, err := services.CreateReminder(userID, input)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, reminder)
}

func CreateReminderBatch(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var input services.CreateReminderBatchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	reminders, err := services.CreateReminderBatch(userID, input)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, reminders)
}

func ListReminders(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	reminders, err := services.GetReminders(userID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, reminders)
}

func UpdateReminder(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid reminder id")
		return
	}
	var input struct {
		RemindBeforeMinutes int    `json:"remindBeforeMinutes"`
		Channel             string `json:"channel"`
	}
	c.ShouldBindJSON(&input)
	reminder, err := services.UpdateReminder(uint(id), userID, input.RemindBeforeMinutes, input.Channel)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, reminder)
}

func DeleteReminder(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid reminder id")
		return
	}
	if err := services.DeleteReminder(uint(id), userID); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, gin.H{"deleted": true})
}
