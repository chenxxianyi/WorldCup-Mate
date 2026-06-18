package repositories

import (
	"time"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"

	"gorm.io/gorm"
)

type MatchFilter struct {
	Date      string
	TeamID    uint
	GroupID   uint
	GroupName string
	Stage     string
	CityID    uint
	Status    string
	Keyword   string
	Page      int
	PageSize  int
}

func ListMatches(f MatchFilter) ([]models.Match, int64, error) {
	var matches []models.Match
	var total int64
	q := database.DB.Model(&models.Match{})

	if f.Date != "" {
		start, err := time.Parse("2006-01-02", f.Date)
		if err == nil {
			end := start.AddDate(0, 0, 1)
			q = q.Where("kickoff_time_utc >= ? AND kickoff_time_utc < ?", start, end)
		}
	}
	if f.TeamID > 0 {
		q = q.Where("home_team_id = ? OR away_team_id = ?", f.TeamID, f.TeamID)
	}
	if f.GroupID > 0 {
		q = q.Where("group_id = ?", f.GroupID)
	}
	if f.GroupName != "" {
		q = q.Joins("LEFT JOIN `groups` ON `groups`.`id` = `matches`.`group_id`").
			Where("`groups`.`name` = ?", f.GroupName)
	}
	if f.Stage != "" {
		if f.Stage == "knockout" {
			q = q.Where("stage NOT IN ?", []string{"group", "group_stage"})
		} else {
			q = q.Where("stage = ?", f.Stage)
		}
	}
	if f.CityID > 0 {
		q = q.Where("city_id = ?", f.CityID)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.Keyword != "" {
		keyword := "%" + f.Keyword + "%"
		q = q.Joins("LEFT JOIN teams home_team_search ON home_team_search.id = matches.home_team_id").
			Joins("LEFT JOIN teams away_team_search ON away_team_search.id = matches.away_team_id").
			Joins("LEFT JOIN cities city_search ON city_search.id = matches.city_id").
			Joins("LEFT JOIN stadia stadium_search ON stadium_search.id = matches.stadium_id").
			Where(`home_team_search.name LIKE ? OR home_team_search.name_en LIKE ?
				OR away_team_search.name LIKE ? OR away_team_search.name_en LIKE ?
				OR city_search.name LIKE ? OR city_search.name_en LIKE ?
				OR stadium_search.name LIKE ? OR stadium_search.name_en LIKE ?`,
				keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	q.Count(&total)
	err := q.Preload("HomeTeam").Preload("AwayTeam").Preload("Stadium").Preload("City").Preload("Group").
		Order("kickoff_time_utc ASC").
		Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).
		Find(&matches).Error
	return matches, total, err
}

func GetMatchByID(id uint) (*models.Match, error) {
	var match models.Match
	err := database.DB.Preload("HomeTeam").Preload("AwayTeam").
		Preload("Stadium").Preload("City").Preload("Group").
		First(&match, id).Error
	return &match, err
}

func GetMatchByExternal(provider, externalID string) (*models.Match, error) {
	var match models.Match
	err := database.DB.Where("external_provider = ? AND external_id = ?", provider, externalID).
		First(&match).Error
	return &match, err
}

func GetMatchesByDateRange(start, end time.Time) ([]models.Match, error) {
	var matches []models.Match
	err := database.DB.Preload("HomeTeam").Preload("AwayTeam").Preload("Stadium").Preload("City").
		Where("kickoff_time_utc >= ? AND kickoff_time_utc < ?", start, end).
		Order("kickoff_time_utc ASC").Find(&matches).Error
	return matches, err
}

func CountMatchesInSyncWindow(now time.Time) (int64, error) {
	var count int64
	start := now.Add(-150 * time.Minute)
	end := now.Add(15 * time.Minute)
	err := database.DB.Model(&models.Match{}).
		Where("(status = ? OR (status != ? AND kickoff_time_utc >= ? AND kickoff_time_utc <= ?))", "live", "finished", start, end).
		Count(&count).Error
	return count, err
}

func GetRecommendedMatches() ([]models.Match, error) {
	var matches []models.Match
	err := database.DB.Preload("HomeTeam").Preload("AwayTeam").Preload("Stadium").Preload("City").
		Where("status = ? OR (status != ? AND kickoff_time_utc >= ?)", "live", "finished", time.Now().UTC()).
		Order("CASE WHEN status = 'live' THEN 0 ELSE 1 END, kickoff_time_utc ASC").
		Limit(3).
		Find(&matches).Error
	return matches, err
}

func GetUpcomingMatches(limit int) ([]models.Match, error) {
	var matches []models.Match
	if limit <= 0 {
		limit = 10
	}
	err := database.DB.Preload("HomeTeam").Preload("AwayTeam").Preload("Stadium").Preload("City").Preload("Group").
		Where("status IN ?", []string{"scheduled", "upcoming"}).
		Where("kickoff_time_utc > ?", time.Now().UTC()).
		Order("kickoff_time_utc ASC").
		Limit(limit).
		Find(&matches).Error
	return matches, err
}

func GetMatchesByTeamID(teamID uint) ([]models.Match, error) {
	var matches []models.Match
	err := database.DB.Preload("HomeTeam").Preload("AwayTeam").Preload("Stadium").Preload("City").
		Where("home_team_id = ? OR away_team_id = ?", teamID, teamID).
		Order("kickoff_time_utc ASC").Find(&matches).Error
	return matches, err
}

func CreateMatch(match *models.Match) error {
	return database.DB.Create(match).Error
}

func UpdateMatch(match *models.Match) error {
	return database.DB.Save(match).Error
}

func DeleteMatch(id uint) error {
	return database.DB.Delete(&models.Match{}, id).Error
}

func DeleteSeedDemoMatches() error {
	return database.DB.Where("(external_provider = ? OR external_provider IS NULL) AND (external_id = ? OR external_id IS NULL) AND match_no IN ?", "", "", []int{1, 2, 3, 4}).
		Delete(&models.Match{}).Error
}

func ListMatchesForLineupAlert(startWindow, endWindow time.Time) ([]models.Match, error) {
	var matches []models.Match
	err := database.DB.Preload("HomeTeam").Preload("AwayTeam").
		Where("kickoff_time_utc >= ? AND kickoff_time_utc <= ?", startWindow, endWindow).
		Where("status NOT IN ?", []string{"finished", "cancelled"}).
		Order("kickoff_time_utc ASC").
		Find(&matches).Error
	return matches, err
}

func ListFinishedMatchesSince(since time.Time) ([]models.Match, error) {
	var matches []models.Match
	err := database.DB.Preload("HomeTeam").Preload("AwayTeam").
		Where("status = ?", "finished").
		Where("home_score IS NOT NULL AND away_score IS NOT NULL").
		Where("updated_at >= ?", since).
		Order("kickoff_time_utc DESC").
		Find(&matches).Error
	return matches, err
}

func CountMatches() int64 {
	var count int64
	database.DB.Model(&models.Match{}).Count(&count)
	return count
}

func GetMatchByDB(db *gorm.DB, id uint) (*models.Match, error) {
	var match models.Match
	err := db.First(&match, id).Error
	return &match, err
}

func ListMatchesInWindow(start, end time.Time) ([]models.Match, error) {
	var matches []models.Match
	orderExpr := `
		CASE
			WHEN status = 'live' THEN 0
			WHEN status = 'scheduled' THEN 1
			WHEN status = 'finished' THEN 2
			ELSE 3
		END ASC,
		CASE
			WHEN status = 'finished' THEN UNIX_TIMESTAMP(kickoff_time_utc) * -1
			ELSE UNIX_TIMESTAMP(kickoff_time_utc)
		END ASC`
	err := database.DB.Preload("HomeTeam").Preload("AwayTeam").Preload("Stadium").Preload("City").Preload("Group").
		Where("kickoff_time_utc >= ? AND kickoff_time_utc < ?", start, end).
		Order(orderExpr).
		Find(&matches).Error
	return matches, err
}

func CountMatchesByStageAndStatus() (map[string]map[string]int64, error) {
	type result struct {
		Stage  string
		Status string
		Count  int64
	}
	var results []result
	err := database.DB.Model(&models.Match{}).
		Select("stage, status, count(*) as count").
		Group("stage, status").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	stageStatus := make(map[string]map[string]int64)
	for _, r := range results {
		if stageStatus[r.Stage] == nil {
			stageStatus[r.Stage] = make(map[string]int64)
		}
		stageStatus[r.Stage][r.Status] = r.Count
	}
	return stageStatus, nil
}
