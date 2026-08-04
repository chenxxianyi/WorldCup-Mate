package repositories

import (
	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
)

func ListTeams(continent, keyword string, groupID uint, page, pageSize int) ([]models.Team, int64, error) {
	var teams []models.Team
	var total int64
	q := database.DB.Model(&models.Team{})
	if continent != "" {
		q = q.Where("continent = ?", continent)
	}
	if groupID > 0 {
		q = q.Where("group_id = ?", groupID)
	}
	if keyword != "" {
		q = q.Where("name LIKE ? OR name_en LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	q.Count(&total)
	err := q.Preload("Group").Offset((page - 1) * pageSize).Limit(pageSize).Find(&teams).Error
	return teams, total, err
}

// TeamFilter carries optional filters for multi-competition team queries.
// Zero values mean "no filter", preserving the legacy behavior.
// Note: teams have no competition_id column; filter clubs via TeamType/Country.
type TeamFilter struct {
	Continent string
	Keyword   string
	GroupID   uint
	TeamType  string
	Country   string
	Page      int
	PageSize  int
}

func ListTeamsFiltered(f TeamFilter) ([]models.Team, int64, error) {
	var teams []models.Team
	var total int64
	q := database.DB.Model(&models.Team{})
	if f.Continent != "" {
		q = q.Where("continent = ?", f.Continent)
	}
	if f.GroupID > 0 {
		q = q.Where("group_id = ?", f.GroupID)
	}
	if f.TeamType != "" {
		q = q.Where("team_type = ?", f.TeamType)
	}
	if f.Country != "" {
		q = q.Where("country = ?", f.Country)
	}
	if f.Keyword != "" {
		q = q.Where("name LIKE ? OR name_en LIKE ?", "%"+f.Keyword+"%", "%"+f.Keyword+"%")
	}
	q.Count(&total)
	err := q.Preload("Group").Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).Find(&teams).Error
	return teams, total, err
}

func GetTeamByID(id uint) (*models.Team, error) {
	var team models.Team
	err := database.DB.Preload("Group").First(&team, id).Error
	return &team, err
}

func CreateTeam(team *models.Team) error {
	return database.DB.Create(team).Error
}

func UpdateTeam(team *models.Team) error {
	return database.DB.Save(team).Error
}

func DeleteTeam(id uint) error {
	return database.DB.Delete(&models.Team{}, id).Error
}

func GetTeamsByGroupID(groupID uint) ([]models.Team, error) {
	var teams []models.Team
	err := database.DB.Where("group_id = ?", groupID).Order("id ASC").Find(&teams).Error
	return teams, err
}
