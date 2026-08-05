package handlers_test

import (
	"net/http"
	"testing"

	"worldcup-mate/internal/testutil"
)

func TestAuthLoginFlow(t *testing.T) {
	testutil.Setup(t)
	defer testutil.ResetDB(t)
	testutil.CreateUser(t, "user@test.dev", "user", "active")

	t.Run("login succeeds with correct password", func(t *testing.T) {
		w := perform(t, http.MethodPost, "/api/auth/login", "", map[string]string{
			"email": "user@test.dev", "password": "Password123!",
		})
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		code, _, data := read(t, w)
		if code != 0 {
			t.Fatalf("business code = %d", code)
		}
		if data["token"] == "" || data["refresh_token"] == "" {
			t.Fatal("login must return access + refresh tokens")
		}
	})

	t.Run("wrong password returns 401", func(t *testing.T) {
		w := perform(t, http.MethodPost, "/api/auth/login", "", map[string]string{
			"email": "user@test.dev", "password": "WrongPass123!",
		})
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", w.Code)
		}
	})

	t.Run("disabled user cannot log in", func(t *testing.T) {
		testutil.CreateUser(t, "disabled@test.dev", "user", "disabled")
		w := perform(t, http.MethodPost, "/api/auth/login", "", map[string]string{
			"email": "disabled@test.dev", "password": "Password123!",
		})
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", w.Code)
		}
	})

	t.Run("disabled user token is rejected by the JWT middleware", func(t *testing.T) {
		// A token minted while the account was active must die immediately
		// after the account is disabled (ADM-06 + SEC-04).
		user := testutil.CreateUser(t, "late-disable@test.dev", "user", "active")
		token := testutil.TokenFor(t, user)
		databaseUpdateStatus(t, user.ID, "disabled")

		w := perform(t, http.MethodGet, "/api/user/profile", token, nil)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401 (token must die on disable)", w.Code)
		}
	})

	t.Run("register enforces password policy", func(t *testing.T) {
		w := perform(t, http.MethodPost, "/api/auth/register", "", map[string]string{
			"username": "newbie", "email": "newbie@test.dev", "password": "weak",
		})
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400 (weak password)", w.Code)
		}
	})
}
