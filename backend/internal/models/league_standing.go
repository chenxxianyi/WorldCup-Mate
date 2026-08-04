package models

import "time"

// LeagueStanding stores league standings (TOTAL / HOME / AWAY) for a
// competition season. It is a brand-new table, independent from
// GroupStanding which keeps serving the World Cup group tables untouched.
type LeagueStanding struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CompetitionID  uint      `gorm:"uniqueIndex:idx_league_standing;not null" json:"competition_id"`
	Season         int       `gorm:"uniqueIndex:idx_league_standing;not null" json:"season"`
	TeamID         uint      `gorm:"uniqueIndex:idx_league_standing;not null" json:"team_id"`
	Team           Team      `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Type           string    `gorm:"size:10;uniqueIndex:idx_league_standing;not null;default:total" json:"type"` // total | home | away
	Position       int       `json:"position"`
	Played         int       `json:"played"`
	Won            int       `json:"won"`
	Drawn          int       `json:"drawn"`
	Lost           int       `json:"lost"`
	GoalsFor       int       `json:"goals_for"`
	GoalsAgainst   int       `json:"goals_against"`
	GoalDifference int       `json:"goal_difference"`
	Points         int       `json:"points"`
	Zone           string    `gorm:"size:30" json:"zone"` // champions_league | europa_league | relegation ...
	UpdatedAt      time.Time `json:"updated_at"`
}
