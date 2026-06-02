package repositories

import (
	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
)

func AddFavoriteTeam(userID, teamID uint) error {
	fav := models.UserFavoriteTeam{UserID: userID, TeamID: teamID}
	return database.DB.Where("user_id = ? AND team_id = ?", userID, teamID).FirstOrCreate(&fav).Error
}

func RemoveFavoriteTeam(userID, teamID uint) error {
	return database.DB.Where("user_id = ? AND team_id = ?", userID, teamID).Delete(&models.UserFavoriteTeam{}).Error
}

func IsTeamFavorited(userID, teamID uint) bool {
	var count int64
	database.DB.Model(&models.UserFavoriteTeam{}).Where("user_id = ? AND team_id = ?", userID, teamID).Count(&count)
	return count > 0
}

func GetFavoriteTeams(userID uint) ([]models.UserFavoriteTeam, error) {
	var favs []models.UserFavoriteTeam
	err := database.DB.Preload("Team").Where("user_id = ?", userID).Find(&favs).Error
	return favs, err
}

func AddFavoriteMatch(userID, matchID uint) error {
	fav := models.UserFavoriteMatch{UserID: userID, MatchID: matchID}
	return database.DB.Where("user_id = ? AND match_id = ?", userID, matchID).FirstOrCreate(&fav).Error
}

func RemoveFavoriteMatch(userID, matchID uint) error {
	return database.DB.Where("user_id = ? AND match_id = ?", userID, matchID).Delete(&models.UserFavoriteMatch{}).Error
}

func IsMatchFavorited(userID, matchID uint) bool {
	var count int64
	database.DB.Model(&models.UserFavoriteMatch{}).Where("user_id = ? AND match_id = ?", userID, matchID).Count(&count)
	return count > 0
}

func GetFavoriteMatches(userID uint) ([]models.UserFavoriteMatch, error) {
	var favs []models.UserFavoriteMatch
	err := database.DB.Preload("Match").Preload("Match.HomeTeam").Preload("Match.AwayTeam").
		Where("user_id = ?", userID).Find(&favs).Error
	return favs, err
}
