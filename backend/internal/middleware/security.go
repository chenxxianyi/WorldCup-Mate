package middleware

import "github.com/gin-gonic/gin"

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		if len(c.Request.URL.Path) >= len("/uploads/") && c.Request.URL.Path[:len("/uploads/")] == "/uploads/" {
			c.Header("Cache-Control", "public, max-age=86400")
		}
		c.Next()
	}
}
