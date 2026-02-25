package routes

import (
	"api-gateway/handlers"
	"api-gateway/middleware"
	"api-gateway/service"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(
	router *gin.Engine,
	dataHandler *handlers.DataHandler,
	healthHandler *handlers.HealthHandler,
	tokenHandler *handlers.TokenHandler,
	tokenService *service.TokenService,
	apiKey string,
) {
	router.Use(middleware.CORS())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to API Gateway"})
	})

	// Health check endpoints (no authentication required)
	router.GET("/health", healthHandler.Check)
	router.GET("/ping", healthHandler.Ping)

	// Admin dashboard static files and pages
	router.Static("/admin/assets", "./templates/assets")
	router.LoadHTMLGlob("templates/*.html")

	router.GET("/admin", func(c *gin.Context) { c.HTML(200, "login.html", nil) })
	router.GET("/admin/login", func(c *gin.Context) { c.HTML(200, "login.html", nil) })
	router.GET("/admin/dashboard", func(c *gin.Context) { c.HTML(200, "dashboard.html", nil) })
	router.GET("/admin/tokens", func(c *gin.Context) { c.HTML(200, "dashboard.html", nil) })

	// Admin API routes (session-authenticated)
	if tokenHandler != nil && tokenService != nil {
		adminAPI := router.Group("/api/v1/admin")
		{
			adminAPI.POST("/auth/login", tokenHandler.Login)
			adminAPI.POST("/auth/logout", tokenHandler.Logout)

			protected := adminAPI.Group("")
			protected.Use(middleware.AdminAuthMiddleware(tokenService))
			{
				protected.GET("/auth/me", tokenHandler.GetCurrentUser)

				// Token management
				protected.GET("/tokens", tokenHandler.ListTokens)
				protected.POST("/tokens", tokenHandler.CreateToken)
				protected.GET("/tokens/:id", tokenHandler.GetToken)
				protected.PUT("/tokens/:id", tokenHandler.UpdateToken)
				protected.DELETE("/tokens/:id", tokenHandler.DeleteToken)
				protected.PATCH("/tokens/:id/disable", tokenHandler.DisableToken)
				protected.PATCH("/tokens/:id/enable", tokenHandler.EnableToken)
				protected.GET("/tokens/:id/logs", tokenHandler.GetTokenUsageLogs)

				// Analytics
				protected.GET("/analytics/dashboard", tokenHandler.GetDashboardStats)
				protected.GET("/analytics/tokens/:id", tokenHandler.GetTokenAnalytics)
				protected.GET("/analytics/endpoints", tokenHandler.GetEndpointStats)
				protected.GET("/analytics/daily", tokenHandler.GetDailyUsage)

				// Audit logs
				protected.GET("/audit-logs", tokenHandler.GetAuditLogs)
			}
		}
	}

	// ── Unified data endpoint (token-authenticated) ──────────────────────────
	api := router.Group("/api/v1")
	api.Use(middleware.CombinedAuth(tokenService))
	{
		data := api.Group("/data")
		{
			data.GET("", dataHandler.GetAll)
			data.GET("/metadata", dataHandler.GetMetadata)
			data.GET("/:terminal_id", dataHandler.GetByID)
			data.PUT("/:terminal_id", dataHandler.Update)
		}
	}
}
