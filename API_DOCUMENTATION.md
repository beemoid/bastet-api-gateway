# API Gateway - Complete API Documentation

## Overview

This API Gateway serves as middleware between on-premise databases (ticket_master, machine_master, and token_management) and cloud applications, providing a secure and standardized REST API for ATM monitoring, ticket management, and token-based access control.

**Base URL:** `http://localhost:8080`
**API Version:** v1
**API Prefix:** `/api/v1`
**Swagger UI:** `http://localhost:8080/swagger/index.html`
**Admin Dashboard:** `http://localhost:8080/admin`

---

## Key Features

- **Hybrid Adaptive Metadata System** - Automatically discovers new field values from database
- **Intelligent Caching** - 1-hour cache for optimal performance
- **Thread-Safe Operations** - Production-ready concurrent request handling
- **Comprehensive Error Handling** - Standardized error responses
- **Real-time Database Queries** - Always up-to-date with actual data
- **Token Management System** - Scoped API tokens with rate limiting, IP whitelisting, and expiration
- **Admin Dashboard** - Web UI for managing tokens, viewing analytics, and audit logs
- **Gzip Compression** - Automatic response compression for reduced bandwidth
- **Analytics & Monitoring** - Dashboard stats, endpoint analytics, daily usage tracking
- **Audit Logging** - Complete history of all administrative actions

---

## Authentication

The API supports two authentication methods for `/api/v1` endpoints.

### API Token (Recommended)

Token-based authentication with scopes, rate limiting, usage tracking, and expiration. Tokens are created and managed via the admin dashboard or admin API.

**Header:**
```
X-API-Token: tok_live_abc123xyz456
```

**Features:**
- Scoped permissions
- Configurable rate limits (per minute, hour, day)
- IP whitelisting (optional)
- Token expiration dates
- Usage analytics and audit logging
- Token revocation support

### API Key (Legacy)

Simple key-based authentication using a shared secret configured in `.env`.

**Header:**
```
X-API-Key: your_api_key_here
```

**Note:** Health checks (`/health`, `/ping`), Swagger docs (`/swagger/*`), and admin dashboard pages (`/admin/*`) do not require authentication.

**Configuration:**
Set your credentials in the `.env` file:
```env
API_KEY=your_secure_api_key_here
JWT_SECRET=your_jwt_secret_key
```

---

## Response Format

All API responses follow a standardized format for consistency.

### Success Response
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": { ... },
  "total": 100
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

---

## Endpoints

### 1. Health Check Endpoints

#### 1.1 Health Check
Check the overall health of the API and all database connections.

**Endpoint:** `GET /health`
**Authentication:** Not required

**Response (200 OK):**
```json
{
  "status": "healthy",
  "message": "API Gateway is running",
  "services": {
    "ticket_database": "connected",
    "machine_database": "connected",
    "token_database": "connected"
  }
}
```

**Response (503 Service Unavailable):**
```json
{
  "status": "unhealthy",
  "message": "Database connection failed",
  "error": "connection timeout"
}
```

#### 1.2 Ping
Simple ping endpoint to check if the API is running.

**Endpoint:** `GET /ping`
**Authentication:** Not required

**Response (200 OK):**
```json
{
  "message": "pong"
}
```

---

### 2. Ticket Endpoints

#### 2.1 Get All Tickets
Retrieve all tickets from the ticket_master database.

**Endpoint:** `GET /api/v1/tickets`
**Authentication:** Required

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Tickets retrieved successfully",
  "data": [
    {
      "terminal_id": "ATM-001",
      "terminal_name": "Main Branch ATM",
      "priority": "1.High",
      "mode": "Off-line",
      "initial_problem": "Cash dispenser jam",
      "current_problem": "Card reader error",
      "p_duration": "2h 30m",
      "incident_start_datetime": "2024-01-15 10:30:00",
      "count": 5,
      "status": "0.NEW",
      "remarks": "Waiting for technician",
      "balance": 1000000,
      "condition": "Critical",
      "tickets_no": "TKT-2024-001",
      "tickets_duration": 150.5,
      "open_time": "2024-01-15 08:00:00",
      "close_time": null,
      "problem_history": "Card reader issue resolved",
      "mode_history": "Online->Offline->Online",
      "dsp_flm": "FLM-001",
      "dsp_slm": "SLM-001",
      "last_withdrawal": "2024-01-15T09:30:00Z",
      "export_name": "ATM_Report_Jan2024"
    }
  ],
  "total": 1
}
```

**Note:** Fields with NULL values in database will show as `null` in JSON response.

#### 2.2 Get Ticket by Terminal ID
Retrieve a specific ticket by its terminal ID.

**Endpoint:** `GET /api/v1/tickets/:id`
**Authentication:** Required

**Parameters:**
- `id` (path, string) - Terminal ID (e.g., "ATM-001")

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Ticket retrieved successfully",
  "data": {
    "terminal_id": "ATM-001",
    "terminal_name": "Main Branch ATM",
    "priority": "1.High",
    "status": "0.NEW",
    ...
  }
}
```

**Response (404 Not Found):**
```json
{
  "success": false,
  "message": "Ticket not found"
}
```

#### 2.3 Get Ticket by Number
Retrieve a ticket by its unique ticket number.

**Endpoint:** `GET /api/v1/tickets/number/:number`
**Authentication:** Required

**Parameters:**
- `number` (path, string) - Ticket number (e.g., "TKT-2024-001")

**Response:** Same as Get Ticket by Terminal ID

