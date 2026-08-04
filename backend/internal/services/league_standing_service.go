package services

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"worldcup-mate/internal/models"
	"worldcup-mate/internal/providers/footballdata"
	"worldcup-mate/internal/repositories"
)

// League standings live in the new league_standings table, completely
// independent from GroupStanding / best-third logic (World Cup, untouched).

// Zone rules per league code (approximation for the 2025-26 season;
// adjustable — order matters: first matching range wins).
type zoneRange struct {
	from, to int
	zone     string
}

var leagueZoneRules = map[string][]zoneRange{
	"PL":  {{1, 4, "champions_league"}, {5, 6, "europa_league"}, {18, 20, "relegation"}},
	"PD":  {{1, 4, "champions_league"}, {5, 6, "europa_league"}, {18, 20, "relegation"}},
	"SA":  {{1, 4, "champions_league"}, {5, 6, "europa_league"}, {18, 20, "relegation"}},
	"BL1": {{1, 4, "champions_league"}, {5, 6, "europa_league"}, {16, 18, "relegation"}},
	"FL1": {{1, 3, "champions_league"}, {4, 6, "europa_league"}, {16, 18, "relegation"}},
}

func zoneForPosition(code string, position int) string {
	rules := leagueZoneRules[strings.ToUpper(code)]
	for _, r := range rules {
		if position >= r.from && position <= r.to {
			return r.zone
		}
	}
	return ""
}

// SyncLeagueStandings fetches the official standings (TOTAL / HOME / AWAY)
// from football-data and upserts them into league_standings.
func SyncLeagueStandings(ctx context.Context, code string, season int) error {
	cfg := GetLeagueSyncConfig()
	if !IsLeagueSyncEnabled() {
		return errors.New("league sync is disabled")
	}
	competition, err := repositories.GetCompetitionByCode(code)
	if err != nil {
		return err
	}
	if season <= 0 {
		season = competition.Season
	}

	client := footballdata.NewClient(cfg.BaseURL, cfg.APIKey)
	data, err := client.CompetitionStandings(ctx, code, season)
	if err != nil {
		return err
	}

	for _, group := range data.Standings {
		standingType := strings.ToLower(group.Type)
		if standingType == "" {
			standingType = "total"
		}
		for _, row := range group.Table {
			team, err := findOrCreateClub(row.Team)
			if err != nil {
				continue
			}
			standing := &models.LeagueStanding{
				CompetitionID:  competition.ID,
				Season:         season,
				TeamID:         team.ID,
				Type:           standingType,
				Position:       row.Position,
				Played:         row.PlayedGames,
				Won:            row.Won,
				Drawn:          row.Draw,
				Lost:           row.Lost,
				GoalsFor:       row.GoalsFor,
				GoalsAgainst:   row.GoalsAgainst,
				GoalDifference: row.GoalDifference,
				Points:         row.Points,
				UpdatedAt:      time.Now().UTC(),
			}
			if standingType == "total" {
				standing.Zone = zoneForPosition(code, row.Position)
			}
			if err := repositories.UpsertLeagueStanding(standing); err != nil {
				continue
			}
		}
	}
	return nil
}

// RecalculateLeagueStanding rebuilds the TOTAL league table from finished
// matches (points → goal difference → goals for). Used by the admin endpoint
// as a fallback when the official standings are unavailable.
func RecalculateLeagueStanding(competitionID uint, season int) error {
	// Rebuild from scratch: drop stale rows first so teams without any
	// finished matches don't keep old positions/zones.
	if err := repositories.DeleteLeagueStandings(competitionID, season, "total"); err != nil {
		return err
	}

	matches, _, err := repositories.ListMatches(repositories.MatchFilter{
		CompetitionID: competitionID,
		Season:        season,
		Status:        "finished",
		Page:          1,
		PageSize:      10000,
	})
	if err != nil {
		return err
	}

	stats := make(map[uint]*models.LeagueStanding)
	for _, m := range matches {
		if m.HomeScore == nil || m.AwayScore == nil {
			continue
		}
		home, ok := stats[m.HomeTeamID]
		if !ok {
			home = &models.LeagueStanding{CompetitionID: competitionID, Season: season, TeamID: m.HomeTeamID, Type: "total"}
			stats[m.HomeTeamID] = home
		}
		away, ok := stats[m.AwayTeamID]
		if !ok {
			away = &models.LeagueStanding{CompetitionID: competitionID, Season: season, TeamID: m.AwayTeamID, Type: "total"}
			stats[m.AwayTeamID] = away
		}
		home.Played++
		away.Played++
		home.GoalsFor += *m.HomeScore
		home.GoalsAgainst += *m.AwayScore
		away.GoalsFor += *m.AwayScore
		away.GoalsAgainst += *m.HomeScore
		if *m.HomeScore > *m.AwayScore {
			home.Won++
			home.Points += 3
			away.Lost++
		} else if *m.HomeScore < *m.AwayScore {
			away.Won++
			away.Points += 3
			home.Lost++
		} else {
			home.Drawn++
			away.Drawn++
			home.Points++
			away.Points++
		}
	}

	standings := make([]*models.LeagueStanding, 0, len(stats))
	code := ""
	if competition, err := repositories.GetCompetitionByID(competitionID); err == nil {
		code = competition.Code
	}
	for _, s := range stats {
		s.GoalDifference = s.GoalsFor - s.GoalsAgainst
		standings = append(standings, s)
	}
	sort.Slice(standings, func(i, j int) bool {
		if standings[i].Points != standings[j].Points {
			return standings[i].Points > standings[j].Points
		}
		if standings[i].GoalDifference != standings[j].GoalDifference {
			return standings[i].GoalDifference > standings[j].GoalDifference
		}
		if standings[i].GoalsFor != standings[j].GoalsFor {
			return standings[i].GoalsFor > standings[j].GoalsFor
		}
		return standings[i].TeamID < standings[j].TeamID
	})

	for i, s := range standings {
		s.Position = i + 1
		s.Zone = zoneForPosition(code, s.Position)
		s.UpdatedAt = time.Now().UTC()
		if err := repositories.UpsertLeagueStanding(s); err != nil {
			return err
		}
	}
	return nil
}
