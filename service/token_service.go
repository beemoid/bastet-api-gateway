package service

import (
	"api-gateway/models"
	"api-gateway/repository"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// TokenService handles business logic for token management
type TokenService struct {
	repo   *repository.TokenRepository
	logger *logrus.Logger
}

// NewTokenService creates a new token service instance
func NewTokenService(repo *repository.TokenRepository, logger *logrus.Logger) *TokenService {
	return &TokenService{
		repo:   repo,
		logger: logger,
	}
}

// ============================================================================
// Admin Authentication
// ============================================================================

// Login authenticates an admin user and creates a session
func (s *TokenService) Login(username, password, ipAddress, userAgent string) (*models.LoginResponse, error) {
	admin, err := s.repo.GetAdminByUsername(username)
	if err != nil {
		s.logger.Warnf("Login attempt failed for username: %s", username)
		return &models.LoginResponse{
			Success: false,
			Message: "Invalid username or password",
		}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password))
	if err != nil {
		s.logger.Warnf("Invalid password for username: %s", username)
		return &models.LoginResponse{
			Success: false,
			Message: "Invalid username or password",
		}, nil
	}

	sessionToken, err := s.generateSecureToken(64)
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %v", err)
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	session := &models.AdminSession{
		SessionToken: sessionToken,
		AdminUserID:  admin.ID,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		ExpiresAt:    expiresAt,
	}

	err = s.repo.CreateSession(session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	_ = s.repo.UpdateAdminLastLogin(admin.ID, ipAddress)

	s.logger.Infof("Admin user '%s' logged in successfully", username)

	return &models.LoginResponse{
		Success:      true,
		Message:      "Login successful",
		SessionToken: sessionToken,
		User:         admin,
		ExpiresAt:    expiresAt,
	}, nil
}

// Logout deletes a session
func (s *TokenService) Logout(sessionToken string) error {
	return s.repo.DeleteSession(sessionToken)
}

// ValidateSession validates a session token and returns the admin user
func (s *TokenService) ValidateSession(sessionToken string) (*models.AdminUser, error) {
	session, err := s.repo.GetSessionByToken(sessionToken)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired session")
	}

	_ = s.repo.UpdateSessionAccess(session.ID)

	admin, err := s.repo.GetAdminByID(session.AdminUserID)
	if err != nil {
		return nil, fmt.Errorf("admin user not found")
	}

	return admin, nil
}

// ============================================================================
// API Token Management
// ============================================================================

// CreateAPIToken creates a new API token
func (s *TokenService) CreateAPIToken(req *models.CreateTokenRequest, createdBy int) (*models.APIToken, error) {
	tokenValue, err := s.generateAPIToken(req.Environment)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	prefix := s.extractTokenPrefix(tokenValue)

	scopesJSON, _ := repository.ConvertToJSON(req.Scopes)
	ipWhitelistJSON, _ := repository.ConvertToJSON(req.IPWhitelist)
	allowedOriginsJSON, _ := repository.ConvertToJSON(req.AllowedOrigins)

	if req.RateLimitPerMinute == 0 {
		req.RateLimitPerMinute = 100
	}
	if req.RateLimitPerHour == 0 {
		req.RateLimitPerHour = 5000
	}
	if req.RateLimitPerDay == 0 {
		req.RateLimitPerDay = 100000
	}

	token := &models.APIToken{
		Token:              tokenValue,
		Name:               req.Name,
		Description:        req.Description,
		TokenPrefix:        prefix,
		Scopes:             scopesJSON,
		Environment:        req.Environment,
		IsActive:           true,
		IPWhitelist:        ipWhitelistJSON,
		AllowedOrigins:     allowedOriginsJSON,
		RateLimitPerMinute: req.RateLimitPerMinute,
		RateLimitPerHour:   req.RateLimitPerHour,
		RateLimitPerDay:    req.RateLimitPerDay,
	}

	if req.ExpiresAt != nil {
		token.ExpiresAt = models.NullTime{NullTime: sql.NullTime{Valid: true, Time: *req.ExpiresAt}}
	}

	id, err := s.repo.CreateAPIToken(token, createdBy)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %v", err)
	}
	token.ID = id

	newValuesJSON, _ := json.Marshal(map[string]interface{}{
		"name": token.Name, "environment": token.Environment, "scopes": req.Scopes,
	})
	_ = s.repo.CreateAuditLog(&models.AuditLog{
		AdminUserID: &createdBy, Action: "create_token",
		ResourceType: "token", ResourceID: &id,
		NewValues: string(newValuesJSON),
		Description: fmt.Sprintf("Created API token: %s", token.Name),
	})

	s.logger.Infof("Created new API token: %s (ID: %d)", token.Name, id)
	return token, nil
}