#### 2.4 Get Tickets by Status
Retrieve all tickets with a specific status.

**Endpoint:** `GET /api/v1/tickets/status/:status`
**Authentication:** Required

**Parameters:**
- `status` (path, string) - Status value

**Valid Status Values:**
- `0.NEW` - New ticket
- `1.Req FD ke HD` - Request FD to HD
- `2.Kirim FLM` - Send FLM
- `21.Req Replenish` - Request Replenish
- `3.SLM` - SLM Machine
- `4.SLM-Net` - SLM Network
- `5.Menunggu Update` - Waiting for Update
- `6.Follow-up Sales team` - Follow-up Sales team
- `8.Wait transaction` - Wait transaction

**Note:** Use the Metadata endpoint to get the current list of all status values in the database.

**Example:**
```
GET /api/v1/tickets/status/0.NEW
```

**Response:** Same as Get All Tickets (filtered list)

#### 2.5 Get Tickets by Terminal
Retrieve all tickets associated with a specific terminal.

**Endpoint:** `GET /api/v1/tickets/terminal/:terminal_id`
**Authentication:** Required

**Parameters:**
- `terminal_id` (path, string) - Terminal ID (e.g., "ATM-001")

**Response:** Same as Get All Tickets (filtered list)

#### 2.6 Get Ticket Metadata - Adaptive System
Retrieve all valid values for ticket fields directly from the database.

**Endpoint:** `GET /api/v1/tickets/metadata`
**Authentication:** Required

**Features:**
- Fully adaptive - automatically discovers new values from database
- Cached for performance - results cached for 1 hour
- Self-documenting - shows which values are documented
- Real-time accuracy - always reflects actual database content

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Metadata retrieved successfully from database",
  "last_updated": "2024-01-15T10:30:00Z",
  "statuses": [
    {
      "code": "0.NEW",
      "description": "New ticket",
      "is_documented": true
    },
    {
      "code": "9.Emergency",
      "description": "9.Emergency",
      "is_documented": false
    }
  ],
  "modes": [
    {
      "code": "Closed",
      "description": "Terminal is closed",
      "is_documented": true
    }
  ],
  "priorities": [
    {
      "code": "1.High",
      "description": "High priority",
      "is_documented": true
    }
  ]
}
```

**Understanding the Response:**
- `code`: The actual value stored in the database
- `description`: Human-readable description
- `is_documented`:
  - `true` = Value has official documentation
  - `false` = New value discovered from database (not yet documented)

**How Adaptation Works:**
1. First request queries database for DISTINCT values
2. Results cached for 1 hour for performance
3. New values added to database appear automatically after cache expires
4. No code changes ever needed when new values are added

#### 2.7 Create Ticket
Create a new ticket in the system.

**Endpoint:** `POST /api/v1/tickets`
**Authentication:** Required

**Request Body:**
```json
{
  "terminal_id": "ATM-001",
  "terminal_name": "Main Branch ATM",
  "priority": "1.High",
  "mode": "Off-line",
  "initial_problem": "Card reader not responding",
  "current_problem": "Card reader error",
  "p_duration": "30m",
  "incident_start_datetime": "2024-01-15 10:30:00",
  "status": "0.NEW",
  "remarks": "Urgent",
  "condition": "Critical",
  "tickets_no": "TKT-2024-001",
  "export_name": "Jan2024_Report"
}
```

**Field Validation:**
- `terminal_id` (required, string) - Terminal identifier
- `terminal_name` (required, string) - Terminal name
- `priority` (required, string) - Priority level (1.High, 2.Middle, 3.Low, 4.Minimum)
- `mode` (required, string) - Terminal mode (Closed, In Service, nan, Off-line, Supervisor)
- `initial_problem` (required, string) - Initial problem description
- `current_problem` (optional, string) - Current problem description
- `p_duration` (optional, string) - Problem duration
- `incident_start_datetime` (optional, string) - Incident start time
- `status` (optional, string) - Ticket status
- `remarks` (optional, string) - Additional remarks
- `condition` (optional, string) - Condition status
- `tickets_no` (optional, string) - Ticket number
- `export_name` (optional, string) - Export name

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Ticket created successfully",
  "data": {
    "terminal_id": "ATM-001",
    "tickets_no": "TKT-2024-001",
    ...
  }
}
```

**Response (409 Conflict):**
```json
{
  "success": false,
  "message": "ticket with this number already exists"
}
```

#### 2.8 Update Ticket
Update an existing ticket.

**Endpoint:** `PUT /api/v1/tickets/:id`
**Authentication:** Required

**Parameters:**
- `id` (path, string) - Terminal ID to update

**Request Body (all fields optional):**
```json
{
  "priority": "2.Middle",
  "mode": "In Service",
  "current_problem": "Issue resolved",
  "status": "6.Follow-up Sales team",
  "remarks": "Technician dispatched",
  "condition": "Normal",
  "close_time": "2024-01-15 18:00:00",
  "problem_history": "Card reader replaced",
  "mode_history": "Off-line->In Service"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Ticket updated successfully",
  "data": {
    "terminal_id": "ATM-001",
    "status": "6.Follow-up Sales team",
    ...
  }
}
```

---

### 3. Machine Endpoints

#### 3.1 Get All Machines
Retrieve all machines/terminals from the machine_master database.

