package models

import "time"

type GroupStanding struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	GroupID            uint      `gorm:"uniqueIndex:idx_group_team;not null" json:"group_id"`
	Group              Group     `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	TeamID             uint      `gorm:"uniqueIndex:idx_group_team;not null" json:"team_id"`
	Team               Team      `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Played             int       `gorm:"default:0" json:"played"`
	Won                int       `gorm:"default:0" json:"won"`
	Drawn              int       `gorm:"default:0" json:"drawn"`
	Lost               int       `gorm:"default:0" json:"lost"`
	GoalsFor           int       `gorm:"default:0" json:"goals_for"`
	GoalsAgainst       int       `gorm:"default:0" json:"goals_against"`
	GoalDifference     int       `gorm:"default:0" json:"goal_difference"`
	Points             int       `gorm:"default:0" json:"points"`
	Rank               int       `gorm:"default:0" json:"rank"`
	QualificationStatus string   `gorm:"size:20;default:unknown" json:"qualification_status"`
	UpdatedAt          time.Time `json:"updated_at"`
}
