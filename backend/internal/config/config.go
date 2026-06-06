package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                        string
	AppPort                       string
	JWTSecret                     string
	MySQLDSN                      string
	RedisAddr                     string
	RedisPass                     string
	RedisDB                       string
	DataSyncEnabled               bool
	DataSyncProvider              string
	DataSyncLiveIntervalSeconds   int
	DataSyncIdleIntervalMinutes   int
	DataSyncFullIntervalHours     int
	FootballDataAPIKey            string
	FootballDataBaseURL           string
	PlayerSyncEnabled             bool
	PlayerSyncProvider            string
	APIFootballKey                string
	APIFootballBaseURL            string
	PlayerSyncIntervalHours       int
	PlayerSyncOnStartup           bool
	LineupSyncEnabled             bool
	LineupSyncPrimaryProvider     string
	LineupSyncEnhancedProvider    string
	LineupSyncPregameMinutes      int
	LineupSyncLiveIntervalSeconds int
	SMTPHost                      string
	SMTPPort                      int
	SMTPUsername                  string
	SMTPPassword                  string
	SMTPFrom                      string
	AIProvider                    string
	AIBaseURL                     string
	AIAPIKey                      string
	AIModel                       string
	AITimeoutSeconds              int
	AIDailyLimitUser              int
	AITemperature                 float64
	AIMaxTokens                   int
	AICacheEnabled                bool
	LineupAlertEnabled            bool
	LineupAlertScanIntervalSeconds int
	LineupAlertWindowBeforeMinutes int
	LineupAlertWindowAfterMinutes  int
	PostMatchSummaryEnabled        bool
	PostMatchSummaryScanIntervalSeconds int
	PostMatchSummaryLookbackHours      int
	PostMatchSummaryAutoGenerate       bool
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AppEnv:                        getEnv("APP_ENV", "development"),
		AppPort:                       getEnv("APP_PORT", "8080"),
		JWTSecret:                     getEnv("JWT_SECRET", "default_secret"),
		MySQLDSN:                      getEnv("MYSQL_DSN", "xxladmin:XXLadmin_2021!@tcp(127.0.0.1:3310)/worldcup_mate?charset=utf8mb4&parseTime=True&loc=Local"),
		RedisAddr:                     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPass:                     getEnv("REDIS_PASSWORD", ""),
		RedisDB:                       getEnv("REDIS_DB", "0"),
		DataSyncEnabled:               getBoolEnv("DATA_SYNC_ENABLED", false),
		DataSyncProvider:              getEnv("DATA_SYNC_PROVIDER", "football-data"),
		DataSyncLiveIntervalSeconds:   getIntEnv("DATA_SYNC_LIVE_INTERVAL_SECONDS", 120),
		DataSyncIdleIntervalMinutes:   getIntEnv("DATA_SYNC_IDLE_INTERVAL_MINUTES", 30),
		DataSyncFullIntervalHours:     getIntEnv("DATA_SYNC_FULL_INTERVAL_HOURS", 6),
		FootballDataAPIKey:            getEnv("FOOTBALL_DATA_API_KEY", ""),
		FootballDataBaseURL:           getEnv("FOOTBALL_DATA_BASE_URL", "https://api.football-data.org/v4"),
		PlayerSyncEnabled:             getBoolEnv("PLAYER_SYNC_ENABLED", false),
		PlayerSyncProvider:            getEnv("PLAYER_SYNC_PROVIDER", "api-football"),
		APIFootballKey:                getEnv("API_FOOTBALL_KEY", ""),
		APIFootballBaseURL:            getEnv("API_FOOTBALL_BASE_URL", "https://v3.football.api-sports.io"),
		PlayerSyncIntervalHours:       getIntEnv("PLAYER_SYNC_INTERVAL_HOURS", 168),
		PlayerSyncOnStartup:           getBoolEnv("PLAYER_SYNC_ON_STARTUP", false),
		LineupSyncEnabled:             getBoolEnv("LINEUP_SYNC_ENABLED", false),
		LineupSyncPrimaryProvider:     getEnv("LINEUP_SYNC_PRIMARY_PROVIDER", "football-data"),
		LineupSyncEnhancedProvider:    getEnv("LINEUP_SYNC_ENHANCED_PROVIDER", "api-football"),
		LineupSyncPregameMinutes:      getIntEnv("LINEUP_SYNC_PREGAME_MINUTES", 90),
		LineupSyncLiveIntervalSeconds: getIntEnv("LINEUP_SYNC_LIVE_INTERVAL_SECONDS", 900),
		SMTPHost:                      getEnv("SMTP_HOST", ""),
		SMTPPort:                      getIntEnv("SMTP_PORT", 587),
		SMTPUsername:                  getEnv("SMTP_USERNAME", ""),
		SMTPPassword:                  getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:                      getEnv("SMTP_FROM", ""),
		AIProvider:                    getEnv("AI_PROVIDER", "openai"),
		AIBaseURL:                     getEnv("AI_BASE_URL", "https://api.openai.com/v1"),
		AIAPIKey:                      getEnv("AI_API_KEY", ""),
		AIModel:                       getEnv("AI_MODEL", "gpt-4o-mini"),
		AITimeoutSeconds:              getIntEnv("AI_TIMEOUT_SECONDS", 60),
		AIDailyLimitUser:              getIntEnv("AI_DAILY_LIMIT_USER", 50),
		AITemperature:                 getFloatEnv("AI_TEMPERATURE", 0.7),
		AIMaxTokens:                   getIntEnv("AI_MAX_TOKENS", 1200),
		AICacheEnabled:                getBoolEnv("AI_CACHE_ENABLED", true),
		LineupAlertEnabled:            getBoolEnv("LINEUP_ALERT_ENABLED", true),
		LineupAlertScanIntervalSeconds: getIntEnv("LINEUP_ALERT_SCAN_INTERVAL_SECONDS", 120),
		LineupAlertWindowBeforeMinutes: getIntEnv("LINEUP_ALERT_WINDOW_BEFORE_MINUTES", 120),
		LineupAlertWindowAfterMinutes:  getIntEnv("LINEUP_ALERT_WINDOW_AFTER_MINUTES", 30),
		PostMatchSummaryEnabled:        getBoolEnv("POST_MATCH_SUMMARY_ENABLED", true),
		PostMatchSummaryScanIntervalSeconds: getIntEnv("POST_MATCH_SUMMARY_SCAN_INTERVAL_SECONDS", 600),
		PostMatchSummaryLookbackHours:      getIntEnv("POST_MATCH_SUMMARY_LOOKBACK_HOURS", 24),
		PostMatchSummaryAutoGenerate:       getBoolEnv("POST_MATCH_SUMMARY_AUTO_GENERATE", true),
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

func getFloatEnv(key string, fallback float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	parsed, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fallback
	}
	return parsed
}
