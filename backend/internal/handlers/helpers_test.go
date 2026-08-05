package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/testutil"
)

// perform sends a JSON request and returns the recorder.
func perform(t *testing.T, method, path, token string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", testutil.AuthHeader(token))
	}
	w := httptest.NewRecorder()
	testutil.Router.ServeHTTP(w, req)
	return w
}

// read decodes the response envelope {code, message, data}.
func read(t *testing.T, w *httptest.ResponseRecorder) (code int, message string, data map[string]any) {
	t.Helper()
	var envelope struct {
		Code    int            `json:"code"`
		Message string         `json:"message"`
		Data    map[string]any `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode response %q: %v", w.Body.String(), err)
	}
	return envelope.Code, envelope.Message, envelope.Data
}

// databaseUpdateStatus flips a user's status directly (simulating the
// admin disable action without going through the admin endpoint).
func databaseUpdateStatus(t *testing.T, userID uint, status string) {
	t.Helper()
	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Update("status", status).Error; err != nil {
		t.Fatalf("update status: %v", err)
	}
}
