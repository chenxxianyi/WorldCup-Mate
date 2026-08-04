package middleware

import (
	"strings"

	"worldcup-mate/internal/repositories"
	"worldcup-mate/internal/utils"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(c, 401, "missing authorization header")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			utils.Error(c, 401, "invalid authorization format")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			utils.Error(c, 401, "invalid or expired token")
			c.Abort()
			return
		}

		// ADM-06: disabled accounts lose access immediately (DB lookup per
		// request; fine at this scale, cache later if needed).
		// SEC-04: tokens issued before the last password change are rejected,
		// closing the post-change access-token window.
		user, err := repositories.GetUserByID(claims.UserID)
		if err != nil || user.Status == "disabled" {
			utils.Error(c, 401, "invalid or expired token")
			c.Abort()
			return
		}
		if user.PasswordChangedAt != nil && claims.IssuedAt != nil &&
			claims.IssuedAt.Time.Before(*user.PasswordChangedAt) {
			utils.Error(c, 401, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
