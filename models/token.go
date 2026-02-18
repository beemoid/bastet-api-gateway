package models

import (
	"time"
)

// ============================================================================
// Database Models
// ============================================================================

// AdminUser represents an admin user who can access the token management dashboard
type AdminUser struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username" binding:"required,min=3,max=100"`
	Email        string    `json:"email" db:"email" binding:"required,email"`
	PasswordHash string    `json:"-" db:"password_hash"` // Never expose password hash in JSON
	FullName     string    `json:"full_name,omitempty" db:"full_name"`
	Role         string    `json:"role" db:"role" binding:"required,oneof=super_admin admin viewer"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	LastLoginAt  NullTime  `json:"last_login_at,omitempty" db:"last_login_at"`
	LastLoginIP  string    `json:"last_login_ip,omitempty" db:"last_login_ip"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy    *int      `json:"created_by,omitempty" db:"created_by"`
}

// APIToken represents an API token with scopes, permissions, and analytics
type APIToken struct {
	ID          int       `json:"id" db:"id"`
	Token       string    `json:"token" db:"token"` // Only shown once during creation
	Name        string    `json:"name" db:"name" binding:"required,min=3,max=200"`
	Description string    `json:"description,omitempty" db:"description"`
	TokenPrefix string    `json:"token_prefix" db:"token_prefix"`

	// Permissions & Scopes (stored as JSON in database)
	Scopes      string `json:"scopes,omitempty" db:"scopes"`           // JSON array
	Permissions string `json:"permissions,omitempty" db:"permissions"` // JSON object

	// Environment & Status
	Environment string `json:"environment" db:"environment" binding:"required,oneof=production staging development test"`
	IsActive    bool   `json:"is_active" db:"is_active"`

	// Security
	IPWhitelist    string `json:"ip_whitelist,omitempty" db:"ip_whitelist"`       // JSON array
	AllowedOrigins string `json:"allowed_origins,omitempty" db:"allowed_origins"` // JSON array

	// Rate Limiting
	RateLimitPerMinute int `json:"rate_limit_per_minute" db:"rate_limit_per_minute"`
	RateLimitPerHour   int `json:"rate_limit_per_hour" db:"rate_limit_per_hour"`
	RateLimitPerDay    int `json:"rate_limit_per_day" db:"rate_limit_per_day"`

	// Expiration
	ExpiresAt NullTime `json:"expires_at,omitempty" db:"expires_at"`

	// Usage Statistics
	LastUsedAt       NullTime `json:"last_used_at,omitempty" db:"last_used_at"`
	LastUsedIP       string   `json:"last_used_ip,omitempty" db:"last_used_ip"`
	LastUsedEndpoint string   `json:"last_used_endpoint,omitempty" db:"last_used_endpoint"`
	TotalRequests    int64    `json:"total_requests" db:"total_requests"`

	// Metadata
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy    *int      `json:"created_by,omitempty" db:"created_by"`
	RevokedAt    NullTime  `json:"revoked_at,omitempty" db:"revoked_at"`
	RevokedBy    *int      `json:"revoked_by,omitempty" db:"revoked_by"`
	RevokedReason string   `json:"revoked_reason,omitempty" db:"revoked_reason"`
}

