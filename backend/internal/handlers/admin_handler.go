package handlers

import (
	"errors"
	"strconv"
	"time"

	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"
	"worldcup-mate/internal/services"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminLogin(c *gin.Context) {
	var input services.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	token, user, err := services.Login(input)
	if err != nil {
		utils.Error(c, 401, err.Error())
		return
	}
	if user.Role != "admin" {
		utils.Error(c, 403, "admin access required")
		return
	}
	refreshToken, err := services.IssueRefreshToken(user.ID)
	if err != nil {
		utils.Error(c, 500, "failed to issue refresh token")
		return
	}
	utils.Success(c, gin.H{"token": token, "refresh_token": refreshToken, "user": user})
}

func AdminDashboard(c *gin.Context) {
	utils.Success(c, services.GetDashboard())
}

// adminIdentity returns (adminID, adminEmail) from the JWT context.
func adminIdentity(c *gin.Context) (uint, string) {
	userID, _ := c.Get("user_id")
	id, _ := userID.(uint)
	if id > 0 {
		if u, err := repositories.GetUserByID(id); err == nil {
			return id, u.Email
		}
	}
	return id, ""
}

// ---------- Teams (ADM-02) ----------

func AdminListTeams(c *gin.Context) {
	var q services.TeamQuery
	c.ShouldBindQuery(&q) // optional teamType/country
	q.Page = utils.GetPage(c)
	q.PageSize = utils.GetPageSize(c)
	teams, total, err := services.ListTeams(q)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Paginated(c, teams, total, q.Page, q.PageSize)
}

func AdminCreateTeam(c *gin.Context) {
	var input TeamCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	if exists, err := repositories.TeamExistsByCode(input.FIFACode, input.ExternalCode, 0); err != nil {
		utils.Error(c, 500, err.Error())
		return
	} else if exists {
		utils.Error(c, 409, utils.ErrTeamCodeConflict.Error())
		return
	}

	team := &models.Team{
		Name:         input.Name,
		NameEn:       input.NameEn,
		FIFACode:     input.FIFACode,
		ExternalCode: input.ExternalCode,
		TeamType:     input.TeamType,
		FlagURL:      input.FlagURL,
		Continent:    input.Continent,
		Country:      input.Country,
		Venue:        input.Venue,
		GroupID:      input.GroupID,
		Coach:        input.Coach,
		Description:  input.Description,
	}
	if team.TeamType == "" {
		team.TeamType = "national"
	}
	if err := repositories.CreateTeam(team); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "team", strconv.FormatUint(uint64(team.ID), 10), "create", nil, team)
	utils.Success(c, team)
}

func AdminUpdateTeam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid team id")
		return
	}
	var input TeamUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	t, err := repositories.GetTeamByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, 404, "team not found")
			return
		}
		utils.Error(c, 500, err.Error())
		return
	}
	before := *t
	if input.Name != nil {
		t.Name = *input.Name
	}
	if input.NameEn != nil {
		t.NameEn = *input.NameEn
	}
	if input.FIFACode != nil {
		t.FIFACode = input.FIFACode
	}
	if input.ExternalCode != nil {
		t.ExternalCode = input.ExternalCode
	}
	if input.TeamType != nil {
		t.TeamType = *input.TeamType
	}
	if input.FlagURL != nil {
		t.FlagURL = *input.FlagURL
	}
	if input.Continent != nil {
		t.Continent = *input.Continent
	}
	if input.Country != nil {
		t.Country = *input.Country
	}
	if input.Venue != nil {
		t.Venue = *input.Venue
	}
	if input.GroupID != nil {
		t.GroupID = input.GroupID
	}
	if input.Coach != nil {
		t.Coach = *input.Coach
	}
	if input.Description != nil {
		t.Description = *input.Description
	}
	if exists, err := repositories.TeamExistsByCode(t.FIFACode, t.ExternalCode, t.ID); err != nil {
		utils.Error(c, 500, err.Error())
		return
	} else if exists {
		utils.Error(c, 409, utils.ErrTeamCodeConflict.Error())
		return
	}
	if err := repositories.UpdateTeam(t); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "team", strconv.FormatUint(uint64(t.ID), 10), "update", before, t)
	utils.Success(c, t)
}

