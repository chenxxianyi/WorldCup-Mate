package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"worldcup-mate/internal/models"
	"worldcup-mate/internal/utils"

	"gorm.io/gorm"
)

func Seed() {
	seedCompetitions()
	seedAdmin()
	seedGroups()
	seedCities()
	seedStadiums()
	seedTeams()
	seedMatches()
}

// seedCompetitions adds competition metadata (World Cup + top-5 leagues).
// It only inserts new rows and never touches existing data.
func seedCompetitions() {
	competitions := []models.Competition{
		{Code: "WC", Name: "世界杯", NameEn: "World Cup", Country: "World", Format: "cup", Season: 2026, Status: "active", SortOrder: 0},
		{Code: "PL", Name: "英超", NameEn: "Premier League", Country: "England", Format: "league", Season: 2025, Status: "active", SortOrder: 1},
		{Code: "PD", Name: "西甲", NameEn: "Primera Division", Country: "Spain", Format: "league", Season: 2025, Status: "active", SortOrder: 2},
		{Code: "BL1", Name: "德甲", NameEn: "Bundesliga", Country: "Germany", Format: "league", Season: 2025, Status: "active", SortOrder: 3},
		{Code: "SA", Name: "意甲", NameEn: "Serie A", Country: "Italy", Format: "league", Season: 2025, Status: "active", SortOrder: 4},
		{Code: "FL1", Name: "法甲", NameEn: "Ligue 1", Country: "France", Format: "league", Season: 2025, Status: "active", SortOrder: 5},
	}
	for _, comp := range competitions {
		var count int64
		DB.Model(&models.Competition{}).Where("code = ?", comp.Code).Count(&count)
		if count == 0 {
			DB.Create(&comp)
		}
	}
}

func EnsureNoDefaultAdminPassword(appEnv string) error {
	if !strings.EqualFold(appEnv, "production") {
		return nil
	}
	var admin models.User
	err := DB.Where("email = ?", "admin@worldcup.local").First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if utils.CheckPassword("admin123456", admin.PasswordHash) {
		return fmt.Errorf("default admin password is not allowed in production")
	}
	return nil
}

func seedAdmin() {
	if os.Getenv("APP_ENV") == "production" {
		log.Println("skip default admin seed in production")
		return
	}
	var count int64
	DB.Model(&models.User{}).Where("email = ?", "admin@worldcup.local").Count(&count)
	if count > 0 {
		return
	}
	hash, _ := utils.HashPassword("admin123456")
	admin := models.User{
		Username:     "admin",
		Email:        "admin@worldcup.local",
		PasswordHash: hash,
		Role:         "admin",
		Timezone:     "Asia/Shanghai",
		Language:     "zh-CN",
	}
	if err := DB.Create(&admin).Error; err != nil {
		log.Printf("seed admin failed: %v", err)
	} else {
		log.Println("seed admin created: admin@worldcup.local / admin123456")
	}
}

func seedGroups() {
	groups := []string{"Group A", "Group B", "Group C", "Group D", "Group E", "Group F",
		"Group G", "Group H", "Group I", "Group J", "Group K", "Group L"}
	for _, name := range groups {
		var count int64
		DB.Model(&models.Group{}).Where("name = ?", name).Count(&count)
		if count == 0 {
			DB.Create(&models.Group{Name: name, Stage: "group"})
		}
	}
}