// GetAllTokens retrieves all API tokens (with masked token values)
func (s *TokenService) GetAllTokens() ([]*models.APIToken, error) {
	tokens, err := s.repo.GetAllAPITokens()
	if err != nil {
		return nil, err
	}
	for _, token := range tokens {
		token.SanitizeForList()
	}
	return tokens, nil
}

// GetTokenByID retrieves a token by ID (with masked token value)
func (s *TokenService) GetTokenByID(id int) (*models.APIToken, error) {
	token, err := s.repo.GetAPITokenByID(id)
	if err != nil {
		return nil, err
	}
	token.SanitizeForList()
	return token, nil
}

// UpdateToken updates an existing API token
func (s *TokenService) UpdateToken(id int, req *models.UpdateTokenRequest, updatedBy int) (*models.APIToken, error) {
	oldToken, err := s.repo.GetAPITokenByID(id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Scopes != nil {
		j, _ := repository.ConvertToJSON(req.Scopes)
		updates["scopes"] = j
	}
	if req.IPWhitelist != nil {
		j, _ := repository.ConvertToJSON(req.IPWhitelist)
		updates["ip_whitelist"] = j
	}
	if req.AllowedOrigins != nil {
		j, _ := repository.ConvertToJSON(req.AllowedOrigins)
		updates["allowed_origins"] = j
	}
	if req.RateLimitPerMinute != nil {
		updates["rate_limit_per_minute"] = *req.RateLimitPerMinute
	}
	if req.RateLimitPerHour != nil {
		updates["rate_limit_per_hour"] = *req.RateLimitPerHour
	}
	if req.RateLimitPerDay != nil {
		updates["rate_limit_per_day"] = *req.RateLimitPerDay
	}
	if req.ExpiresAt != nil {
		updates["expires_at"] = *req.ExpiresAt
	}

	if len(updates) == 0 {
		return s.GetTokenByID(id)
	}

	err = s.repo.UpdateAPIToken(id, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update token: %v", err)
	}

	oldJSON, _ := json.Marshal(map[string]string{"name": oldToken.Name})
	newJSON, _ := json.Marshal(updates)
	_ = s.repo.CreateAuditLog(&models.AuditLog{
		AdminUserID: &updatedBy, Action: "update_token",
		ResourceType: "token", ResourceID: &id,
		OldValues: string(oldJSON), NewValues: string(newJSON),
		Description: fmt.Sprintf("Updated API token: %s", oldToken.Name),
	})

	return s.GetTokenByID(id)
}

// DisableToken disables a token
func (s *TokenService) DisableToken(id int, disabledBy int) error {
	err := s.repo.DisableToken(id)
	if err != nil {
		return err
	}
	_ = s.repo.CreateAuditLog(&models.AuditLog{
		AdminUserID: &disabledBy, Action: "disable_token",
		ResourceType: "token", ResourceID: &id,
		Description: fmt.Sprintf("Disabled API token ID: %d", id),
	})
	return nil
}

// EnableToken enables a token
func (s *TokenService) EnableToken(id int, enabledBy int) error {
	err := s.repo.EnableToken(id)
	if err != nil {
		return err
	}
	_ = s.repo.CreateAuditLog(&models.AuditLog{
		AdminUserID: &enabledBy, Action: "enable_token",
		ResourceType: "token", ResourceID: &id,
		Description: fmt.Sprintf("Enabled API token ID: %d", id),
	})
	return nil
}

// DeleteToken deletes a token permanently
func (s *TokenService) DeleteToken(id int, deletedBy int) error {
	token, err := s.repo.GetAPITokenByID(id)
	if err != nil {
		return err
	}
	err = s.repo.DeleteToken(id)
	if err != nil {
		return err
	}
	_ = s.repo.CreateAuditLog(&models.AuditLog{
		AdminUserID: &deletedBy, Action: "delete_token",
		ResourceType: "token", ResourceID: &id,
		Description: fmt.Sprintf("Deleted API token: %s", token.Name),
	})
	return nil
}

// ============================================================================
// Token Validation & Usage Tracking
// ============================================================================

// ValidateAPIToken validates a token and checks all security constraints
func (s *TokenService) ValidateAPIToken(tokenValue, ipAddress string) (*models.APIToken, error) {
	token, err := s.repo.GetAPITokenByToken(tokenValue)
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	if !token.IsValid() {
		if token.IsRevoked() {
			return nil, fmt.Errorf("token has been revoked")
		}
		if token.IsExpired() {
			return nil, fmt.Errorf("token has expired")
		}
		return nil, fmt.Errorf("token is disabled")
	}

	if token.IPWhitelist != "" && token.IPWhitelist != "[]" {
		var whitelist []string
		if err := json.Unmarshal([]byte(token.IPWhitelist), &whitelist); err == nil && len(whitelist) > 0 {
			allowed := false
			for _, ip := range whitelist {
				if ipAddress == ip {
					allowed = true
					break
				}
			}
			if !allowed {
				return nil, fmt.Errorf("IP address not whitelisted")
			}
		}
	}

	return token, nil
}

// CheckRateLimit checks if token has exceeded rate limits
func (s *TokenService) CheckRateLimit(tokenID int, rateLimits map[string]int) (bool, string, error) {
	now := time.Now()

	checks := []struct {
		windowType string
		truncate   time.Duration
		duration   time.Duration
	}{
		{"minute", time.Minute, time.Minute},
		{"hour", time.Hour, time.Hour},
	}

	for _, check := range checks {
		limit, ok := rateLimits[check.windowType]
		if !ok || limit <= 0 {
			continue
		}

		windowStart := now.Truncate(check.truncate)
		windowEnd := windowStart.Add(check.duration)

		count, err := s.repo.GetRateLimitCount(tokenID, check.windowType, windowStart)
		if err != nil {
			return false, "", err
		}
		if count >= limit {
			return false, fmt.Sprintf("Rate limit exceeded (per %s)", check.windowType), nil
		}

		_ = s.repo.IncrementRateLimit(tokenID, check.windowType, windowStart, windowEnd)
	}

	// Day check
	if limit, ok := rateLimits["day"]; ok && limit > 0 {
		windowStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		windowEnd := windowStart.Add(24 * time.Hour)

		count, err := s.repo.GetRateLimitCount(tokenID, "day", windowStart)
		if err != nil {
			return false, "", err
		}
		if count >= limit {
			return false, "Rate limit exceeded (per day)", nil
		}
		_ = s.repo.IncrementRateLimit(tokenID, "day", windowStart, windowEnd)
	}

	return true, "", nil
}

// LogTokenUsage logs API token usage
func (s *TokenService) LogTokenUsage(log *models.TokenUsageLog) {
	if err := s.repo.CreateUsageLog(log); err != nil {
		s.logger.Errorf("Failed to log token usage: %v", err)
	}
	if err := s.repo.UpdateTokenUsage(log.TokenID, log.IPAddress, log.Endpoint); err != nil {
		s.logger.Warnf("Failed to update token usage: %v", err)
	}
}

// ============================================================================
// Analytics
// ============================================================================

// GetTokenAnalytics retrieves analytics for a specific token
func (s *TokenService) GetTokenAnalytics(tokenID int, days int) (*models.TokenAnalytics, error) {
	return s.repo.GetTokenAnalytics(tokenID, days)
}

// GetDashboardStats retrieves overall dashboard statistics
func (s *TokenService) GetDashboardStats() (*models.TokenDashboardStats, error) {
	stats, err := s.repo.GetDashboardStats()
	if err != nil {
		return nil, err
	}
	recentLogs, err := s.repo.GetRecentUsageLogs(20)
	if err == nil {
		stats.RecentActivity = recentLogs
	}
	return stats, nil
}

// GetEndpointStats retrieves endpoint statistics
func (s *TokenService) GetEndpointStats(days int, limit int) ([]*models.EndpointStats, error) {
	return s.repo.GetEndpointStats(days, limit)
}

// GetDailyUsage retrieves daily usage for charts
func (s *TokenService) GetDailyUsage(tokenID *int, days int) ([]*models.DailyUsage, error) {
	return s.repo.GetDailyUsage(tokenID, days)
}

// GetUsageLogsByTokenID retrieves usage logs for a specific token
func (s *TokenService) GetUsageLogsByTokenID(tokenID int, limit int) ([]*models.TokenUsageLog, error) {
	return s.repo.GetUsageLogsByTokenID(tokenID, limit)
}

// GetAuditLogs retrieves audit logs
func (s *TokenService) GetAuditLogs(limit int) ([]*models.AuditLog, error) {
	return s.repo.GetAuditLogs(limit)
}

// ============================================================================
// Helper Functions
// ============================================================================

func (s *TokenService) generateAPIToken(environment string) (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	tokenValue := base64.URLEncoding.EncodeToString(randomBytes)

	var prefix string
	switch environment {
	case "production":
		prefix = "tok_live"
	case "staging":
		prefix = "tok_stage"
	case "development":
		prefix = "tok_dev"
	case "test":
		prefix = "tok_test"
	default:
		prefix = "tok"
	}
	return fmt.Sprintf("%s_%s", prefix, tokenValue), nil
}

func (s *TokenService) generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (s *TokenService) extractTokenPrefix(token string) string {
	underscoreCount := 0
	for i, char := range token {
		if char == '_' {
			underscoreCount++
			if underscoreCount == 2 {
				return token[:i]
			}
		}
	}
	return "tok"
}

// HashPassword hashes a password using bcrypt
func (s *TokenService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
