-- ============================================================================
-- API Token Management System - Database Schema
-- ============================================================================
-- Database: token_management
-- Purpose: Manage API tokens, usage analytics, and admin authentication
-- Version: 1.0
-- ============================================================================
-- Run with: sqlcmd -S localhost -U your_username -P your_password -i 001_create_token_management_schema.sql
-- ============================================================================

-- Create database
IF DB_ID('token_management') IS NULL
BEGIN
    CREATE DATABASE token_management;
    PRINT 'Database token_management created.';
END
ELSE
BEGIN
    PRINT 'Database token_management already exists.';
END
GO

USE token_management;
GO

-- ============================================================================
-- Table: admin_users
-- ============================================================================
IF OBJECT_ID('admin_users', 'U') IS NULL
BEGIN
    CREATE TABLE admin_users (
        id INT IDENTITY(1,1) PRIMARY KEY,
        username NVARCHAR(100) NOT NULL UNIQUE,
        email NVARCHAR(255) NOT NULL UNIQUE,
        password_hash NVARCHAR(255) NOT NULL,
        full_name NVARCHAR(200),
        role NVARCHAR(50) NOT NULL DEFAULT 'admin', -- super_admin, admin, viewer
        is_active BIT NOT NULL DEFAULT 1,
        last_login_at DATETIME2,
        last_login_ip NVARCHAR(45),
        created_at DATETIME2 NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME2 NOT NULL DEFAULT GETDATE(),
        created_by INT,
        INDEX idx_username (username),
        INDEX idx_email (email),
        INDEX idx_role (role),
        INDEX idx_is_active (is_active)
    );
    PRINT 'Table admin_users created.';
END
GO

-- ============================================================================
-- Table: api_tokens
-- ============================================================================
IF OBJECT_ID('api_tokens', 'U') IS NULL
BEGIN
    CREATE TABLE api_tokens (
        id INT IDENTITY(1,1) PRIMARY KEY,
        token NVARCHAR(255) NOT NULL UNIQUE,
        name NVARCHAR(200) NOT NULL,
        description NVARCHAR(500),
        token_prefix NVARCHAR(20) NOT NULL,

        -- Permissions & Scopes
        scopes NVARCHAR(MAX),
        permissions NVARCHAR(MAX),

        -- Environment & Status
        environment NVARCHAR(50) NOT NULL DEFAULT 'production',
        is_active BIT NOT NULL DEFAULT 1,

        -- Security
        ip_whitelist NVARCHAR(MAX),
        allowed_origins NVARCHAR(MAX),

        -- Rate Limiting
        rate_limit_per_minute INT DEFAULT 100,
        rate_limit_per_hour INT DEFAULT 5000,
        rate_limit_per_day INT DEFAULT 100000,

        -- Expiration
        expires_at DATETIME2,

        -- Usage Statistics
        last_used_at DATETIME2,
        last_used_ip NVARCHAR(45),
        last_used_endpoint NVARCHAR(500),
        total_requests BIGINT NOT NULL DEFAULT 0,

        -- Metadata
        created_at DATETIME2 NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME2 NOT NULL DEFAULT GETDATE(),
        created_by INT,
        revoked_at DATETIME2,
        revoked_by INT,
        revoked_reason NVARCHAR(500),

        -- Foreign Keys
        CONSTRAINT fk_api_tokens_created_by FOREIGN KEY (created_by) REFERENCES admin_users(id),
        CONSTRAINT fk_api_tokens_revoked_by FOREIGN KEY (revoked_by) REFERENCES admin_users(id),

        -- Indexes
        INDEX idx_token (token),
        INDEX idx_is_active (is_active),
        INDEX idx_environment (environment),
        INDEX idx_expires_at (expires_at),
        INDEX idx_created_by (created_by),
        INDEX idx_token_prefix (token_prefix)
    );
    PRINT 'Table api_tokens created.';
END
GO

