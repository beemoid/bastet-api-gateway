package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIKeyAuth validates the API key in request headers
// This middleware protects endpoints from unauthorized access
// The API key should be sent in the "X-API-Key" header
func APIKeyAuth(expectedKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth check if no API key is configured
		if expectedKey == "" {
			c.Next()
			return
		}

		// Get API key from header
		apiKey := c.GetHeader("X-API-Key")

		// Validate API key
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Missing API key",
				"message": "Please provide X-API-Key header",
			})
			c.Abort()
			return
		}

		if apiKey != expectedKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid API key",
				"message": "The provided API key is not valid",
			})
			c.Abort()
			return
		}

		// API key is valid, proceed to next handler
		c.Next()
	}
}
