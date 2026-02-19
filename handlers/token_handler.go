package handlers

import (
	"api-gateway/models"
	"api-gateway/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// TokenHandler handles HTTP requests for token management
type TokenHandler struct {
	service *service.TokenService
	logger  *logrus.Logger
}

// NewTokenHandler creates a new token handler instance
func NewTokenHandler(service *service.TokenService, logger *logrus.Logger) *TokenHandler {
	return &TokenHandler{
		service: service,
		logger:  logger,
	}
}

// ============================================================================
// Admin Authentication Endpoints
// ============================================================================

// Login handles POST /api/v1/admin/auth/login
// @Summary Admin Login
// @Description Authenticate admin user and return session token
// @Tags Admin Auth
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Login Credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.LoginResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/auth/login [post]
func (h *TokenHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	resp, err := h.service.Login(req.Username, req.Password, ipAddress, userAgent)
	if err != nil {
		h.logger.Errorf("Login error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Login failed",
		})
		return
	}

	if resp.Success {
		// Set session cookie
		c.SetCookie("session_token", resp.SessionToken, 86400, "/", "", false, true)
	}

	status := http.StatusOK
	if !resp.Success {
		status = http.StatusUnauthorized
	}
	c.JSON(status, resp)
}

// Logout handles POST /api/v1/admin/auth/logout
// @Summary Admin Logout
// @Description Invalidate current session
// @Tags Admin Auth
// @Accept json
// @Produce json
// @Param X-Session-Token header string false "Session Token"
// @Success 200 {object} map[string]interface{}
// @Router /admin/auth/logout [post]
func (h *TokenHandler) Logout(c *gin.Context) {
	sessionToken := c.GetHeader("X-Session-Token")
	if sessionToken == "" {
		sessionToken, _ = c.Cookie("session_token")
	}

	if sessionToken != "" {
		_ = h.service.Logout(sessionToken)
	}

	// Clear cookie
	c.SetCookie("session_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}

// GetCurrentUser handles GET /api/v1/admin/auth/me
// @Summary Get Current User
// @Description Get details of currently logged in admin
// @Tags Admin Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /admin/auth/me [get]
func (h *TokenHandler) GetCurrentUser(c *gin.Context) {
	adminID, _ := c.Get("admin_id")
	adminUsername, _ := c.Get("admin_username")
	adminRole, _ := c.Get("admin_role")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":       adminID,
			"username": adminUsername,
			"role":     adminRole,
		},
	})
}

// ============================================================================
// Token CRUD Endpoints
// ============================================================================

// ListTokens handles GET /api/v1/admin/tokens
// @Summary List API Tokens
// @Description Get all API tokens
// @Tags Token Management
// @Accept json
// @Produce json
// @Success 200 {object} models.TokenListResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/tokens [get]
func (h *TokenHandler) ListTokens(c *gin.Context) {
	tokens, err := h.service.GetAllTokens()
	if err != nil {
		h.logger.Errorf("Error listing tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to list tokens",
		})
		return
	}

	if tokens == nil {
		tokens = []*models.APIToken{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tokens retrieved successfully",
		"data":    tokens,
		"total":   len(tokens),
	})
}

// CreateToken handles POST /api/v1/admin/tokens
// @Summary Create API Token
// @Description Create a new API token
// @Tags Token Management
// @Accept json
// @Produce json
// @Param token body models.CreateTokenRequest true "Token Details"
// @Success 201 {object} models.CreateTokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/tokens [post]
func (h *TokenHandler) CreateToken(c *gin.Context) {
	var req models.CreateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	adminID := c.GetInt("admin_id")

	token, err := h.service.CreateAPIToken(&req, adminID)
	if err != nil {
		h.logger.Errorf("Error creating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create token: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Token created successfully",
		"data":    token,
		"warning": "Save this token securely - it won't be shown again!",
	})
}

// GetToken handles GET /api/v1/admin/tokens/:id
// @Summary Get API Token
// @Description Get details of a specific API token
// @Tags Token Management
// @Accept json
// @Produce json
// @Param id path int true "Token ID"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/tokens/{id} [get]
func (h *TokenHandler) GetToken(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid token ID",
		})
		return
	}

	token, err := h.service.GetTokenByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Token not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    token,
	})
}

// UpdateToken handles PUT /api/v1/admin/tokens/:id
// @Summary Update API Token
// @Description Update details of an existing API token
// @Tags Token Management
// @Accept json
// @Produce json
// @Param id path int true "Token ID"
// @Param token body models.UpdateTokenRequest true "Update Details"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/tokens/{id} [put]
func (h *TokenHandler) UpdateToken(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid token ID",
		})
		return
	}

	var req models.UpdateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	adminID := c.GetInt("admin_id")

	token, err := h.service.UpdateToken(id, &req, adminID)
	if err != nil {
		h.logger.Errorf("Error updating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token updated successfully",
		"data":    token,
	})
}

