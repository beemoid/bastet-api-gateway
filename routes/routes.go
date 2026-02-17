package routes

import (
	"api-gateway/handlers"
	"api-gateway/middleware"

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
	apiKey string,
) {
	// Apply global middleware
	router.Use(middleware.CORS())

	// Health check endpoints (no authentication required)
	router.GET("/health", healthHandler.Check)
	router.GET("/ping", healthHandler.Ping)

	// Swagger documentation endpoint (no authentication required)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes group
	api := router.Group("/api/v1")
	{
		// Apply API key authentication to all API routes
		if apiKey != "" {
			api.Use(middleware.APIKeyAuth(apiKey))
		}

		// Ticket routes
		tickets := api.Group("/tickets")
		{
			tickets.GET("", ticketHandler.GetAll)                         // GET /api/v1/tickets - list all tickets
			tickets.GET("/:id", ticketHandler.GetByID)                    // GET /api/v1/tickets/:id - get ticket by ID
			tickets.GET("/number/:number", ticketHandler.GetByNumber)     // GET /api/v1/tickets/number/:number - get ticket by number
			tickets.GET("/status/:status", ticketHandler.GetByStatus)     // GET /api/v1/tickets/status/:status - filter by status
			tickets.GET("/terminal/:terminal_id", ticketHandler.GetByTerminal) // GET /api/v1/tickets/terminal/:terminal_id - filter by terminal
			tickets.POST("", ticketHandler.Create)                        // POST /api/v1/tickets - create new ticket
			tickets.PUT("/:id", ticketHandler.Update)                     // PUT /api/v1/tickets/:id - update ticket
		}

		// Machine routes
		machines := api.Group("/machines")
		{
			machines.GET("", machineHandler.GetAll)                          // GET /api/v1/machines - list all machines
			machines.GET("/search", machineHandler.Search)                   // GET /api/v1/machines/search - search with filters
			machines.GET("/:terminal_id", machineHandler.GetByTerminalID)    // GET /api/v1/machines/:terminal_id - get machine by terminal ID
			machines.GET("/status/:status", machineHandler.GetByStatus)      // GET /api/v1/machines/status/:status - filter by status
			machines.GET("/branch/:branch_code", machineHandler.GetByBranch) // GET /api/v1/machines/branch/:branch_code - filter by branch
			machines.PATCH("/status", machineHandler.UpdateStatus)           // PATCH /api/v1/machines/status - update machine status
		}
	}
}
