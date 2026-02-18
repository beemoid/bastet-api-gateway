package main

import (
	"api-gateway/config"
	"api-gateway/database"
	"api-gateway/handlers"
	"api-gateway/middleware"
	"api-gateway/repository"
	"api-gateway/routes"
	"api-gateway/service"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	_ "api-gateway/docs" // Import generated swagger docs
)

// @title API Gateway for On-Premise to Cloud Communication
// @version 1.0
// @description This API Gateway serves as middleware between on-premise databases and cloud applications, providing RESTful APIs for ticket and machine management.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Proprietary
// @license.url http://www.example.com/license

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @description API key for authentication. Required for all endpoints except /health and /ping.

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	logger.Info("Starting API Gateway...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode based on configuration
	gin.SetMode(cfg.Server.GinMode)

	// Initialize database connections (non-fatal: app keeps running if a DB is unavailable)
	dbManager := database.NewDBManager(
		cfg.TicketDB.GetDSN(),
		cfg.MachineDB.GetDSN(),
		cfg.TokenDB.GetDSN(),
		logger,
	)
	defer dbManager.Close()

	// Initialize repositories
	ticketRepo := repository.NewTicketRepository(dbManager.TicketDB, logger)
	machineRepo := repository.NewMachineRepository(dbManager.MachineDB, logger)

	// Initialize services
	ticketService := service.NewTicketService(ticketRepo, logger)
	machineService := service.NewMachineService(machineRepo, logger)

	// Initialize handlers
	ticketHandler := handlers.NewTicketHandler(ticketService, logger)
	machineHandler := handlers.NewMachineHandler(machineService, logger)
	healthHandler := handlers.NewHealthHandler(dbManager, logger)

	// Initialize token management (if token DB is available)
	var tokenHandler *handlers.TokenHandler
	var tokenService *service.TokenService

	if dbManager.TokenDB != nil {
		tokenRepo := repository.NewTokenRepository(dbManager.TokenDB, logger)
		tokenService = service.NewTokenService(tokenRepo, logger)
		tokenHandler = handlers.NewTokenHandler(tokenService, logger)
		logger.Info("Token management system initialized")
	} else {
		logger.Warn("Token management system not available (no database connection)")
	}

	// Create Gin router
	router := gin.New()

	// Apply global middleware
	router.Use(gin.Recovery())                       // Recover from panics
	router.Use(middleware.Logger(logger))             // Custom logger middleware
	router.Use(gzip.Gzip(gzip.DefaultCompression))   // Compress responses (1-5MB → ~200-500KB)

	// Setup routes
	routes.SetupRoutes(
		router,
		ticketHandler,
		machineHandler,
		healthHandler,
		tokenHandler,
		tokenService,
		cfg.Security.APIKey,
	)

	// Setup graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logger.Info("Shutting down API Gateway...")

		if err := dbManager.Close(); err != nil {
			logger.Errorf("Error during shutdown: %v", err)
		}

		os.Exit(0)
	}()

	// Start server
	address := fmt.Sprintf(":%s", cfg.Server.Port)
	logger.Infof("API Gateway listening on %s", address)
	logger.Infof("Environment: %s", cfg.Server.GinMode)
	if tokenHandler != nil {
		logger.Infof("Admin Dashboard: http://localhost:%s/admin", cfg.Server.Port)
	}

	if err := router.Run(address); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
