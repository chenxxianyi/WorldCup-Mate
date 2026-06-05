package main

import (
	"log"
	"time"

	"worldcup-mate/internal/config"
	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
)

type demoPlayer struct {
	name     string
	nameEn   string
	number   int
	position string
	label    string
	grid     string
}

func main() {
	cfg := config.Load()
	database.InitMySQL(cfg.MySQLDSN)

	if err := database.DB.AutoMigrate(
		&models.MatchLineup{},
		&models.MatchLineupPlayer{},
		&models.ExternalMatchMapping{},
	); err != nil {
		log.Fatalf("migrate lineup tables failed: %v", err)
	}

	var match models.Match
	if err := database.DB.Preload("HomeTeam").Preload("AwayTeam").Order("kickoff_time_utc ASC, id ASC").First(&match).Error; err != nil {
		log.Fatalf("load demo match failed: %v", err)
	}

	now := time.Now().UTC()
	seedTeamLineup(match.ID, match.HomeTeamID, "home", "4-3-3", "哈维尔·阿吉雷", "demo", &now, []demoPlayer{
		{"吉列尔莫·奥乔亚", "G. Ochoa", 13, "GK", "门将", "1:1"},
		{"豪尔赫·桑切斯", "J. Sanchez", 2, "DF", "后卫", "2:4"},
		{"塞萨尔·蒙特斯", "C. Montes", 3, "DF", "后卫", "2:3"},
		{"约翰·巴斯克斯", "J. Vasquez", 5, "DF", "后卫", "2:2"},
		{"赫苏斯·加利亚多", "J. Gallardo", 23, "DF", "后卫", "2:1"},
		{"埃德森·阿尔瓦雷斯", "E. Alvarez", 4, "MF", "中场", "3:2"},
		{"路易斯·查韦斯", "L. Chavez", 18, "MF", "中场", "3:1"},
		{"路易斯·罗莫", "L. Romo", 7, "MF", "中场", "3:3"},
		{"亚历克西斯·维加", "A. Vega", 10, "FW", "前锋", "4:1"},
		{"圣地亚哥·希门尼斯", "S. Gimenez", 11, "FW", "前锋", "4:2"},
		{"罗伯托·阿尔瓦拉多", "R. Alvarado", 25, "FW", "前锋", "4:3"},
	}, []demoPlayer{
		{"劳尔·希门尼斯", "R. Jimenez", 9, "FW", "前锋", ""},
		{"奥贝林·皮内达", "O. Pineda", 17, "MF", "中场", ""},
		{"伊斯雷尔·雷耶斯", "I. Reyes", 15, "DF", "后卫", ""},
	})

	seedTeamLineup(match.ID, match.AwayTeamID, "away", "3-4-2-1", "杰西·马什", "demo", &now, []demoPlayer{
		{"马克西姆·克雷波", "M. Crepeau", 16, "GK", "门将", "1:1"},
		{"阿利斯泰尔·约翰斯顿", "A. Johnston", 2, "DF", "后卫", "2:3"},
		{"德里克·科内柳斯", "D. Cornelius", 13, "DF", "后卫", "2:2"},
		{"莫伊塞·邦比托", "M. Bombito", 15, "DF", "后卫", "2:1"},
		{"阿方索·戴维斯", "A. Davies", 19, "MF", "中场", "3:1"},
		{"斯蒂芬·尤斯塔基奥", "S. Eustaquio", 7, "MF", "中场", "3:2"},
		{"伊斯梅尔·科内", "I. Kone", 8, "MF", "中场", "3:3"},
		{"泰琼·布坎南", "T. Buchanan", 11, "MF", "中场", "3:4"},
		{"乔纳森·戴维", "J. David", 10, "FW", "前锋", "4:2"},
		{"赛尔·拉林", "C. Larin", 17, "FW", "前锋", "4:1"},
		{"雅各布·沙费尔伯格", "J. Shaffelburg", 14, "FW", "前锋", "5:1"},
	}, []demoPlayer{
		{"乔纳森·奥索里奥", "J. Osorio", 21, "MF", "中场", ""},
		{"泰勒·米勒", "K. Miller", 4, "DF", "后卫", ""},
		{"卢卡斯·卡瓦利尼", "L. Cavallini", 9, "FW", "前锋", ""},
	})

	log.Printf("seeded demo lineups for match %d: %s vs %s", match.ID, match.HomeTeam.Name, match.AwayTeam.Name)
}

func seedTeamLineup(matchID, teamID uint, side, formation, coach, source string, syncedAt *time.Time, starters, subs []demoPlayer) {
	lineup := models.MatchLineup{
		MatchID:       matchID,
		TeamID:        teamID,
		Side:          side,
		Formation:     formation,
		CoachName:     coach,
		Source:        "api-football",
		SourceMatchID: source,
		Status:        "available",
		LastSyncedAt:  syncedAt,
	}
	var existing models.MatchLineup
	if err := database.DB.Where("match_id = ? AND team_id = ?", matchID, teamID).First(&existing).Error; err == nil {
		lineup.ID = existing.ID
		lineup.CreatedAt = existing.CreatedAt
	}
	if err := database.DB.Save(&lineup).Error; err != nil {
		log.Fatalf("save %s lineup failed: %v", side, err)
	}

	if err := database.DB.Unscoped().Where("match_lineup_id = ?", lineup.ID).Delete(&models.MatchLineupPlayer{}).Error; err != nil {
		log.Fatalf("clear %s lineup players failed: %v", side, err)
	}

	rows := make([]models.MatchLineupPlayer, 0, len(starters)+len(subs))
	for i, p := range starters {
		rows = append(rows, playerRow(lineup.ID, matchID, teamID, "start_xi", i+1, p))
	}
	for i, p := range subs {
		rows = append(rows, playerRow(lineup.ID, matchID, teamID, "substitute", i+1, p))
	}
	if err := database.DB.Create(&rows).Error; err != nil {
		log.Fatalf("insert %s lineup players failed: %v", side, err)
	}
}

func playerRow(lineupID, matchID, teamID uint, role string, order int, p demoPlayer) models.MatchLineupPlayer {
	return models.MatchLineupPlayer{
		MatchLineupID: lineupID,
		MatchID:       matchID,
		TeamID:        teamID,
		Source:        "demo",
		Name:          p.name,
		NameEn:        p.nameEn,
		ShirtNumber:   p.number,
		Position:      p.position,
		PositionLabel: p.label,
		Role:          role,
		Grid:          p.grid,
		SortOrder:     order,
	}
}
