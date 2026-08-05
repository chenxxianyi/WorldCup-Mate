package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/middleware"
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// setupDB wires an in-memory DB with one user whose PasswordChangedAt is
// pinned relative to the current second (fixed offset, not wall-clock).
func setupDB(t *testing.T, pwChanged *time.Time) {
	t.Helper()
	// Unique shared-cache DB per test (all :memory: connections otherwise
	// see the same database).
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
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	database.DB = db
	utils.SetJWTSecret("jwt-test-secret")

	hash, _ := utils.HashPassword("Password123!")
	if err := db.Create(&models.User{
		Username: "jwtuser", Email: "jwt@test.dev", PasswordHash: hash,
		Role: "user", Status: "active", PasswordChangedAt: pwChanged,
	}).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
}

// runAuth runs the JWTAuth middleware against a request carrying `token`.
func runAuth(t *testing.T, token string) int {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/api/user/profile", nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	c.Request = req

	middleware.JWTAuth()(c)
	if c.IsAborted() {
		return w.Code
	}
	return http.StatusOK
}

func TestJWTAuthPasswordChangedWindow(t *testing.T) {
	// Pin times to the CURRENT second (a fixed past iat would make the
	// token expired before the middleware can evaluate it).
	base := time.Now().UTC().Truncate(time.Second)

	// Password changed at base+900ms — inside the same second as `base`.
	pw := base.Add(900 * time.Millisecond)
	setupDB(t, &pw)

	// A token minted inside the SAME second (whether before or after the
	// change, at second precision they are indistinguishable) is accepted.
	tok, err := utils.GenerateTokenAt(1, "user", base.Add(100*time.Millisecond))
	if err != nil {
		t.Fatalf("token: %v", err)
	}
	if code := runAuth(t, tok); code != http.StatusOK {
		t.Fatalf("same-second token: status = %d, want 200", code)
	}

	// A token issued the second BEFORE the change must be rejected.
	tokPrior, err := utils.GenerateTokenAt(1, "user", base.Add(-1*time.Second))
	if err != nil {
		t.Fatalf("token: %v", err)
	}
	if code := runAuth(t, tokPrior); code != http.StatusUnauthorized {
		t.Fatalf("pre-change token: status = %d, want 401", code)
	}
}

func TestJWTAuthLegacyUserWithoutPasswordChangedAt(t *testing.T) {
	setupDB(t, nil) // legacy account, never changed password
	tok, err := utils.GenerateTokenAt(1, "user", time.Now().Add(-time.Minute))
	if err != nil {
		t.Fatalf("token: %v", err)
	}
	if code := runAuth(t, tok); code != http.StatusOK {
		t.Fatalf("legacy user token: status = %d, want 200", code)
	}
}