-- ============================================================================
-- Table: token_usage_logs
-- ============================================================================
IF OBJECT_ID('token_usage_logs', 'U') IS NULL
BEGIN
    CREATE TABLE token_usage_logs (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        token_id INT NOT NULL,

        -- Request Details
        method NVARCHAR(10) NOT NULL,
        endpoint NVARCHAR(500) NOT NULL,
        full_url NVARCHAR(1000),

        -- Response Details
        status_code INT NOT NULL,
        response_time_ms INT,

        -- Request Context
        ip_address NVARCHAR(45) NOT NULL,
        user_agent NVARCHAR(500),
        referer NVARCHAR(500),

        -- Metadata
        request_id NVARCHAR(100),
        request_body_size INT,
        response_body_size INT,

        -- Error Tracking
        error_message NVARCHAR(MAX),
        error_code NVARCHAR(100),

        -- Timestamps
        created_at DATETIME2 NOT NULL DEFAULT GETDATE(),

        -- Foreign Key
        CONSTRAINT fk_token_usage_logs_token_id FOREIGN KEY (token_id) REFERENCES api_tokens(id) ON DELETE CASCADE,

        -- Indexes
        INDEX idx_token_id (token_id),
        INDEX idx_created_at (created_at),
        INDEX idx_endpoint (endpoint),
        INDEX idx_status_code (status_code),
        INDEX idx_ip_address (ip_address),
        INDEX idx_request_id (request_id)
    );
    PRINT 'Table token_usage_logs created.';
END
GO

-- ============================================================================
-- Table: token_rate_limits
-- ============================================================================
IF OBJECT_ID('token_rate_limits', 'U') IS NULL
BEGIN
    CREATE TABLE token_rate_limits (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        token_id INT NOT NULL,

        -- Time Windows
        window_type NVARCHAR(20) NOT NULL,
        window_start DATETIME2 NOT NULL,
        window_end DATETIME2 NOT NULL,

        -- Counters
        request_count INT NOT NULL DEFAULT 0,

        -- Metadata
        created_at DATETIME2 NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME2 NOT NULL DEFAULT GETDATE(),

        -- Foreign Key
        CONSTRAINT fk_token_rate_limits_token_id FOREIGN KEY (token_id) REFERENCES api_tokens(id) ON DELETE CASCADE,

        -- Unique constraint
        CONSTRAINT uq_token_window UNIQUE (token_id, window_type, window_start),

        -- Indexes
        INDEX idx_token_id_window (token_id, window_type, window_start),
        INDEX idx_window_end (window_end)
    );
    PRINT 'Table token_rate_limits created.';
END
GO

-- ============================================================================
-- Table: admin_sessions
-- ============================================================================
IF OBJECT_ID('admin_sessions', 'U') IS NULL
BEGIN
    CREATE TABLE admin_sessions (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        session_token NVARCHAR(255) NOT NULL UNIQUE,
        admin_user_id INT NOT NULL,

        -- Session Details
        ip_address NVARCHAR(45),
        user_agent NVARCHAR(500),

        -- Expiration
        expires_at DATETIME2 NOT NULL,

        -- Metadata
        created_at DATETIME2 NOT NULL DEFAULT GETDATE(),
        last_accessed_at DATETIME2 NOT NULL DEFAULT GETDATE(),

        -- Foreign Key
        CONSTRAINT fk_admin_sessions_user_id FOREIGN KEY (admin_user_id) REFERENCES admin_users(id) ON DELETE CASCADE,

        -- Indexes
        INDEX idx_session_token (session_token),
        INDEX idx_admin_user_id (admin_user_id),
        INDEX idx_expires_at (expires_at)
    );
    PRINT 'Table admin_sessions created.';
END
GO

