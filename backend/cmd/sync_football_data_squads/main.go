package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"worldcup-mate/internal/config"
	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/providers/footballdata"
	"worldcup-mate/internal/repositories"
	"worldcup-mate/internal/services"
)

type fallbackReport struct {
	TeamID       uint   `json:"team_id"`
	Name         string `json:"name"`
	FIFACode     string `json:"fifa_code"`
	ExternalID   int64  `json:"external_id,omitempty"`
	ExternalName string `json:"external_name,omitempty"`
	Status       string `json:"status"`
	Players      int    `json:"players"`
	Error        string `json:"error,omitempty"`
}

func main() {
	includeExisting := flag.Bool("include-existing", false, "also import teams that already have active players")
	flag.Parse()

	cfg := config.Load()
	database.InitMySQL(cfg.MySQLDSN)

	client := footballdata.NewClient(cfg.FootballDataBaseURL, cfg.FootballDataAPIKey)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	data, err := client.CompetitionTeams(ctx, "WC", 2026)
	if err != nil {
		log.Fatalf("fetch football-data squads failed: %v", err)
	}

	var teams []models.Team
	if err := database.DB.Where("deleted_at IS NULL").Order("id ASC").Find(&teams).Error; err != nil {
		log.Fatalf("list local teams failed: %v", err)
	}

	externalByKey := indexExternalTeams(data.Teams)
	reports := []fallbackReport{}
	for _, team := range teams {
		report := fallbackReport{
			TeamID:   team.ID,
			Name:     team.Name,
			FIFACode: team.FIFACode,
		}
		if isPlaceholder(team) {
			report.Status = "skipped"
			report.Error = "placeholder team"
			reports = append(reports, report)
			continue
		}
		if countActivePlayers(team.ID) > 0 && !*includeExisting {
			report.Status = "skipped"
			report.Players = countActivePlayers(team.ID)
			reports = append(reports, report)
			continue
		}

		external, ok := matchExternalTeam(team, externalByKey)
		if !ok {
			report.Status = "failed"
			report.Error = "football-data team not found"
			reports = append(reports, report)
			continue
		}
		report.ExternalID = external.ID
		report.ExternalName = external.Name

		players := mapSquad(team.ID, external)
		if len(players) == 0 {
			report.Status = "failed"
			report.Error = "football-data squad empty"
			reports = append(reports, report)
			continue
		}
		for i := range players {
			if _, err := repositories.UpsertPlayer(&players[i]); err != nil {
				report.Status = "failed"
				report.Error = err.Error()
				break
			}
			report.Players++
		}
		if report.Status == "" {
			report.Status = "success"
		}
		reports = append(reports, report)
	}

	payload, err := json.MarshalIndent(reports, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(payload))
}

func indexExternalTeams(teams []footballdata.Team) map[string]footballdata.Team {
	index := map[string]footballdata.Team{}
	for _, team := range teams {
		keys := []string{team.TLA, team.Name, team.ShortName}
		for _, key := range keys {
			key = normalizeKey(key)
			if key != "" {
				index[key] = team
			}
		}
	}
	return index
}

func matchExternalTeam(team models.Team, index map[string]footballdata.Team) (footballdata.Team, bool) {
	candidates := []string{team.FIFACode, team.NameEn, team.Name}
	aliases := map[string][]string{
		"CUR": {"CUW", "Curacao", "Curaçao"},
		"CUW": {"CUW", "Curacao", "Curaçao"},
		"CPV": {"Cape Verde"},
		"IRQ": {"Iraq"},
		"CRO": {"Croatia"},
	}
	if extra := aliases[strings.ToUpper(strings.TrimSpace(team.FIFACode))]; len(extra) > 0 {
		candidates = append(extra, candidates...)
	}
	for _, candidate := range candidates {
		if external, ok := index[normalizeKey(candidate)]; ok {
			return external, true
		}
	}
	return footballdata.Team{}, false
}

func mapSquad(teamID uint, team footballdata.Team) []models.Player {
	now := time.Now().UTC()
	players := make([]models.Player, 0, len(team.Squad))
	for _, item := range team.Squad {
		if strings.TrimSpace(item.Name) == "" {
			continue
		}
		nameEn := strings.TrimSpace(item.Name)
		sourcePlayerID := strconv.FormatInt(item.ID, 10)
		position, positionLabel := normalizeFootballDataPosition(item.Position)
		players = append(players, models.Player{
			TeamID:         teamID,
			Name:           services.LocalizePlayerName("football-data", sourcePlayerID, nameEn),
			NameEn:         nameEn,
			ShirtNumber:    0,
			Position:       position,
			PositionLabel:  positionLabel,
			Source:         "football-data",
			SourcePlayerID: sourcePlayerID,
			ExternalTeamID: strconv.FormatInt(team.ID, 10),
			IsActive:       true,
			LastSyncedAt:   &now,
		})
	}
	return players
}

func normalizeFootballDataPosition(value string) (string, string) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "goalkeeper":
		return "GK", "门将"
	case "defence", "defender":
		return "DF", "后卫"
	case "midfield", "midfielder":
		return "MF", "中场"
	case "offence", "attacker", "forward":
		return "FW", "前锋"
	default:
		return strings.TrimSpace(value), ""
	}
}

func countActivePlayers(teamID uint) int {
	var count int64
	database.DB.Model(&models.Player{}).
		Where("team_id = ? AND is_active = ?", teamID, true).
		Count(&count)
	return int(count)
}

func isPlaceholder(team models.Team) bool {
	code := strings.TrimSpace(strings.ToUpper(team.FIFACode))
	name := strings.TrimSpace(strings.ToUpper(team.NameEn))
	return code == "" || code == "TBD" || name == "" || name == "TBD"
}

func normalizeKey(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	replacer := strings.NewReplacer("ç", "c", "ã", "a", "á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u")
	value = replacer.Replace(value)
	return strings.Join(strings.Fields(value), " ")
}