// TokenUsageLog tracks every API request for analytics and audit
type TokenUsageLog struct {
	ID      int64 `json:"id" db:"id"`
	TokenID int   `json:"token_id" db:"token_id"`

	// Request Details
	Method   string `json:"method" db:"method"`
	Endpoint string `json:"endpoint" db:"endpoint"`
	FullURL  string `json:"full_url,omitempty" db:"full_url"`

	// Response Details
	StatusCode     int `json:"status_code" db:"status_code"`
	ResponseTimeMs int `json:"response_time_ms,omitempty" db:"response_time_ms"`

	// Request Context
	IPAddress string `json:"ip_address" db:"ip_address"`
	UserAgent string `json:"user_agent,omitempty" db:"user_agent"`
	Referer   string `json:"referer,omitempty" db:"referer"`

	// Metadata
	RequestID        string `json:"request_id,omitempty" db:"request_id"`
	RequestBodySize  int    `json:"request_body_size,omitempty" db:"request_body_size"`
	ResponseBodySize int    `json:"response_body_size,omitempty" db:"response_body_size"`

	// Error Tracking
	ErrorMessage string `json:"error_message,omitempty" db:"error_message"`
	ErrorCode    string `json:"error_code,omitempty" db:"error_code"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// TokenRateLimit tracks rate limit counters for each token
type TokenRateLimit struct {
	ID      int64 `json:"id" db:"id"`
	TokenID int   `json:"token_id" db:"token_id"`

	// Time Windows
	WindowType  string    `json:"window_type" db:"window_type"` // 'minute', 'hour', 'day'
	WindowStart time.Time `json:"window_start" db:"window_start"`
	WindowEnd   time.Time `json:"window_end" db:"window_end"`

	// Counters
	RequestCount int `json:"request_count" db:"request_count"`

	// Metadata
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// AdminSession manages admin user sessions for dashboard authentication
type AdminSession struct {
	ID             int64     `json:"id" db:"id"`
	SessionToken   string    `json:"session_token" db:"session_token"`
	AdminUserID    int       `json:"admin_user_id" db:"admin_user_id"`
	IPAddress      string    `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent      string    `json:"user_agent,omitempty" db:"user_agent"`
	ExpiresAt      time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	LastAccessedAt time.Time `json:"last_accessed_at" db:"last_accessed_at"`
}

// AuditLog tracks all administrative actions for compliance and security
type AuditLog struct {
	ID          int64  `json:"id" db:"id"`
	AdminUserID *int   `json:"admin_user_id,omitempty" db:"admin_user_id"`
	Action      string `json:"action" db:"action"`
	ResourceType string `json:"resource_type" db:"resource_type"`
	ResourceID  *int   `json:"resource_id,omitempty" db:"resource_id"`
	OldValues   string `json:"old_values,omitempty" db:"old_values"` // JSON
	NewValues   string `json:"new_values,omitempty" db:"new_values"` // JSON
	IPAddress   string `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent   string `json:"user_agent,omitempty" db:"user_agent"`
	Description string `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// ============================================================================
// Request/Response Models
// ============================================================================

// LoginRequest represents admin login credentials
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse contains session token and user info
type LoginResponse struct {
	Success      bool       `json:"success"`
	Message      string     `json:"message"`
	SessionToken string     `json:"session_token,omitempty"`
	User         *AdminUser `json:"user,omitempty"`
	ExpiresAt    time.Time  `json:"expires_at,omitempty"`
}

// CreateTokenRequest represents a request to create a new API token
type CreateTokenRequest struct {
	Name               string   `json:"name" binding:"required,min=3,max=200"`
	Description        string   `json:"description"`
	Environment        string   `json:"environment" binding:"required,oneof=production staging development test"`
	Scopes             []string `json:"scopes"`             // Will be converted to JSON
	IPWhitelist        []string `json:"ip_whitelist"`       // Will be converted to JSON
	AllowedOrigins     []string `json:"allowed_origins"`    // Will be converted to JSON
	RateLimitPerMinute int      `json:"rate_limit_per_minute"`
	RateLimitPerHour   int      `json:"rate_limit_per_hour"`
	RateLimitPerDay    int      `json:"rate_limit_per_day"`
	ExpiresAt          *time.Time `json:"expires_at"`
}

// CreateTokenResponse contains the newly created token (only shown once)
type CreateTokenResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Token   *APIToken `json:"token,omitempty"`
	Warning string    `json:"warning,omitempty"` // "Save this token - it won't be shown again"
}

// UpdateTokenRequest represents a request to update an existing token
type UpdateTokenRequest struct {
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	Scopes             []string   `json:"scopes"`
	IPWhitelist        []string   `json:"ip_whitelist"`
	AllowedOrigins     []string   `json:"allowed_origins"`
	RateLimitPerMinute *int       `json:"rate_limit_per_minute"`
	RateLimitPerHour   *int       `json:"rate_limit_per_hour"`
	RateLimitPerDay    *int       `json:"rate_limit_per_day"`
	ExpiresAt          *time.Time `json:"expires_at"`
}