**Endpoint:** `GET /api/v1/machines`
**Authentication:** Required

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Machines retrieved successfully",
  "data": [
    {
      "terminal_id": "ATM-001",
      "store": "Main Branch",
      "store_code": "STR-001",
      "store_name": "Jakarta Main Branch",
      "date_of_activation": "2023-01-15T00:00:00Z",
      "status": "Active",
      "std": 1,
      "gps": "-6.200000,106.816666",
      "lat": -6.200000,
      "lon": 106.816666,
      "province": "DKI Jakarta",
      "city_regency": "Jakarta Pusat",
      "district": "Menteng"
    }
  ],
  "total": 1
}
```

#### 3.2 Get Machine by Terminal ID
Retrieve a specific machine by its terminal ID.

**Endpoint:** `GET /api/v1/machines/:terminal_id`
**Authentication:** Required

**Parameters:**
- `terminal_id` (path, string) - Terminal ID (e.g., "ATM-001")

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Machine retrieved successfully",
  "data": {
    "terminal_id": "ATM-001",
    "store_name": "Jakarta Main Branch",
    "status": "Active",
    "province": "DKI Jakarta",
    ...
  }
}
```

**Response (404 Not Found):**
```json
{
  "success": false,
  "message": "Machine not found"
}
```

#### 3.3 Get Machines by Status
Retrieve all machines with a specific operational status.

**Endpoint:** `GET /api/v1/machines/status/:status`
**Authentication:** Required

**Parameters:**
- `status` (path, string) - Status value (e.g., "Active", "Inactive", "Maintenance", "Offline")

**Example:**
```
GET /api/v1/machines/status/Active
```

**Response:** Same as Get All Machines (filtered list)

#### 3.4 Get Machines by Branch
Retrieve all machines for a specific branch/store.

**Endpoint:** `GET /api/v1/machines/branch/:branch_code`
**Authentication:** Required

**Parameters:**
- `branch_code` (path, string) - Store code (e.g., "STR-001")

**Response:** Same as Get All Machines (filtered list)

#### 3.5 Search Machines
Search machines using multiple filter criteria.

**Endpoint:** `GET /api/v1/machines/search`
**Authentication:** Required

**Query Parameters (all optional):**
- `status` (string) - Filter by status
- `store_code` (string) - Filter by store code
- `province` (string) - Filter by province
- `city_regency` (string) - Filter by city/regency
- `district` (string) - Search by district (partial match)

**Example:**
```
GET /api/v1/machines/search?status=Active&province=DKI Jakarta&city_regency=Jakarta Pusat
```

**Response:** Same as Get All Machines (filtered list)

#### 3.6 Get Machine Metadata - Adaptive System
Retrieve all valid values for machine fields directly from the database.

**Endpoint:** `GET /api/v1/machines/metadata`
**Authentication:** Required

**Features:**
- Fully adaptive - automatically discovers new values from database
- Cached for performance - results cached for 1 hour
- Geographic mapping - FLM providers mapped to service areas
- Real-time accuracy - always reflects actual database content

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Machine metadata retrieved successfully from database",
  "last_updated": "2024-01-15T10:30:00Z",
  "slms": [
    {
      "code": "KGP - WINCOR DW",
      "description": "KGP - WINCOR DW",
      "is_documented": true
    }
  ],
  "flms": [
    {
      "code": "AVT - BANDUNG",
      "description": "AVT - BANDUNG",
      "area": "BANDUNG",
      "is_documented": true
    }
  ],
  "nets": [
    {
      "code": "NOSAIRIS",
      "description": "NOSAIRIS",
      "is_documented": true
    }
  ],
  "flm_names": [
    {
      "code": "AVT",
      "description": "AVT",
      "is_documented": true
    }
  ]
}
```

**Understanding the Response:**
- **SLMs** (Second Level Maintenance): All SLM provider/type combinations
- **FLMs** (Field Maintenance): All FLM providers with their service areas
- **NETs**: Network providers
- **FLM Names**: Provider names (AVT, ABS, BRS, TAG)

**FLM Geographic Areas:**
The system maps 58 FLM providers to their service areas including:
- **AVT**: 36 locations (Bandung, Jakarta, Surabaya, Bali, Medan, etc.)
- **ABS**: 13 locations (Jakarta, Bandung, Surabaya, Bali, etc.)
- **BRS**: 7 locations (Bandung, Bogor, Surabaya, etc.)
- **TAG**: 1 location (Cimone)

#### 3.7 Update Machine Status
Update the operational status of a machine.

**Endpoint:** `PATCH /api/v1/machines/status`
**Authentication:** Required

**Request Body:**
```json
{
  "terminal_id": "ATM-001",
  "status": "Maintenance",
  "gps": "-6.200000,106.816666",
  "lat": -6.200000,
  "lon": 106.816666
}
```

**Field Validation:**
- `terminal_id` (required, string) - Terminal to update
- `status` (required, string) - New operational status
- `gps` (optional, string) - GPS coordinates
- `lat` (optional, float64) - Latitude
- `lon` (optional, float64) - Longitude

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Machine status updated successfully",
  "data": {
    "terminal_id": "ATM-001",
    "status": "Maintenance",
    ...
  }
}
```

---

### 4. Admin Authentication Endpoints

Admin endpoints use session-based authentication. Login to receive a session token, then pass it via the `X-Session-Token` header or `session_token` cookie.

#### 4.1 Admin Login
Authenticate an admin user and receive a session token.

**Endpoint:** `POST /api/v1/admin/auth/login`
**Authentication:** Not required

**Request Body:**
```json
{
  "username": "admin",
  "password": "your_password"
}
```

