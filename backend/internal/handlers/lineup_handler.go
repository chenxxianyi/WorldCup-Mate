package handlers

import (
	"errors"
	"strconv"

	"worldcup-mate/internal/services"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMatchLineups(c *gin.Context) {
	matchID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid match id")
		return
	}

	lineups, err := services.GetMatchLineups(uint(matchID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, 404, "match not found")
			return
		}
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, lineups)
}
