package handlers

import (
	"strconv"

	"worldcup-mate/internal/services"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
)

func ListTeams(c *gin.Context) {
	var q services.TeamQuery
	c.ShouldBindQuery(&q)
	q.Page = utils.GetPage(c)
	q.PageSize = utils.GetPageSize(c)

	teams, total, err := services.ListTeams(q)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Paginated(c, teams, total, q.Page, q.PageSize)
}

func GetTeamDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid team id")
		return
	}
	team, err := services.GetTeamByID(uint(id))
	if err != nil {
		utils.Error(c, 404, "team not found")
		return
	}
	utils.Success(c, team)
}

func GetTeamMatches(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid team id")
		return
	}
	matches, err := services.GetTeamMatches(uint(teamID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, matches)
}