**Field Validation:**
- `username` (required, string, min 3 chars)
- `password` (required, string, min 6 chars)

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "session_token": "sess_abc123xyz456",
  "user": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "full_name": "Administrator",
    "role": "super_admin",
    "is_active": true,
    "last_login_at": "2024-01-15T10:30:00Z",
    "created_at": "2024-01-01T00:00:00Z"
  },
  "expires_at": "2024-01-16T10:30:00Z"
}
```

**Response (401 Unauthorized):**
```json
{
  "success": false,
  "message": "Invalid credentials"
}
```

**Note:** The session token is also set as an HTTP-only cookie (`session_token`).

#### 4.2 Admin Logout
Invalidate the current admin session.

**Endpoint:** `POST /api/v1/admin/auth/logout`
**Authentication:** Session token (header or cookie)

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

#### 4.3 Get Current User
Get details of the currently logged-in admin.

**Endpoint:** `GET /api/v1/admin/auth/me`
**Authentication:** Session token required

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "admin",
    "role": "super_admin"
  }
}
```

---

### 5. Token Management Endpoints

All token management endpoints require admin session authentication via `X-Session-Token` header or `session_token` cookie.

#### 5.1 List All Tokens
Retrieve all API tokens (token values are masked).

**Endpoint:** `GET /api/v1/admin/tokens`
**Authentication:** Admin session required

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Tokens retrieved successfully",
  "data": [
    {
      "id": 1,
      "token": "tok_****xyz4",
      "name": "Production App",
      "description": "Main production API token",
      "token_prefix": "tok",
      "scopes": "[\"tickets:read\",\"tickets:write\",\"machines:read\"]",
      "environment": "production",
      "is_active": true,
      "rate_limit_per_minute": 60,
      "rate_limit_per_hour": 1000,
      "rate_limit_per_day": 10000,
      "expires_at": "2025-12-31T23:59:59Z",
      "last_used_at": "2024-01-15T10:30:00Z",
      "total_requests": 15234,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1
}
```

#### 5.2 Create Token
Create a new API token. The full token value is only shown once.

**Endpoint:** `POST /api/v1/admin/tokens`
**Authentication:** Admin session required

**Request Body:**
```json
{
  "name": "Production App",
  "description": "API token for production application",
  "environment": "production",
  "scopes": ["tickets:read", "tickets:write", "machines:read"],
  "ip_whitelist": ["192.168.1.0/24", "10.0.0.0/8"],
  "allowed_origins": ["https://your-app.com"],
  "rate_limit_per_minute": 60,
  "rate_limit_per_hour": 1000,
  "rate_limit_per_day": 10000,
  "expires_at": "2025-12-31T23:59:59Z"
}
```

**Field Validation:**
- `name` (required, string, 3-200 chars) - Token name
- `environment` (required, string) - One of: `production`, `staging`, `development`, `test`
- `description` (optional, string) - Token description
- `scopes` (optional, array of strings) - Permission scopes
- `ip_whitelist` (optional, array of strings) - Allowed IP addresses/ranges
- `allowed_origins` (optional, array of strings) - Allowed CORS origins
- `rate_limit_per_minute` (optional, int) - Max requests per minute
- `rate_limit_per_hour` (optional, int) - Max requests per hour
- `rate_limit_per_day` (optional, int) - Max requests per day
- `expires_at` (optional, datetime) - Token expiration date

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Token created successfully",
  "data": {
    "id": 1,
    "token": "tok_live_abc123xyz456789",
    "name": "Production App",
    "environment": "production",
    ...
  },
  "warning": "Save this token securely - it won't be shown again!"
}
```

#### 5.3 Get Token
Get details of a specific token.

**Endpoint:** `GET /api/v1/admin/tokens/:id`
**Authentication:** Admin session required

**Parameters:**
- `id` (path, int) - Token ID

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "token": "tok_****xyz4",
    "name": "Production App",
    ...
  }
}
```

#### 5.4 Update Token
Update an existing token's settings.

**Endpoint:** `PUT /api/v1/admin/tokens/:id`
**Authentication:** Admin session required

**Parameters:**
- `id` (path, int) - Token ID

**Request Body (all fields optional):**
```json
{
  "name": "Updated Name",
  "description": "Updated description",
  "scopes": ["tickets:read"],
  "ip_whitelist": ["10.0.0.0/8"],
  "allowed_origins": ["https://new-app.com"],
  "rate_limit_per_minute": 120,
  "rate_limit_per_hour": 2000,
  "rate_limit_per_day": 20000,
  "expires_at": "2026-12-31T23:59:59Z"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Token updated successfully",
  "data": { ... }
}
```

#### 5.5 Delete Token
Permanently delete an API token.

**Endpoint:** `DELETE /api/v1/admin/tokens/:id`
**Authentication:** Admin session required

**Parameters:**
- `id` (path, int) - Token ID

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Token deleted successfully"
}
```

#### 5.6 Disable Token
Temporarily disable a token without deleting it.

**Endpoint:** `PATCH /api/v1/admin/tokens/:id/disable`
**Authentication:** Admin session required

**Parameters:**
- `id` (path, int) - Token ID

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Token disabled successfully"
}
```

#### 5.7 Enable Token
Re-enable a previously disabled token.

**Endpoint:** `PATCH /api/v1/admin/tokens/:id/enable`
**Authentication:** Admin session required

**Parameters:**
- `id` (path, int) - Token ID

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Token enabled successfully"
}
```

#### 5.8 Get Token Usage Logs
Get detailed access logs for a specific token.

**Endpoint:** `GET /api/v1/admin/tokens/:id/logs`
**Authentication:** Admin session required

