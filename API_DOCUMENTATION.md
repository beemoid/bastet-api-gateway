# API Gateway - Complete API Documentation

## Overview

This API Gateway serves as middleware between on-premise databases (ticket_master and machine_master) and cloud applications, providing a secure and standardized REST API for ATM monitoring and ticket management.

**Base URL:** `http://localhost:8080`
**API Version:** v1
**API Prefix:** `/api/v1`
**Swagger UI:** `http://localhost:8080/swagger/index.html`

---

## Key Features

✅ **Hybrid Adaptive Metadata System** - Automatically discovers new field values from database
✅ **Intelligent Caching** - 1-hour cache for optimal performance
✅ **Thread-Safe Operations** - Production-ready concurrent request handling
✅ **Comprehensive Error Handling** - Standardized error responses
✅ **Real-time Database Queries** - Always up-to-date with actual data

---

## Authentication

All API endpoints (except health checks and Swagger documentation) require authentication using an API key.

**Header Required:**
```
X-API-Key: your_api_key_here
```

**Example:**
```bash
curl -H "X-API-Key: your_api_key" http://localhost:8080/api/v1/tickets
```

**Configuration:**
Set your API key in the `.env` file:
```
API_KEY=your_secure_api_key_here
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

#### 2.6 Get Ticket Metadata ⭐ NEW - Adaptive System
Retrieve all valid values for ticket fields directly from the database.

**Endpoint:** `GET /api/v1/tickets/metadata`
**Authentication:** Required

**Features:**
- ✅ **Fully Adaptive** - Automatically discovers new values from database
- ✅ **Cached for Performance** - Results cached for 1 hour
- ✅ **Self-Documenting** - Shows which values are documented
- ✅ **Real-time Accuracy** - Always reflects actual database content

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
      "code": "1.Req FD ke HD",
      "description": "Request FD to HD",
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
    },
    {
      "code": "In Service",
      "description": "Terminal is in service",
      "is_documented": true
    },
    {
      "code": "nan",
      "description": "No mode data available",
      "is_documented": true
    },
    {
      "code": "Off-line",
      "description": "Terminal is offline",
      "is_documented": true
    },
    {
      "code": "Supervisor",
      "description": "Supervisor mode",
      "is_documented": true
    }
  ],
  "priorities": [
    {
      "code": "1.High",
      "description": "High priority",
      "is_documented": true
    },
    {
      "code": "2.Middle",
      "description": "Middle priority",
      "is_documented": true
    },
    {
      "code": "3.Low",
      "description": "Low priority",
      "is_documented": true
    },
    {
      "code": "4.Minimum",
      "description": "Minimum priority",
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

#### 3.6 Get Machine Metadata ⭐ NEW - Adaptive System
Retrieve all valid values for machine fields directly from the database.

**Endpoint:** `GET /api/v1/machines/metadata`
**Authentication:** Required

**Features:**
- ✅ **Fully Adaptive** - Automatically discovers new values from database
- ✅ **Cached for Performance** - Results cached for 1 hour
- ✅ **Geographic Mapping** - FLM providers mapped to service areas
- ✅ **Real-time Accuracy** - Always reflects actual database content

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
    },
    {
      "code": "GPS - NCR",
      "description": "GPS - NCR",
      "is_documented": true
    },
    {
      "code": "NCR",
      "description": "NCR",
      "is_documented": true
    }
  ],
  "flms": [
    {
      "code": "AVT - BANDUNG",
      "description": "AVT - BANDUNG",
      "area": "BANDUNG",
      "is_documented": true
    },
    {
      "code": "AVT - JAKARTA",
      "description": "AVT - JAKARTA",
      "area": "JAKARTA",
      "is_documented": true
    },
    {
      "code": "BRS - SURABAYA",
      "description": "BRS - SURABAYA",
      "area": "SURABAYA",
      "is_documented": true
    }
  ],
  "nets": [
    {
      "code": "NOSAIRIS",
      "description": "NOSAIRIS",
      "is_documented": true
    },
    {
      "code": "SMS",
      "description": "SMS",
      "is_documented": true
    },
    {
      "code": "TANGARA",
      "description": "TANGARA",
      "is_documented": true
    },
    {
      "code": "IFORTE",
      "description": "IFORTE",
      "is_documented": true
    }
  ],
  "flm_names": [
    {
      "code": "AVT",
      "description": "AVT",
      "is_documented": true
    },
    {
      "code": "ABS",
      "description": "ABS",
      "is_documented": true
    },
    {
      "code": "BRS",
      "description": "BRS",
      "is_documented": true
    },
    {
      "code": "TAG",
      "description": "TAG",
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
Client Request → Handler → Service (Check Cache)
                              ↓
                         Cache Hit? Yes → Return Cached Data
                              ↓ No
                         Query Database (SELECT DISTINCT)
                              ↓
                         Combine with Descriptions
                              ↓
                         Update Cache
                              ↓
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
→ Returns: {"code": "9.Emergency", "description": "9.Emergency", "is_documented": false}

-- Value works immediately in all operations
POST /api/v1/tickets {"status": "9.Emergency", ...}  ✅ Works!
GET /api/v1/tickets/status/9.Emergency  ✅ Works!
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
| 409 | Conflict - Resource already exists (duplicate ticket number) |
| 500 | Internal Server Error - Server error |
| 503 | Service Unavailable - Database connection unavailable |

---

## Examples

### cURL Examples

**Get ticket metadata:**
```bash
curl -X GET \
  -H "X-API-Key: your_api_key" \
  http://localhost:8080/api/v1/tickets/metadata
