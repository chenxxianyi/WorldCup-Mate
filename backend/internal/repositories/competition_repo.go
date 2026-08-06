package repositories

import (
	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
)

func ListCompetitions() ([]models.Competition, error) {
	var competitions []models.Competition
	err := database.DB.Order("sort_order ASC, id ASC").Find(&competitions).Error
	return competitions, err
}

// ListActiveCompetitions returns only the competitions enabled for the
// public frontend (admin management can disable a league without deleting
// its data — ADM/competition-config).
func ListActiveCompetitions() ([]models.Competition, error) {
	var competitions []models.Competition
	err := database.DB.Where("status = ?", "active").Order("sort_order ASC, id ASC").Find(&competitions).Error
	return competitions, err
}

func GetCompetitionByID(id uint) (*models.Competition, error) {
	var competition models.Competition
	err := database.DB.First(&competition, id).Error
	return &competition, err
}

func GetCompetitionByCode(code string) (*models.Competition, error) {
	var competition models.Competition
	err := database.DB.Where("code = ?", code).First(&competition).Error
	return &competition, err
}

func CreateCompetition(competition *models.Competition) error {
	return database.DB.Create(competition).Error
}

func UpdateCompetition(competition *models.Competition) error {
	return database.DB.Save(competition).Error
}

func CountCompetitions() int64 {
	var count int64
	database.DB.Model(&models.Competition{}).Count(&count)
	return count
}

// ListCompetitionSeasons returns distinct seasons for a competition,
// ordered descending (most recent first).
func ListCompetitionSeasons(competitionID uint) ([]int, error) {
	var seasons []int
	err := database.DB.Model(&models.Match{}).
		Where("competition_id = ? AND season IS NOT NULL", competitionID).
		Select("DISTINCT season").Order("season DESC").Find(&seasons).Error
	return seasons, err
}

// GetLatestMatchday returns the latest matchday number for the given
// competition and season (league only).
func GetLatestMatchday(competitionID uint, season int) (*int, error) {
	var maxMatchday *int
	err := database.DB.Model(&models.Match{}).
		Where("competition_id = ? AND season = ?", competitionID, season).
		Order("matchday DESC").Limit(1).Pluck("matchday", &maxMatchday).Error
	return maxMatchday, err
}

// CountMatchesByCompetitionAndSeason returns the total match count
// for a competition in the given season.
func CountMatchesByCompetitionAndSeason(competitionID uint, season int) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Match{}).
		Where("competition_id = ? AND season = ?", competitionID, season).Count(&count).Error
	return count, err
}