func AdminDeleteTeam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid team id")
		return
	}
	t, err := repositories.GetTeamByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, 404, "team not found")
			return
		}
		utils.Error(c, 500, err.Error())
		return
	}
	if repositories.CountMatchesByTeam(t.ID) > 0 || repositories.CountStandingsByTeam(t.ID) > 0 || repositories.CountFavoriteTeams(t.ID) > 0 {
		utils.Error(c, 409, utils.ErrTeamInUse.Error())
		return
	}
	if err := repositories.DeleteTeam(t.ID); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "team", strconv.FormatUint(uint64(t.ID), 10), "delete", t, nil)
	utils.Success(c, gin.H{"deleted": true})
}

// ---------- Groups (ADM-03) ----------

func AdminListGroups(c *gin.Context) {
	groups, err := repositories.ListGroups()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, groups)
}

func AdminCreateGroup(c *gin.Context) {
	var input GroupCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	if exists, err := repositories.GroupExistsByName(input.Name, 0); err != nil {
		utils.Error(c, 500, err.Error())
		return
	} else if exists {
		utils.Error(c, 409, utils.ErrGroupNameConflict.Error())
		return
	}
	group := &models.Group{Name: input.Name, Stage: input.Stage}
	if err := repositories.CreateGroup(group); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "group", strconv.FormatUint(uint64(group.ID), 10), "create", nil, group)
	utils.Success(c, group)
}

func AdminUpdateGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid group id")
		return
	}
	var input GroupUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	group, err := repositories.GetGroupByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, 404, "group not found")
			return
		}
		utils.Error(c, 500, err.Error())
		return
	}
	before := *group
	if input.Name != nil {
		if exists, err := repositories.GroupExistsByName(*input.Name, group.ID); err != nil {
			utils.Error(c, 500, err.Error())
			return
		} else if exists {
			utils.Error(c, 409, utils.ErrGroupNameConflict.Error())
			return
		}
		group.Name = *input.Name
	}
	if input.Stage != nil {
		group.Stage = *input.Stage
	}
	if err := repositories.UpdateGroup(group); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "group", strconv.FormatUint(uint64(group.ID), 10), "update", before, group)
	utils.Success(c, group)
}

func AdminDeleteGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid group id")
		return
	}
	group, err := repositories.GetGroupByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, 404, "group not found")
			return
		}
		utils.Error(c, 500, err.Error())
		return
	}
	if repositories.CountTeamsByGroup(group.ID) > 0 || repositories.CountMatchesByGroup(group.ID) > 0 {
		utils.Error(c, 409, utils.ErrGroupInUse.Error())
		return
	}
	if err := repositories.DeleteGroup(group.ID); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "group", strconv.FormatUint(uint64(group.ID), 10), "delete", group, nil)
	utils.Success(c, gin.H{"deleted": true})
}

// ---------- Cities (ADM-04) ----------

func AdminListCities(c *gin.Context) {
	cities, err := repositories.ListCities()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, cities)
}

func AdminCreateCity(c *gin.Context) {
	var input CityCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	city := &models.City{Name: input.Name, NameEn: input.NameEn, Country: input.Country, Timezone: input.Timezone}
	if err := repositories.CreateCity(city); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "city", strconv.FormatUint(uint64(city.ID), 10), "create", nil, city)
	utils.Success(c, city)
}

