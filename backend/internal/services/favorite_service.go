package services

import (
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"
)

func ToggleFavoriteTeam(userID, teamID uint) (bool, error) {
	if repositories.IsTeamFavorited(userID, teamID) {
		return false, repositories.RemoveFavoriteTeam(userID, teamID)
	}
	return true, repositories.AddFavoriteTeam(userID, teamID)
}

func ToggleFavoriteMatch(userID, matchID uint) (bool, error) {
	if repositories.IsMatchFavorited(userID, matchID) {
		return false, repositories.RemoveFavoriteMatch(userID, matchID)
	}
	return true, repositories.AddFavoriteMatch(userID, matchID)
}

func GetFavoriteTeams(userID uint) ([]models.UserFavoriteTeam, error) {
	return repositories.GetFavoriteTeams(userID)
}

func GetFavoriteMatches(userID uint) ([]models.UserFavoriteMatch, error) {
	return repositories.GetFavoriteMatches(userID)
}