-- ============================================================================
-- Table: audit_logs
-- ============================================================================
IF OBJECT_ID('audit_logs', 'U') IS NULL
BEGIN
    CREATE TABLE audit_logs (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        admin_user_id INT,

        -- Action Details
        action NVARCHAR(100) NOT NULL,
        resource_type NVARCHAR(50) NOT NULL,
        resource_id INT,

        -- Changes
        old_values NVARCHAR(MAX),
        new_values NVARCHAR(MAX),

        -- Context
        ip_address NVARCHAR(45),
        user_agent NVARCHAR(500),
        description NVARCHAR(1000),

        -- Metadata
        created_at DATETIME2 NOT NULL DEFAULT GETDATE(),

        -- Foreign Key
        CONSTRAINT fk_audit_logs_user_id FOREIGN KEY (admin_user_id) REFERENCES admin_users(id) ON DELETE SET NULL,

        -- Indexes
        INDEX idx_admin_user_id (admin_user_id),
        INDEX idx_created_at (created_at),
        INDEX idx_action (action),
        INDEX idx_resource (resource_type, resource_id)
    );
    PRINT 'Table audit_logs created.';
END
GO

-- ============================================================================
-- Insert Default Admin User
-- ============================================================================
-- Password: admin123 (CHANGE THIS IN PRODUCTION!)
-- Bcrypt hash for 'admin123' at cost 10
IF NOT EXISTS (SELECT 1 FROM admin_users WHERE username = 'admin')
BEGIN
    INSERT INTO admin_users (username, email, password_hash, full_name, role, is_active)
    VALUES (
        'admin',
        'admin@example.com',
        '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
        'System Administrator',
        'super_admin',
        1
    );
    PRINT 'Default admin user created (admin / admin123).';
END
GO

-- ============================================================================
-- Insert Sample Token (for testing)
-- ============================================================================
IF NOT EXISTS (SELECT 1 FROM api_tokens WHERE token = 'tok_test_sample123456789')
BEGIN
    INSERT INTO api_tokens (
        token, name, description, token_prefix, scopes,
        environment, is_active,
        rate_limit_per_minute, rate_limit_per_hour, rate_limit_per_day,
        created_by
    )
    VALUES (
        'tok_test_sample123456789',
        'Test Token',
        'Sample token for testing purposes',
        'tok_test',
        '["tickets:read", "tickets:write", "machines:read", "machines:write"]',
        'development',
        1,
        100, 5000, 100000,
        1
    );
    PRINT 'Sample test token created.';
END
GO

-- ============================================================================
-- Views
-- ============================================================================

-- View: Token usage summary
IF OBJECT_ID('vw_token_usage_summary', 'V') IS NOT NULL
    DROP VIEW vw_token_usage_summary;
GO

CREATE VIEW vw_token_usage_summary AS
SELECT
    t.id AS token_id,
    t.name AS token_name,
    t.token_prefix,
    t.environment,
    t.is_active,
    COUNT(l.id) AS total_requests,
    COUNT(CASE WHEN l.status_code >= 200 AND l.status_code < 300 THEN 1 END) AS successful_requests,
    COUNT(CASE WHEN l.status_code >= 400 THEN 1 END) AS failed_requests,
    AVG(l.response_time_ms) AS avg_response_time_ms,
    MAX(l.created_at) AS last_request_at,
    COUNT(DISTINCT l.ip_address) AS unique_ips
FROM api_tokens t
LEFT JOIN token_usage_logs l ON t.id = l.token_id
GROUP BY t.id, t.name, t.token_prefix, t.environment, t.is_active;
GO

-- View: Endpoint usage statistics
IF OBJECT_ID('vw_endpoint_usage_stats', 'V') IS NOT NULL
    DROP VIEW vw_endpoint_usage_stats;
GO

CREATE VIEW vw_endpoint_usage_stats AS
SELECT
    l.endpoint,
    l.method,
    COUNT(*) AS request_count,
    COUNT(DISTINCT l.token_id) AS unique_tokens,
    AVG(l.response_time_ms) AS avg_response_time_ms,
    MIN(l.response_time_ms) AS min_response_time_ms,
    MAX(l.response_time_ms) AS max_response_time_ms,
    COUNT(CASE WHEN l.status_code >= 200 AND l.status_code < 300 THEN 1 END) AS successful_requests,
    COUNT(CASE WHEN l.status_code >= 400 THEN 1 END) AS failed_requests
