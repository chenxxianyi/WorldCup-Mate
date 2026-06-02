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
