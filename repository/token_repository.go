package repository

import (
	"api-gateway/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// TokenRepository handles database operations for token management
type TokenRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewTokenRepository creates a new token repository instance
func NewTokenRepository(db *sql.DB, logger *logrus.Logger) *TokenRepository {
	return &TokenRepository{
		db:     db,
		logger: logger,
	}
}

// ============================================================================
// Admin User Operations
// ============================================================================

// GetAdminByUsername retrieves an admin user by username
func (r *TokenRepository) GetAdminByUsername(username string) (*models.AdminUser, error) {
	query := `
		SELECT id, username, email, password_hash, ISNULL(full_name, '') as full_name,
		       role, is_active, last_login_at, ISNULL(last_login_ip, '') as last_login_ip,
		       created_at, updated_at, created_by
		FROM admin_users
		WHERE username = @p1 AND is_active = 1
	`
	row := r.db.QueryRow(query, username)

	var admin models.AdminUser
	var createdBy sql.NullInt64
	err := row.Scan(
		&admin.ID, &admin.Username, &admin.Email, &admin.PasswordHash,
		&admin.FullName, &admin.Role, &admin.IsActive, &admin.LastLoginAt,
		&admin.LastLoginIP, &admin.CreatedAt, &admin.UpdatedAt, &createdBy,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("admin user not found")
		}
		return nil, err
	}
	if createdBy.Valid {
		v := int(createdBy.Int64)
		admin.CreatedBy = &v
	}
	return &admin, nil
}

// GetAdminByID retrieves an admin user by ID
func (r *TokenRepository) GetAdminByID(id int) (*models.AdminUser, error) {
	query := `
		SELECT id, username, email, password_hash, ISNULL(full_name, '') as full_name,
		       role, is_active, last_login_at, ISNULL(last_login_ip, '') as last_login_ip,
		       created_at, updated_at, created_by
		FROM admin_users
		WHERE id = @p1
	`
	row := r.db.QueryRow(query, id)

	var admin models.AdminUser
	var createdBy sql.NullInt64
	err := row.Scan(
		&admin.ID, &admin.Username, &admin.Email, &admin.PasswordHash,
		&admin.FullName, &admin.Role, &admin.IsActive, &admin.LastLoginAt,
		&admin.LastLoginIP, &admin.CreatedAt, &admin.UpdatedAt, &createdBy,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("admin user not found")
		}
		return nil, err
	}
	if createdBy.Valid {
		v := int(createdBy.Int64)
		admin.CreatedBy = &v
	}
	return &admin, nil
}

// UpdateAdminLastLogin updates the last login timestamp and IP
func (r *TokenRepository) UpdateAdminLastLogin(adminID int, ipAddress string) error {
	query := `UPDATE admin_users SET last_login_at = GETDATE(), last_login_ip = @p1 WHERE id = @p2`
	_, err := r.db.Exec(query, ipAddress, adminID)
	return err
}

// ============================================================================
// Session Operations
// ============================================================================

// CreateSession creates a new admin session
func (r *TokenRepository) CreateSession(session *models.AdminSession) error {
	query := `
		INSERT INTO admin_sessions (session_token, admin_user_id, ip_address, user_agent, expires_at)
		VALUES (@p1, @p2, @p3, @p4, @p5)
	`
	_, err := r.db.Exec(query,
		session.SessionToken, session.AdminUserID,
		session.IPAddress, session.UserAgent, session.ExpiresAt,
	)
	return err
}

// GetSessionByToken retrieves a session by token
func (r *TokenRepository) GetSessionByToken(token string) (*models.AdminSession, error) {
	query := `
		SELECT id, session_token, admin_user_id, ISNULL(ip_address, '') as ip_address,
		       ISNULL(user_agent, '') as user_agent, expires_at, created_at, last_accessed_at
		FROM admin_sessions
		WHERE session_token = @p1 AND expires_at > GETDATE()
	`
	row := r.db.QueryRow(query, token)

	var session models.AdminSession
	err := row.Scan(
		&session.ID, &session.SessionToken, &session.AdminUserID,
		&session.IPAddress, &session.UserAgent, &session.ExpiresAt,
		&session.CreatedAt, &session.LastAccessedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found or expired")
		}
		return nil, err
	}
	return &session, nil
}

