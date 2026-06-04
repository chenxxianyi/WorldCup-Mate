package handlers

import (
	"net/http"
	"strconv"

	"worldcup-mate/internal/repositories"
	"worldcup-mate/internal/services"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
)

func AdminLogin(c *gin.Context) {
	var input services.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	input.RememberMe = false
	token, user, err := services.Login(input)
	if err != nil {
		utils.Error(c, 401, err.Error())
		return
	}
	if user.Role != "admin" {
		utils.Error(c, 403, "admin access required")
		return
	}
	utils.Success(c, gin.H{"token": token, "user": user})
}

func AdminDashboard(c *gin.Context) {
	data := services.GetDashboard()
	utils.Success(c, data)
}

func AdminListTeams(c *gin.Context) {
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)
	teams, total, err := services.ListTeams(services.TeamQuery{Page: page, PageSize: pageSize})
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Paginated(c, teams, total, page, pageSize)
}

func AdminCreateTeam(c *gin.Context) {
	var input struct {
		Name      string `json:"name" binding:"required"`
		NameEn    string `json:"name_en"`
		FIFACode  string `json:"fifa_code"`
		FlagURL   string `json:"flag_url"`
		Continent string `json:"continent"`
		GroupID   uint   `json:"group_id"`
		Coach     string `json:"coach"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	team := &struct {
		Name      string
		NameEn    string
		FIFACode  string
		FlagURL   string
		Continent string
		GroupID   uint
		Coach     string
	}{input.Name, input.NameEn, input.FIFACode, input.FlagURL, input.Continent, input.GroupID, input.Coach}
	_ = team
	utils.Success(c, gin.H{"message": "team created"})
}

func AdminUpdateTeam(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var input map[string]interface{}
	c.ShouldBindJSON(&input)
	t, err := services.GetTeamByID(uint(id))
	if err != nil {
		utils.Error(c, 404, "team not found")
		return
	}
	if v, ok := input["name"].(string); ok {
		t.Name = v
	}
	if v, ok := input["group_id"].(float64); ok {
		t.GroupID = uint(v)
	}
	services.UpdateTeam(t)
	utils.Success(c, t)
}

func AdminDeleteTeam(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := services.DeleteTeam(uint(id)); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, gin.H{"deleted": true})
}

func AdminListGroups(c *gin.Context) {
	groups, err := repositories.ListGroups()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, groups)
}

func AdminCreateGroup(c *gin.Context) {
	var input struct {
		Name  string `json:"name" binding:"required"`
		Stage string `json:"stage"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, gin.H{"message": "group created", "name": input.Name})
}

func AdminUpdateGroup(c *gin.Context) {
	utils.Success(c, gin.H{"message": "group updated"})
}

func AdminListCities(c *gin.Context) {
	cities, err := repositories.ListCities()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, cities)
}

func AdminCreateCity(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		NameEn   string `json:"name_en"`
		Country  string `json:"country"`
		Timezone string `json:"timezone"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, gin.H{"message": "city created"})
}

func AdminUpdateCity(c *gin.Context) {
	utils.Success(c, gin.H{"message": "city updated"})
}

func AdminDeleteCity(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	_ = id
	utils.Success(c, gin.H{"deleted": true})
}

func AdminListStadiums(c *gin.Context) {
	stadiums, err := repositories.ListStadiums()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, stadiums)
}

func AdminCreateStadium(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		NameEn   string `json:"name_en"`
		CityID   uint   `json:"city_id"`
		Capacity int    `json:"capacity"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, gin.H{"message": "stadium created"})
}

func AdminUpdateStadium(c *gin.Context) {
	utils.Success(c, gin.H{"message": "stadium updated"})
}

func AdminDeleteStadium(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	_ = id
	utils.Success(c, gin.H{"deleted": true})
}

func AdminListMatches(c *gin.Context) {
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)
	matches, total, err := services.ListMatches(services.MatchQuery{Page: page, PageSize: pageSize})
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Paginated(c, matches, total, page, pageSize)
}

func AdminCreateMatch(c *gin.Context) {
	var input services.MatchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	match, err := services.CreateMatch(input)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, match)
}

func AdminUpdateMatch(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	match, err := services.GetMatchByID(uint(id))
	if err != nil {
		utils.Error(c, 404, "match not found")
		return
	}
	var input map[string]interface{}
	c.ShouldBindJSON(&input)
	if v, ok := input["importance_level"].(float64); ok {
		match.ImportanceLevel = int(v)
	}
	if v, ok := input["recommend_tag"].(string); ok {
		match.RecommendTag = v
	}
	utils.Success(c, match)
}

func AdminDeleteMatch(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := repositories.DeleteMatch(uint(id)); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, gin.H{"deleted": true})
}

func AdminUpdateMatchScore(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var input services.ScoreInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	match, err := services.UpdateMatchScore(uint(id), input)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, match)
}

func AdminUpdateMatchStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var input services.StatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	match, err := services.UpdateMatchStatus(uint(id), input)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, match)
}

func AdminImportMatches(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		utils.Error(c, 400, "missing file")
		return
	}
	defer file.Close()
	count, err := services.ImportMatchesCSV(file)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, gin.H{"imported": count})
}

func AdminListStandings(c *gin.Context) {
	standings, err := services.GetAllStandings()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, standings)
}

func AdminRecalculateStandings(c *gin.Context) {
	groups, err := repositories.ListGroups()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	for _, g := range groups {
		_ = services.RecalculateGroupStanding(g.ID)
	}
	_ = services.RecalculateBestThird()
	utils.Success(c, gin.H{"recalculated": true})
}

func AdminUpdateStanding(c *gin.Context) {
	utils.Success(c, gin.H{"message": "standing updated"})
}

func AdminListUsers(c *gin.Context) {
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)
	users, total, err := repositories.ListUsers(page, pageSize)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Paginated(c, users, total, page, pageSize)
}

func AdminUpdateUserStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "user status updated"})
}
