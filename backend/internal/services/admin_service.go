package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"
)

type DashboardData struct {
	TotalMatches  int64 `json:"total_matches"`
	TotalGroups   int64 `json:"total_groups"`
	TotalTeams    int64 `json:"total_teams"`
	TotalReminders int64 `json:"total_reminders"`
}

func GetDashboard() DashboardData {
	return DashboardData{
		TotalMatches:  repositories.CountMatches(),
		TotalGroups:   repositories.CountGroups(),
		TotalTeams:    CountTeams(),
		TotalReminders: 0,
	}
}

type ScoreInput struct {
	HomeScore int `json:"home_score" binding:"required"`
	AwayScore int `json:"away_score" binding:"required"`
}

func UpdateMatchScore(matchID uint, input ScoreInput) (*models.Match, error) {
	match, err := repositories.GetMatchByID(matchID)
	if err != nil {
		return nil, fmt.Errorf("match not found")
	}
	match.HomeScore = &input.HomeScore
	match.AwayScore = &input.AwayScore
	if input.HomeScore > input.AwayScore {
		match.WinnerTeamID = &match.HomeTeamID
	} else if input.AwayScore > input.HomeScore {
		match.WinnerTeamID = &match.AwayTeamID
	}
	if err := repositories.UpdateMatch(match); err != nil {
		return nil, err
	}
	return match, nil
}

type StatusInput struct {
	Status string `json:"status" binding:"required"`
}

func UpdateMatchStatus(matchID uint, input StatusInput) (*models.Match, error) {
	match, err := repositories.GetMatchByID(matchID)
	if err != nil {
		return nil, fmt.Errorf("match not found")
	}
	match.Status = input.Status
	if err := repositories.UpdateMatch(match); err != nil {
		return nil, err
	}

	if match.Status == "finished" && match.Stage == "group" && match.GroupID != nil {
		_ = RecalculateGroupStanding(*match.GroupID)
	}

	return match, nil
}

type MatchInput struct {
	MatchNo         int    `json:"match_no"`
	HomeTeamID      uint   `json:"home_team_id"`
	AwayTeamID      uint   `json:"away_team_id"`
	GroupID         *uint  `json:"group_id"`
	Stage           string `json:"stage"`
	StadiumID       uint   `json:"stadium_id"`
	CityID          uint   `json:"city_id"`
	KickoffTimeUTC  string `json:"kickoff_time_utc"`
	ImportanceLevel int    `json:"importance_level"`
	RecommendTag    string `json:"recommend_tag"`
}

func CreateMatch(input MatchInput) (*models.Match, error) {
	match := &models.Match{
		MatchNo:         input.MatchNo,
		HomeTeamID:      input.HomeTeamID,
		AwayTeamID:      input.AwayTeamID,
		GroupID:         input.GroupID,
		Stage:           input.Stage,
		StadiumID:       input.StadiumID,
		CityID:          input.CityID,
		ImportanceLevel: input.ImportanceLevel,
		RecommendTag:    input.RecommendTag,
		Status:          "scheduled",
	}
	if err := repositories.CreateMatch(match); err != nil {
		return nil, err
	}
	return match, nil
}

func ImportMatchesCSV(reader io.Reader) (int, error) {
	r := csv.NewReader(reader)
	records, err := r.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("failed to parse CSV: %v", err)
	}

	imported := 0
	for i, row := range records {
		if i == 0 {
			continue
		}
		if len(row) < 8 {
			continue
		}

		homeTeam, err := getTeamByName(row[1])
		if err != nil {
			continue
		}
		awayTeam, err := getTeamByName(row[2])
		if err != nil {
			continue
		}

		match := &models.Match{
			MatchNo:         mustAtoi(row[0]),
			HomeTeamID:      homeTeam.ID,
			AwayTeamID:      awayTeam.ID,
			Stage:           row[4],
			Status:          "scheduled",
			ImportanceLevel: mustAtoi(row[8]),
			RecommendTag:    row[9],
		}

		if row[3] != "" {
			group, err := repositories.GetGroupByName(row[3])
			if err == nil {
				match.GroupID = &group.ID
			}
		}
		if row[5] != "" {
			stadium, err := repositories.GetStadiumByName(row[5])
			if err == nil {
				match.StadiumID = stadium.ID
				match.CityID = stadium.CityID
			}
		}

		if err := repositories.CreateMatch(match); err != nil {
			continue
		}
		imported++
	}
	return imported, nil
}

func getTeamByName(name string) (*models.Team, error) {
	teams, _, err := repositories.ListTeams("", name, 0, 1, 1)
	if err != nil || len(teams) == 0 {
		return nil, fmt.Errorf("team not found: %s", name)
	}
	return &teams[0], nil
}

func mustAtoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
