package middleware

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// JSONLogger emits structured JSON per request (OBS-04).
// Each log line is a single JSON object consumable by log collectors.
func JSONLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start).Milliseconds()
		requestID, _ := c.Get("request_id")
		rid := ""
		if v, ok := requestID.(string); ok {
			rid = v
		}
		record := gin.H{
			"timestamp":  time.Now().UTC().Format(time.RFC3339),
			"level":      "info",
			"event":      "http_request",
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"latency_ms": latency,
			"request_id": rid,
			"ip":         c.ClientIP(),
		}
		if userID, exists := c.Get("user_id"); exists {
			record["user_id"] = userID
		}
		b, err := json.Marshal(record)
		if err != nil {
			log.Printf(`{"event":"http_request","level":"error","error":"marshal failed"}`)
			return
		}
		log.Println(string(b))
	}
}
