package handlers

import (
	"strconv"

	"worldcup-mate/internal/services"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
)

func GetTeamPlayers(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid team id")
		return
	}

	players, err := services.ListPlayersByTeam(uint(teamID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, players)
}
