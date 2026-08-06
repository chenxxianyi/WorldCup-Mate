package handlers

import (
	"worldcup-mate/internal/database"
	"worldcup-mate/internal/repositories"
	"worldcup-mate/internal/services"

	"github.com/gin-gonic/gin"
)

// HealthLive returns 200 immediately to indicate the process is alive (OBS-04).
func HealthLive(c *gin.Context) {
	c.JSON(200, gin.H{"status": "alive", "service": "worldcup-mate"})
}

// HealthReady checks MySQL and Redis are reachable (OBS-04).
func HealthReady(c *gin.Context) {
	checks := map[string]string{}
	healthy := true

	// MySQL
	if database.DB != nil {
		if err := database.DB.Raw("SELECT 1").Error; err != nil {
			checks["mysql"] = err.Error()
			healthy = false
		} else {
			checks["mysql"] = "ok"
		}
	} else {
		checks["mysql"] = "uninitialized"
		healthy = false
	}

	// Redis
	if database.RDB != nil {
		if err := database.RDB.Ping(c.Request.Context()).Err(); err != nil {
			checks["redis"] = err.Error()
			healthy = false
		} else {
			checks["redis"] = "ok"
		}
	} else {
		checks["redis"] = "uninitialized"
		healthy = false
	}

	if !healthy {
		c.JSON(503, gin.H{"status": "not_ready", "checks": checks})
		return
	}

	syncStates, _ := services.GetSyncStates()
	reminderStats := map[string]int64{
		"pending": repositories.CountRemindersByStatus("pending"),
		"sending": repositories.CountRemindersByStatus("sending"),
		"sent":    repositories.CountRemindersByStatus("sent"),
		"failed":  repositories.CountRemindersByStatus("failed"),
	}

	c.JSON(200, gin.H{
		"status":         "ready",
		"service":        "worldcup-mate",
		"checks":         checks,
		"sync_states":    syncStates,
		"reminder_stats": reminderStats,
	})
}
