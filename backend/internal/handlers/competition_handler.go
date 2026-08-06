package handlers

import (
	"strconv"
	"strings"

	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
)

// ListCompetitions returns all competitions ordered by sort_order.
// Public endpoint used by the frontend competition switcher.
func ListCompetitions(c *gin.Context) {
	// Public endpoint: only enabled competitions (the admin endpoint lists
	// all of them).
	competitions, err := repositories.ListActiveCompetitions()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, competitions)
}

// GetCompetitionStandings returns the league standings of a competition.
// Optional query params: season (defaults to the competition's current season), type (total|home|away, default total).
func GetCompetitionStandings(c *gin.Context) {
	code := strings.ToUpper(c.Param("code"))
	competition, err := repositories.GetCompetitionByCode(code)
	if err != nil {
		utils.Error(c, 404, "competition not found")
		return
	}
	// Disabled competitions are hidden from the public API as well (the
	// admin/sync paths share GetCompetitionByCode and are unaffected).
	if competition.Status != "active" {
		utils.Error(c, 404, "competition not found")
		return
	}
	season := competition.Season
	if v, err := strconv.Atoi(c.Query("season")); err == nil && v > 0 {
		season = v
	}
	standingType := c.DefaultQuery("type", "total")
	standings, err := repositories.GetLeagueStandings(competition.ID, season, standingType)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, standings)
}

// CompetitionOverview returns the public overview for a league competition:
// its current season, available seasons, latest matchday, and total match count.
// This is the DATA-09 endpoint used by the league switcher and season selector.
func CompetitionOverview(c *gin.Context) {
	code := strings.ToUpper(c.Param("code"))
	competition, err := repositories.GetCompetitionByCode(code)
	if err != nil {
		utils.Error(c, 404, "competition not found")
		return
	}
	if competition.Status != "active" {
		utils.Error(c, 404, "competition not found")
		return
	}

	seasons, err := repositories.ListCompetitionSeasons(competition.ID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	// Default season is the competition's current season.
	season := competition.Season
	if v, err := strconv.Atoi(c.Query("season")); err == nil && v > 0 {
		season = v
	}

	// Latest matchday for the selected season (league only).
	matchday, err := repositories.GetLatestMatchday(competition.ID, season)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	// Match count for the selected season.
	matchCount, err := repositories.CountMatchesByCompetitionAndSeason(competition.ID, season)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"competition": competition,
		"seasons":     seasons,
		"season":      season,
		"matchday":    matchday,
		"match_count": matchCount,
	})
}

func AdminListCompetitions(c *gin.Context) {
	competitions, err := repositories.ListCompetitions()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, competitions)
}

func AdminCreateCompetition(c *gin.Context) {
	var input struct {
		Code      string `json:"code" binding:"required"`
		Name      string `json:"name" binding:"required"`
		NameEn    string `json:"name_en"`
		Country   string `json:"country"`
		LogoURL   string `json:"logo_url"`
		Format    string `json:"format"`
		Season    int    `json:"season"`
		Status    string `json:"status"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	competition := models.Competition{
		Code:      strings.ToUpper(input.Code),
		Name:      input.Name,
		NameEn:    input.NameEn,
		Country:   input.Country,
		LogoURL:   input.LogoURL,
		Format:    input.Format,
		Season:    input.Season,
		Status:    input.Status,
		SortOrder: input.SortOrder,
	}
	if competition.Format == "" {
		competition.Format = "league"
	}
	if competition.Status == "" {
		competition.Status = "active"
	}
	// Normalize/validate status exactly like the update path.
	competition.Status = strings.ToLower(competition.Status)
	if competition.Status != "active" && competition.Status != "inactive" {
		utils.Error(c, 400, "invalid status, must be active or inactive")
		return
	}
	// Pre-check the code like the team endpoints do; the unique index
	// remains the hard guarantee against races.
	if existing, err := repositories.GetCompetitionByCode(competition.Code); err == nil && existing != nil {
		utils.Error(c, 409, "competition code already exists")
		return
	}
	if err := repositories.CreateCompetition(&competition); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, competition)
}

func AdminUpdateCompetition(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid competition id")
		return
	}
	var input CompetitionUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	competition, err := repositories.GetCompetitionByID(uint(id))
	if err != nil {
		utils.Error(c, 404, "competition not found")
		return
	}
	if input.Name != nil {
		competition.Name = *input.Name
	}
	if input.NameEn != nil {
		competition.NameEn = *input.NameEn
	}
	if input.Country != nil {
		competition.Country = *input.Country
	}
	if input.LogoURL != nil {
		competition.LogoURL = *input.LogoURL
	}
	if input.Format != nil {
		competition.Format = *input.Format
	}
	if input.Status != nil {
		status := strings.ToLower(*input.Status)
		if status != "active" && status != "inactive" {
			utils.Error(c, 400, "invalid status, must be active or inactive")
			return
		}
		competition.Status = status
	}
	if input.Season != nil {
		competition.Season = *input.Season
	}
	if input.SortOrder != nil {
		competition.SortOrder = *input.SortOrder
	}
	if err := repositories.UpdateCompetition(competition); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, competition)
}
