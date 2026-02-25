# API Gateway - Complete API Documentation

## Overview

This API Gateway serves as middleware between on-premise databases (`ticket_master`, `machine_master`, and `token_management`) and cloud applications. It exposes a **single unified endpoint** `/api/v1/data` that always JOINs `ticket_master.dbo.open_ticket` with `machine_master.dbo.machine`, applying vendor-scoped access control from the API token.

**Base URL:** `http://localhost:8080`
**API Version:** v1
**API Prefix:** `/api/v1`
**Swagger UI:** `http://localhost:8080/swagger/index.html`
**Admin Dashboard:** `http://localhost:8080/admin`

---

## Key Features

- **Single Unified Endpoint** – `/api/v1/data` always returns joined ticket + machine data
- **Vendor-Scoped Token System** – Each API token can be restricted to a specific vendor via `filter_column` / `filter_value` (e.g. `mm.[FLM name] = 'AVT'`)
- **Admin / Internal Token Support** – Tokens with `is_super_token = true` bypass all vendor filters and use a customizable admin query (`repository/queries/admin_data_query.go`)
- **Metadata Discovery** – Automatically discovers distinct status, mode, and priority values from the database (1-hour cache)
- **Thread-Safe Operations** – Production-ready concurrent request handling
- **Standardized Error Responses** – Consistent JSON error format across all endpoints
- **Admin Dashboard** – Web UI for managing tokens, viewing analytics, and audit logs
- **Gzip Compression** – Automatic response compression
- **Analytics & Monitoring** – Dashboard stats, endpoint analytics, daily usage tracking
- **Audit Logging** – Complete history of all administrative actions

---

## Authentication

All `/api/v1/data` endpoints require an `X-API-Token` header.

```
X-API-Token: tok_live_your_token_here
```

Two token types exist:

| Token type | Behavior |
|---|---|
| **Vendor token** (`filter_column` set) | Only returns/modifies rows matching `mm.[col] = value` |
| **Admin / Internal token** (`is_super_token = true`) | Full access — uses the customizable admin query |

Admin dashboard endpoints use session-based auth (cookie / `X-Session-Token`).

---

## Endpoint Reference

### Health

#### `GET /health`
Check API and database connection status. No auth required.

**Response 200:**
```json
{
  "status": "healthy",
  "databases": { "ticket_db": "ok", "machine_db": "ok", "token_db": "ok" }
}
```

#### `GET /ping`
Simple liveness check.

**Response 200:**
```json
{ "message": "pong" }
```

---

### Data (`/api/v1/data`)

All data endpoints require `X-API-Token`.

#### `GET /api/v1/data`
Retrieve all joined ticket + machine rows.

- Vendor tokens: rows filtered by their `filter_column` / `filter_value`
- Admin / Internal tokens: all rows via the customizable query in `repository/queries/admin_data_query.go`

**Query parameters:**

| Parameter | Type | Description |
|---|---|---|
| `page` | integer | Page number (omit for all results) |
| `page_size` | integer | Items per page (default: 100, max: 500) |

**Response 200:**
```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": [ { ...DataRow... } ],
  "total": 250,
  "page": 1,
  "page_size": 100,
  "total_pages": 3
}
```

---

#### `GET /api/v1/data/metadata`
Retrieve distinct valid values for `status`, `mode`, and `priority` fields. Cached for 1 hour.

**Response 200:**
```json
{
  "success": true,
  "message": "Metadata retrieved successfully",
  "statuses": [
    { "code": "0.NEW", "description": "New ticket", "is_documented": true }
  ],
  "modes": [
    { "code": "Off-line", "description": "Terminal is offline", "is_documented": true }
  ],
  "priorities": [
    { "code": "1.High", "description": "High priority", "is_documented": true }
  ],
  "last_updated": "2024-01-15T10:30:00Z"
}
```

---

#### `GET /api/v1/data/:terminal_id`
Retrieve a single joined row by terminal ID.

- Vendor tokens: returns 404 if the terminal is outside their scope
- Admin / Internal tokens: returns any terminal

**Response 200:**
```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": { ...DataRow... }
}
```

**Response 404:** Terminal not found or outside vendor scope.

---

#### `PUT /api/v1/data/:terminal_id`
Update ticket fields for a terminal.

- Vendor tokens: returns 403 if the terminal is outside their scope
- Admin / Internal tokens: can update any terminal

**Request body (all fields optional):**
```json
{
  "priority": "1.High",
  "mode": "Off-line",
  "current_problem": "Card reader fixed",
  "status": "2.Kirim FLM",
  "remarks": "Technician dispatched",
  "condition": "Normal",
  "close_time": "2024-01-15 18:00:00",
  "problem_history": "Card reader issue resolved",
  "mode_history": "Online->Offline->Online"
}
```

**Response 200:**
```json
{
  "success": true,
  "message": "Updated successfully",
  "data": { ...DataRow... }
}
```

**Error responses:**

| Status | Condition |
|---|---|
| 400 | No fields provided / invalid JSON |
| 403 | Terminal outside vendor token scope |
| 404 | Terminal not found |
| 500 | Database error |

---

### DataRow Schema

Every data response returns `DataRow` objects with the following fields:

