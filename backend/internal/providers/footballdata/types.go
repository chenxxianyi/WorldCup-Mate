package footballdata

import "time"

type MatchesResponse struct {
	Count   int     `json:"count"`
	Matches []Match `json:"matches"`
}

type Match struct {
	ID       int64     `json:"id"`
	UTCDate  time.Time `json:"utcDate"`
	Status   string    `json:"status"`
	Matchday int       `json:"matchday"`
	Stage    string    `json:"stage"`
	Group    string    `json:"group"`
	Venue    string    `json:"venue"`
	HomeTeam Team      `json:"homeTeam"`
	AwayTeam Team      `json:"awayTeam"`
	Score    Score     `json:"score"`
}

type Team struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	TLA       string `json:"tla"`
	Crest     string `json:"crest"`
}

type Score struct {
	Winner      string    `json:"winner"`
	Duration    string    `json:"duration"`
	FullTime    ScoreLine `json:"fullTime"`
	RegularTime ScoreLine `json:"regularTime"`
	HalfTime    ScoreLine `json:"halfTime"`
}

type ScoreLine struct {
	Home *int `json:"home"`
	Away *int `json:"away"`
}

// StandingsResponse is the response of GET /competitions/{code}/standings.
type StandingsResponse struct {
	Standings []StandingGroup `json:"standings"`
}

type StandingGroup struct {
	Stage string        `json:"stage"`
	Type  string        `json:"type"` // TOTAL | HOME | AWAY
	Table []StandingRow `json:"table"`
}

type StandingRow struct {
	Position       int  `json:"position"`
	Team           Team `json:"team"`
	PlayedGames    int  `json:"playedGames"`
	Won            int  `json:"won"`
	Draw           int  `json:"draw"`
	Lost           int  `json:"lost"`
	Points         int  `json:"points"`
	GoalsFor       int  `json:"goalsFor"`
	GoalsAgainst   int  `json:"goalsAgainst"`
	GoalDifference int  `json:"goalDifference"`
}