func seedCities() {
	cities := []models.City{
		{Name: "墨西哥城", NameEn: "Mexico City", Country: "Mexico", Timezone: "America/Mexico_City"},
		{Name: "纽约", NameEn: "New York", Country: "USA", Timezone: "America/New_York"},
		{Name: "洛杉矶", NameEn: "Los Angeles", Country: "USA", Timezone: "America/Los_Angeles"},
		{Name: "达拉斯", NameEn: "Dallas", Country: "USA", Timezone: "America/Chicago"},
		{Name: "西雅图", NameEn: "Seattle", Country: "USA", Timezone: "America/Los_Angeles"},
		{Name: "迈阿密", NameEn: "Miami", Country: "USA", Timezone: "America/New_York"},
		{Name: "多伦多", NameEn: "Toronto", Country: "Canada", Timezone: "America/Toronto"},
		{Name: "温哥华", NameEn: "Vancouver", Country: "Canada", Timezone: "America/Vancouver"},
	}
	for _, city := range cities {
		var count int64
		DB.Model(&models.City{}).Where("name = ?", city.Name).Count(&count)
		if count == 0 {
			DB.Create(&city)
		}
	}
}

func seedStadiums() {
	type stadiumSeed struct {
		Name   string
		NameEn string
		City   string
		Cap    int
	}
	seeds := []stadiumSeed{
		{"Estadio Azteca", "Estadio Azteca", "墨西哥城", 87000},
		{"MetLife Stadium", "MetLife Stadium", "纽约", 82500},
		{"SoFi Stadium", "SoFi Stadium", "洛杉矶", 70000},
		{"AT&T Stadium", "AT&T Stadium", "达拉斯", 80000},
		{"Lumen Field", "Lumen Field", "西雅图", 69000},
		{"Hard Rock Stadium", "Hard Rock Stadium", "迈阿密", 65000},
		{"BMO Field", "BMO Field", "多伦多", 30000},
		{"BC Place", "BC Place", "温哥华", 54000},
	}
	for _, s := range seeds {
		var count int64
		DB.Model(&models.Stadium{}).Where("name = ?", s.Name).Count(&count)
		if count == 0 {
			var city models.City
			DB.Where("name = ?", s.City).First(&city)
			DB.Create(&models.Stadium{Name: s.Name, NameEn: s.NameEn, CityID: city.ID, Capacity: s.Cap})
		}
	}
}

func seedTeams() {
	type teamSeed struct {
		Name      string
		NameEn    string
		Code      string
		Flag      string
		Continent string
		Group     string
	}
	seeds := []teamSeed{
		{"墨西哥", "Mexico", "MEX", "🇲🇽", "北美洲", "Group A"},
		{"加拿大", "Canada", "CAN", "🇨🇦", "北美洲", "Group A"},
		{"南非", "South Africa", "RSA", "🇿🇦", "非洲", "Group A"},
		{"新西兰", "New Zealand", "NZL", "🇳🇿", "大洋洲", "Group A"},
		{"阿根廷", "Argentina", "ARG", "🇦🇷", "南美洲", "Group B"},
		{"法国", "France", "FRA", "🇫🇷", "欧洲", "Group B"},
		{"美国", "USA", "USA", "🇺🇸", "北美洲", "Group B"},
		{"加纳", "Ghana", "GHA", "🇬🇭", "非洲", "Group B"},
		{"巴西", "Brazil", "BRA", "🇧🇷", "南美洲", "Group C"},
		{"西班牙", "Spain", "ESP", "🇪🇸", "欧洲", "Group C"},
		{"英格兰", "England", "ENG", "🏴󠁧󠁢󠁥󠁮󠁧󠁿", "欧洲", "Group C"},
		{"葡萄牙", "Portugal", "POR", "🇵🇹", "欧洲", "Group C"},
		{"日本", "Japan", "JPN", "🇯🇵", "亚洲", "Group D"},
		{"德国", "Germany", "GER", "🇩🇪", "欧洲", "Group D"},
		{"荷兰", "Netherlands", "NED", "🇳🇱", "欧洲", "Group D"},
		{"比利时", "Belgium", "BEL", "🇧🇪", "欧洲", "Group D"},
	}
	for _, t := range seeds {
		var count int64
		DB.Model(&models.Team{}).Where("fifa_code = ?", t.Code).Count(&count)
		if count == 0 {
			var group models.Group
			DB.Where("name = ?", t.Group).First(&group)
			var codePtr *string
			if t.Code != "" {
				code := t.Code
				codePtr = &code
			}
			groupID := group.ID
			DB.Create(&models.Team{
				Name: t.Name, NameEn: t.NameEn, FIFACode: codePtr,
				FlagURL: t.Flag, Continent: t.Continent, GroupID: &groupID,
			})
		}
	}
}

