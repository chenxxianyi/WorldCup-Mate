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

func AdminGetTeamPlayerMapping(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid team id")
		return
	}

	mapping, err := services.GetTeamPlayerMapping(uint(teamID), c.Query("provider"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, 404, "external team mapping not found")
			return
		}
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, mapping)
}

func AdminUpsertTeamPlayerMapping(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid team id")
		return
	}

	var input services.ExternalTeamMappingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	mapping, err := services.UpsertTeamPlayerMapping(uint(teamID), input)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, mapping)
}

func AdminSyncTeamPlayers(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid team id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	result, err := services.SyncTeamPlayersWithDefault(ctx, uint(teamID), "manual")
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, result)
}

func AdminSyncAllPlayers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Minute)
	defer cancel()

	result, err := services.SyncAllMappedTeamPlayersWithDefault(ctx, "manual")
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, result)
}