func AdminUpdateCity(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid city id")
		return
	}
	var input CityUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	city, err := repositories.GetCityByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, 404, "city not found")
			return
		}
		utils.Error(c, 500, err.Error())
		return
	}
	before := *city
	if input.Name != nil {
		city.Name = *input.Name
	}
	if input.NameEn != nil {
		city.NameEn = *input.NameEn
	}
	if input.Country != nil {
		city.Country = *input.Country
	}
	if input.Timezone != nil {
		city.Timezone = *input.Timezone
	}
	if err := repositories.UpdateCity(city); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "city", strconv.FormatUint(uint64(city.ID), 10), "update", before, city)
	utils.Success(c, city)
}

func AdminDeleteCity(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid city id")
		return
	}
	city, err := repositories.GetCityByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, 404, "city not found")
			return
		}
		utils.Error(c, 500, err.Error())
		return
	}
	if repositories.CountStadiumsByCity(city.ID) > 0 || repositories.CountMatchesByCity(city.ID) > 0 {
		utils.Error(c, 409, utils.ErrCityInUse.Error())
		return
	}
	if err := repositories.DeleteCity(city.ID); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "city", strconv.FormatUint(uint64(city.ID), 10), "delete", city, nil)
	utils.Success(c, gin.H{"deleted": true})
}

// ---------- Stadiums (ADM-04) ----------

func AdminListStadiums(c *gin.Context) {
	stadiums, err := repositories.ListStadiums()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, stadiums)
}

func AdminCreateStadium(c *gin.Context) {
	var input StadiumCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	if _, err := repositories.GetCityByID(input.CityID); err != nil {
		utils.Error(c, 400, "city does not exist")
		return
	}
	stadium := &models.Stadium{Name: input.Name, NameEn: input.NameEn, CityID: input.CityID, Capacity: input.Capacity}
	if err := repositories.CreateStadium(stadium); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "stadium", strconv.FormatUint(uint64(stadium.ID), 10), "create", nil, stadium)
	utils.Success(c, stadium)
}

func AdminUpdateStadium(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid stadium id")
		return
	}
	var input StadiumUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	stadium, err := repositories.GetStadiumByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, 404, "stadium not found")
			return
		}
		utils.Error(c, 500, err.Error())
		return
	}
	before := *stadium
	if input.Name != nil {
		stadium.Name = *input.Name
	}
	if input.NameEn != nil {
		stadium.NameEn = *input.NameEn
	}
	if input.CityID != nil {
		if _, err := repositories.GetCityByID(*input.CityID); err != nil {
			utils.Error(c, 400, "city does not exist")
			return
		}
		stadium.CityID = *input.CityID
	}
	if input.Capacity != nil {
		stadium.Capacity = *input.Capacity
	}
	if err := repositories.UpdateStadium(stadium); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "stadium", strconv.FormatUint(uint64(stadium.ID), 10), "update", before, stadium)
	utils.Success(c, stadium)
}

func AdminDeleteStadium(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid stadium id")
		return
	}
	stadium, err := repositories.GetStadiumByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, 404, "stadium not found")
			return
		}
		utils.Error(c, 500, err.Error())
		return
	}
	if repositories.CountMatchesByStadium(stadium.ID) > 0 {
		utils.Error(c, 409, utils.ErrStadiumInUse.Error())
		return
	}
	if err := repositories.DeleteStadium(stadium.ID); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "stadium", strconv.FormatUint(uint64(stadium.ID), 10), "delete", stadium, nil)
	utils.Success(c, gin.H{"deleted": true})
}

// ---------- Matches (ADM-05) ----------

func AdminListMatches(c *gin.Context) {
	var q services.MatchQuery
	c.ShouldBindQuery(&q) // optional competitionId/season/matchday
	q.Page = utils.GetPage(c)
	q.PageSize = utils.GetPageSize(c)
	matches, total, err := services.ListMatches(q)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Paginated(c, matches, total, q.Page, q.PageSize)
}

func AdminCreateMatch(c *gin.Context) {
	var input services.MatchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	match, err := services.CreateMatch(input)
	if err != nil {
		if errors.Is(err, utils.ErrInvalidTime) {
			utils.Error(c, 400, err.Error())
			return
		}
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "match", strconv.FormatUint(uint64(match.ID), 10), "create", nil, match)
	utils.Success(c, match)
}

func AdminUpdateMatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid match id")
		return
	}
	var input MatchUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	match, err := services.GetMatchByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, 404, "match not found")
			return
		}
		utils.Error(c, 500, err.Error())
		return
	}
	before := *match
	if input.MatchNo != nil {
		match.MatchNo = *input.MatchNo
	}
	if input.HomeTeamID != nil {
		match.HomeTeamID = *input.HomeTeamID
	}
	if input.AwayTeamID != nil {
		match.AwayTeamID = *input.AwayTeamID
	}
	if input.GroupID != nil {
		match.GroupID = input.GroupID
	}
	if input.Stage != nil {
		match.Stage = *input.Stage
	}
	if input.StadiumID != nil {
		match.StadiumID = *input.StadiumID
	}
	if input.CityID != nil {
		match.CityID = *input.CityID
	}
	if input.KickoffTimeUTC != nil {
		parsed, err := time.Parse(time.RFC3339, *input.KickoffTimeUTC)
		if err != nil {
			utils.Error(c, 400, utils.ErrInvalidTime.Error())
			return
		}
		match.KickoffTimeUTC = parsed.UTC()
	}
	if input.ImportanceLevel != nil {
		match.ImportanceLevel = *input.ImportanceLevel
	}
	if input.RecommendTag != nil {
		match.RecommendTag = *input.RecommendTag
	}
	if input.CompetitionID != nil {
		match.CompetitionID = input.CompetitionID
	}
	if input.Season != nil {
		match.Season = input.Season
	}
	if input.Matchday != nil {
		match.Matchday = input.Matchday
	}
	if err := services.UpdateMatch(match); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "match", strconv.FormatUint(uint64(match.ID), 10), "update", before, match)
	utils.Success(c, match)
}

func AdminDeleteMatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid match id")
		return
	}
	match, err := services.GetMatchByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, 404, "match not found")
			return
		}
		utils.Error(c, 500, err.Error())
		return
	}
	if repositories.CountRemindersByMatch(match.ID) > 0 || repositories.CountFavoriteMatches(match.ID) > 0 {
		utils.Error(c, 409, utils.ErrMatchInUse.Error())
		return
	}
	if err := repositories.DeleteMatch(match.ID); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "match", strconv.FormatUint(uint64(match.ID), 10), "delete", match, nil)
	utils.Success(c, gin.H{"deleted": true})
}

func AdminUpdateMatchScore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid match id")
		return
	}
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
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "match", strconv.FormatUint(uint64(match.ID), 10), "score", nil, match)
	utils.Success(c, match)
}

func AdminUpdateMatchStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid match id")
		return
	}
	var input services.StatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	match, err := services.UpdateMatchStatus(uint(id), input)
	if err != nil {
		if errors.Is(err, utils.ErrInvalidStatusTransition) {
			utils.Error(c, 409, err.Error())
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, 404, "match not found")
			return
		}
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "match", strconv.FormatUint(uint64(match.ID), 10), "status", nil, match)
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

// ---------- Standings (ADM-06) ----------