// TokenListResponse contains a list of tokens (without full token value)
type TokenListResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    []APIToken  `json:"data"`
	Total   int         `json:"total"`
}

// TokenResponse contains a single token response
type TokenResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    *APIToken `json:"data,omitempty"`
}

// TokenAnalyticsRequest filters for analytics queries
type TokenAnalyticsRequest struct {
	TokenID   *int       `form:"token_id"`
	StartDate *time.Time `form:"start_date"`
	EndDate   *time.Time `form:"end_date"`
	Limit     int        `form:"limit"`
}

// TokenAnalytics contains detailed usage statistics for a token
type TokenAnalytics struct {
	TokenID            int       `json:"token_id"`
	TokenName          string    `json:"token_name"`
	TotalRequests      int64     `json:"total_requests"`
	SuccessfulRequests int64     `json:"successful_requests"`
	FailedRequests     int64     `json:"failed_requests"`
	ClientErrors       int64     `json:"client_errors"`
	ServerErrors       int64     `json:"server_errors"`
	AvgResponseTimeMs  float64   `json:"avg_response_time_ms"`
	MaxResponseTimeMs  int       `json:"max_response_time_ms"`
	UniqueIPs          int       `json:"unique_ips"`
	UniqueEndpoints    int       `json:"unique_endpoints"`
	LastUsedAt         NullTime  `json:"last_used_at"`
}

// TokenDashboardStats contains overall token system statistics
type TokenDashboardStats struct {
	TotalTokens        int     `json:"total_tokens"`
	ActiveTokens       int     `json:"active_tokens"`
	TotalRequests24h   int64   `json:"total_requests_24h"`
	SuccessRate        float64 `json:"success_rate"`
	AvgResponseTimeMs  float64 `json:"avg_response_time_ms"`
	TopTokens          []TokenAnalytics  `json:"top_tokens"`
	RecentActivity     []*TokenUsageLog  `json:"recent_activity"`
}

// EndpointStats contains statistics per endpoint
type EndpointStats struct {
	Endpoint           string  `json:"endpoint"`
	Method             string  `json:"method"`
	RequestCount       int64   `json:"request_count"`
	UniqueTokens       int     `json:"unique_tokens"`
	AvgResponseTimeMs  float64 `json:"avg_response_time_ms"`
	SuccessfulRequests int64   `json:"successful_requests"`
	FailedRequests     int64   `json:"failed_requests"`
}

// DailyUsage contains daily aggregated usage data
type DailyUsage struct {
	Date               string  `json:"date"`
	TokenID            int     `json:"token_id"`
	TokenName          string  `json:"token_name"`
	RequestCount       int64   `json:"request_count"`
	SuccessfulRequests int64   `json:"successful_requests"`
	FailedRequests     int64   `json:"failed_requests"`
	AvgResponseTimeMs  float64 `json:"avg_response_time_ms"`
}

// ============================================================================
// Helper Functions
// ============================================================================

// MaskToken returns a masked version of the token for display (shows only prefix and last 4 chars)
func (t *APIToken) MaskToken() string {
	if len(t.Token) <= 12 {
		return t.TokenPrefix + "_" + "****"
	}
	return t.TokenPrefix + "_****" + t.Token[len(t.Token)-4:]
}

// IsExpired checks if the token has expired
func (t *APIToken) IsExpired() bool {
	if !t.ExpiresAt.Valid {
		return false // No expiration set
	}
	return time.Now().After(t.ExpiresAt.Time)
}

// IsRevoked checks if the token has been revoked
func (t *APIToken) IsRevoked() bool {
	return t.RevokedAt.Valid
}

// IsValid checks if token is active, not expired, and not revoked
func (t *APIToken) IsValid() bool {
	return t.IsActive && !t.IsExpired() && !t.IsRevoked()
}

// SanitizeForList removes sensitive data for list responses
func (t *APIToken) SanitizeForList() {
	if len(t.Token) > 0 {
		t.Token = t.MaskToken()
	}
}
