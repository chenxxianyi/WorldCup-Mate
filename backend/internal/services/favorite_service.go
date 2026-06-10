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

func AddFavoriteTeam(userID, teamID uint) error {
	return repositories.AddFavoriteTeam(userID, teamID)
}

func RemoveFavoriteTeam(userID, teamID uint) error {
	return repositories.RemoveFavoriteTeam(userID, teamID)
}

func ToggleFavoriteMatch(userID, matchID uint) (bool, error) {
	if repositories.IsMatchFavorited(userID, matchID) {
		return false, repositories.RemoveFavoriteMatch(userID, matchID)
	}
	return true, repositories.AddFavoriteMatch(userID, matchID)
}

func AddFavoriteMatch(userID, matchID uint) error {
	return repositories.AddFavoriteMatch(userID, matchID)
}

func RemoveFavoriteMatch(userID, matchID uint) error {
	return repositories.RemoveFavoriteMatch(userID, matchID)
}

func GetFavoriteTeams(userID uint) ([]models.UserFavoriteTeam, error) {
	return repositories.GetFavoriteTeams(userID)
}

func GetFavoriteMatches(userID uint) ([]models.UserFavoriteMatch, error) {
	return repositories.GetFavoriteMatches(userID)
}
