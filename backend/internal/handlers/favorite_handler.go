package handlers

import (
	"strconv"

	"worldcup-mate/internal/services"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
)

func AddFavoriteTeam(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	teamID, err := strconv.ParseUint(c.Param("teamId"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid team id")
		return
	}
	err = services.AddFavoriteTeam(userID, uint(teamID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, gin.H{"added": true})
}

func RemoveFavoriteTeam(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	teamID, err := strconv.ParseUint(c.Param("teamId"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid team id")
		return
	}
	err = services.RemoveFavoriteTeam(userID, uint(teamID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, gin.H{"removed": true})
}

func ListFavoriteTeams(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	favs, err := services.GetFavoriteTeams(userID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, favs)
}

func AddFavoriteMatch(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	matchID, err := strconv.ParseUint(c.Param("matchId"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid match id")
		return
	}
	err = services.AddFavoriteMatch(userID, uint(matchID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, gin.H{"added": true})
}

func RemoveFavoriteMatch(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	matchID, err := strconv.ParseUint(c.Param("matchId"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid match id")
		return
	}
	err = services.RemoveFavoriteMatch(userID, uint(matchID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, gin.H{"removed": true})
}

func ListFavoriteMatches(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	favs, err := services.GetFavoriteMatches(userID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, favs)
}
