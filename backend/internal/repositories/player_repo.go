package repositories

import (
	"errors"
	"time"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ListPlayersByTeam(teamID uint) ([]models.Player, error) {
	var players []models.Player
	err := database.DB.
		Where("team_id = ? AND is_active = ?", teamID, true).
		Order("shirt_number ASC, id ASC").
		Find(&players).Error
	return players, err
}

func UpsertPlayer(player *models.Player) (string, error) {
	var existing models.Player
	err := database.DB.
		Where("team_id = ? AND source = ? AND source_player_id = ?", player.TeamID, player.Source, player.SourcePlayerID).
		Select("id").
		First(&existing).Error
	existed := err == nil
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "skipped", err
	}

	err = database.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "team_id"},
			{Name: "source"},
			{Name: "source_player_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"name_en",
			"shirt_number",
			"position",
			"position_label",
			"photo_url",
			"club",
			"external_team_id",
			"is_active",
			"last_synced_at",
		}),
	}).Create(player).Error
	if err != nil {
		return "skipped", err
	}
	if existed {
		return "updated", nil
	}
	return "created", nil
}

func MarkMissingPlayersInactive(teamID uint, source string, activeSourcePlayerIDs []string, syncedAt time.Time) (int64, error) {
	q := database.DB.Model(&models.Player{}).
		Where("team_id = ? AND source = ? AND is_active = ?", teamID, source, true)
	if len(activeSourcePlayerIDs) > 0 {
		q = q.Where("source_player_id NOT IN ?", activeSourcePlayerIDs)
	}
	result := q.Updates(map[string]interface{}{
		"is_active":      false,
		"last_synced_at": syncedAt,
	})
	return result.RowsAffected, result.Error
}
