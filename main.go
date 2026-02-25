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
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	logger.Info("Starting API Gateway...")

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	gin.SetMode(cfg.Server.GinMode)

	dbManager := database.NewDBManager(
		cfg.TicketDB.GetDSN(),
		cfg.MachineDB.GetDSN(),
		cfg.TokenDB.GetDSN(),
		logger,
	)
	defer dbManager.Close()

	// Initialize unified data repository (uses ticket_master; cross-db JOIN to machine_master)
	dataRepo := repository.NewDataRepository(dbManager.TicketDB, logger)
	dataService := service.NewDataService(dataRepo, logger)
	dataHandler := handlers.NewDataHandler(dataService, logger)

	healthHandler := handlers.NewHealthHandler(dbManager, logger)

	// Token management (optional — requires token DB)
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

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger))
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	routes.SetupRoutes(
		router,
		dataHandler,
		healthHandler,
		tokenHandler,
		tokenService,
		cfg.Security.APIKey,
	)

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