// UpdateSessionAccess updates the last accessed timestamp
func (r *TokenRepository) UpdateSessionAccess(sessionID int64) error {
	query := `UPDATE admin_sessions SET last_accessed_at = GETDATE() WHERE id = @p1`
	_, err := r.db.Exec(query, sessionID)
	return err
}

// DeleteSession deletes a session (logout)
func (r *TokenRepository) DeleteSession(token string) error {
	query := `DELETE FROM admin_sessions WHERE session_token = @p1`
	_, err := r.db.Exec(query, token)
	return err
}

// ============================================================================
// API Token Operations
// ============================================================================

// CreateAPIToken creates a new API token and returns its ID
func (r *TokenRepository) CreateAPIToken(token *models.APIToken, createdBy int) (int, error) {
	query := `
		INSERT INTO api_tokens (
			token, name, description, token_prefix, scopes, permissions,
			environment, is_active, ip_whitelist, allowed_origins,
			rate_limit_per_minute, rate_limit_per_hour, rate_limit_per_day,
			expires_at, created_by,
			vendor_name, filter_column, filter_value, is_super_token
		)
		OUTPUT INSERTED.id
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10,
		        @p11, @p12, @p13, @p14, @p15,
		        @p16, @p17, @p18, @p19)
	`

	var expiresAt interface{}
	if token.ExpiresAt.Valid {
		expiresAt = token.ExpiresAt.Time
	}

	var vendorName, filterColumn, filterValue interface{}
	if token.VendorName != "" {
		vendorName = token.VendorName
	}
	if token.FilterColumn != "" {
		filterColumn = token.FilterColumn
	}
	if token.FilterValue != "" {
		filterValue = token.FilterValue
	}

	var id int
	err := r.db.QueryRow(query,
		token.Token, token.Name, token.Description, token.TokenPrefix,
		token.Scopes, token.Permissions, token.Environment, token.IsActive,
		token.IPWhitelist, token.AllowedOrigins,
		token.RateLimitPerMinute, token.RateLimitPerHour, token.RateLimitPerDay,
		expiresAt, createdBy,
		vendorName, filterColumn, filterValue, token.IsSuperToken,
	).Scan(&id)

	return id, err
}

// scanToken scans a row into an APIToken struct
func (r *TokenRepository) scanToken(row interface{ Scan(dest ...interface{}) error }) (*models.APIToken, error) {
	var t models.APIToken
	var description, scopes, permissions, ipWhitelist, allowedOrigins sql.NullString
	var lastUsedIP, lastUsedEndpoint, revokedReason sql.NullString
	var createdBy, revokedBy sql.NullInt64

	err := row.Scan(
		&t.ID, &t.Token, &t.Name, &description, &t.TokenPrefix,
		&scopes, &permissions, &t.Environment, &t.IsActive,
		&ipWhitelist, &allowedOrigins,
		&t.RateLimitPerMinute, &t.RateLimitPerHour, &t.RateLimitPerDay,
		&t.ExpiresAt, &t.LastUsedAt, &lastUsedIP, &lastUsedEndpoint,
		&t.TotalRequests, &t.CreatedAt, &t.UpdatedAt, &createdBy,
		&t.RevokedAt, &revokedBy, &revokedReason,
		&t.VendorName, &t.FilterColumn, &t.FilterValue, &t.IsSuperToken,
	)
	if err != nil {
		return nil, err
	}

	t.Description = description.String
	t.Scopes = scopes.String
	t.Permissions = permissions.String
	t.IPWhitelist = ipWhitelist.String
	t.AllowedOrigins = allowedOrigins.String
	t.LastUsedIP = lastUsedIP.String
	t.LastUsedEndpoint = lastUsedEndpoint.String
	t.RevokedReason = revokedReason.String
	if createdBy.Valid {
		v := int(createdBy.Int64)
		t.CreatedBy = &v
	}
	if revokedBy.Valid {
		v := int(revokedBy.Int64)
		t.RevokedBy = &v
	}
	return &t, nil
}

