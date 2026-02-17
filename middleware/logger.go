package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger creates a middleware that logs HTTP requests
// Logs method, path, status code, latency, and client IP
func Logger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Get request details
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		// Determine log level based on status code
		entry := logger.WithFields(logrus.Fields{
			"status":  statusCode,
			"method":  method,
			"path":    path,
			"ip":      clientIP,
			"latency": latency,
		})

		// Log based on status code
		switch {
		case statusCode >= 500:
			entry.Error("Server error")
		case statusCode >= 400:
			entry.Warn("Client error")
		default:
			entry.Info("Request processed")
		}
	}
}
