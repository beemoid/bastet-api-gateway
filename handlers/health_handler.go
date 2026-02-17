package handlers

import (
	"api-gateway/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	dbManager *database.DBManager
	logger    *logrus.Logger
}

// NewHealthHandler creates a new health handler instance
func NewHealthHandler(dbManager *database.DBManager, logger *logrus.Logger) *HealthHandler {
	return &HealthHandler{
		dbManager: dbManager,
		logger:    logger,
	}
}

// Check handles GET /health - performs health check on the API and databases
// @Summary Health check
// @Description Check the health status of the API and database connections
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "API is healthy"
// @Failure 503 {object} map[string]interface{} "Service unavailable"
// @Router /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	// Check database health
	err := h.dbManager.HealthCheck()
	if err != nil {
		h.logger.Errorf("Health check failed: %v", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "unhealthy",
			"message": "Database connection failed",
			"error":   err.Error(),
		})
		return
	}

	// All checks passed
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "API Gateway is running",
		"services": gin.H{
			"ticket_database":  "connected",
			"machine_database": "connected",
		},
	})
}

// Ping handles GET /ping - simple ping endpoint
// @Summary Ping
// @Description Simple ping endpoint to check if the API is running
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "pong"
// @Router /ping [get]
func (h *HealthHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