const tokenSelectQuery = `
	SELECT id, token, name, description, token_prefix, scopes, permissions,
	       environment, is_active, ip_whitelist, allowed_origins,
	       rate_limit_per_minute, rate_limit_per_hour, rate_limit_per_day,
	       expires_at, last_used_at, last_used_ip, last_used_endpoint,
	       total_requests, created_at, updated_at, created_by,
	       revoked_at, revoked_by, revoked_reason,
	       ISNULL(vendor_name, '') as vendor_name,
	       ISNULL(filter_column, '') as filter_column,
	       ISNULL(filter_value, '') as filter_value,
	       ISNULL(is_super_token, 0) as is_super_token
	FROM api_tokens
`

// GetAPITokenByToken retrieves a token by its value
func (r *TokenRepository) GetAPITokenByToken(tokenValue string) (*models.APIToken, error) {
	query := tokenSelectQuery + ` WHERE token = @p1`
	row := r.db.QueryRow(query, tokenValue)
	token, err := r.scanToken(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("token not found")
		}
		return nil, err
	}
	return token, nil
}

// GetAPITokenByID retrieves a token by ID
func (r *TokenRepository) GetAPITokenByID(id int) (*models.APIToken, error) {
	query := tokenSelectQuery + ` WHERE id = @p1`
	row := r.db.QueryRow(query, id)
	token, err := r.scanToken(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("token not found")
		}
		return nil, err
	}
	return token, nil
}

// GetAllAPITokens retrieves all API tokens
func (r *TokenRepository) GetAllAPITokens() ([]*models.APIToken, error) {
	query := tokenSelectQuery + ` ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*models.APIToken
	for rows.Next() {
		token, err := r.scanToken(rows)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, rows.Err()
}

// UpdateAPIToken updates an existing API token
func (r *TokenRepository) UpdateAPIToken(id int, updates map[string]interface{}) error {
	query := "UPDATE api_tokens SET "
	args := []interface{}{}
	paramNum := 1

	for key, value := range updates {
		if paramNum > 1 {
			query += ", "
		}
		query += fmt.Sprintf("%s = @p%d", key, paramNum)
		args = append(args, value)
		paramNum++
	}

	query += fmt.Sprintf(" WHERE id = @p%d", paramNum)
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	return err
}

// UpdateTokenUsage updates token usage statistics
func (r *TokenRepository) UpdateTokenUsage(tokenID int, ipAddress, endpoint string) error {
	query := `
		UPDATE api_tokens
		SET last_used_at = GETDATE(), last_used_ip = @p1,
		    last_used_endpoint = @p2, total_requests = total_requests + 1
		WHERE id = @p3
	`
	_, err := r.db.Exec(query, ipAddress, endpoint, tokenID)
	return err
}

// DisableToken disables a token
func (r *TokenRepository) DisableToken(id int) error {
	_, err := r.db.Exec(`UPDATE api_tokens SET is_active = 0 WHERE id = @p1`, id)
	return err
}

// EnableToken enables a token
func (r *TokenRepository) EnableToken(id int) error {
	_, err := r.db.Exec(`UPDATE api_tokens SET is_active = 1 WHERE id = @p1`, id)
	return err
}

// RevokeToken revokes a token permanently
func (r *TokenRepository) RevokeToken(id int, revokedBy int, reason string) error {
	query := `
		UPDATE api_tokens
		SET revoked_at = GETDATE(), revoked_by = @p1, revoked_reason = @p2, is_active = 0
		WHERE id = @p3
	`
	_, err := r.db.Exec(query, revokedBy, reason, id)
	return err
}

// DeleteToken deletes a token permanently
func (r *TokenRepository) DeleteToken(id int) error {
	_, err := r.db.Exec(`DELETE FROM api_tokens WHERE id = @p1`, id)
	return err
}

// ============================================================================
// Token Usage Logs
// ============================================================================

// CreateUsageLog creates a new usage log entry
func (r *TokenRepository) CreateUsageLog(log *models.TokenUsageLog) error {
	// Set created_at if not already set
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO token_usage_logs (
			token_id, method, endpoint, full_url, status_code, response_time_ms,
			ip_address, user_agent, referer, request_id, request_body_size,
			response_body_size, error_message, error_code, created_at
		)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14, @p15)
	`
	_, err := r.db.Exec(query,
		log.TokenID, log.Method, log.Endpoint, log.FullURL,
		log.StatusCode, log.ResponseTimeMs, log.IPAddress, log.UserAgent,
		log.Referer, log.RequestID, log.RequestBodySize, log.ResponseBodySize,
		log.ErrorMessage, log.ErrorCode, log.CreatedAt,
	)
	return err
}