```

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
  -H "X-API-Key: your_api_key" \
  "http://localhost:8080/api/v1/tickets/status/0.NEW"
```

**Update ticket:**
```bash
curl -X PUT \
  -H "X-API-Key: your_api_key" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "6.Follow-up Sales team",
    "remarks": "Technician dispatched",
    "condition": "Normal"
  }' \
  http://localhost:8080/api/v1/tickets/ATM-001
```

**Get machine metadata:**
```bash
curl -X GET \
  -H "X-API-Key: your_api_key" \
  http://localhost:8080/api/v1/machines/metadata
```

**Search machines:**
```bash
curl -X GET \
  -H "X-API-Key: your_api_key" \
  "http://localhost:8080/api/v1/machines/search?status=Active&province=DKI%20Jakarta"
```

**Update machine status:**
```bash
curl -X PATCH \
  -H "X-API-Key: your_api_key" \
  -H "Content-Type: application/json" \
  -d '{
    "terminal_id": "ATM-001",
    "status": "Maintenance"
  }' \
  http://localhost:8080/api/v1/machines/status
```

### JavaScript/Fetch Examples

**Get ticket metadata:**
```javascript
fetch('http://localhost:8080/api/v1/tickets/metadata', {
  headers: {
    'X-API-Key': 'your_api_key'
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
    'X-API-Key': 'your_api_key',
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

---

## Best Practices

### 1. Use Metadata Endpoints
Always use metadata endpoints to discover valid field values rather than hard-coding them:

```javascript
// ✅ Good: Dynamic values from metadata
const metadata = await fetch('/api/v1/tickets/metadata');
const statuses = metadata.statuses.map(s => s.code);

// ❌ Bad: Hard-coded values
const statuses = ['Open', 'Closed']; // Will be outdated
```

### 2. Handle Nullable Fields
Many fields can be NULL in the database. Always handle null values:

```javascript
const priority = ticket.priority || 'Not Set';
const closeTime = ticket.close_time || 'Still Open';
```

### 3. Cache Metadata Client-Side
Since metadata is cached for 1 hour server-side, you can also cache it client-side:

```javascript
// Cache metadata for 1 hour
const CACHE_DURATION = 60 * 60 * 1000; // 1 hour
let cachedMetadata = null;
let cacheTime = 0;

async function getMetadata() {
  if (cachedMetadata && Date.now() - cacheTime < CACHE_DURATION) {
    return cachedMetadata;
  }

  const response = await fetch('/api/v1/tickets/metadata');
  cachedMetadata = await response.json();
  cacheTime = Date.now();
  return cachedMetadata;
}
```

### 4. Use Proper Status Codes
Always check the `success` field and HTTP status code:

```javascript
const response = await fetch('/api/v1/tickets');
const data = await response.json();

if (data.success) {
  // Handle success
  console.log('Tickets:', data.data);
} else {
  // Handle error
  console.error('Error:', data.message);
}
```

### 5. Handle Undocumented Values
Use the `is_documented` flag to identify new values:

```javascript
metadata.statuses.forEach(status => {
  if (!status.is_documented) {
    console.warn(`New status discovered: ${status.code}`);
    // Alert admin to document this value
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

---

## Versioning

**Current API version:** v1

All endpoints are prefixed with `/api/v1` to support future versioning.

**Version Strategy:**
- Breaking changes will increment the major version (v2, v3, etc.)
- Backward-compatible changes will be added to existing versions
- Old versions will be maintained for minimum 6 months after new version release

---

## Rate Limiting

Currently, no rate limiting is implemented.

**Recommendations for Production:**
- Implement rate limiting per API key
- Suggested limit: 1000 requests per hour per key
- Return `429 Too Many Requests` when limit exceeded
- Include `X-RateLimit-*` headers in responses

---

## Support & Contact

For API support, bug reports, or feature requests:

- **GitHub Issues**: [Create an issue](https://github.com/your-org/api-gateway/issues)
- **Email**: support@your-org.com
- **Documentation**: This file and Swagger UI at `/swagger/index.html`

---

## Changelog

### Version 1.0.0 (2024-01-15)

**Features:**
- ✅ Complete CRUD operations for tickets and machines
- ✅ Hybrid adaptive metadata system with database queries
- ✅ Intelligent 1-hour caching for optimal performance
- ✅ Thread-safe concurrent request handling
- ✅ Comprehensive Swagger/OpenAPI documentation
- ✅ Support for NULL database values with proper JSON marshaling
- ✅ Float64 support for ticket duration field
- ✅ Geographic area mapping for FLM providers
- ✅ Standardized error responses
- ✅ API key authentication
- ✅ Health check endpoints
- ✅ Advanced search and filtering

**Database Support:**
- ✅ SQL Server (ticket_master and machine_master databases)
- ✅ Proper handling of database column names with spaces
- ✅ Support for nullable fields

**Performance:**
- ✅ Metadata caching (1-hour TTL)
- ✅ Efficient database queries with proper indexing
- ✅ Thread-safe operations with RWMutex

---

*Last Updated: 2024-01-15*
*API Gateway Version: 1.0.0*
*Documentation Version: 1.0.0*