**Parameters:**
- `id` (path, int) - Token ID

**Query Parameters:**
- `limit` (optional, int, default 100) - Number of log entries to return

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "token_id": 1,
      "method": "GET",
      "endpoint": "/api/v1/tickets",
      "full_url": "/api/v1/tickets?status=0.NEW",
      "status_code": 200,
      "response_time_ms": 45,
      "ip_address": "192.168.1.100",
      "user_agent": "curl/7.88.0",
      "request_id": "req_abc123",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "total": 1
}
```

---

### 6. Analytics Endpoints

All analytics endpoints require admin session authentication.

#### 6.1 Dashboard Statistics
Get overview statistics for the admin dashboard.

**Endpoint:** `GET /api/v1/admin/analytics/dashboard`
**Authentication:** Admin session required

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "total_tokens": 10,
    "active_tokens": 8,
    "total_requests_24h": 15234,
    "success_rate": 98.5,
    "avg_response_time_ms": 45.2,
    "top_tokens": [
      {
        "token_id": 1,
        "token_name": "Production App",
        "total_requests": 8000,
        "successful_requests": 7900,
        "failed_requests": 100,
        "avg_response_time_ms": 42.5
      }
    ],
    "recent_activity": [ ... ]
  }
}
```

#### 6.2 Token Analytics
Get detailed analytics for a specific token.

**Endpoint:** `GET /api/v1/admin/analytics/tokens/:id`
**Authentication:** Admin session required

**Parameters:**
- `id` (path, int) - Token ID

**Query Parameters:**
- `days` (optional, int, default 7) - Number of days to analyze

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "token_id": 1,
    "token_name": "Production App",
    "total_requests": 8000,
    "successful_requests": 7900,
    "failed_requests": 100,
    "client_errors": 80,
    "server_errors": 20,
    "avg_response_time_ms": 42.5,
    "max_response_time_ms": 500,
    "unique_ips": 5,
    "unique_endpoints": 12,
    "last_used_at": "2024-01-15T10:30:00Z"
  }
}
```

#### 6.3 Endpoint Statistics
Get usage statistics grouped by endpoint.

**Endpoint:** `GET /api/v1/admin/analytics/endpoints`
**Authentication:** Admin session required

**Query Parameters:**
- `days` (optional, int, default 7) - Number of days to analyze
- `limit` (optional, int, default 20) - Number of endpoints to return

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "endpoint": "/api/v1/tickets",
      "method": "GET",
      "request_count": 5000,
      "unique_tokens": 5,
      "avg_response_time_ms": 40.2,
      "successful_requests": 4950,
      "failed_requests": 50
    }
  ]
}
```

#### 6.4 Daily Usage
Get daily aggregated usage data.

**Endpoint:** `GET /api/v1/admin/analytics/daily`
**Authentication:** Admin session required

**Query Parameters:**
- `days` (optional, int, default 30) - Number of days to return
- `token_id` (optional, int) - Filter by specific token

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "date": "2024-01-15",
      "token_id": 1,
      "token_name": "Production App",
      "request_count": 1500,
      "successful_requests": 1480,
      "failed_requests": 20,
      "avg_response_time_ms": 43.1
    }
  ]
}
```

---

### 7. Audit Logs

#### 7.1 Get Audit Logs
Retrieve administrative audit logs for compliance tracking.

**Endpoint:** `GET /api/v1/admin/audit-logs`
**Authentication:** Admin session required

**Query Parameters:**
- `limit` (optional, int, default 100) - Number of log entries to return

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "admin_user_id": 1,
      "action": "token.created",
      "resource_type": "api_token",
      "resource_id": 5,
      "old_values": null,
      "new_values": "{\"name\":\"New Token\",\"environment\":\"production\"}",
      "ip_address": "192.168.1.100",
      "user_agent": "Mozilla/5.0 ...",
      "description": "Created API token 'New Token'",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "total": 1
}
```

---

## Data Models

### Ticket Model (OpenTicket)

| Field | Type | Nullable | Description |
|-------|------|----------|-------------|
| terminal_id | string | No | Terminal identifier |
| terminal_name | string | No | Terminal name |
| priority | string | Yes | Priority level (1.High, 2.Middle, 3.Low, 4.Minimum) |
| mode | string | Yes | Terminal mode (Closed, In Service, nan, Off-line, Supervisor) |
| initial_problem | string | Yes | Initial problem description |
| current_problem | string | Yes | Current problem description |
| p_duration | string | Yes | Problem duration |
| incident_start_datetime | string | Yes | Incident start timestamp |
| count | int | No | Count value |
| status | string | Yes | Ticket status (0.NEW, 1.Req FD ke HD, etc.) |
| remarks | string | Yes | Remarks/notes |
| balance | int | No | Balance amount |
| condition | string | Yes | Condition status |
| tickets_no | string | Yes | Ticket number |
| tickets_duration | float64 | No | Ticket duration in minutes (supports decimals) |
| open_time | string | Yes | Ticket open time |
| close_time | string | Yes | Ticket close time |
| problem_history | string | Yes | Problem history |
| mode_history | string | Yes | Mode change history |
| dsp_flm | string | Yes | DSP FLM identifier |
| dsp_slm | string | Yes | DSP SLM identifier |
| last_withdrawal | timestamp | Yes | Last withdrawal timestamp |
| export_name | string | Yes | Export name for reports |

### Machine Model (ATMI)

