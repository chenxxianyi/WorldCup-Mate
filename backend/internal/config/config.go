package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                      string
	AppPort                     string
	JWTSecret                   string
	MySQLDSN                    string
	RedisAddr                   string
	RedisPass                   string
	RedisDB                     string
	DataSyncEnabled             bool
	DataSyncProvider            string
	DataSyncLiveIntervalSeconds int
	DataSyncIdleIntervalMinutes int
	DataSyncFullIntervalHours   int
	FootballDataAPIKey          string
	FootballDataBaseURL         string
	SyncCompetitions            string
	LeagueSyncIntervalMinutes   int
	SMTPHost                    string
	SMTPPort                    int
	SMTPUsername                string
	SMTPPassword                string
	SMTPFrom                    string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AppEnv:                      getEnv("APP_ENV", "development"),
		AppPort:                     getEnv("APP_PORT", "8080"),
		JWTSecret:                   getEnv("JWT_SECRET", "default_secret"),
		MySQLDSN:                    getEnv("MYSQL_DSN", "xxladmin:XXLadmin_2021!@tcp(127.0.0.1:3310)/worldcup_mate?charset=utf8mb4&parseTime=True&loc=Local"),
		RedisAddr:                   getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPass:                   getEnv("REDIS_PASSWORD", ""),
		RedisDB:                     getEnv("REDIS_DB", "0"),
		DataSyncEnabled:             getBoolEnv("DATA_SYNC_ENABLED", false),
		DataSyncProvider:            getEnv("DATA_SYNC_PROVIDER", "football-data"),
		DataSyncLiveIntervalSeconds: getIntEnv("DATA_SYNC_LIVE_INTERVAL_SECONDS", 120),
		DataSyncIdleIntervalMinutes: getIntEnv("DATA_SYNC_IDLE_INTERVAL_MINUTES", 30),
		DataSyncFullIntervalHours:   getIntEnv("DATA_SYNC_FULL_INTERVAL_HOURS", 6),
		FootballDataAPIKey:          getEnv("FOOTBALL_DATA_API_KEY", ""),
		FootballDataBaseURL:         getEnv("FOOTBALL_DATA_BASE_URL", "https://api.football-data.org/v4"),
		SyncCompetitions:            getEnv("SYNC_COMPETITIONS", ""),
		LeagueSyncIntervalMinutes:   getIntEnv("LEAGUE_SYNC_INTERVAL_MINUTES", 30),
		SMTPHost:                    getEnv("SMTP_HOST", ""),
		SMTPPort:                    getIntEnv("SMTP_PORT", 587),
		SMTPUsername:                getEnv("SMTP_USERNAME", ""),
		SMTPPassword:                getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:                    getEnv("SMTP_FROM", ""),
	}
}

func (c *Config) ValidateProduction() error {
	if !strings.EqualFold(c.AppEnv, "production") {
		return nil
	}
	if strings.TrimSpace(c.JWTSecret) == "" || c.JWTSecret == "default_secret" {
		return fmt.Errorf("JWT_SECRET must be configured with a strong non-default value in production")
	}
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters in production")
	}
	if strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS")) == "" {
		return fmt.Errorf("CORS_ALLOWED_ORIGINS must be configured in production")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getBoolEnv(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return parsed
}

func getIntEnv(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(v)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
