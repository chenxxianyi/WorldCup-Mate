package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"worldcup-mate/internal/config"
	"worldcup-mate/internal/database"
	"worldcup-mate/internal/services"
)

func main() {
	teamID := flag.Uint("team-id", 0, "local teams.id value")
	externalTeamID := flag.String("external-team-id", "", "provider team id")
	externalTeamName := flag.String("external-team-name", "", "provider team name")
	flag.Parse()

	if *teamID == 0 || *externalTeamID == "" {
		log.Fatal("team-id and external-team-id are required")
	}

	cfg := config.Load()
	database.InitMySQL(cfg.MySQLDSN)
	services.ConfigurePlayerSync(services.PlayerSyncConfig{
		Enabled:            true,
		Provider:           cfg.PlayerSyncProvider,
		APIFootballKey:     cfg.APIFootballKey,
		APIFootballBaseURL: cfg.APIFootballBaseURL,
		Interval:           time.Duration(cfg.PlayerSyncIntervalHours) * time.Hour,
	})

	if _, err := services.UpsertTeamPlayerMapping(uint(*teamID), services.ExternalTeamMappingInput{
		ExternalTeamID:   *externalTeamID,
		ExternalTeamName: *externalTeamName,
		Provider:         cfg.PlayerSyncProvider,
	}); err != nil {
		log.Fatalf("upsert mapping failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	result, err := services.SyncTeamPlayersWithDefault(ctx, uint(*teamID), "manual-cli")
	if err != nil {
		log.Fatalf("sync failed: %v", err)
	}

	payload, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(payload))
}