| Field | Type | Nullable | Description |
|-------|------|----------|-------------|
| terminal_id | string | No | Unique terminal identifier |
| store | string | No | Store name |
| store_code | string | No | Store code |
| store_name | string | No | Full store name |
| date_of_activation | timestamp | Yes | Activation date |
| status | string | No | Operational status |
| std | int | No | Standard value |
| gps | string | No | GPS coordinates |
| lat | float64 | No | Latitude |
| lon | float64 | No | Longitude |
| province | string | No | Province name |
| city_regency | string | No | City or regency name |
| district | string | No | District name |

### API Token Model

| Field | Type | Nullable | Description |
|-------|------|----------|-------------|
| id | int | No | Token ID |
| token | string | No | Token value (masked after creation) |
| name | string | No | Token name (3-200 chars) |
| description | string | Yes | Token description |
| token_prefix | string | No | Token prefix (e.g., "tok") |
| scopes | string (JSON) | Yes | Permission scopes array |
| permissions | string (JSON) | Yes | Permission object |
| environment | string | No | Environment: production, staging, development, test |
| is_active | bool | No | Whether token is active |
| ip_whitelist | string (JSON) | Yes | Allowed IP addresses |
| allowed_origins | string (JSON) | Yes | Allowed CORS origins |
| rate_limit_per_minute | int | No | Max requests per minute |
| rate_limit_per_hour | int | No | Max requests per hour |
| rate_limit_per_day | int | No | Max requests per day |
| expires_at | timestamp | Yes | Token expiration date |
| last_used_at | timestamp | Yes | Last usage timestamp |
| last_used_ip | string | Yes | Last used IP address |
| last_used_endpoint | string | Yes | Last used endpoint |
| total_requests | int64 | No | Total request count |
| created_at | timestamp | No | Creation timestamp |
| updated_at | timestamp | No | Last update timestamp |
| created_by | int | Yes | Admin user ID who created |
| revoked_at | timestamp | Yes | Revocation timestamp |
| revoked_by | int | Yes | Admin user ID who revoked |
| revoked_reason | string | Yes | Revocation reason |

### Admin User Model

| Field | Type | Nullable | Description |
|-------|------|----------|-------------|
| id | int | No | Admin user ID |
| username | string | No | Username (min 3 chars) |
| email | string | No | Email address |
| full_name | string | Yes | Full name |
| role | string | No | Role: super_admin, admin, viewer |
| is_active | bool | No | Whether user is active |
| last_login_at | timestamp | Yes | Last login timestamp |
| last_login_ip | string | Yes | Last login IP address |
| created_at | timestamp | No | Creation timestamp |
| updated_at | timestamp | No | Last update timestamp |

### Token Usage Log Model

| Field | Type | Nullable | Description |
|-------|------|----------|-------------|
| id | int64 | No | Log entry ID |
| token_id | int | No | Associated token ID |
| method | string | No | HTTP method (GET, POST, etc.) |
| endpoint | string | No | Request endpoint path |
| full_url | string | Yes | Full request URL with query params |
| status_code | int | No | HTTP response status code |
| response_time_ms | int | Yes | Response time in milliseconds |
| ip_address | string | No | Client IP address |
| user_agent | string | Yes | Client user agent |
| referer | string | Yes | Request referer |
| request_id | string | Yes | Unique request ID |
| request_body_size | int | Yes | Request body size in bytes |
| response_body_size | int | Yes | Response body size in bytes |
| error_message | string | Yes | Error message (if failed) |
| error_code | string | Yes | Error code (if failed) |
| created_at | timestamp | No | Log timestamp |

### Audit Log Model

| Field | Type | Nullable | Description |
|-------|------|----------|-------------|
| id | int64 | No | Audit log ID |
| admin_user_id | int | Yes | Admin who performed the action |
| action | string | No | Action performed (e.g., "token.created") |
| resource_type | string | No | Resource type (e.g., "api_token") |
| resource_id | int | Yes | ID of affected resource |
| old_values | string (JSON) | Yes | Previous values before change |
| new_values | string (JSON) | Yes | New values after change |
| ip_address | string | Yes | Admin's IP address |
| user_agent | string | Yes | Admin's user agent |
| description | string | Yes | Human-readable description |
| created_at | timestamp | No | Audit log timestamp |

---

## Hybrid Adaptive Metadata System

### How It Works

The API Gateway uses a **hybrid database-driven approach** that combines:

1. **Real-time Database Queries** - Queries DISTINCT values from actual database
2. **Optional Documentation** - Provides descriptions for known values
3. **Intelligent Caching** - Caches results for 1 hour for performance
4. **Graceful Degradation** - Unknown values still work (use code as description)

### Architecture

```
Client Request -> Handler -> Service (Check Cache)
                              |
                         Cache Hit? Yes -> Return Cached Data
                              | No
                         Query Database (SELECT DISTINCT)
                              |
                         Combine with Descriptions
                              |
                         Update Cache
                              |
                         Return Results
```

### Performance Characteristics

| Metric | Value |
|--------|-------|
| Cache Duration | 1 hour (configurable) |
| Cache Hit Speed | 1-5ms |
| Database Query | 50-100ms |
| Thread Safety | Yes (RWMutex) |
| Expected Cache Hit Rate | ~95% |

### Adaptation Process

**When new values are added to the database:**

1. **Immediate**: Value works in all API operations (create, update, filter)
2. **Within 1 hour**: Value appears in metadata endpoint (after cache expires)
3. **No code changes needed**: System automatically discovers new values
4. **Optional documentation**: Can add description to code later

