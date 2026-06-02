package middleware

import (
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(string) != "admin" {
			utils.Error(c, 403, "admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}
