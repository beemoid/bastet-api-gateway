package middleware

import (
	"api-gateway/service"
	"net/http"
	"time"

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

// CombinedAuth validates generated API tokens from the token management system.
// Accepts X-API-Token header with tokens created via the admin dashboard.
func CombinedAuth(tokenService *service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if tokenService == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"error":   "Token service unavailable",
				"message": "Token management system is not configured",
			})
			c.Abort()
			return
		}

		// Extract token from header
		apiToken := c.GetHeader("X-API-Token")
		if apiToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Missing authentication",
				"message": "Please provide X-API-Token header",
			})
			c.Abort()
			return
		}

		// Validate token
		token, err := tokenService.ValidateAPIToken(apiToken, c.ClientIP())
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   err.Error(),
				"message": "Invalid API token",
			})
			c.Abort()
			return
		}

		// Check rate limits
		rateLimits := map[string]int{
			"minute": token.RateLimitPerMinute,
			"hour":   token.RateLimitPerHour,
			"day":    token.RateLimitPerDay,
		}

		allowed, message, err := tokenService.CheckRateLimit(token.ID, rateLimits)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
				"message": "Rate limit check failed",
			})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "Please slow down your requests",
				"message": message,
			})
			c.Abort()
			return
		}

		// Store token info in context
		c.Set("auth_type", "token")
		c.Set("token_id", token.ID)
		c.Set("token_name", token.Name)
		c.Set("token_scopes", token.Scopes)
		// Vendor filter context – read by handlers to scope DB queries
		c.Set("token_is_super", token.IsSuperToken)
		c.Set("token_vendor_name", token.VendorName)
		c.Set("token_filter_column", token.FilterColumn)
		c.Set("token_filter_value", token.FilterValue)

		// Process request
		startTime := time.Now()
		c.Next()

		// Log usage after request completes
		statusCode := c.Writer.Status()
		logUsage(tokenService, token.ID, c, startTime, statusCode, "")
	}
}