**Example:**
```sql
-- Someone adds new status to database
INSERT INTO open_ticket (status, ...) VALUES ('9.Emergency', ...)

-- After cache expires (max 1 hour):
GET /api/v1/tickets/metadata
-- Returns: {"code": "9.Emergency", "description": "9.Emergency", "is_documented": false}

-- Value works immediately in all operations
POST /api/v1/tickets {"status": "9.Emergency", ...}
GET /api/v1/tickets/status/9.Emergency
```

---

## Token Management System

### Overview

The token management system provides enterprise-grade API access control:

- **Create tokens** with granular scopes and rate limits
- **Monitor usage** with real-time analytics and daily breakdowns
- **Control access** with IP whitelisting and environment separation
- **Track everything** with comprehensive audit logging
- **Manage lifecycle** with expiration, disable/enable, and revocation

### Admin Roles

| Role | Permissions |
|------|-------------|
| `super_admin` | Full access to all admin features |
| `admin` | Token management, analytics, audit logs |
| `viewer` | Read-only access to dashboard and analytics |

### Token Lifecycle

```
Created -> Active -> [Disabled] -> [Re-enabled] -> [Expired/Revoked/Deleted]
```

### Rate Limiting

Rate limits are enforced per-token at three levels:
- **Per minute** - Burst protection
- **Per hour** - Sustained usage control
- **Per day** - Daily quota enforcement

When a rate limit is exceeded, the API returns `429 Too Many Requests`:
```json
{
  "success": false,
  "message": "Rate limit exceeded: minute limit (60 requests)",
  "error": "Please slow down your requests"
}
```

---

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid input data |
| 401 | Unauthorized - Missing or invalid API key/token |
| 404 | Not Found - Resource not found |
| 409 | Conflict - Resource already exists (duplicate ticket number) |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error - Server error |
| 503 | Service Unavailable - Database connection or token service unavailable |

---

## Examples

### cURL Examples

**Get ticket metadata:**
```bash
curl -X GET \
  -H "X-API-Token: tok_live_abc123" \
  http://localhost:8080/api/v1/tickets/metadata
```

**Get all tickets:**
```bash
curl -X GET \
  -H "X-API-Token: tok_live_abc123" \
  http://localhost:8080/api/v1/tickets
```

**Create a ticket:**
```bash
curl -X POST \
  -H "X-API-Token: tok_live_abc123" \
  -H "Content-Type: application/json" \
  -d '{
    "terminal_id": "ATM-001",
    "terminal_name": "Main Branch ATM",
    "priority": "1.High",
    "mode": "Off-line",
    "initial_problem": "Card reader not responding",
    "status": "0.NEW",
    "tickets_no": "TKT-2024-001"
  }' \
  http://localhost:8080/api/v1/tickets
```

**Get tickets by status:**
```bash
curl -X GET \
  -H "X-API-Token: tok_live_abc123" \
  "http://localhost:8080/api/v1/tickets/status/0.NEW"
```

**Update ticket:**
```bash
curl -X PUT \
  -H "X-API-Token: tok_live_abc123" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "6.Follow-up Sales team",
    "remarks": "Technician dispatched",
    "condition": "Normal"
  }' \
  http://localhost:8080/api/v1/tickets/ATM-001
```

**Search machines:**
```bash
curl -X GET \
  -H "X-API-Token: tok_live_abc123" \
  "http://localhost:8080/api/v1/machines/search?status=Active&province=DKI%20Jakarta"
```

**Admin login:**
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "your_password"
  }' \
  http://localhost:8080/api/v1/admin/auth/login
```

**Create API token (admin):**
```bash
curl -X POST \
  -H "X-Session-Token: sess_abc123" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Production App",
    "environment": "production",
    "scopes": ["tickets:read", "tickets:write", "machines:read"],
    "rate_limit_per_minute": 60,
    "rate_limit_per_hour": 1000,
    "rate_limit_per_day": 10000
  }' \
  http://localhost:8080/api/v1/admin/tokens
```

**Get dashboard analytics (admin):**
```bash
curl -X GET \
  -H "X-Session-Token: sess_abc123" \
  http://localhost:8080/api/v1/admin/analytics/dashboard
```

**Get audit logs (admin):**
```bash
curl -X GET \
  -H "X-Session-Token: sess_abc123" \
  "http://localhost:8080/api/v1/admin/audit-logs?limit=50"
```

### JavaScript/Fetch Examples

**Get ticket metadata:**
```javascript
fetch('http://localhost:8080/api/v1/tickets/metadata', {
  headers: {
    'X-API-Token': 'tok_live_abc123'
  }
})
  .then(response => response.json())
  .then(data => {
    console.log('Statuses:', data.statuses);
    console.log('Modes:', data.modes);
    console.log('Priorities:', data.priorities);
  });
```

**Create ticket:**
```javascript
fetch('http://localhost:8080/api/v1/tickets', {
  method: 'POST',
  headers: {
    'X-API-Token': 'tok_live_abc123',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    terminal_id: 'ATM-001',
    terminal_name: 'Main Branch ATM',
    priority: '1.High',
    mode: 'Off-line',
    initial_problem: 'Card reader issue',
    status: '0.NEW',
    tickets_no: 'TKT-2024-001'
  })
})
  .then(response => response.json())
  .then(data => console.log('Created:', data));
```

**Admin login and token creation:**
```javascript
// Login
const loginResp = await fetch('http://localhost:8080/api/v1/admin/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ username: 'admin', password: 'your_password' })
});
const { session_token } = await loginResp.json();

