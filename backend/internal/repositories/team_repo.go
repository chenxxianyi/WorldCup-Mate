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
	err := q.Preload("Group").Offset((page-1)*pageSize).Limit(pageSize).Find(&teams).Error
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
