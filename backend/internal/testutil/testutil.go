// Package testutil provides an in-memory MySQL-free test harness for
// handler-level tests (QA-02): a GORM database backed by pure-Go SQLite
// (no CGO), a miniredis stand-in, seeded users and JWT helpers.
package testutil

import (
	"strings"
	"testing"
	"time"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/routes"
	"worldcup-mate/internal/services"
	"worldcup-mate/internal/utils"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// model list mirrors cmd/server/main.go AutoMigrate.
var allModels = []any{
	&models.User{},
	&models.Group{},
	&models.Team{},
	&models.Stadium{},
	&models.City{},
	&models.Match{},
	&models.GroupStanding{},
	&models.UserFavoriteTeam{},
	&models.UserFavoriteMatch{},
	&models.Reminder{},
	&models.Notification{},
	&models.Competition{},
	&models.LeagueStanding{},
	&models.AdminAuditLog{},
	&models.RefreshToken{},
	&models.SyncState{},
	&models.FeaturedConfig{},
}

// Router is the engine built by the last Setup call (helpers access it).
var Router *gin.Engine

// Setup replaces the global DB/Redis with in-memory stand-ins and returns
// a fully wired router. Tests must not run in parallel (global singletons).
func Setup(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	// In-memory SQLite. A unique name per test avoids the shared-cache
	// trap where every ":memory:" connection sees the same database.
	name := "mem" + strings.ReplaceAll(t.Name(), "/", "_")
	db, err := gorm.Open(sqlite.Open("file:"+name+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(allModels...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	database.DB = db

	// miniredis for rate limiting + refresh tokens.
	mr := miniredis.RunT(t)
	database.RDB = redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = database.RDB.Close() })

	utils.SetJWTSecret("qa2-test-secret")

	Router = routes.Setup()
	return Router
}

// CreateUser inserts a user with the given role/status and returns it.
func CreateUser(t *testing.T, email, role, status string) *models.User {
	t.Helper()
	hash, err := utils.HashPassword("Password123!")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	now := time.Now()
	username := "tester"
	if at := strings.IndexByte(email, '@'); at > 0 {
		username = email[:at]
	}
	user := models.User{
		Username:          username,
		Email:             email,
		PasswordHash:      hash,
		Role:              role,
		Status:            status,
		Timezone:          "Asia/Shanghai",
		Language:          "zh",
		NotificationEmail: email,
		PasswordChangedAt: &now,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	return &user
}

// TokenFor issues an access token for the user (30min TTL).
func TokenFor(t *testing.T, user *models.User) string {
	t.Helper()
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		t.Fatalf("token: %v", err)
	}
	return token
}

// AuthHeader builds the Authorization header value.
func AuthHeader(token string) string { return "Bearer " + token }

// RefreshTokenFor issues a refresh token for the user via the service.
func RefreshTokenFor(t *testing.T, user *models.User) string {
	t.Helper()
	raw, err := services.IssueRefreshToken(user.ID)
	if err != nil {
		t.Fatalf("refresh token: %v", err)
	}
	return raw
}

// SetupServices wires an in-memory DB and miniredis (no router) for
// service-level tests. Tests must not run in parallel (global singletons).
func SetupServices(t *testing.T) {
	t.Helper()
	name := "svc" + strings.ReplaceAll(t.Name(), "/", "_")
	db, err := gorm.Open(sqlite.Open("file:"+name+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(allModels...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	database.DB = db

	mr := miniredis.RunT(t)
	database.RDB = redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = database.RDB.Close() })
}

// CreateTeam inserts a minimal club team.
func CreateTeam(t *testing.T, name string) *models.Team {
	t.Helper()
	team := models.Team{Name: name, NameEn: name, TeamType: "club"}
	if err := database.DB.Create(&team).Error; err != nil {
		t.Fatalf("seed team: %v", err)
	}
	return &team
}

// CreateMatch inserts a minimal scheduled match.
func CreateMatch(t *testing.T, home, away uint) *models.Match {
	t.Helper()
	m := models.Match{
		MatchNo:        int(time.Now().UnixNano() % 1e9),
		Stage:          "Group Stage",
		Status:         "scheduled",
		HomeTeamID:     home,
		AwayTeamID:     away,
		KickoffTimeUTC: time.Now().UTC().Add(24 * time.Hour),
	}
	if err := database.DB.Create(&m).Error; err != nil {
		t.Fatalf("seed match: %v", err)
	}
	return &m
}

// ResetDB wipes all rows (between test groups sharing one DB instance).
func ResetDB(t *testing.T) {
	t.Helper()
	for _, m := range allModels {
		if err := database.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(m).Error; err != nil {
			t.Fatalf("reset %T: %v", m, err)
		}
	}
}