FROM token_usage_logs l
WHERE l.created_at >= DATEADD(day, -7, GETDATE())
GROUP BY l.endpoint, l.method;
GO

-- View: Daily token usage
IF OBJECT_ID('vw_daily_token_usage', 'V') IS NOT NULL
    DROP VIEW vw_daily_token_usage;
GO

CREATE VIEW vw_daily_token_usage AS
SELECT
    t.id AS token_id,
    t.name AS token_name,
    CAST(l.created_at AS DATE) AS usage_date,
    COUNT(*) AS request_count,
    COUNT(CASE WHEN l.status_code >= 200 AND l.status_code < 300 THEN 1 END) AS successful_requests,
    COUNT(CASE WHEN l.status_code >= 400 THEN 1 END) AS failed_requests,
    AVG(l.response_time_ms) AS avg_response_time_ms
FROM api_tokens t
LEFT JOIN token_usage_logs l ON t.id = l.token_id
WHERE l.created_at >= DATEADD(day, -30, GETDATE())
GROUP BY t.id, t.name, CAST(l.created_at AS DATE);
GO

-- ============================================================================
-- Stored Procedures
-- ============================================================================

-- Procedure: Clean up old usage logs (retention policy)
IF OBJECT_ID('sp_cleanup_old_usage_logs', 'P') IS NOT NULL
    DROP PROCEDURE sp_cleanup_old_usage_logs;
GO

CREATE PROCEDURE sp_cleanup_old_usage_logs
    @retention_days INT = 90
AS
BEGIN
    SET NOCOUNT ON;
    DECLARE @cutoff_date DATETIME2 = DATEADD(day, -@retention_days, GETDATE());
    DELETE FROM token_usage_logs WHERE created_at < @cutoff_date;
    SELECT @@ROWCOUNT AS rows_deleted;
END;
GO

-- Procedure: Clean up expired rate limit windows
IF OBJECT_ID('sp_cleanup_expired_rate_limits', 'P') IS NOT NULL
    DROP PROCEDURE sp_cleanup_expired_rate_limits;
GO

CREATE PROCEDURE sp_cleanup_expired_rate_limits
AS
BEGIN
    SET NOCOUNT ON;
    DELETE FROM token_rate_limits WHERE window_end < DATEADD(hour, -1, GETDATE());
    SELECT @@ROWCOUNT AS rows_deleted;
END;
GO

-- Procedure: Clean up expired sessions
IF OBJECT_ID('sp_cleanup_expired_sessions', 'P') IS NOT NULL
    DROP PROCEDURE sp_cleanup_expired_sessions;
GO

CREATE PROCEDURE sp_cleanup_expired_sessions
AS
BEGIN
    SET NOCOUNT ON;
    DELETE FROM admin_sessions WHERE expires_at < GETDATE();
    SELECT @@ROWCOUNT AS rows_deleted;
END;
GO

-- Procedure: Get token analytics
IF OBJECT_ID('sp_get_token_analytics', 'P') IS NOT NULL
    DROP PROCEDURE sp_get_token_analytics;
GO

CREATE PROCEDURE sp_get_token_analytics
    @token_id INT = NULL,
    @days INT = 7
