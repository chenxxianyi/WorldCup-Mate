package services

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"
	"worldcup-mate/internal/utils"

	"gorm.io/gorm"
)

type DashboardData struct {
	TotalMatches      int64 `json:"total_matches"`
	TotalGroups       int64 `json:"total_groups"`
	TotalTeams        int64 `json:"total_teams"`
	TotalUsers        int64 `json:"total_users"`
	TotalCompetitions int64 `json:"total_competitions"`
	TotalReminders    int64 `json:"total_reminders"`
}

func GetDashboard() DashboardData {
	return DashboardData{
		TotalMatches:      repositories.CountMatches(),
		TotalGroups:       repositories.CountGroups(),
		TotalTeams:        CountTeams(),
		TotalUsers:        repositories.CountUsers(),
		TotalCompetitions: repositories.CountCompetitions(),
		TotalReminders:    repositories.CountAllReminders(),
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
	} else {
		match.WinnerTeamID = nil
	}
	if err := repositories.UpdateMatch(match); err != nil {
		return nil, err
	}
	// Score correction on a finished group match must keep standings in sync.
	if match.Status == "finished" && match.Stage == "group" && match.GroupID != nil {
		if err := RecalculateGroupStanding(*match.GroupID); err != nil {
			return nil, err
		}
	}
	return match, nil
}

type StatusInput struct {
	Status string `json:"status" binding:"required"`
}

// matchStatusTransitions defines the legal state machine (ADM-05):
//   scheduled -> live | finished | postponed | cancelled
//   live      -> finished | postponed
//   postponed -> scheduled | live
//   cancelled / finished -> (terminal)
var matchStatusTransitions = map[string][]string{
	"scheduled": {"live", "finished", "postponed", "cancelled"},
	"live":      {"finished", "postponed"},
	"postponed": {"scheduled", "live"},
	"cancelled": {},
	"finished":  {},
}

func canTransition(from, to string) bool {
	allowed, ok := matchStatusTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

func UpdateMatchStatus(matchID uint, input StatusInput) (*models.Match, error) {
	match, err := repositories.GetMatchByID(matchID)
	if err != nil {
		return nil, err // includes gorm.ErrRecordNotFound
	}
	if !canTransition(match.Status, input.Status) {
		return nil, utils.ErrInvalidStatusTransition
	}
	if input.Status == "finished" && (match.HomeScore == nil || match.AwayScore == nil) {
		return nil, errors.New("cannot finish a match without a score")
	}

	// Persist status atomically; standings recalculation runs afterwards and
	// its failure is logged (keeps the status change from being rolled back
	// by a non-critical recalculation issue).
	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(match).Update("status", input.Status).Error
	}); err != nil {
		return nil, err
	}
	match.Status = input.Status

	if match.Status == "finished" && match.Stage == "group" && match.GroupID != nil {
		if err := RecalculateGroupStanding(*match.GroupID); err != nil {
			return nil, err
		}
		if match.Stage == "group" {
			if err := RecalculateBestThird(); err != nil {
				return nil, err
			}
		}
	}
	return match, nil
}

// UpdateMatch persists an already-mutated match (admin edit).
func UpdateMatch(match *models.Match) error {
	return repositories.UpdateMatch(match)
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
	CompetitionID   *uint  `json:"competition_id"`
	Season          *int   `json:"season"`
	Matchday        *int   `json:"matchday"`
}

func CreateMatch(input MatchInput) (*models.Match, error) {
	var kickoff time.Time
	if input.KickoffTimeUTC != "" {
		parsed, err := time.Parse(time.RFC3339, input.KickoffTimeUTC)
		if err != nil {
			return nil, utils.ErrInvalidTime
		}
		kickoff = parsed.UTC()
	}
	match := &models.Match{
		MatchNo:         input.MatchNo,
		HomeTeamID:      input.HomeTeamID,
		AwayTeamID:      input.AwayTeamID,
		GroupID:         input.GroupID,
		Stage:           input.Stage,
		StadiumID:       input.StadiumID,
		CityID:          input.CityID,
		KickoffTimeUTC:  kickoff,
		ImportanceLevel: input.ImportanceLevel,
		RecommendTag:    input.RecommendTag,
		CompetitionID:   input.CompetitionID,
		Season:          input.Season,
		Matchday:        input.Matchday,
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
