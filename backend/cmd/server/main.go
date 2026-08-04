package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"worldcup-mate/internal/config"
	"worldcup-mate/internal/database"
	"worldcup-mate/internal/jobs"
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/routes"
	"worldcup-mate/internal/services"
	"worldcup-mate/internal/utils"
)

func main() {
	cfg := config.Load()
	if err := cfg.ValidateProduction(); err != nil {
		log.Fatalf("config validation failed: %v", err)
	}

	utils.SetJWTSecret(cfg.JWTSecret)
	utils.InitEmail(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPFrom)
	services.ConfigureMatchSync(services.MatchSyncConfig{
		Enabled:             cfg.DataSyncEnabled,
		Provider:            cfg.DataSyncProvider,
		FootballDataAPIKey:  cfg.FootballDataAPIKey,
		FootballDataBaseURL: cfg.FootballDataBaseURL,
		LiveInterval:        time.Duration(cfg.DataSyncLiveIntervalSeconds) * time.Second,
		IdleInterval:        time.Duration(cfg.DataSyncIdleIntervalMinutes) * time.Minute,
		FullInterval:        time.Duration(cfg.DataSyncFullIntervalHours) * time.Hour,
	})
	services.ConfigureLeagueSync(
		cfg.SyncCompetitions,
		time.Duration(cfg.LeagueSyncIntervalMinutes)*time.Minute,
		cfg.FootballDataBaseURL,
		cfg.FootballDataAPIKey,
	)

	database.InitMySQL(cfg.MySQLDSN)
	database.InitRedis(cfg.RedisAddr, cfg.RedisPass, cfg.RedisDB)

	// Auto migrate
	err := database.DB.AutoMigrate(
		&models.User{},
		&models.Group{},
		&models.Team{},
		&models.City{},
		&models.Stadium{},
		&models.Match{},
		&models.UserFavoriteTeam{},
		&models.UserFavoriteMatch{},
		&models.Reminder{},
		&models.GroupStanding{},
		&models.Notification{},
		&models.SyncState{},
		&models.Competition{},
		&models.LeagueStanding{},
		&models.AdminAuditLog{},
		&models.RefreshToken{},
	)
	if err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}
	if err := database.EnsureNoDefaultAdminPassword(cfg.AppEnv); err != nil {
		log.Fatalf("security validation failed: %v", err)
	}

	// Seed data
	database.Seed()

	// Update team names with localized Chinese versions (after seed, before sync)
	services.UpdateAllTeamNames()
	if updated, err := services.ApplyOfficialVenueMappings(); err != nil {
		log.Printf("official venue mapping failed: %v", err)
	} else if updated > 0 {
		log.Printf("official venue mapping updated %d matches", updated)
	}

	// Setup routes and start
	r := routes.Setup()

	// Serve uploaded files
	if err := os.MkdirAll(filepath.Join(cfg.UploadDir, "avatars"), 0755); err != nil {
		log.Printf("[warn] failed to create upload dir: %v", err)
	}
	services.SetUploadRoot(cfg.UploadDir)
	r.Static("/uploads", cfg.UploadDir)
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("server listen failed on %s: %v", addr, err)
	}

	// Start background workers only after the HTTP port is available.
	jobs.StartReminderScanner()
	jobs.StartMatchSyncer()
	jobs.StartLeagueSyncer()

	log.Printf("Server starting on %s", addr)
	if err := r.RunListener(listener); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
