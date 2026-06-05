package footballdata

import "time"

type MatchesResponse struct {
	Count   int     `json:"count"`
	Matches []Match `json:"matches"`
}

type TeamsResponse struct {
	Count int    `json:"count"`
	Teams []Team `json:"teams"`
}

type Match struct {
	ID       int64     `json:"id"`
	UTCDate  time.Time `json:"utcDate"`
	Status   string    `json:"status"`
	Matchday int       `json:"matchday"`
	Stage    string    `json:"stage"`
	Group    string    `json:"group"`
	HomeTeam Team      `json:"homeTeam"`
	AwayTeam Team      `json:"awayTeam"`
	Score    Score     `json:"score"`
}

type Team struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	ShortName string        `json:"shortName"`
	TLA       string        `json:"tla"`
	Crest     string        `json:"crest"`
	Squad     []SquadPlayer `json:"squad"`
}

type SquadPlayer struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Position    string `json:"position"`
	DateOfBirth string `json:"dateOfBirth"`
	Nationality string `json:"nationality"`
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

type MatchDetailResponse struct {
	ID       int64           `json:"id"`
	Status   string          `json:"status"`
	HomeTeam MatchLineupTeam `json:"homeTeam"`
	AwayTeam MatchLineupTeam `json:"awayTeam"`
}

type MatchLineupTeam struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	TLA       string         `json:"tla"`
	Formation string         `json:"formation"`
	Lineup    []LineupPlayer `json:"lineup"`
	Bench     []LineupPlayer `json:"bench"`
}

type LineupPlayer struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Position    string `json:"position"`
	ShirtNumber int    `json:"shirtNumber"`
}
