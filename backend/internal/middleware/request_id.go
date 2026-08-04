package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"

	"github.com/gin-gonic/gin"
)

const RequestIDHeader = "X-Request-ID"

// Only accept well-formed upstream IDs (letters, digits, dot, dash,
// underscore, 8-64 chars); anything else is replaced with a fresh ID.
var requestIDPattern = regexp.MustCompile(`^[A-Za-z0-9._-]{8,64}$`)

// RequestID injects a request ID for tracing: accepts a trusted upstream
// X-Request-ID header, otherwise generates one. The ID is stored on the
// gin context ("request_id"), echoed in the response header, and used by
// the logger and error responses (API-02).
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(RequestIDHeader)
		if !requestIDPattern.MatchString(id) {
			id = generateRequestID()
		}
		c.Set("request_id", id)
		c.Header(RequestIDHeader, id)
		c.Next()
	}
}

func generateRequestID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