// DeleteToken handles DELETE /api/v1/admin/tokens/:id
// @Summary Delete API Token
// @Description Permanently delete an API token
// @Tags Token Management
// @Accept json
// @Produce json
// @Param id path int true "Token ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/tokens/{id} [delete]
func (h *TokenHandler) DeleteToken(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid token ID",
		})
		return
	}

	adminID := c.GetInt("admin_id")

	err = h.service.DeleteToken(id, adminID)
	if err != nil {
		h.logger.Errorf("Error deleting token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token deleted successfully",
	})
}

// DisableToken handles PATCH /api/v1/admin/tokens/:id/disable
// @Summary Disable API Token
// @Description Temporarily disable an API token
// @Tags Token Management
// @Accept json
// @Produce json
// @Param id path int true "Token ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/tokens/{id}/disable [patch]
func (h *TokenHandler) DisableToken(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid token ID",
		})
		return
	}

	adminID := c.GetInt("admin_id")

	err = h.service.DisableToken(id, adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to disable token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token disabled successfully",
	})
}

// EnableToken handles PATCH /api/v1/admin/tokens/:id/enable
// @Summary Enable API Token
// @Description Re-enable a disabled API token
// @Tags Token Management
// @Accept json
// @Produce json
// @Param id path int true "Token ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/tokens/{id}/enable [patch]
func (h *TokenHandler) EnableToken(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid token ID",
		})
		return
	}

	adminID := c.GetInt("admin_id")

	err = h.service.EnableToken(id, adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to enable token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token enabled successfully",
	})
}

// ============================================================================
// Analytics Endpoints
// ============================================================================

// GetDashboardStats handles GET /api/v1/admin/analytics/dashboard
// @Summary Get Dashboard Stats
// @Description Get overview statistics for the admin dashboard
// @Tags Analytics
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/analytics/dashboard [get]
func (h *TokenHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.service.GetDashboardStats()
	if err != nil {
		h.logger.Errorf("Error getting dashboard stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get dashboard stats",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// GetTokenAnalytics handles GET /api/v1/admin/analytics/tokens/:id
// @Summary Get Token Analytics
// @Description Get detailed analytics for a specific token
// @Tags Analytics
// @Accept json
// @Produce json
// @Param id path int true "Token ID"
// @Param days query int false "Number of days (default 7)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/analytics/tokens/{id} [get]
func (h *TokenHandler) GetTokenAnalytics(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid token ID",
		})
		return
	}

	days := 7
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	analytics, err := h.service.GetTokenAnalytics(id, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get token analytics",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    analytics,
	})
}

// GetEndpointStats handles GET /api/v1/admin/analytics/endpoints
// @Summary Get Endpoint Stats
// @Description Get usage statistics by endpoint
// @Tags Analytics
// @Accept json
// @Produce json
// @Param days query int false "Number of days (default 7)"
// @Param limit query int false "Limit results (default 20)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/analytics/endpoints [get]
func (h *TokenHandler) GetEndpointStats(c *gin.Context) {
	days := 7
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	stats, err := h.service.GetEndpointStats(days, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get endpoint stats",
		})
		return
	}

	if stats == nil {
		stats = []*models.EndpointStats{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// GetDailyUsage handles GET /api/v1/admin/analytics/daily
// @Summary Get Daily Usage
// @Description Get daily request volume
// @Tags Analytics
// @Accept json
// @Produce json
// @Param days query int false "Number of days (default 30)"
// @Param token_id query int false "Filter by Token ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/analytics/daily [get]
func (h *TokenHandler) GetDailyUsage(c *gin.Context) {
	days := 30
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	var tokenID *int
	if t := c.Query("token_id"); t != "" {
		if parsed, err := strconv.Atoi(t); err == nil {
			tokenID = &parsed
		}
	}

	usage, err := h.service.GetDailyUsage(tokenID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get daily usage",
		})
		return
	}

	if usage == nil {
		usage = []*models.DailyUsage{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    usage,
	})
}

// GetTokenUsageLogs handles GET /api/v1/admin/tokens/:id/logs
// @Summary Get Token Logs
// @Description Get access logs for a specific token
// @Tags Analytics
// @Accept json
// @Produce json
// @Param id path int true "Token ID"
// @Param limit query int false "Limit results (default 100)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/tokens/{id}/logs [get]
func (h *TokenHandler) GetTokenUsageLogs(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid token ID",
		})
		return
	}

	limit := 100
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	logs, err := h.service.GetUsageLogsByTokenID(id, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get usage logs",
		})
		return
	}

	if logs == nil {
		logs = []*models.TokenUsageLog{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    logs,
		"total":   len(logs),
	})
}

// GetAuditLogs handles GET /api/v1/admin/audit-logs
// @Summary Get Audit Logs
// @Description Get administrative audit logs
// @Tags Admin Auth
// @Accept json
// @Produce json
// @Param limit query int false "Limit results (default 100)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/audit-logs [get]
func (h *TokenHandler) GetAuditLogs(c *gin.Context) {
	limit := 100
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	logs, err := h.service.GetAuditLogs(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get audit logs",
		})
		return
	}

	if logs == nil {
		logs = []*models.AuditLog{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    logs,
		"total":   len(logs),
	})
}
