package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "internal server error",
					"data":    nil,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
