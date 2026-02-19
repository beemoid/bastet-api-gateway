package routes

import (
	"api-gateway/handlers"
	"api-gateway/middleware"
	"api-gateway/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configures all API routes
// Groups routes by resource and applies appropriate middleware
func SetupRoutes(
	router *gin.Engine,
	ticketHandler *handlers.TicketHandler,
	machineHandler *handlers.MachineHandler,
	healthHandler *handlers.HealthHandler,
	tokenHandler *handlers.TokenHandler,
	tokenService *service.TokenService,
	apiKey string,
) {
	// Apply global middleware
	router.Use(middleware.CORS())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to API Gateway",
		})
	})

	// Health check endpoints (no authentication required)
	router.GET("/health", healthHandler.Check)
	router.GET("/ping", healthHandler.Ping)

	// Swagger documentation endpoint (no authentication required)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Serve admin dashboard static files
	router.Static("/admin/assets", "./templates/assets")
	router.LoadHTMLGlob("templates/*.html")

	// Admin dashboard pages (HTML)
	router.GET("/admin", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})
	router.GET("/admin/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})
	router.GET("/admin/dashboard", func(c *gin.Context) {
		c.HTML(200, "dashboard.html", nil)
	})
	router.GET("/admin/tokens", func(c *gin.Context) {
		c.HTML(200, "dashboard.html", nil)
	})

	// Admin routes (no API key required — uses session authentication)
	if tokenHandler != nil && tokenService != nil {
		adminAPI := router.Group("/api/v1/admin")
		{
			// Auth routes (no session required)
			adminAPI.POST("/auth/login", tokenHandler.Login)
			adminAPI.POST("/auth/logout", tokenHandler.Logout)

			// Protected admin routes (session required)
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

	// API v1 routes group (accepts X-API-Key or X-API-Token)
	api := router.Group("/api/v1")
	{
		api.Use(middleware.CombinedAuth(tokenService))

		// Ticket routes
		tickets := api.Group("/tickets")
		{
			tickets.GET("", ticketHandler.GetAll)
			tickets.GET("/metadata", ticketHandler.GetMetadata)
			tickets.GET("/:id", ticketHandler.GetByID)
			tickets.GET("/number/:number", ticketHandler.GetByNumber)
			tickets.GET("/status/:status", ticketHandler.GetByStatus)
			tickets.GET("/terminal/:terminal_id", ticketHandler.GetByTerminal)
			tickets.POST("", ticketHandler.Create)
			tickets.PUT("/:id", ticketHandler.Update)
		}

		// Machine routes
		machines := api.Group("/machines")
		{
			machines.GET("", machineHandler.GetAll)
			machines.GET("/metadata", machineHandler.GetMetadata)
			machines.GET("/search", machineHandler.Search)
			machines.GET("/:terminal_id", machineHandler.GetByTerminalID)
			machines.GET("/status/:status", machineHandler.GetByStatus)
			machines.GET("/branch/:branch_code", machineHandler.GetByBranch)
			machines.PATCH("/status", machineHandler.UpdateStatus)
		}
	}
}
