package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		requestID, _ := c.Get("request_id")
		log.Printf("[%s] %s %s %d %v request_id=%v",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Writer.Status(),
			time.Since(start),
			requestID,
		)
	}
}
