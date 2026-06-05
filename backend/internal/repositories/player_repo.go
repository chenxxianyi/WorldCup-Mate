package repositories

import (
	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
)

func ListPlayersByTeam(teamID uint) ([]models.Player, error) {
	var players []models.Player
	err := database.DB.
		Where("team_id = ?", teamID).
		Order("shirt_number ASC, id ASC").
		Find(&players).Error
	return players, err
}
