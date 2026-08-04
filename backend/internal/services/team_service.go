package services

import (
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"
)

type TeamQuery struct {
	Continent string `form:"continent"`
	GroupID   uint   `form:"groupId"`
	Keyword   string `form:"keyword"`
	TeamType  string `form:"teamType"` // optional: national | club
	Country   string `form:"country"`  // optional
	Page      int
	PageSize  int
}

func ListTeams(q TeamQuery) ([]models.Team, int64, error) {
	if q.TeamType != "" || q.Country != "" {
		return repositories.ListTeamsFiltered(repositories.TeamFilter{
			Continent: q.Continent,
			Keyword:   q.Keyword,
			GroupID:   q.GroupID,
			TeamType:  q.TeamType,
			Country:   q.Country,
			Page:      q.Page,
			PageSize:  q.PageSize,
		})
	}
	return repositories.ListTeams(q.Continent, q.Keyword, q.GroupID, q.Page, q.PageSize)
}

func GetTeamByID(id uint) (*models.Team, error) {
	return repositories.GetTeamByID(id)
}

func GetTeamMatches(teamID uint) ([]models.Match, error) {
	return repositories.GetMatchesByTeamID(teamID)
}

func CreateTeam(team *models.Team) error {
	return repositories.CreateTeam(team)
}

func UpdateTeam(team *models.Team) error {
	return repositories.UpdateTeam(team)
}

func DeleteTeam(id uint) error {
	return repositories.DeleteTeam(id)
}

func CountTeams() int64 {
	var count int64
	repositories.ListTeams("", "", 0, 1, 1)
	_, total, _ := repositories.ListTeams("", "", 0, 1, 1)
	count = total
	return count
}