AS
BEGIN
    SET NOCOUNT ON;
    DECLARE @start_date DATETIME2 = DATEADD(day, -@days, GETDATE());

    SELECT
        t.id AS token_id,
        t.name AS token_name,
        t.environment,
        t.is_active,
        COUNT(l.id) AS total_requests,
        COUNT(CASE WHEN l.status_code >= 200 AND l.status_code < 300 THEN 1 END) AS successful_requests,
        COUNT(CASE WHEN l.status_code >= 400 AND l.status_code < 500 THEN 1 END) AS client_errors,
        COUNT(CASE WHEN l.status_code >= 500 THEN 1 END) AS server_errors,
        AVG(l.response_time_ms) AS avg_response_time_ms,
        MAX(l.response_time_ms) AS max_response_time_ms,
        COUNT(DISTINCT l.ip_address) AS unique_ips,
        COUNT(DISTINCT l.endpoint) AS unique_endpoints,
        MAX(l.created_at) AS last_used_at
    FROM api_tokens t
    LEFT JOIN token_usage_logs l ON t.id = l.token_id AND l.created_at >= @start_date
    WHERE (@token_id IS NULL OR t.id = @token_id)
    GROUP BY t.id, t.name, t.environment, t.is_active
    ORDER BY total_requests DESC;
END;
GO

-- ============================================================================
-- Triggers
-- ============================================================================

-- Trigger: Update api_tokens.updated_at on modification
IF OBJECT_ID('trg_api_tokens_update_timestamp', 'TR') IS NOT NULL
    DROP TRIGGER trg_api_tokens_update_timestamp;
GO

CREATE TRIGGER trg_api_tokens_update_timestamp
ON api_tokens
AFTER UPDATE
AS
BEGIN
    SET NOCOUNT ON;
    UPDATE api_tokens
    SET updated_at = GETDATE()
    FROM api_tokens t
    INNER JOIN inserted i ON t.id = i.id;
END;
GO

-- Trigger: Update admin_users.updated_at on modification
IF OBJECT_ID('trg_admin_users_update_timestamp', 'TR') IS NOT NULL
    DROP TRIGGER trg_admin_users_update_timestamp;
GO

CREATE TRIGGER trg_admin_users_update_timestamp
ON admin_users
AFTER UPDATE
AS
BEGIN
    SET NOCOUNT ON;
    UPDATE admin_users
    SET updated_at = GETDATE()
    FROM admin_users u
    INNER JOIN inserted i ON u.id = i.id;
END;
GO

-- ============================================================================
-- Additional Indexes
-- ============================================================================

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_token_usage_logs_composite')
    CREATE INDEX idx_token_usage_logs_composite
    ON token_usage_logs(token_id, created_at DESC, status_code);
GO

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_token_usage_logs_endpoint_date')
    CREATE INDEX idx_token_usage_logs_endpoint_date
    ON token_usage_logs(endpoint, created_at DESC);
GO

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_audit_logs_composite')
    CREATE INDEX idx_audit_logs_composite
    ON audit_logs(admin_user_id, created_at DESC, action);
GO

-- ============================================================================
-- Documentation
-- ============================================================================

EXEC sp_addextendedproperty
    @name = N'MS_Description',
    @value = N'Admin users who can access the token management dashboard',
    @level0type = N'SCHEMA', @level0name = 'dbo',
    @level1type = N'TABLE',  @level1name = 'admin_users';
GO

EXEC sp_addextendedproperty
    @name = N'MS_Description',
    @value = N'API tokens with scopes, permissions, and analytics',
    @level0type = N'SCHEMA', @level0name = 'dbo',
    @level1type = N'TABLE',  @level1name = 'api_tokens';
GO

EXEC sp_addextendedproperty
    @name = N'MS_Description',
    @value = N'Complete audit trail of all API requests for analytics and compliance',
    @level0type = N'SCHEMA', @level0name = 'dbo',
    @level1type = N'TABLE',  @level1name = 'token_usage_logs';
GO

-- ============================================================================
-- DONE
-- ============================================================================
PRINT '============================================';
PRINT 'Token Management Schema created successfully!';
PRINT '============================================';
PRINT 'Default admin login: admin / admin123';
PRINT 'IMPORTANT: Change the default password in production!';
PRINT '';
PRINT 'Next steps:';
PRINT '  1. Start the API Gateway: go run main.go';
PRINT '  2. Open http://localhost:8080/admin/login';
PRINT '  3. Login with admin / admin123';
PRINT '  4. Create your first API token';
GO