// GetRecentUsageLogs retrieves recent usage logs
func (r *TokenRepository) GetRecentUsageLogs(limit int) ([]*models.TokenUsageLog, error) {
	query := `
		SELECT TOP (@p1) id, token_id, method, endpoint, ISNULL(full_url, '') as full_url,
		       status_code, ISNULL(response_time_ms, 0) as response_time_ms,
		       ip_address, ISNULL(user_agent, '') as user_agent,
		       ISNULL(referer, '') as referer, ISNULL(request_id, '') as request_id,
		       ISNULL(request_body_size, 0) as request_body_size,
		       ISNULL(response_body_size, 0) as response_body_size,
		       ISNULL(error_message, '') as error_message,
		       ISNULL(error_code, '') as error_code, created_at
		FROM token_usage_logs
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.TokenUsageLog
	for rows.Next() {
		var l models.TokenUsageLog
		err := rows.Scan(
			&l.ID, &l.TokenID, &l.Method, &l.Endpoint, &l.FullURL,
			&l.StatusCode, &l.ResponseTimeMs, &l.IPAddress, &l.UserAgent,
			&l.Referer, &l.RequestID, &l.RequestBodySize, &l.ResponseBodySize,
			&l.ErrorMessage, &l.ErrorCode, &l.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &l)
	}
	return logs, rows.Err()
}

// GetUsageLogsByTokenID retrieves usage logs for a specific token
func (r *TokenRepository) GetUsageLogsByTokenID(tokenID int, limit int) ([]*models.TokenUsageLog, error) {
	query := `
		SELECT TOP (@p2) id, token_id, method, endpoint, ISNULL(full_url, '') as full_url,
		       status_code, ISNULL(response_time_ms, 0) as response_time_ms,
		       ip_address, ISNULL(user_agent, '') as user_agent,
		       ISNULL(referer, '') as referer, ISNULL(request_id, '') as request_id,
		       ISNULL(request_body_size, 0) as request_body_size,
		       ISNULL(response_body_size, 0) as response_body_size,
		       ISNULL(error_message, '') as error_message,
		       ISNULL(error_code, '') as error_code, created_at
		FROM token_usage_logs
		WHERE token_id = @p1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, tokenID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.TokenUsageLog
	for rows.Next() {
		var l models.TokenUsageLog
		err := rows.Scan(
			&l.ID, &l.TokenID, &l.Method, &l.Endpoint, &l.FullURL,
			&l.StatusCode, &l.ResponseTimeMs, &l.IPAddress, &l.UserAgent,
			&l.Referer, &l.RequestID, &l.RequestBodySize, &l.ResponseBodySize,
			&l.ErrorMessage, &l.ErrorCode, &l.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &l)
	}
	return logs, rows.Err()
}

// ============================================================================
// Rate Limiting
// ============================================================================

// GetRateLimitCount gets the current request count for a rate limit window
func (r *TokenRepository) GetRateLimitCount(tokenID int, windowType string, windowStart time.Time) (int, error) {
	query := `
		SELECT ISNULL(request_count, 0)
		FROM token_rate_limits
		WHERE token_id = @p1 AND window_type = @p2 AND window_start = @p3
	`
	var count int
	err := r.db.QueryRow(query, tokenID, windowType, windowStart).Scan(&count)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return count, err
}

// IncrementRateLimit increments or creates a rate limit counter
func (r *TokenRepository) IncrementRateLimit(tokenID int, windowType string, windowStart, windowEnd time.Time) error {
	// Use MERGE for upsert
	query := `
		MERGE token_rate_limits AS target
		USING (SELECT @p1 AS token_id, @p2 AS window_type, @p3 AS window_start) AS source
		ON target.token_id = source.token_id
		   AND target.window_type = source.window_type
		   AND target.window_start = source.window_start
		WHEN MATCHED THEN
			UPDATE SET request_count = request_count + 1, updated_at = GETDATE()
		WHEN NOT MATCHED THEN
			INSERT (token_id, window_type, window_start, window_end, request_count)
			VALUES (@p1, @p2, @p3, @p4, 1);
	`
	_, err := r.db.Exec(query, tokenID, windowType, windowStart, windowEnd)
	return err
}

