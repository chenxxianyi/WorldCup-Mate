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
	HomePossession    *int           `json:"home_possession"`
	AwayPossession    *int           `json:"away_possession"`
	HomeShots         *int           `json:"home_shots"`
	AwayShots         *int           `json:"away_shots"`
	HomeShotsOnTarget *int           `json:"home_shots_on_target"`
	AwayShotsOnTarget *int           `json:"away_shots_on_target"`
	HomeCorners       *int           `json:"home_corners"`
	AwayCorners       *int           `json:"away_corners"`
	HomeOffsides      *int           `json:"home_offsides"`
	AwayOffsides      *int           `json:"away_offsides"`
	HomeYellowCards   *int           `json:"home_yellow_cards"`
	AwayYellowCards   *int           `json:"away_yellow_cards"`
	HomeRedCards      *int           `json:"home_red_cards"`
	AwayRedCards      *int           `json:"away_red_cards"`
	HomeFouls         *int           `json:"home_fouls"`
	AwayFouls         *int           `json:"away_fouls"`
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
