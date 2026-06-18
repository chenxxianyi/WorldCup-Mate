package repositories

import (
	"errors"
	"time"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"

	"gorm.io/gorm"
)

func GetLineupsByMatch(matchID uint) ([]models.MatchLineup, error) {
	var lineups []models.MatchLineup
	err := database.DB.
		Preload("Team").
		Preload("Players", func(db *gorm.DB) *gorm.DB {
			return db.Order("role ASC, sort_order ASC, shirt_number ASC, id ASC")
		}).
		Preload("Players.Player").
		Where("match_id = ?", matchID).
		Order("FIELD(side, 'home', 'away'), id ASC").
		Find(&lineups).Error
	return lineups, err
}

func UpsertMatchLineup(lineup *models.MatchLineup) (*models.MatchLineup, error) {
	var existing models.MatchLineup
	err := database.DB.Where("match_id = ? AND team_id = ?", lineup.MatchID, lineup.TeamID).First(&existing).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if err := database.DB.Create(lineup).Error; err != nil {
			return nil, err
		}
		return lineup, nil
	}

	existing.Side = lineup.Side
	existing.Formation = lineup.Formation
	existing.CoachName = lineup.CoachName
	existing.CoachNameEn = lineup.CoachNameEn
	existing.Source = lineup.Source
	existing.SourceMatchID = lineup.SourceMatchID
	existing.ExternalTeamID = lineup.ExternalTeamID
	existing.Status = lineup.Status
	existing.LastSyncedAt = lineup.LastSyncedAt
	if err := database.DB.Save(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

func ReplaceLineupPlayers(lineupID uint, players []models.MatchLineupPlayer) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("match_lineup_id = ?", lineupID).Delete(&models.MatchLineupPlayer{}).Error; err != nil {
			return err
		}
		if len(players) == 0 {
			return nil
		}
		for i := range players {
			players[i].MatchLineupID = lineupID
		}
		return tx.Create(&players).Error
	})
}

func ListMatchesNeedingLineupSync(now time.Time, pregameWindow time.Duration) ([]models.Match, error) {
	if pregameWindow <= 0 {
		pregameWindow = 90 * time.Minute
	}
	start := now.UTC().Add(-72 * time.Hour)
	end := now.UTC().Add(pregameWindow)

	var matches []models.Match
	err := database.DB.
		Preload("HomeTeam").
		Preload("AwayTeam").
		Where("status = ? OR kickoff_time_utc BETWEEN ? AND ?", "live", start, end).
		Order("kickoff_time_utc ASC").
		Find(&matches).Error
	return matches, err
}

func GetPlayerBySource(teamID uint, source, sourcePlayerID string) (*models.Player, error) {
	var player models.Player
	err := database.DB.
		Where("team_id = ? AND source = ? AND source_player_id = ? AND is_active = ?", teamID, source, sourcePlayerID, true).
		First(&player).Error
	return &player, err
}

func FindPlayerForLineup(teamID uint, nameEn string, shirtNumber int, position string) (*models.Player, error) {
	var player models.Player
	q := database.DB.Where("team_id = ? AND is_active = ?", teamID, true)
	if nameEn != "" {
		err := q.Where("LOWER(name_en) = LOWER(?)", nameEn).First(&player).Error
		if err == nil {
			return &player, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	if shirtNumber > 0 && position != "" {
		err := q.Where("shirt_number = ? AND position = ?", shirtNumber, position).First(&player).Error
		if err == nil {
			return &player, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	return nil, gorm.ErrRecordNotFound
}