// ============================================================================
// Analytics
// ============================================================================

// GetTokenAnalytics retrieves analytics for a specific token
func (r *TokenRepository) GetTokenAnalytics(tokenID int, days int) (*models.TokenAnalytics, error) {
	query := `
		SELECT
			t.id AS token_id,
			t.name AS token_name,
			COUNT(l.id) AS total_requests,
			COUNT(CASE WHEN l.status_code >= 200 AND l.status_code < 300 THEN 1 END) AS successful_requests,
			COUNT(CASE WHEN l.status_code >= 400 THEN 1 END) AS failed_requests,
			COUNT(CASE WHEN l.status_code >= 400 AND l.status_code < 500 THEN 1 END) AS client_errors,
			COUNT(CASE WHEN l.status_code >= 500 THEN 1 END) AS server_errors,
			ISNULL(AVG(CAST(l.response_time_ms AS FLOAT)), 0) AS avg_response_time_ms,
			ISNULL(MAX(l.response_time_ms), 0) AS max_response_time_ms,
			COUNT(DISTINCT l.ip_address) AS unique_ips,
			COUNT(DISTINCT l.endpoint) AS unique_endpoints,
			MAX(l.created_at) AS last_used_at
		FROM api_tokens t
		LEFT JOIN token_usage_logs l ON t.id = l.token_id
			AND l.created_at >= DATEADD(day, -@p2, GETDATE())
		WHERE t.id = @p1
		GROUP BY t.id, t.name
	`
	row := r.db.QueryRow(query, tokenID, days)

	var a models.TokenAnalytics
	err := row.Scan(
		&a.TokenID, &a.TokenName, &a.TotalRequests,
		&a.SuccessfulRequests, &a.FailedRequests, &a.ClientErrors,
		&a.ServerErrors, &a.AvgResponseTimeMs, &a.MaxResponseTimeMs,
		&a.UniqueIPs, &a.UniqueEndpoints, &a.LastUsedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// GetDashboardStats retrieves overall dashboard statistics
func (r *TokenRepository) GetDashboardStats() (*models.TokenDashboardStats, error) {
	var stats models.TokenDashboardStats

	// Get token counts
	err := r.db.QueryRow(`
		SELECT
			COUNT(*) AS total_tokens,
			COUNT(CASE WHEN is_active = 1 AND (expires_at IS NULL OR expires_at > GETDATE())
				AND revoked_at IS NULL THEN 1 END) AS active_tokens
		FROM api_tokens
	`).Scan(&stats.TotalTokens, &stats.ActiveTokens)
	if err != nil {
		return nil, err
	}

	// Get 24h request stats
	err = r.db.QueryRow(`
		SELECT
			COUNT(*),
			CASE WHEN COUNT(*) > 0 THEN
				CAST(COUNT(CASE WHEN status_code >= 200 AND status_code < 300 THEN 1 END) AS FLOAT) / COUNT(*) * 100
			ELSE 0 END,
			ISNULL(AVG(CAST(response_time_ms AS FLOAT)), 0)
		FROM token_usage_logs
		WHERE created_at >= DATEADD(hour, -24, GETDATE())
	`).Scan(&stats.TotalRequests24h, &stats.SuccessRate, &stats.AvgResponseTimeMs)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetEndpointStats retrieves statistics per endpoint
func (r *TokenRepository) GetEndpointStats(days int, limit int) ([]*models.EndpointStats, error) {
	query := `
		SELECT TOP (@p2)
			endpoint, method, COUNT(*) AS request_count,
			COUNT(DISTINCT token_id) AS unique_tokens,
			ISNULL(AVG(CAST(response_time_ms AS FLOAT)), 0) AS avg_response_time_ms,
			COUNT(CASE WHEN status_code >= 200 AND status_code < 300 THEN 1 END) AS successful_requests,
			COUNT(CASE WHEN status_code >= 400 THEN 1 END) AS failed_requests
		FROM token_usage_logs
		WHERE created_at >= DATEADD(day, -@p1, GETDATE())
		GROUP BY endpoint, method
		ORDER BY request_count DESC
	`
	rows, err := r.db.Query(query, days, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []*models.EndpointStats
	for rows.Next() {
		var s models.EndpointStats
		err := rows.Scan(&s.Endpoint, &s.Method, &s.RequestCount,
			&s.UniqueTokens, &s.AvgResponseTimeMs,
			&s.SuccessfulRequests, &s.FailedRequests)
		if err != nil {
			return nil, err
		}
		stats = append(stats, &s)
	}
	return stats, rows.Err()
}

// GetDailyUsage retrieves daily usage stats for charting
func (r *TokenRepository) GetDailyUsage(tokenID *int, days int) ([]*models.DailyUsage, error) {
	query := `
		SELECT
			CONVERT(VARCHAR(10), l.created_at, 120) AS usage_date,
			t.id AS token_id, t.name AS token_name,
			COUNT(*) AS request_count,
			COUNT(CASE WHEN l.status_code >= 200 AND l.status_code < 300 THEN 1 END) AS successful_requests,
			COUNT(CASE WHEN l.status_code >= 400 THEN 1 END) AS failed_requests,
			ISNULL(AVG(CAST(l.response_time_ms AS FLOAT)), 0) AS avg_response_time_ms
		FROM api_tokens t
		INNER JOIN token_usage_logs l ON t.id = l.token_id
		WHERE l.created_at >= DATEADD(day, -@p1, GETDATE())
	`
	args := []interface{}{days}

	if tokenID != nil {
		query += ` AND t.id = @p2`
		args = append(args, *tokenID)
	}

	query += `
		GROUP BY CONVERT(VARCHAR(10), l.created_at, 120), t.id, t.name
		ORDER BY usage_date DESC
	`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usage []*models.DailyUsage
	for rows.Next() {
		var u models.DailyUsage
		err := rows.Scan(&u.Date, &u.TokenID, &u.TokenName,
			&u.RequestCount, &u.SuccessfulRequests, &u.FailedRequests,
			&u.AvgResponseTimeMs)
		if err != nil {
			return nil, err
		}
		usage = append(usage, &u)
	}
	return usage, rows.Err()
}

// ============================================================================
// Audit Logging
// ============================================================================

// CreateAuditLog creates a new audit log entry
func (r *TokenRepository) CreateAuditLog(log *models.AuditLog) error {
	query := `
		INSERT INTO audit_logs (
			admin_user_id, action, resource_type, resource_id,
			old_values, new_values, ip_address, user_agent, description
		)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9)
	`
	_, err := r.db.Exec(query,
		log.AdminUserID, log.Action, log.ResourceType, log.ResourceID,
		log.OldValues, log.NewValues, log.IPAddress, log.UserAgent,
		log.Description,
	)
	return err
}

// GetAuditLogs retrieves audit logs with limit
func (r *TokenRepository) GetAuditLogs(limit int) ([]*models.AuditLog, error) {
	query := `
		SELECT TOP (@p1) id, admin_user_id, action, resource_type, resource_id,
		       ISNULL(old_values, '') as old_values, ISNULL(new_values, '') as new_values,
		       ISNULL(ip_address, '') as ip_address, ISNULL(user_agent, '') as user_agent,
		       ISNULL(description, '') as description, created_at
		FROM audit_logs
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.AuditLog
	for rows.Next() {
		var l models.AuditLog
		var adminUserID, resourceID sql.NullInt64
		err := rows.Scan(
			&l.ID, &adminUserID, &l.Action, &l.ResourceType, &resourceID,
			&l.OldValues, &l.NewValues, &l.IPAddress, &l.UserAgent,
			&l.Description, &l.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if adminUserID.Valid {
			v := int(adminUserID.Int64)
			l.AdminUserID = &v
		}
		if resourceID.Valid {
			v := int(resourceID.Int64)
			l.ResourceID = &v
		}
		logs = append(logs, &l)
	}
	return logs, rows.Err()
}

// ============================================================================
// Helper Functions
// ============================================================================

// ConvertToJSON converts a value to JSON string
func ConvertToJSON(data interface{}) (string, error) {
	if data == nil {
		return "[]", nil
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
