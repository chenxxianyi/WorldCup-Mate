package handlers

import (
	"strconv"

	"worldcup-mate/internal/services"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
)

func ListGroups(c *gin.Context) {
	groups, err := services.GetAllStandings()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, groups)
}

func GetGroupDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid group id")
		return
	}
	standings, err := services.GetStandingsByGroupID(uint(id))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, standings)
}

func GetGroupStandings(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid group id")
		return
	}
	standings, err := services.GetStandingsByGroupID(uint(id))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, standings)
}