func seedMatches() {
	if os.Getenv("DATA_SYNC_ENABLED") == "true" {
		return
	}

	// Delete existing matches with zero kickoff time (from old seed)
	DB.Where("kickoff_time_utc = ?", time.Time{}).Delete(&models.Match{})

	var count int64
	DB.Model(&models.Match{}).Count(&count)
	if count > 0 {
		return
	}

	getTeam := func(code string) models.Team {
		var t models.Team
		DB.Where("fifa_code = ?", code).First(&t)
		return t
	}
	getCity := func(name string) models.City {
		var c models.City
		DB.Where("name = ?", name).First(&c)
		return c
	}
	getStadium := func(name string) models.Stadium {
		var s models.Stadium
		DB.Where("name = ?", name).First(&s)
		return s
	}
	getGroup := func(name string) models.Group {
		var g models.Group
		DB.Where("name = ?", name).First(&g)
		return g
	}

	mex := getTeam("MEX")
	can := getTeam("CAN")
	arg := getTeam("ARG")
	fra := getTeam("FRA")
	bra := getTeam("BRA")
	esp := getTeam("ESP")
	jpn := getTeam("JPN")
	ger := getTeam("GER")

	cityMex := getCity("墨西哥城")
	cityNY := getCity("纽约")
	cityLA := getCity("洛杉矶")
	citySea := getCity("西雅图")

	stadAz := getStadium("Estadio Azteca")
	stadMet := getStadium("MetLife Stadium")
	stadSoFi := getStadium("SoFi Stadium")
	stadLum := getStadium("Lumen Field")

	grpA := getGroup("Group A")
	grpB := getGroup("Group B")
	grpC := getGroup("Group C")
	grpD := getGroup("Group D")

	now := time.Now().UTC()
	today9am := time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, time.UTC)   // 09:00 CST = 01:00 UTC
	today21pm := time.Date(now.Year(), now.Month(), now.Day(), 13, 0, 0, 0, time.UTC) // 21:00 CST = 13:00 UTC
	tomorrow18pm := today9am.Add(33 * time.Hour)
	dayAfter3am := today9am.Add(42 * time.Hour)

	matches := []models.Match{
		{MatchNo: 1, HomeTeamID: mex.ID, AwayTeamID: can.ID, GroupID: &grpA.ID, Stage: "group", StadiumID: stadAz.ID, CityID: cityMex.ID, KickoffTimeUTC: today9am, Status: "scheduled", ImportanceLevel: 3, RecommendTag: "揭幕战"},
		{MatchNo: 2, HomeTeamID: arg.ID, AwayTeamID: fra.ID, GroupID: &grpB.ID, Stage: "group", StadiumID: stadMet.ID, CityID: cityNY.ID, KickoffTimeUTC: today21pm, Status: "scheduled", ImportanceLevel: 3, RecommendTag: "焦点大战"},
		{MatchNo: 3, HomeTeamID: bra.ID, AwayTeamID: esp.ID, GroupID: &grpC.ID, Stage: "group", StadiumID: stadSoFi.ID, CityID: cityLA.ID, KickoffTimeUTC: tomorrow18pm, Status: "scheduled", ImportanceLevel: 3, RecommendTag: "热门比赛"},
		{MatchNo: 4, HomeTeamID: jpn.ID, AwayTeamID: ger.ID, GroupID: &grpD.ID, Stage: "group", StadiumID: stadLum.ID, CityID: citySea.ID, KickoffTimeUTC: dayAfter3am, Status: "scheduled", ImportanceLevel: 2, RecommendTag: ""},
	}
	for _, m := range matches {
		DB.Create(&m)
	}
	log.Println("seed matches created")
}
