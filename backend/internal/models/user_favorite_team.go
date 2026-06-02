package models

import "time"

type UserFavoriteTeam struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex:idx_user_team;not null" json:"user_id"`
	TeamID    uint      `gorm:"uniqueIndex:idx_user_team;not null" json:"team_id"`
	Team      Team      `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