| Field | Source | Type | Description |
|---|---|---|---|
| `terminal_id` | `op.[Terminal ID]` | string | Primary identifier |
| `terminal_name` | `op.[Terminal Name]` | string | |
| `priority` | `op.[Priority]` | string | e.g. `1.High` |
| `mode` | `op.[Mode]` | string | e.g. `Off-line` |
| `initial_problem` | `op.[Initial Problem]` | string | |
| `current_problem` | `op.[Current Problem]` | string | |
| `p_duration` | `op.[P Duration]` | string | |
| `incident_start_datetime` | `op.[Incident start datetime]` | string | |
| `count` | `op.[Count]` | integer | |
| `status` | `op.[Status]` | string | e.g. `0.NEW` |
| `remarks` | `op.[Remarks]` | string | |
| `balance` | `op.[Balance]` | integer | |
| `condition` | `op.[Condition]` | string | |
| `tickets_no` | `op.[Tickets No]` | string | |
| `tickets_duration` | `op.[Tickets Duration]` | number | |
| `open_time` | `op.[Open Time]` | string | |
| `close_time` | `op.[Close Time]` | string | |
| `problem_history` | `op.[Problem History]` | string | |
| `mode_history` | `op.[Mode History]` | string | |
| `dsp_flm` | `op.[DSP FLM]` | string | |
| `dsp_slm` | `op.[DSP SLM]` | string | |
| `last_withdrawal` | `op.[Last Withdrawal]` | string | |
| `export_name` | `op.[Export Name]` | string | |
| `flm_name` | `mm.[FLM name]` | string | From machine_master JOIN |
| `flm` | `mm.[FLM]` | string | From machine_master JOIN |
| `slm` | `mm.[SLM]` | string | From machine_master JOIN |
| `net` | `mm.[Net]` | string | From machine_master JOIN |

Nullable fields are omitted from the JSON response when NULL.

---

### Admin Auth (`/api/v1/admin/auth`)

#### `POST /api/v1/admin/auth/login`
Authenticate as admin and receive a session token.

**Request:**
```json
{ "username": "admin", "password": "yourpassword" }
```

**Response 200:**
```json
{
  "success": true,
  "session_token": "sess_abc123...",
  "expires_at": "2024-01-16T10:30:00Z",
  "user": { "id": 1, "username": "admin", "role": "super_admin", ... }
}
```

#### `POST /api/v1/admin/auth/logout`
Invalidate the current session.

#### `GET /api/v1/admin/auth/me`
Get details of the currently authenticated admin.

---

### Token Management (`/api/v1/admin/tokens`)

All token management endpoints require admin session auth.

#### `GET /api/v1/admin/tokens`
List all API tokens.

#### `POST /api/v1/admin/tokens`
Create a new API token.

**Request:**
```json
{
  "name": "AVT Vendor Token",
  "description": "Token for AVT vendor access",
  "environment": "production",
  "vendor_name": "AVT",
  "filter_column": "flm_name",
  "filter_value": "AVT",
  "is_super_token": false,
  "rate_limit_per_minute": 60,
  "rate_limit_per_hour": 1000,
  "rate_limit_per_day": 10000,
  "expires_at": "2025-12-31T23:59:59Z"
}
```

**Filter column options:**

| `filter_column` | SQL expression | Example `filter_value` |
|---|---|---|
| `flm_name` | `mm.[FLM name]` | `AVT` |
| `flm` | `mm.[FLM]` | `AVT - BANDUNG` |
| `slm` | `mm.[SLM]` | `KGP - WINCOR DW` |
| `net` | `mm.[Net]` | `NOSAIRIS` |
| `terminal_id` | `op.[Terminal ID]` | `ATM-001` |
| `status` | `op.[Status]` | `0.NEW` |
| `priority` | `op.[Priority]` | `1.High` |

Set `is_super_token: true` to create an Admin / Internal token (bypasses all vendor filters).

#### `GET /api/v1/admin/tokens/:id`
Get a specific token by ID.

#### `PUT /api/v1/admin/tokens/:id`
Update token details.

#### `DELETE /api/v1/admin/tokens/:id`
Permanently delete a token.

#### `PATCH /api/v1/admin/tokens/:id/disable`
Temporarily disable a token.

#### `PATCH /api/v1/admin/tokens/:id/enable`
Re-enable a disabled token.

#### `GET /api/v1/admin/tokens/:id/logs`
Get access logs for a specific token.

---

### Analytics (`/api/v1/admin/analytics`)

All analytics endpoints require admin session auth.

#### `GET /api/v1/admin/analytics/dashboard`
Overview statistics for the admin dashboard.

#### `GET /api/v1/admin/analytics/daily?days=30&token_id=1`
Daily request volume. Optionally filter by token ID.

#### `GET /api/v1/admin/analytics/endpoints?days=7&limit=20`
Usage statistics broken down by endpoint.

#### `GET /api/v1/admin/analytics/tokens/:id?days=7`
Detailed analytics for a specific token.

#### `GET /api/v1/admin/audit-logs?limit=100`
Administrative audit log entries.

---

## Error Response Format

All errors follow this structure:

```json
{
  "success": false,
  "message": "Human-readable description",
  "error": "detailed error information"
}
```

---

## Customizing the Admin Query

Admin / Internal tokens (`is_super_token = true`) use a dedicated query defined in:

```
repository/queries/admin_data_query.go
```

Edit the `AdminDataQuery` constant in that file to change which columns or ordering the admin view returns. The column order in the SELECT must match the `scanDataRow()` function in `repository/data_repository.go`.

---

## Environment Variables

| Variable | Description |
|---|---|
| `TICKET_DB_SERVER` | MS SQL Server host for ticket_master |
| `TICKET_DB_NAME` | Database name (ticket_master) |
| `TICKET_DB_USER` / `TICKET_DB_PASSWORD` | Credentials |
| `MACHINE_DB_SERVER` | MS SQL Server host for machine_master |
| `MACHINE_DB_NAME` | Database name (machine_master) |
| `TOKEN_DB_SERVER` | MS SQL Server host for token management |
| `TOKEN_DB_NAME` | Database name for tokens/sessions/analytics |
| `API_KEY` | Fallback static API key (legacy) |
| `PORT` | Server port (default: 8080) |
| `GIN_MODE` | `debug` or `release` |
