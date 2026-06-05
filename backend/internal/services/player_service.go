package services

import (
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"
)

func ListPlayersByTeam(teamID uint) ([]models.Player, error) {
	return repositories.ListPlayersByTeam(teamID)
}
