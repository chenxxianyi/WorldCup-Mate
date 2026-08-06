package handlers

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"worldcup-mate/internal/repositories"
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

// HomeAggregate returns the homepage aggregation: upcoming matches,
// live matches, competitions, and the latest sync state (DATA-10).
func HomeAggregate(c *gin.Context) {
	now := time.Now().UTC()

	upcoming, _ := services.GetUpcomingMatches(false)
	live, _ := services.GetLiveMatches(false)
	competitions, _ := repositories.ListActiveCompetitions()

	// Latest sync states for credibility display (DATA-11).
	syncStates, _ := repositories.GetSyncStates()

	// Upcoming matches limited to 10 for the homepage.
	if len(upcoming) > 10 {
		upcoming = upcoming[:10]
	}

	utils.Success(c, gin.H{
		"upcoming_matches": upcoming,
		"live_matches":     live,
		"competitions":     competitions,
		"sync_states":      syncStates,
		"synced_at":        now.Format(time.RFC3339),
	})
}

func GetTodayMatches(c *gin.Context) {
	matches, err := services.GetTodayMatches(isWorldCupQuery(c), c.Query("tz"))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, matches)
}

func GetTomorrowMatches(c *gin.Context) {
	matches, err := services.GetTomorrowMatches(isWorldCupQuery(c), c.Query("tz"))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, matches)
}

func GetUpcomingMatches(c *gin.Context) {
	matches, err := services.GetUpcomingMatches(isWorldCupQuery(c))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, matches)
}

func GetLiveMatches(c *gin.Context) {
	matches, err := services.GetLiveMatches(isWorldCupQuery(c))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, matches)
}

func GetRecommendedMatches(c *gin.Context) {
	matches, err := services.GetRecommendedMatches(isWorldCupQuery(c))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, matches)
}

// isWorldCupQuery reports whether the client requested World Cup-only
// data (worldCup=true, used by the default World Cup view).
func isWorldCupQuery(c *gin.Context) bool {
	return c.Query("worldCup") == "true"
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
	// League sync: ?code=PL&season=2025 triggers a single league; otherwise
	// the legacy World Cup sync path runs unchanged.
	if code := strings.ToUpper(c.Query("code")); code != "" && code != "WC" {
		season, _ := strconv.Atoi(c.Query("season"))
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
		defer cancel()
		result, err := services.SyncLeague(ctx, code, season, "manual")
		if err != nil {
			// REL-06: another worker (manual or scheduled) holds the lock.
			if errors.Is(err, services.ErrSyncAlreadyRunning) {
				utils.Error(c, 409, err.Error())
				return
			}
			utils.Error(c, 500, err.Error())
			return
		}
		utils.Success(c, result)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	result, err := services.SyncMatchesWithDefault(ctx, "manual")
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, result)
}
