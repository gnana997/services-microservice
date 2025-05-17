package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks if the request is authenticated
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// TODO: Implement token validation
		// For now, we'll just log the token
		log.Printf("Received token: %s", token)

		c.Next()
	}
}



