# API Gateway - Complete API Documentation

## Overview

This document provides detailed information about all available API endpoints in the API Gateway.

**Base URL:** `http://localhost:8080`
**API Version:** v1
**API Prefix:** `/api/v1`

---

## Authentication

All API endpoints (except health checks) require authentication using an API key.

**Header Required:**
```
X-API-Key: your_api_key_here
```

**Example:**
```bash
curl -H "X-API-Key: your_api_key" http://localhost:8080/api/v1/tickets
```

---

## Response Format

All API responses follow a standardized format:

### Success Response
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": { ... }
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
Check the overall health of the API and database connections.

**Endpoint:** `GET /health`
**Authentication:** Not required

**Response (200 OK):**
```json
{
  "status": "healthy",
  "message": "API Gateway is running",
  "services": {
    "ticket_database": "connected",
    "machine_database": "connected"
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
Retrieve all tickets from the system.

**Endpoint:** `GET /api/v1/tickets`
**Authentication:** Required

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Tickets retrieved successfully",
  "data": [
    {
      "id": 1,
      "ticket_number": "TKT-2024-001",
      "terminal_id": "ATM-001",
      "description": "Card reader malfunction",
      "priority": "High",
      "status": "Open",
      "category": "Hardware",
      "reported_by": "John Doe",
      "assigned_to": "Tech Team A",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z",
      "resolved_at": null,
      "resolution_notes": null
    }
  ],
  "total": 1
}
```

#### 2.2 Get Ticket by ID
Retrieve a specific ticket by its ID.

**Endpoint:** `GET /api/v1/tickets/:id`
**Authentication:** Required

**Parameters:**
- `id` (path, integer) - Ticket ID

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Ticket retrieved successfully",
  "data": {
    "id": 1,
    "ticket_number": "TKT-2024-001",
    "terminal_id": "ATM-001",
    "description": "Card reader malfunction",
    "priority": "High",
    "status": "Open",
    "category": "Hardware",
    "reported_by": "John Doe",
    "assigned_to": "Tech Team A",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "resolved_at": null,
    "resolution_notes": null
  }
}
```

**Response (404 Not Found):**
```json
{
  "success": false,
  "message": "Ticket not found",
  "data": null
}
```

#### 2.3 Get Ticket by Number
Retrieve a ticket by its unique ticket number.

**Endpoint:** `GET /api/v1/tickets/number/:number`
**Authentication:** Required

**Parameters:**
- `number` (path, string) - Ticket number (e.g., "TKT-2024-001")

**Response:** Same as Get Ticket by ID

#### 2.4 Get Tickets by Status
Retrieve all tickets with a specific status.

**Endpoint:** `GET /api/v1/tickets/status/:status`
**Authentication:** Required

**Parameters:**
- `status` (path, string) - Status value: `Open`, `InProgress`, `Pending`, `Resolved`

**Response:** Same as Get All Tickets (filtered list)

#### 2.5 Get Tickets by Terminal
Retrieve all tickets associated with a specific terminal.

**Endpoint:** `GET /api/v1/tickets/terminal/:terminal_id`
**Authentication:** Required

**Parameters:**
- `terminal_id` (path, string) - Terminal ID (e.g., "ATM-001")

**Response:** Same as Get All Tickets (filtered list)

#### 2.6 Create Ticket
Create a new ticket in the system.

**Endpoint:** `POST /api/v1/tickets`
**Authentication:** Required

**Request Body:**
```json
{
  "ticket_number": "TKT-2024-001",
  "terminal_id": "ATM-001",
  "description": "Card reader not responding to card insertions",
  "priority": "High",
  "category": "Hardware",
  "reported_by": "John Doe",
  "assigned_to": "Tech Team A"
}
```

**Field Validation:**
- `ticket_number` (required, string) - Unique ticket identifier
- `terminal_id` (required, string) - Associated terminal
- `description` (required, string) - Detailed issue description
- `priority` (required, string) - Must be one of: `Low`, `Medium`, `High`, `Critical`
- `category` (required, string) - Ticket category
- `reported_by` (required, string) - Name of person reporting
- `assigned_to` (optional, string) - Assigned technician/team

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Ticket created successfully",
  "data": {
    "id": 1,
    "ticket_number": "TKT-2024-001",
    "terminal_id": "ATM-001",
    "description": "Card reader not responding to card insertions",
    "priority": "High",
    "status": "Open",
    "category": "Hardware",
    "reported_by": "John Doe",
    "assigned_to": "Tech Team A",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "resolved_at": null,
    "resolution_notes": null
  }
}
```

**Response (409 Conflict):**
```json
{
  "success": false,
  "message": "ticket with this number already exists",
  "data": null
}
```

#### 2.7 Update Ticket
Update an existing ticket.

**Endpoint:** `PUT /api/v1/tickets/:id`
**Authentication:** Required

**Parameters:**
- `id` (path, integer) - Ticket ID to update

**Request Body (all fields optional):**
```json
{
  "status": "Resolved",
  "priority": "Medium",
  "assigned_to": "Tech Team B",
  "description": "Updated description",
  "resolution_notes": "Replaced card reader module. Tested successfully."
}
```

