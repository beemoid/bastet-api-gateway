package middleware

import (
	"api-gateway/models"
	"api-gateway/service"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TokenAuthMiddleware validates API tokens and logs usage
func TokenAuthMiddleware(tokenService *service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Extract token from header
		tokenValue := c.GetHeader("X-API-Token")
		if tokenValue == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Missing API token",
				Error:   "X-API-Token header is required",
			})
			c.Abort()
			return
		}

		// Get client IP
		clientIP := c.ClientIP()

		// Validate token
		token, err := tokenService.ValidateAPIToken(tokenValue, clientIP)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Invalid API token",
				Error:   err.Error(),
			})
			c.Abort()

			// Still log failed attempt
			logUsage(tokenService, -1, c, startTime, http.StatusUnauthorized, err.Error())
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
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Message: "Rate limit check failed",
				Error:   err.Error(),
			})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, models.ErrorResponse{
				Success: false,
				Message: message,
				Error:   "Please slow down your requests",
			})
			c.Abort()

			// Log rate limit exceeded
			logUsage(tokenService, token.ID, c, startTime, http.StatusTooManyRequests, message)
			return
		}

		// Store token info in context for handlers
		c.Set("token_id", token.ID)
		c.Set("token_name", token.Name)
		c.Set("token_scopes", token.Scopes)
		// Vendor filter context – read by handlers to scope DB queries
		c.Set("token_is_super", token.IsSuperToken)
		c.Set("token_vendor_name", token.VendorName)
		c.Set("token_filter_column", token.FilterColumn)
		c.Set("token_filter_value", token.FilterValue)

		// Process request
		c.Next()

		// Log successful request after processing
		statusCode := c.Writer.Status()
		logUsage(tokenService, token.ID, c, startTime, statusCode, "")
	}
}

// AdminAuthMiddleware validates admin session tokens
func AdminAuthMiddleware(tokenService *service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract session token from header or cookie
		sessionToken := c.GetHeader("X-Session-Token")
		if sessionToken == "" {
			// Try cookie
			sessionToken, _ = c.Cookie("session_token")
		}

		if sessionToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authentication required",
			})
			c.Abort()
			return
		}

		// Validate session
		admin, err := tokenService.ValidateSession(sessionToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired session",
			})
			c.Abort()
			return
		}

		// Store admin info in context
		c.Set("admin_id", admin.ID)
		c.Set("admin_username", admin.Username)
		c.Set("admin_role", admin.Role)

		c.Next()
	}
}

// RequireRole checks if admin has required role
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminRole, exists := c.Get("admin_role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Access denied",
			})
			c.Abort()
			return
		}

		roleStr := adminRole.(string)
		allowed := false
		for _, role := range allowedRoles {
			if roleStr == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// logUsage creates a usage log entry
func logUsage(tokenService *service.TokenService, tokenID int, c *gin.Context, startTime time.Time, statusCode int, errorMsg string) {
	// Generate request ID if not exists
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	// Calculate response time
	responseTimeMs := int(time.Since(startTime).Milliseconds())

	// Create usage log
	log := &models.TokenUsageLog{
		TokenID:        tokenID,
		Method:         c.Request.Method,
		Endpoint:       c.Request.URL.Path,
		FullURL:        c.Request.URL.String(),
		StatusCode:     statusCode,
		ResponseTimeMs: responseTimeMs,
		IPAddress:      c.ClientIP(),
		UserAgent:      c.Request.UserAgent(),
		Referer:        c.Request.Referer(),
		RequestID:      requestID,
		ErrorMessage:   errorMsg,
	}

	// Log asynchronously to avoid blocking request
	go tokenService.LogTokenUsage(log)
}

// CORSForAdmin configures CORS for admin dashboard
func CORSForAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Session-Token, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// ScopeChecker checks if token has required scope
func ScopeChecker(requiredScopes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenScopes, exists := c.Get("token_scopes")
		if !exists {
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Success: false,
				Message: "Access denied",
				Error:   "No scopes found for token",
			})
			c.Abort()
			return
		}

		// Parse scopes JSON
		scopesJSON := tokenScopes.(string)
		if scopesJSON == "" || scopesJSON == "[]" {
			// No scopes defined - allow all (backward compatibility)
			c.Next()
			return
		}

		// Check if token has required scopes
		// For now, simplified check - in production, parse JSON array
		for _, required := range requiredScopes {
			if !strings.Contains(scopesJSON, required) {
				c.JSON(http.StatusForbidden, models.ErrorResponse{
					Success: false,
					Message: "Insufficient permissions",
					Error:   "Token does not have required scope: " + required,
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
