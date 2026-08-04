package models

import (
	"time"

	"gorm.io/gorm"
)

type Match struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	MatchNo          int            `gorm:"not null" json:"match_no"`
	HomeTeamID       uint           `gorm:"index" json:"home_team_id"`
	AwayTeamID       uint           `gorm:"index" json:"away_team_id"`
	HomeTeam         Team           `gorm:"foreignKey:HomeTeamID" json:"home_team,omitempty"`
	AwayTeam         Team           `gorm:"foreignKey:AwayTeamID" json:"away_team,omitempty"`
	GroupID          *uint          `gorm:"index" json:"group_id"`
	Group            *Group         `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Stage            string         `gorm:"size:30;default:group;index" json:"stage"`
	CompetitionID    *uint          `gorm:"index" json:"competition_id"` // nullable: NULL = World Cup / legacy data
	Season           *int           `gorm:"index" json:"season"`         // starting year, e.g. 2025 for 2025-26 season
	Matchday         *int           `gorm:"index" json:"matchday"`       // league round, league competitions only
	StadiumID        uint           `gorm:"index" json:"stadium_id"`
	Stadium          Stadium        `gorm:"foreignKey:StadiumID" json:"stadium,omitempty"`
	CityID           uint           `gorm:"index" json:"city_id"`
	City             City           `gorm:"foreignKey:CityID" json:"city,omitempty"`
	KickoffTimeUTC   time.Time      `gorm:"not null;index" json:"kickoff_time_utc"`
	HomeScore        *int           `json:"home_score"`
	AwayScore        *int           `json:"away_score"`
	Status           string         `gorm:"size:20;default:scheduled;index" json:"status"`
	StatusDetail     string         `gorm:"size:30" json:"status_detail"`
	LiveMinute       *int           `json:"live_minute"`
	WinnerTeamID     *uint          `json:"winner_team_id"`
	ExternalProvider string         `gorm:"size:30;index" json:"external_provider"`
	ExternalID       string         `gorm:"size:64;index" json:"external_id"`
	LastSyncedAt     *time.Time     `json:"last_synced_at"`
	ImportanceLevel  int            `gorm:"default:0" json:"importance_level"`
	RecommendTag     string         `gorm:"size:50" json:"recommend_tag"`
	RecommendReason  string         `gorm:"size:255" json:"recommend_reason"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}
