package main

import (
	"fmt"
	"log"
	"net"
	"os"
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
	services.ConfigurePlayerSync(services.PlayerSyncConfig{
		Enabled:            cfg.PlayerSyncEnabled,
		Provider:           cfg.PlayerSyncProvider,
		APIFootballKey:     cfg.APIFootballKey,
		APIFootballBaseURL: cfg.APIFootballBaseURL,
		Interval:           time.Duration(cfg.PlayerSyncIntervalHours) * time.Hour,
		SyncOnStartup:      cfg.PlayerSyncOnStartup,
	})
	if err := services.ConfigureAI(services.AIServiceConfig{
		Provider:       cfg.AIProvider,
		BaseURL:        cfg.AIBaseURL,
		APIKey:         cfg.AIAPIKey,
		Model:          cfg.AIModel,
		TimeoutSeconds: cfg.AITimeoutSeconds,
		DailyLimitUser: cfg.AIDailyLimitUser,
		Temperature:    cfg.AITemperature,
		MaxTokens:      cfg.AIMaxTokens,
		CacheEnabled:   cfg.AICacheEnabled,
	}); err != nil {
		log.Fatalf("AI config failed: %v", err)
	}

	database.InitMySQL(cfg.MySQLDSN)
	database.InitRedis(cfg.RedisAddr, cfg.RedisPass, cfg.RedisDB)

	if err := services.EnsurePlayerSourceUniqueness(); err != nil {
		log.Fatalf("player source uniqueness cleanup failed: %v", err)
	}

	// Auto migrate
	err := database.DB.AutoMigrate(
		&models.User{},
		&models.Group{},
		&models.Team{},
		&models.Player{},
		&models.ExternalTeamMapping{},
		&models.City{},
		&models.Stadium{},
		&models.Match{},
		&models.UserFavoriteTeam{},
		&models.UserFavoriteMatch{},
		&models.Reminder{},
		&models.GroupStanding{},
		&models.Notification{},
		&models.SyncState{},
		&models.AIConversation{},
		&models.AIMessage{},
		&models.AIGeneratedContent{},
		&models.AIUsageLog{},
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
	os.MkdirAll("uploads/avatars", 0755)
	r.Static("/uploads", "uploads")
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("server listen failed on %s: %v", addr, err)
	}

	// Start background workers only after the HTTP port is available.
	jobs.StartReminderScanner()
	jobs.StartMatchSyncer()
	jobs.StartPlayerSyncer()

	log.Printf("Server starting on %s", addr)
	if err := r.RunListener(listener); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
