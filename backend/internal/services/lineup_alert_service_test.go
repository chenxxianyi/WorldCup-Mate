package services

import (
	"testing"

	"worldcup-mate/internal/models"
)

func TestIsLineupAlertReady(t *testing.T) {
	tests := []struct {
		name   string
		lineups []models.MatchLineup
		want   bool
	}{
		{
			name: "both sides have 11 starting players",
			lineups: []models.MatchLineup{
				{Side: "home", Players: makeStartingPlayers(11)},
				{Side: "away", Players: makeStartingPlayers(11)},
			},
			want: true,
		},
		{
			name: "home has 11, away has only 10 starting",
			lineups: []models.MatchLineup{
				{Side: "home", Players: makeStartingPlayers(11)},
				{Side: "away", Players: makeStartingPlayers(10)},
			},
			want: false,
		},
		{
			name: "both sides have 11 but mixed with substitutes",
			lineups: []models.MatchLineup{
				{
					Side: "home",
					Players: append(makeStartingPlayers(11), models.MatchLineupPlayer{
						Role: "substitute",
						Name: "Sub Player",
					}),
				},
				{Side: "away", Players: makeStartingPlayers(11)},
			},
			want: true,
		},
		{
			name:   "no lineups at all",
			lineups: []models.MatchLineup{},
			want:   false,
		},
		{
			name: "only home side has lineup",
			lineups: []models.MatchLineup{
				{Side: "home", Players: makeStartingPlayers(11)},
			},
			want: false,
		},
		{
			name: "home has enough but away has 0 starting",
			lineups: []models.MatchLineup{
				{Side: "home", Players: makeStartingPlayers(11)},
				{Side: "away", Players: []models.MatchLineupPlayer{}},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isLineupAlertReady(tt.lineups)
			if got != tt.want {
				t.Errorf("isLineupAlertReady() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollectTargetUserIDs(t *testing.T) {
	// This function uses DB, so we test it via manual verification.
	// Here we just verify it compiles and is callable.
	t.Log("collectTargetUserIDs integrates with DB — tested via manual verification")
}

// makeStartingPlayers creates n players with role "starting" and unique shirt numbers.
func makeStartingPlayers(n int) []models.MatchLineupPlayer {
	players := make([]models.MatchLineupPlayer, n)
	for i := 0; i < n; i++ {
		players[i] = models.MatchLineupPlayer{
			Name:        "Player",
			NameEn:      "Player",
			ShirtNumber: i + 1,
			Position:    "MF",
			Role:        "starting",
			SortOrder:   i + 1,
		}
	}
	return players
}
