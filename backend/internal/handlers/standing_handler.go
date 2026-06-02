package handlers

import (
	"worldcup-mate/internal/services"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
)

func ListStandings(c *gin.Context) {
	standings, err := services.GetAllStandings()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, standings)
}

func GetBestThird(c *gin.Context) {
	standings, err := services.GetBestThird()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, standings)
}