// AdminRecalculateLeagueStanding rebuilds the TOTAL league table from
// finished matches. Body: { competition_id (required), season (optional,
// defaults to the competition's current season) }.
func AdminRecalculateLeagueStanding(c *gin.Context) {
	var input struct {
		CompetitionID uint `json:"competition_id" binding:"required"`
		Season        int  `json:"season"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	if input.Season <= 0 {
		if competition, err := repositories.GetCompetitionByID(input.CompetitionID); err == nil {
			input.Season = competition.Season
		}
	}
	if input.Season <= 0 {
		utils.Error(c, 400, "season is required")
		return
	}
	if err := services.RecalculateLeagueStanding(input.CompetitionID, input.Season); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "standing", "league", "recalculate", nil, gin.H{"competition_id": input.CompetitionID, "season": input.Season})
	utils.Success(c, gin.H{"recalculated": true, "competition_id": input.CompetitionID, "season": input.Season})
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
		if err := services.RecalculateGroupStanding(g.ID); err != nil {
			utils.Error(c, 500, err.Error())
			return
		}
	}
	if err := services.RecalculateBestThird(); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "standing", "all-groups", "recalculate", nil, gin.H{"recalculated": true})
	utils.Success(c, gin.H{"recalculated": true})
}

func AdminUpdateStanding(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid standing id")
		return
	}
	var input StandingUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	if input.Reason == "" {
		utils.Error(c, 400, "reason is required for manual standings correction")
		return
	}
	standing, err := repositories.GetGroupStandingByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, 404, "standing not found")
			return
		}
		utils.Error(c, 500, err.Error())
		return
	}
	before := *standing
	if input.Position != nil {
		standing.Rank = *input.Position
	}
	if input.Points != nil {
		standing.Points = *input.Points
	}
	if input.Played != nil {
		standing.Played = *input.Played
	}
	if input.Won != nil {
		standing.Won = *input.Won
	}
	if input.Drawn != nil {
		standing.Drawn = *input.Drawn
	}
	if input.Lost != nil {
		standing.Lost = *input.Lost
	}
	if input.GoalsFor != nil {
		standing.GoalsFor = *input.GoalsFor
	}
	if input.GoalsAgainst != nil {
		standing.GoalsAgainst = *input.GoalsAgainst
	}
	if input.GoalDifference != nil {
		standing.GoalDifference = *input.GoalDifference
	}
	if err := repositories.UpdateGroupStanding(standing); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "standing", strconv.FormatUint(uint64(standing.ID), 10), "update", before, standing)
	utils.Success(c, standing)
}

// ---------- Users (ADM-06) ----------

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
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid user id")
		return
	}
	var input struct {
		Status string `json:"status" binding:"required,oneof=active disabled"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, 400, "status must be active or disabled")
		return
	}
	user, err := repositories.GetUserByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, 404, "user not found")
			return
		}
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	if user.ID == adminID {
		utils.Error(c, 400, "cannot change your own status")
		return
	}
	before := *user
	user.Status = input.Status
	if err := repositories.UpdateUser(user); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	// SEC-04: disabling a user must kill every active session immediately.
	if user.Status == "disabled" {
		_ = services.RevokeAllRefreshTokens(user.ID)
	}
	services.RecordAudit(adminID, adminEmail, "user", strconv.FormatUint(uint64(user.ID), 10), "status", before, user)
	utils.Success(c, user)
}

// ---------- Sync history (ADM-13) ----------

func AdminListSyncHistory(c *gin.Context) {
	limit := utils.GetPageSize(c)
	if limit > 200 {
		limit = 200
	}
	states, err := repositories.ListSyncHistory(limit)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, states)
}

// ---------- Reminder operations (ADM-13) ----------

func AdminListReminders(c *gin.Context) {
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)
	reminders, total, err := repositories.ListRemindersWithStats(page, pageSize)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Paginated(c, reminders, total, page, pageSize)
}

func AdminRetryReminder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, 400, "invalid reminder id")
		return
	}
	reminder, err := repositories.GetReminderByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, 404, "reminder not found")
			return
		}
		utils.Error(c, 500, err.Error())
		return
	}
	// Requeue a failed/sending reminder for a fresh attempt.
	now := time.Now().UTC()
	reminder.Status = "pending"
	reminder.RetryCount = 0
	reminder.ClaimToken = ""
	reminder.ClaimedAt = nil
	reminder.NextRetryAt = &now
	reminder.LastError = ""
	if err := repositories.UpdateReminder(reminder); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	adminID, adminEmail := adminIdentity(c)
	services.RecordAudit(adminID, adminEmail, "reminder", strconv.FormatUint(uint64(reminder.ID), 10), "retry", nil, reminder)
	utils.Success(c, gin.H{"retried": true, "reminder_id": reminder.ID})
}
