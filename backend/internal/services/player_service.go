package services

import (
	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"
)

func ListPlayersByTeam(teamID uint) ([]models.Player, error) {
	return repositories.ListPlayersByTeam(teamID)
}

func EnsurePlayerSourceUniqueness() error {
	if !database.DB.Migrator().HasTable(&models.Player{}) {
		return nil
	}
	if err := database.DB.Exec(`
		DELETE p1 FROM players p1
		INNER JOIN players p2
			ON p1.team_id = p2.team_id
			AND p1.source = p2.source
			AND p1.source_player_id = p2.source_player_id
			AND p1.id < p2.id
		WHERE p1.source <> ''
			AND p1.source_player_id <> ''
	`).Error; err != nil {
		return err
	}
	return database.DB.Exec(`
		UPDATE players
		SET source = 'legacy',
			source_player_id = CONCAT('legacy-', id)
		WHERE source = ''
			OR source_player_id = ''
	`).Error
}
