package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS returns a CORS middleware with secure defaults
// Allows the cloud app to make cross-origin requests to this API
func CORS() gin.HandlerFunc {
	config := cors.Config{
		// Allow all origins in development, restrict in production
		AllowOrigins: []string{"*"}, // TODO: Replace with actual cloud app URL in production

		// Allow common HTTP methods
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},

		// Allow common headers plus custom auth headers
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-API-Key",
			"X-API-Token",
			"X-Session-Token",
		},

		// Expose custom headers to the client
		ExposeHeaders: []string{"Content-Length"},

		// Allow credentials (cookies, authorization headers)
		AllowCredentials: true,

		// Cache preflight requests for 12 hours
		MaxAge: 12 * time.Hour,
	}

	return cors.New(config)
}
