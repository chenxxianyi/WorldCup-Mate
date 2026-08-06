package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
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
		&models.FeaturedConfig{},
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
	// A cancelable process context allows graceful shutdown of all workers.
	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := &http.Server{Handler: r}

	// Listen for OS signals and shut down gracefully on SIGINT/SIGTERM (OBS-04).
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Printf("received signal %v, shutting down...", sig)

		// Stop background workers first.
		jobs.StopReminderScanner()
		jobs.StopMatchSyncer()
		jobs.StopLeagueSyncer()
		cancel() // cancels in-flight worker contexts

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("graceful shutdown error: %v", err)
		}
		os.Exit(0)
	}()

	jobs.StartReminderScanner(appCtx)
	jobs.StartMatchSyncer(appCtx)
	jobs.StartLeagueSyncer(appCtx)

	log.Printf("Server starting on %s", addr)
	if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
