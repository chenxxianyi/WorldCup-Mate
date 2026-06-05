package handlers

import (
	"context"
	"errors"
	"strconv"
	"time"

	"worldcup-mate/internal/services"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminSyncMatchLineups(c *gin.Context) {
	matchID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid match id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	result, err := services.SyncMatchLineups(ctx, uint(matchID), "manual")
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, result)
}

func AdminSyncLiveWindowLineups(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Minute)
	defer cancel()

	result, err := services.SyncLiveWindowLineups(ctx, "manual")
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, result)
}

func AdminGetMatchExternalMapping(c *gin.Context) {
	matchID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid match id")
		return
	}

	mapping, err := services.GetMatchExternalMapping(uint(matchID), c.Query("provider"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, 404, "external match mapping not found")
			return
		}
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, mapping)
}

func AdminUpsertMatchExternalMapping(c *gin.Context) {
	matchID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid match id")
		return
	}

	var input services.ExternalMatchMappingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	mapping, err := services.UpsertMatchExternalMapping(uint(matchID), input)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, mapping)
}
