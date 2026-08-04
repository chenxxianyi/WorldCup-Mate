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
	competitions, err := repositories.ListCompetitions()
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
	if err := repositories.CreateCompetition(&competition); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, competition)
}

func AdminUpdateCompetition(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	competition, err := repositories.GetCompetitionByID(uint(id))
	if err != nil {
		utils.Error(c, 404, "competition not found")
		return
	}
	if v, ok := input["name"].(string); ok {
		competition.Name = v
	}
	if v, ok := input["name_en"].(string); ok {
		competition.NameEn = v
	}
	if v, ok := input["country"].(string); ok {
		competition.Country = v
	}
	if v, ok := input["logo_url"].(string); ok {
		competition.LogoURL = v
	}
	if v, ok := input["format"].(string); ok {
		competition.Format = v
	}
	if v, ok := input["status"].(string); ok {
		competition.Status = v
	}
	if v, ok := input["season"].(float64); ok {
		competition.Season = int(v)
	}
	if v, ok := input["sort_order"].(float64); ok {
		competition.SortOrder = int(v)
	}
	if err := repositories.UpdateCompetition(competition); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, competition)
}
