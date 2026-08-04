package repositories

import (
	"errors"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"

	"gorm.io/gorm"
)

func GetStandingsByGroupID(groupID uint) ([]models.GroupStanding, error) {
	var standings []models.GroupStanding
	err := database.DB.Preload("Team").Where("group_id = ?", groupID).
		Order("`rank` ASC").Find(&standings).Error
	return standings, err
}

func GetAllStandings() ([]models.GroupStanding, error) {
	var standings []models.GroupStanding
	err := database.DB.Preload("Team").Preload("Group").
		Order("group_id ASC, `rank` ASC").Find(&standings).Error
	return standings, err
}

func GetBestThird() ([]models.GroupStanding, error) {
	var standings []models.GroupStanding
	err := database.DB.Preload("Team").Preload("Group").
		Where("`rank` = ?", 3).
		Order("points DESC, goal_difference DESC, goals_for DESC").
		Find(&standings).Error
	return standings, err
}

func UpsertStanding(standing *models.GroupStanding) error {
	var existing models.GroupStanding
	err := database.DB.Where("group_id = ? AND team_id = ?", standing.GroupID, standing.TeamID).First(&existing).Error
	if err != nil {
		return database.DB.Create(standing).Error
	}
	standing.ID = existing.ID
	return database.DB.Save(standing).Error
}

func DeleteStandingsByGroupID(groupID uint) error {
	return database.DB.Where("group_id = ?", groupID).Delete(&models.GroupStanding{}).Error
}

// GetLeagueStandings queries the new league_standings table.
// competitionID is required; season/standingType are optional (0/"" = no filter).
func GetLeagueStandings(competitionID uint, season int, standingType string) ([]models.LeagueStanding, error) {
	var standings []models.LeagueStanding
	q := database.DB.Preload("Team").Where("competition_id = ?", competitionID)
	if season > 0 {
		q = q.Where("season = ?", season)
	}
	if standingType != "" {
		q = q.Where("type = ?", standingType)
	}
	err := q.Order("`position` ASC").Find(&standings).Error
	return standings, err
}

// UpsertLeagueStanding inserts or updates a league standing row keyed by
// (competition_id, season, type, team_id).
func UpsertLeagueStanding(standing *models.LeagueStanding) error {
	var existing models.LeagueStanding
	err := database.DB.Where(
		"competition_id = ? AND season = ? AND type = ? AND team_id = ?",
		standing.CompetitionID, standing.Season, standing.Type, standing.TeamID,
	).First(&existing).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return database.DB.Create(standing).Error
	}
	standing.ID = existing.ID
	return database.DB.Save(standing).Error
}

// DeleteLeagueStandings removes league standing rows of a competition season.
// standingType "total"/"home"/"away", or "" for all types.
func DeleteLeagueStandings(competitionID uint, season int, standingType string) error {
	q := database.DB.Where("competition_id = ? AND season = ?", competitionID, season)
	if standingType != "" {
		q = q.Where("type = ?", standingType)
	}
	return q.Delete(&models.LeagueStanding{}).Error
}