**Field Validation:**
- `status` (optional, string) - Must be one of: `Open`, `InProgress`, `Pending`, `Resolved`
- `priority` (optional, string) - Must be one of: `Low`, `Medium`, `High`, `Critical`
- `assigned_to` (optional, string) - Technician/team name
- `description` (optional, string) - Updated description
- `resolution_notes` (optional, string) - Resolution details

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Ticket updated successfully",
  "data": {
    "id": 1,
    "ticket_number": "TKT-2024-001",
    "status": "Resolved",
    "priority": "Medium",
    ...
  }
}
```

---

### 3. Machine Endpoints

#### 3.1 Get All Machines
Retrieve all machines/terminals from the system.

**Endpoint:** `GET /api/v1/machines`
**Authentication:** Required

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Machines retrieved successfully",
  "data": [
    {
      "id": 1,
      "terminal_id": "ATM-001",
      "terminal_name": "Main Branch ATM 1",
      "location": "Jakarta - Main Branch",
      "branch_code": "JKT-001",
      "ip_address": "192.168.1.100",
      "model": "NCR 6634",
      "manufacturer": "NCR Corporation",
      "serial_number": "SN123456789",
      "status": "Active",
      "last_ping_time": "2024-01-15T10:25:00Z",
      "install_date": "2023-01-15T00:00:00Z",
      "warranty_exp": "2026-01-15T00:00:00Z",
      "notes": "Regular maintenance scheduled",
      "created_at": "2023-01-15T00:00:00Z",
      "updated_at": "2024-01-15T10:25:00Z"
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
    "id": 1,
    "terminal_id": "ATM-001",
    "terminal_name": "Main Branch ATM 1",
    ...
  }
}
```

#### 3.3 Get Machines by Status
Retrieve all machines with a specific operational status.

**Endpoint:** `GET /api/v1/machines/status/:status`
**Authentication:** Required

**Parameters:**
- `status` (path, string) - Status value: `Active`, `Inactive`, `Maintenance`, `Offline`

**Response:** Same as Get All Machines (filtered list)

#### 3.4 Get Machines by Branch
Retrieve all machines for a specific branch.

**Endpoint:** `GET /api/v1/machines/branch/:branch_code`
**Authentication:** Required

**Parameters:**
- `branch_code` (path, string) - Branch code (e.g., "JKT-001")

**Response:** Same as Get All Machines (filtered list)

#### 3.5 Search Machines
Search machines using multiple filter criteria.

**Endpoint:** `GET /api/v1/machines/search`
**Authentication:** Required

**Query Parameters (all optional):**
- `status` (string) - Filter by status
- `branch_code` (string) - Filter by branch code
- `location` (string) - Search by location (partial match)

**Example:**
```
GET /api/v1/machines/search?status=Active&branch_code=JKT-001&location=Jakarta
```

**Response:** Same as Get All Machines (filtered list)

#### 3.6 Update Machine Status
Update the operational status of a machine.

**Endpoint:** `PATCH /api/v1/machines/status`
**Authentication:** Required

**Request Body:**
```json
{
  "terminal_id": "ATM-001",
  "status": "Maintenance",
  "last_ping_time": "2024-01-15T10:30:00Z",
  "notes": "Scheduled maintenance in progress"
}
```

**Field Validation:**
- `terminal_id` (required, string) - Terminal to update
- `status` (required, string) - Must be one of: `Active`, `Inactive`, `Maintenance`, `Offline`
- `last_ping_time` (optional, string) - ISO 8601 timestamp
- `notes` (optional, string) - Additional notes

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Machine status updated successfully",
  "data": {
    "id": 1,
    "terminal_id": "ATM-001",
    "status": "Maintenance",
    ...
  }
}
```

---

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid input data |
| 401 | Unauthorized - Missing or invalid API key |
| 404 | Not Found - Resource not found |
| 409 | Conflict - Resource already exists |
| 500 | Internal Server Error - Server error |
| 503 | Service Unavailable - Service/database unavailable |

---

## Rate Limiting

Currently, no rate limiting is implemented. Consider adding rate limiting in production.

---

## Examples

### cURL Examples

**Get all tickets:**
```bash
curl -X GET \
  -H "X-API-Key: your_api_key" \
  http://localhost:8080/api/v1/tickets
```

**Create a ticket:**
```bash
curl -X POST \
  -H "X-API-Key: your_api_key" \
  -H "Content-Type: application/json" \
  -d '{
    "ticket_number": "TKT-2024-001",
    "terminal_id": "ATM-001",
    "description": "Card reader issue",
    "priority": "High",
    "category": "Hardware",
    "reported_by": "John Doe"
  }' \
  http://localhost:8080/api/v1/tickets
```

**Update machine status:**
```bash
curl -X PATCH \
  -H "X-API-Key: your_api_key" \
  -H "Content-Type: application/json" \
  -d '{
    "terminal_id": "ATM-001",
    "status": "Maintenance",
    "notes": "Scheduled maintenance"
  }' \
  http://localhost:8080/api/v1/machines/status
```

---

## Versioning

Current API version: **v1**

All endpoints are prefixed with `/api/v1` to support future versioning.

---

## Support

For API support or questions, contact the development team.
