-- ============================================================================
-- Migration 002: Add Vendor Filter & Super Token Support
-- ============================================================================
-- Purpose: Extend api_tokens to support per-token vendor data filtering.
--          Each token can be scoped to a specific vendor by specifying which
--          column and value to filter the data source by.
--          Super-tokens bypass all filters and can read/write all data.
-- ============================================================================

USE token_management;
GO

-- Add vendor_name column (human-readable label, e.g. 'AVT', 'VENDOR_B')
IF NOT EXISTS (
    SELECT 1 FROM sys.columns
    WHERE object_id = OBJECT_ID('api_tokens') AND name = 'vendor_name'
)
BEGIN
    ALTER TABLE api_tokens
    ADD vendor_name NVARCHAR(200) NULL;
    PRINT 'Column vendor_name added to api_tokens.';
END
GO

-- Add filter_column (the DB column used for filtering, e.g. 'mm.[FLM name]')
IF NOT EXISTS (
    SELECT 1 FROM sys.columns
    WHERE object_id = OBJECT_ID('api_tokens') AND name = 'filter_column'
)
BEGIN
    ALTER TABLE api_tokens
    ADD filter_column NVARCHAR(200) NULL;
    PRINT 'Column filter_column added to api_tokens.';
END
GO

-- Add filter_value (the value to match, e.g. 'AVT')
IF NOT EXISTS (
    SELECT 1 FROM sys.columns
    WHERE object_id = OBJECT_ID('api_tokens') AND name = 'filter_value'
)
BEGIN
    ALTER TABLE api_tokens
    ADD filter_value NVARCHAR(500) NULL;
    PRINT 'Column filter_value added to api_tokens.';
END
GO

-- Add is_super_token (bypasses all filters, full read/write access)
IF NOT EXISTS (
    SELECT 1 FROM sys.columns
    WHERE object_id = OBJECT_ID('api_tokens') AND name = 'is_super_token'
)
BEGIN
    ALTER TABLE api_tokens
    ADD is_super_token BIT NOT NULL DEFAULT 0;
    PRINT 'Column is_super_token added to api_tokens.';
END
GO

-- Index for fast vendor lookups
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_vendor_name' AND object_id = OBJECT_ID('api_tokens'))
    CREATE INDEX idx_vendor_name ON api_tokens(vendor_name);
GO

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_is_super_token' AND object_id = OBJECT_ID('api_tokens'))
    CREATE INDEX idx_is_super_token ON api_tokens(is_super_token);
GO

-- ============================================================================
-- Update sample token: grant super access to the default test token
-- ============================================================================
UPDATE api_tokens
SET is_super_token = 1, vendor_name = NULL, filter_column = NULL, filter_value = NULL
WHERE token = 'tok_test_sample123456789';
GO

-- ============================================================================
-- Sample vendor-scoped token for AVT
-- (Only shows tickets where mm.[FLM name] = 'AVT')
-- ============================================================================
IF NOT EXISTS (SELECT 1 FROM api_tokens WHERE token = 'tok_live_avt_sample')
BEGIN
    INSERT INTO api_tokens (
        token, name, description, token_prefix, scopes,
        environment, is_active,
        vendor_name, filter_column, filter_value, is_super_token,
        rate_limit_per_minute, rate_limit_per_hour, rate_limit_per_day,
        created_by
    )
    VALUES (
        'tok_live_avt_sample',
        'AVT Vendor Token',
        'Scoped to AVT vendor - only sees AVT tickets via mm.[FLM name] = AVT',
        'tok_live',
        '["tickets:read", "machines:read"]',
        'production',
        1,
        'AVT', 'flm_name', 'AVT', 0,
        60, 3000, 50000,
        1
    );
    PRINT 'Sample AVT vendor token created.';
END
GO

-- ============================================================================
-- DONE
-- ============================================================================
PRINT '============================================';
PRINT 'Migration 002 applied successfully!';
PRINT '============================================';
PRINT 'New columns on api_tokens:';
PRINT '  - vendor_name    : human-readable label (e.g. AVT)';
PRINT '  - filter_column  : DB column key used for WHERE clause (e.g. flm_name)';
PRINT '  - filter_value   : value to match (e.g. AVT)';
PRINT '  - is_super_token : 1 = bypass all filters, full access';
GO
