package repositories

import (
	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
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
