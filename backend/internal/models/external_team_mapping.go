package models

import "time"

type ExternalTeamMapping struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	TeamID           uint      `gorm:"uniqueIndex:idx_provider_team;not null" json:"team_id"`
	Team             Team      `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Provider         string    `gorm:"uniqueIndex:idx_provider_team;size:50;not null" json:"provider"`
	ExternalTeamID   string    `gorm:"size:100;index;not null" json:"external_team_id"`
	ExternalTeamName string    `gorm:"size:100" json:"external_team_name"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
