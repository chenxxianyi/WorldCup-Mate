package services

import (
	"time"

	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"
)

type MatchQuery struct {
	Date          string `form:"date"`
	TeamID        uint   `form:"teamId"`
	GroupID       uint   `form:"groupId"`
	GroupName     string `form:"groupName"`
	Stage         string `form:"stage"`
	CityID        uint   `form:"cityId"`
	Status        string `form:"status"`
	Keyword       string `form:"keyword"`
	CompetitionID uint   `form:"competitionId"` // optional: 0 = legacy behavior
	WorldCupOnly  bool   `form:"worldCup"`      // World Cup records use a NULL competition_id for backwards compatibility
	Season        int    `form:"season"`
	Matchday      int    `form:"matchday"`
	Page          int
	PageSize      int
}

func ListMatches(q MatchQuery) ([]models.Match, int64, error) {
	return repositories.ListMatches(repositories.MatchFilter{
		Date:          q.Date,
		TeamID:        q.TeamID,
		GroupID:       q.GroupID,
		GroupName:     q.GroupName,
		Stage:         q.Stage,
		CityID:        q.CityID,
		Status:        q.Status,
		Keyword:       q.Keyword,
		CompetitionID: q.CompetitionID,
		WorldCupOnly:  q.WorldCupOnly,
		Season:        q.Season,
		Matchday:      q.Matchday,
		Page:          q.Page,
		PageSize:      q.PageSize,
	})
}

func GetMatchByID(id uint) (*models.Match, error) {
	return repositories.GetMatchByID(id)
}

// dayRangeInTZ returns the [start, end) UTC range of the calendar day that
// contains `now` in the given IANA timezone (UTC when tz is empty/invalid).
// DST transitions produce 23/25-hour days — handled naturally by local
// midnight math (DATA-05D).
func dayRangeInTZ(now time.Time, tz string) (start, end time.Time) {
	loc := time.UTC
	if tz != "" {
		if l, err := time.LoadLocation(tz); err == nil {
			loc = l
		}
	}
	local := now.In(loc)
	start = time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc).UTC()
	// Next local midnight (handles DST 23/25-hour days and month ends via
	// time.Date normalization; AddDate on the UTC value would add 24h).
	end = time.Date(local.Year(), local.Month(), local.Day()+1, 0, 0, 0, 0, loc).UTC()
	return start, end
}

func GetTodayMatches(worldCup bool, tz string) ([]models.Match, error) {
	start, end := dayRangeInTZ(time.Now().UTC(), tz)
	return repositories.GetMatchesByDateRange(start, end, worldCup)
}

func GetTomorrowMatches(worldCup bool, tz string) ([]models.Match, error) {
	_, end := dayRangeInTZ(time.Now().UTC(), tz)
	// "Tomorrow" is the calendar day containing `end` — NOT end+24h, which
	// drifts by an hour across DST transitions (DATA-05D).
	_, nextEnd := dayRangeInTZ(end, tz)
	return repositories.GetMatchesByDateRange(end, nextEnd, worldCup)
}

func GetUpcomingMatches(worldCup bool) ([]models.Match, error) {
	var matches []models.Match
	all, _, err := repositories.ListMatches(repositories.MatchFilter{
		Status:       "scheduled",
		WorldCupOnly: worldCup,
		Page:         1,
		PageSize:     10,
	})
	if err != nil {
		return nil, err
	}
	matches = all
	return matches, nil
}

func GetLiveMatches(worldCup bool) ([]models.Match, error) {
	matches, _, err := repositories.ListMatches(repositories.MatchFilter{
		Status:       "live",
		WorldCupOnly: worldCup,
		Page:         1,
		PageSize:     100,
	})
	return matches, err
}

func GetRecommendedMatches(worldCup bool) ([]models.Match, error) {
	return repositories.GetRecommendedMatches(worldCup)
}

func GetMatchesByTeamID(teamID uint) ([]models.Match, error) {
	return repositories.GetMatchesByTeamID(teamID)
}

func GetMatchesByGroupID(groupID uint) ([]models.Match, error) {
	matches, _, err := repositories.ListMatches(repositories.MatchFilter{
		GroupID:  groupID,
		Page:     1,
		PageSize: 100,
	})
	return matches, err
}

func GetMatchesByStage(stage string) ([]models.Match, error) {
	matches, _, err := repositories.ListMatches(repositories.MatchFilter{
		Stage:    stage,
		Page:     1,
		PageSize: 100,
	})
	return matches, err
}

type TournamentProgress struct {
	StageName    string          `json:"stage_name"`
	TotalMatches int64           `json:"total_matches"`
	Completed    int64           `json:"completed"`
	Live         int64           `json:"live"`
	Scheduled    int64           `json:"scheduled"`
	Progress     float64         `json:"progress"`
	StageDetails []StageProgress `json:"stage_details"`
}

type StageProgress struct {
	Stage     string  `json:"stage"`
	StageName string  `json:"stage_name"`
	Total     int64   `json:"total"`
	Completed int64   `json:"completed"`
	Live      int64   `json:"live"`
	Scheduled int64   `json:"scheduled"`
	Progress  float64 `json:"progress"`
}

func GetTournamentProgress() (*TournamentProgress, error) {
	data, err := repositories.CountMatchesByStageAndStatus()
	if err != nil {
		return nil, err
	}

	stageNames := map[string]string{
		"group":         "小组赛",
		"round_of_32":   "1/8决赛",
		"round_of_16":   "1/8决赛",
		"quarter_final": "1/4决赛",
		"semi_final":    "半决赛",
		"third_place":   "三四名决赛",
		"final":         "决赛",
	}

	stageOrder := []string{"group", "round_of_16", "round_of_32", "quarter_final", "semi_final", "third_place", "final"}

	progress := &TournamentProgress{
		StageName: "小组赛阶段",
	}

	var totalAll, completedAll int64

	for _, stage := range stageOrder {
		statuses, ok := data[stage]
		if !ok {
			continue
		}

		var total, completed, live, scheduled int64
		for status, count := range statuses {
			total += count
			switch status {
			case "finished":
				completed = count
			case "live":
				live = count
			case "scheduled":
				scheduled = count
			}
		}

		totalAll += total
		completedAll += completed

		var stageProgress float64
		if total > 0 {
			stageProgress = float64(completed) / float64(total) * 100
		}

		stageName := stageNames[stage]
		if stageName == "" {
			stageName = stage
		}

		progress.StageDetails = append(progress.StageDetails, StageProgress{
			Stage:     stage,
			StageName: stageName,
			Total:     total,
			Completed: completed,
			Live:      live,
			Scheduled: scheduled,
			Progress:  stageProgress,
		})
	}

	progress.TotalMatches = totalAll
	progress.Completed = completedAll
	progress.Live = 0
	progress.Scheduled = totalAll - completedAll
	for _, s := range progress.StageDetails {
		progress.Live += s.Live
	}

	if totalAll > 0 {
		progress.Progress = float64(completedAll) / float64(totalAll) * 100
	}

	if len(progress.StageDetails) > 0 {
		for _, s := range progress.StageDetails {
			if s.Total > 0 && s.Completed < s.Total {
				progress.StageName = s.StageName + "阶段"
				break
			}
		}
	}

	return progress, nil
}
