package handlers

import (
	"context"
	"strconv"
	"time"

	"worldcup-mate/internal/services"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
)

func ListMatches(c *gin.Context) {
	var q services.MatchQuery
	c.ShouldBindQuery(&q)
	q.Page = utils.GetPage(c)
	q.PageSize = utils.GetPageSize(c)

	matches, total, err := services.ListMatches(q)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Paginated(c, matches, total, q.Page, q.PageSize)
}

func GetTodayMatches(c *gin.Context) {
	matches, err := services.GetTodayMatches()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, matches)
}

func GetTomorrowMatches(c *gin.Context) {
	matches, err := services.GetTomorrowMatches()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, matches)
}

func GetUpcomingMatches(c *gin.Context) {
	matches, err := services.GetUpcomingMatches()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, matches)
}

func GetLiveMatches(c *gin.Context) {
	matches, err := services.GetLiveMatches()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, matches)
}

func GetRecommendedMatches(c *gin.Context) {
	matches, err := services.GetRecommendedMatches()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, matches)
}

func GetMatchDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid match id")
		return
	}
	match, err := services.GetMatchByID(uint(id))
	if err != nil {
		utils.Error(c, 404, "match not found")
		return
	}
	utils.Success(c, match)
}

func GetTimeline(c *gin.Context) {
	var q services.TimelineQuery
	c.ShouldBindQuery(&q)
	matches, err := services.GetTimeline(q)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, matches)
}

func GetMatchesByTeam(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("teamId"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid team id")
		return
	}
	matches, err := services.GetMatchesByTeamID(uint(teamID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, matches)
}

func GetMatchesByGroup(c *gin.Context) {
	groupID, err := strconv.ParseUint(c.Param("groupId"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid group id")
		return
	}
	matches, err := services.GetMatchesByGroupID(uint(groupID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, matches)
}

func GetMatchesByStage(c *gin.Context) {
	stage := c.Param("stage")
	matches, err := services.GetMatchesByStage(stage)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, matches)
}

func GetTournamentProgress(c *gin.Context) {
	progress, err := services.GetTournamentProgress()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, progress)
}

func GetSyncStatus(c *gin.Context) {
	states, err := services.GetSyncStates()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, states)
}

func AdminSyncMatches(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	result, err := services.SyncMatchesWithDefault(ctx, "manual")
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, result)
}
