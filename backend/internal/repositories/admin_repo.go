package repositories

import (
	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
)

// Reference checks and admin CRUD helpers (ADM-01/02/03/04/05/06).
// Kept in one file so existing repositories stay untouched.

// ---- Teams ----

// TeamExistsByCode reports whether fifa_code or external_code is already
// taken by another team (excludeID excludes the team being updated).
func TeamExistsByCode(fifaCode, externalCode *string, excludeID uint) (bool, error) {
	q := database.DB.Model(&models.Team{})
	if fifaCode != nil && *fifaCode != "" {
		q = q.Where("fifa_code = ?", *fifaCode)
	}
	if externalCode != nil && *externalCode != "" {
		q = q.Or("external_code = ?", *externalCode)
	}
	var count int64
	if err := q.Where("id != ?", excludeID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func CountMatchesByTeam(teamID uint) int64 {
	var count int64
	database.DB.Model(&models.Match{}).
		Where("home_team_id = ? OR away_team_id = ?", teamID, teamID).
		Count(&count)
	return count
}

func CountStandingsByTeam(teamID uint) int64 {
	var count int64
	database.DB.Model(&models.GroupStanding{}).Where("team_id = ?", teamID).Count(&count)
	database.DB.Model(&models.LeagueStanding{}).Where("team_id = ?", teamID).Count(&count)
	return count
}

func CountFavoriteTeams(teamID uint) int64 {
	var count int64
	database.DB.Model(&models.UserFavoriteTeam{}).Where("team_id = ?", teamID).Count(&count)
	return count
}

func CountUsers() int64 {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	return count
}

func CountAllReminders() int64 {
	var count int64
	database.DB.Model(&models.Reminder{}).Count(&count)
	return count
}

// ---- Groups ----

func GroupExistsByName(name string, excludeID uint) (bool, error) {
	var count int64
	err := database.DB.Model(&models.Group{}).
		Where("name = ? AND id != ?", name, excludeID).
		Count(&count).Error
	return count > 0, err
}

func CountTeamsByGroup(groupID uint) int64 {
	var count int64
	database.DB.Model(&models.Team{}).Where("group_id = ?", groupID).Count(&count)
	return count
}

func CountMatchesByGroup(groupID uint) int64 {
	var count int64
	database.DB.Model(&models.Match{}).Where("group_id = ?", groupID).Count(&count)
	return count
}

// ---- Cities & Stadiums ----

func CountStadiumsByCity(cityID uint) int64 {
	var count int64
	database.DB.Model(&models.Stadium{}).Where("city_id = ?", cityID).Count(&count)
	return count
}

func CountMatchesByCity(cityID uint) int64 {
	var count int64
	database.DB.Model(&models.Match{}).Where("city_id = ?", cityID).Count(&count)
	return count
}

func CountMatchesByStadium(stadiumID uint) int64 {
	var count int64
	database.DB.Model(&models.Match{}).Where("stadium_id = ?", stadiumID).Count(&count)
	return count
}

// ---- Matches ----

func CountRemindersByMatch(matchID uint) int64 {
	var count int64
	database.DB.Model(&models.Reminder{}).Where("match_id = ?", matchID).Count(&count)
	return count
}

func CountFavoriteMatches(matchID uint) int64 {
	var count int64
	database.DB.Model(&models.UserFavoriteMatch{}).Where("match_id = ?", matchID).Count(&count)
	return count
}

// ---- Standings ----

func GetGroupStandingByID(id uint) (*models.GroupStanding, error) {
	var standing models.GroupStanding
	err := database.DB.Preload("Team").First(&standing, id).Error
	return &standing, err
}

func UpdateGroupStanding(standing *models.GroupStanding) error {
	return database.DB.Save(standing).Error
}

func DeleteGroup(id uint) error {
	return database.DB.Delete(&models.Group{}, id).Error
}