// Create token
const tokenResp = await fetch('http://localhost:8080/api/v1/admin/tokens', {
  method: 'POST',
  headers: {
    'X-Session-Token': session_token,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    name: 'My App',
    environment: 'production',
    rate_limit_per_minute: 60
  })
});
const { data: newToken, warning } = await tokenResp.json();
console.log('New token:', newToken.token); // Save this!
console.log(warning); // "Save this token securely - it won't be shown again!"
```

---

## Best Practices

### 1. Use Metadata Endpoints
Always use metadata endpoints to discover valid field values rather than hard-coding them:

```javascript
// Good: Dynamic values from metadata
const metadata = await fetch('/api/v1/tickets/metadata');
const statuses = metadata.statuses.map(s => s.code);

// Bad: Hard-coded values
const statuses = ['Open', 'Closed']; // Will be outdated
```

### 2. Handle Nullable Fields
Many fields can be NULL in the database. Always handle null values:

```javascript
const priority = ticket.priority || 'Not Set';
const closeTime = ticket.close_time || 'Still Open';
```

### 3. Use API Tokens Over API Keys
Prefer `X-API-Token` over `X-API-Key` for better security, analytics, and rate limiting:

```bash
# Recommended
curl -H "X-API-Token: tok_live_abc123" http://localhost:8080/api/v1/tickets

# Legacy (still supported)
curl -H "X-API-Key: your_key" http://localhost:8080/api/v1/tickets
```

### 4. Handle Rate Limiting
When using API tokens, handle `429` responses gracefully:

```javascript
const response = await fetch('/api/v1/tickets', {
  headers: { 'X-API-Token': token }
});

if (response.status === 429) {
  // Back off and retry after a delay
  await new Promise(resolve => setTimeout(resolve, 5000));
  // Retry request...
}
```

### 5. Use Proper Status Codes
Always check the `success` field and HTTP status code:

```javascript
const response = await fetch('/api/v1/tickets');
const data = await response.json();

if (data.success) {
  console.log('Tickets:', data.data);
} else {
  console.error('Error:', data.message);
}
```

### 6. Handle Undocumented Values
Use the `is_documented` flag to identify new values:

```javascript
metadata.statuses.forEach(status => {
  if (!status.is_documented) {
    console.warn(`New status discovered: ${status.code}`);
  }
});
```

---

## Database Schema Reference

### Ticket Table (ticket_master.dbo.open_ticket)

**Database Column Names** (with spaces):
- `[Terminal ID]`
- `[Terminal Name]`
- `[Priority]`
- `[Mode]`
- `[Initial Problem]`
- `[Current Problem]`
- `[P-Duration]`
- `[Incident start datetime]`
- `[Count]`
- `[Status]`
- `[Remarks]`
- `[Balance]`
- `[Condition]`
- `[Tickets no]`
- `[Tickets duration]` - Stored as float64 (supports decimal values like 4.6, 150.5)
- `[Open time]`
- `[Close time]`
- `[Problem History]`
- `[Mode History]`
- `[DSP FLM]`
- `[DSP SLM]`
- `[Last Withdrawal]`
- `[Export Name]`

**Note:** Many fields are nullable in the database.

### Machine Table (machine_master.dbo.atmi)

**Database Column Names**:
- `terminal_id`
- `store`
- `store_code`
- `store_name`
- `date_of_activation`
- `status`
- `std`
- `gps`
- `lat`
- `lon`
- `province`
- `city/regency` (note: has slash in column name)
- `district`
- `slm`
- `flm`
- `net`
- `flm_name`

### Token Management Tables (token_management database)

- `admin_users` - Admin user accounts
- `admin_sessions` - Active admin sessions
- `api_tokens` - API token definitions and metadata
- `token_usage_logs` - Per-request usage tracking
- `token_rate_limits` - Rate limit counters by time window
- `audit_logs` - Administrative action audit trail

---

## Versioning

**Current API version:** v1

All endpoints are prefixed with `/api/v1` to support future versioning.

**Version Strategy:**
- Breaking changes will increment the major version (v2, v3, etc.)
- Backward-compatible changes will be added to existing versions
- Old versions will be maintained for minimum 6 months after new version release

---

## Changelog

### Version 1.1.0 (Current)

**New Features:**
- Token management system with admin dashboard
- Admin authentication with session management
- API token creation with scopes, rate limits, and IP whitelisting
- Usage analytics (dashboard stats, endpoint stats, daily usage)
- Audit logging for all administrative actions
- Combined authentication (X-API-Token with rate limiting and usage tracking)
- Gzip response compression
- Token disable/enable lifecycle management

**Improvements:**
- Health check now reports token_database status
- Non-fatal database connections (app starts even if a DB is unavailable)
- Asynchronous usage log recording (non-blocking)

### Version 1.0.0

**Features:**
- Complete CRUD operations for tickets and machines
- Hybrid adaptive metadata system with database queries
- Intelligent 1-hour caching for optimal performance
- Thread-safe concurrent request handling
- Comprehensive Swagger/OpenAPI documentation
- Support for NULL database values with proper JSON marshaling
- Float64 support for ticket duration field
- Geographic area mapping for FLM providers
- Standardized error responses
- API key authentication
- Health check endpoints
- Advanced search and filtering

**Database Support:**
- SQL Server (ticket_master and machine_master databases)
- Proper handling of database column names with spaces
- Support for nullable fields

**Performance:**
- Metadata caching (1-hour TTL)
- Efficient database queries with proper indexing
- Thread-safe operations with RWMutex

---

*API Gateway Version: 1.1.0*
*Documentation Version: 1.1.0*
