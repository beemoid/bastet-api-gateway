# API Migration Guide: v1 to v2

## Overview

This guide helps you migrate from API Gateway v1 to v2. Version 2 introduces significant improvements including enhanced token management, better analytics, and improved data structures.

**Migration Timeline:**
- **v2 Launch Date:** March 1, 2024
- **v1 Support Period:** 6 months (until September 1, 2024)
- **v1 Deprecation Date:** September 1, 2024
- **v1 Sunset Date:** October 1, 2024

After the sunset date, v1 endpoints will return `410 Gone` status code.

---

## Table of Contents

1. [Breaking Changes](#breaking-changes)
2. [New Features](#new-features)
3. [Field Mappings](#field-mappings)
4. [Authentication Changes](#authentication-changes)
5. [Response Format Changes](#response-format-changes)
6. [Code Examples](#code-examples)
7. [Testing Your Migration](#testing-your-migration)
8. [FAQ](#faq)

---

## Breaking Changes

### 1. Field Name Changes

Several field names have been improved for clarity and consistency.

#### Ticket Fields

| v1 Field Name | v2 Field Name | Reason |
|---------------|---------------|--------|
| `tickets_no` | `ticket_number` | Improved clarity |
| `dsp_flm` | `flm_identifier` | Consistency with naming conventions |
| `dsp_slm` | `slm_identifier` | Consistency with naming conventions |
| `p_duration` | `problem_duration` | Clarity |
| `incident_start_datetime` | `incident_started_at` | Consistency with timestamp naming |
| `open_time` | `opened_at` | Consistency |
| `close_time` | `closed_at` | Consistency |

#### Machine Fields

| v1 Field Name | v2 Field Name | Reason |
|---------------|---------------|--------|
| `city_regency` | `city` | Simplified naming |
| `std` | `standard` | Clarity |

### 2. Response Structure Changes

v2 introduces a more structured response format with metadata and pagination.

**v1 Response:**
```json
{
  "success": true,
  "message": "Tickets retrieved successfully",
  "data": [...],
  "total": 100
}
```

**v2 Response:**
```json
{
  "success": true,
  "message": "Tickets retrieved successfully",
  "version": "2.0",
  "data": [...],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total_pages": 5,
    "total_items": 100,
    "has_next": true,
    "has_prev": false
  },
  "meta": {
    "request_id": "req_abc123xyz",
    "timestamp": "2024-01-15T10:30:00Z",
    "response_time_ms": 45
  }
}
```

### 3. Pagination Support

v2 adds mandatory pagination for list endpoints to improve performance.

**v1 (No pagination):**
```
GET /api/v1/tickets
```

**v2 (With pagination):**
```
GET /api/v2/tickets?page=1&per_page=20
```

**Default values if not specified:**
- `page`: 1
- `per_page`: 20
- `max per_page`: 100

### 4. Authentication Changes

#### v1: Simple API Key
```
Header: X-API-Key: your_simple_api_key
```

#### v2: Token-based Authentication with Metadata
```
Header: X-API-Token: tok_live_abc123xyz456
```

**New Token Features:**
- ✅ Scoped permissions (read-only, read-write, admin)
- ✅ Token expiration dates
- ✅ Usage analytics and rate limiting
- ✅ IP whitelisting (optional)
- ✅ Token rotation support
- ✅ Detailed audit logs

### 5. Error Response Changes

**v1 Error:**
```json
{
  "success": false,
  "message": "Ticket not found"
}
```

**v2 Error (Enhanced):**
```json
{
  "success": false,
  "message": "Ticket not found",
  "error": {
    "code": "TICKET_NOT_FOUND",
    "type": "NotFoundError",
    "details": "No ticket found with terminal_id: ATM-001",
    "request_id": "req_abc123xyz",
    "timestamp": "2024-01-15T10:30:00Z",
    "documentation_url": "https://docs.api.example.com/errors/TICKET_NOT_FOUND"
  }
}
```

---

## New Features

### 1. Advanced Token Management

v2 introduces a complete token management system:

- **Token Dashboard:** Web UI to manage all API tokens
- **Scoped Permissions:** Fine-grained access control
- **Usage Analytics:** Real-time monitoring of token usage
- **Rate Limiting:** Per-token rate limits
- **Token Rotation:** Secure token rotation without downtime
- **Audit Logs:** Complete history of all API calls

**Access Token Management:**
```
Dashboard: http://localhost:8080/admin/tokens
API: /api/v2/admin/tokens
```

### 2. Analytics Endpoints

New endpoints for monitoring and insights:

```
GET /api/v2/analytics/dashboard           - Overall system metrics
GET /api/v2/analytics/critical-terminals  - Terminals needing attention
GET /api/v2/analytics/flm-workload        - Maintenance provider workload
GET /api/v2/analytics/area-analysis       - Geographic performance
```

### 3. Enhanced Metadata

Metadata endpoints now include additional information:

```json
{
  "statuses": [
    {
      "code": "0.NEW",
      "description": "New ticket",
      "is_documented": true,
      "usage_count": 1234,          // NEW in v2
      "last_used": "2024-01-15T10:30:00Z"  // NEW in v2
    }
  ]
}
```

### 4. Webhook Support

v2 adds webhook notifications for events:

```
POST /api/v2/webhooks/register
{
  "url": "https://your-app.com/webhook",
  "events": ["ticket.created", "ticket.updated", "machine.offline"],
  "secret": "your_webhook_secret"
}
```

### 5. Batch Operations

New batch endpoints for efficiency:

```
POST /api/v2/tickets/batch
{
  "operations": [
    {"action": "create", "data": {...}},
    {"action": "update", "id": "ATM-001", "data": {...}}
  ]
}
```

---

## Field Mappings

### Complete Ticket Field Mapping

```javascript
// v1 to v2 field conversion
const convertTicketV1ToV2 = (v1Ticket) => {
  return {
    // Renamed fields
    ticket_number: v1Ticket.tickets_no,
    flm_identifier: v1Ticket.dsp_flm,
    slm_identifier: v1Ticket.dsp_slm,
    problem_duration: v1Ticket.p_duration,
    incident_started_at: v1Ticket.incident_start_datetime,
    opened_at: v1Ticket.open_time,
    closed_at: v1Ticket.close_time,

    // Unchanged fields
    terminal_id: v1Ticket.terminal_id,
    terminal_name: v1Ticket.terminal_name,
    priority: v1Ticket.priority,
    mode: v1Ticket.mode,
    initial_problem: v1Ticket.initial_problem,
    current_problem: v1Ticket.current_problem,
    count: v1Ticket.count,
    status: v1Ticket.status,
    remarks: v1Ticket.remarks,
    balance: v1Ticket.balance,
    condition: v1Ticket.condition,
    tickets_duration: v1Ticket.tickets_duration,
    problem_history: v1Ticket.problem_history,
    mode_history: v1Ticket.mode_history,
    last_withdrawal: v1Ticket.last_withdrawal,
    export_name: v1Ticket.export_name
  };
};
```

### Complete Machine Field Mapping

```javascript
// v1 to v2 field conversion
const convertMachineV1ToV2 = (v1Machine) => {
  return {
    // Renamed fields
    city: v1Machine.city_regency,
    standard: v1Machine.std,

    // Unchanged fields
    terminal_id: v1Machine.terminal_id,
    store: v1Machine.store,
    store_code: v1Machine.store_code,
    store_name: v1Machine.store_name,
    date_of_activation: v1Machine.date_of_activation,
    status: v1Machine.status,
    gps: v1Machine.gps,
    lat: v1Machine.lat,
    lon: v1Machine.lon,
    province: v1Machine.province,
    district: v1Machine.district
  };
};
```

---

## Authentication Changes

### v1 Authentication (Deprecated)

```bash
# Simple API key in header
curl -H "X-API-Key: your_simple_key" \
  http://localhost:8080/api/v1/tickets
```

**Limitations:**
- ❌ No expiration
- ❌ No permissions/scopes
- ❌ No usage tracking
- ❌ Single key for all applications

### v2 Authentication (Recommended)

```bash
# Token-based with scopes and expiration
curl -H "X-API-Token: tok_live_abc123xyz456" \
  http://localhost:8080/api/v2/tickets
```

**Benefits:**
- ✅ Token expiration
- ✅ Scoped permissions (read, write, admin)
- ✅ Usage analytics
- ✅ Multiple tokens per application
- ✅ Token rotation support
- ✅ IP whitelisting

### Obtaining v2 Tokens

**Step 1: Access Token Dashboard**
```
URL: http://localhost:8080/admin/tokens
Login: admin / your_admin_password
```

**Step 2: Create New Token**
```
Click "Create Token"
- Name: "Production App"
- Scopes: ["tickets:read", "tickets:write", "machines:read"]
- Expires: 2025-12-31
- IP Whitelist: 192.168.1.0/24 (optional)
```

**Step 3: Save Token**
```
Token generated: tok_live_abc123xyz456
⚠️ Save this token - it won't be shown again!
```

**Step 4: Use Token**
```bash
curl -H "X-API-Token: tok_live_abc123xyz456" \
  http://localhost:8080/api/v2/tickets
```

---

## Response Format Changes

### List Endpoints

#### v1 Response Format

```json
GET /api/v1/tickets

{
  "success": true,
  "message": "Tickets retrieved successfully",
  "data": [
    {
      "terminal_id": "ATM-001",
      "tickets_no": "TKT-2024-001",
      "status": "0.NEW"
    }
  ],
  "total": 100
}
```

#### v2 Response Format

```json
GET /api/v2/tickets?page=1&per_page=20

{
  "success": true,
  "message": "Tickets retrieved successfully",
  "version": "2.0",
  "data": [
    {
      "terminal_id": "ATM-001",
      "ticket_number": "TKT-2024-001",  // Renamed
      "status": "0.NEW",
      "created_at": "2024-01-15T10:30:00Z",  // New field
      "updated_at": "2024-01-15T12:00:00Z"   // New field
    }
  ],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total_pages": 5,
    "total_items": 100,
    "has_next": true,
    "has_prev": false
  },
  "meta": {
    "request_id": "req_abc123xyz",
    "timestamp": "2024-01-15T14:30:00Z",
    "response_time_ms": 45,
    "api_version": "2.0"
  }
}
```

### Error Responses

#### v1 Error Format

```json
{
  "success": false,
  "message": "Ticket not found"
}
```

#### v2 Error Format

```json
{
  "success": false,
  "message": "Ticket not found",
  "error": {
    "code": "TICKET_NOT_FOUND",
    "type": "NotFoundError",
    "details": "No ticket found with terminal_id: ATM-001",
    "request_id": "req_abc123xyz",
    "timestamp": "2024-01-15T10:30:00Z",
    "documentation_url": "https://docs.example.com/errors/TICKET_NOT_FOUND",
    "suggestions": [
      "Verify the terminal_id is correct",
      "Check if the ticket was closed",
      "Use GET /api/v2/tickets to list all tickets"
    ]
  }
}
```

---

## Code Examples

### JavaScript/Node.js Migration

#### Before (v1):

```javascript
// v1 client code
const apiKey = 'your_simple_api_key';

async function getAllTickets() {
  const response = await fetch('http://localhost:8080/api/v1/tickets', {
    headers: {
      'X-API-Key': apiKey
    }
  });

  const result = await response.json();

  if (result.success) {
    return result.data;  // Returns all tickets (no pagination)
  }

  throw new Error(result.message);
}

// Usage
const tickets = await getAllTickets();
console.log(`Total tickets: ${tickets.length}`);
```

#### After (v2):

```javascript
// v2 client code
const apiToken = 'tok_live_abc123xyz456';

async function getAllTickets(page = 1, perPage = 20) {
  const response = await fetch(
    `http://localhost:8080/api/v2/tickets?page=${page}&per_page=${perPage}`,
    {
      headers: {
        'X-API-Token': apiToken
      }
    }
  );

  const result = await response.json();

  if (result.success) {
    return {
      tickets: result.data,
      pagination: result.pagination,
      meta: result.meta
    };
  }

  // Enhanced error handling
  throw new Error(
    `${result.error.code}: ${result.error.details} (Request ID: ${result.error.request_id})`
  );
}

// Usage with pagination
const { tickets, pagination, meta } = await getAllTickets(1, 20);
console.log(`Page ${pagination.page} of ${pagination.total_pages}`);
console.log(`Total tickets: ${pagination.total_items}`);
console.log(`Request ID: ${meta.request_id}`);
console.log(`Response time: ${meta.response_time_ms}ms`);

// Field name changes - convert ticket data
const ticketNumber = tickets[0].ticket_number;  // v2: ticket_number (was tickets_no in v1)
const flmId = tickets[0].flm_identifier;        // v2: flm_identifier (was dsp_flm in v1)
const openedAt = tickets[0].opened_at;          // v2: opened_at (was open_time in v1)
```

### Python Migration

#### Before (v1):

```python
# v1 client code
import requests

API_KEY = 'your_simple_api_key'
BASE_URL = 'http://localhost:8080/api/v1'

def get_all_tickets():
    headers = {'X-API-Key': API_KEY}
    response = requests.get(f'{BASE_URL}/tickets', headers=headers)

    result = response.json()

    if result['success']:
        return result['data']

    raise Exception(result['message'])

# Usage
tickets = get_all_tickets()
print(f'Total tickets: {len(tickets)}')
```

#### After (v2):

```python
# v2 client code
import requests

API_TOKEN = 'tok_live_abc123xyz456'
BASE_URL = 'http://localhost:8080/api/v2'

def get_all_tickets(page=1, per_page=20):
    headers = {'X-API-Token': API_TOKEN}
    params = {'page': page, 'per_page': per_page}

    response = requests.get(f'{BASE_URL}/tickets', headers=headers, params=params)
    result = response.json()

    if result['success']:
        return {
            'tickets': result['data'],
            'pagination': result['pagination'],
            'meta': result['meta']
        }

    # Enhanced error handling
    error = result.get('error', {})
    raise Exception(
        f"{error.get('code')}: {error.get('details')} "
        f"(Request ID: {error.get('request_id')})"
    )

# Usage with pagination
response = get_all_tickets(page=1, per_page=20)
tickets = response['tickets']
pagination = response['pagination']
meta = response['meta']

print(f"Page {pagination['page']} of {pagination['total_pages']}")
print(f"Total tickets: {pagination['total_items']}")
print(f"Request ID: {meta['request_id']}")
print(f"Response time: {meta['response_time_ms']}ms")

# Field name changes - access ticket data
ticket_number = tickets[0]['ticket_number']    # v2: ticket_number (was tickets_no in v1)
flm_id = tickets[0]['flm_identifier']         # v2: flm_identifier (was dsp_flm in v1)
opened_at = tickets[0]['opened_at']            # v2: opened_at (was open_time in v1)
```

### cURL Migration

#### Before (v1):

```bash
# Get all tickets (v1)
curl -H "X-API-Key: your_simple_key" \
  http://localhost:8080/api/v1/tickets

# Create ticket (v1)
curl -X POST \
  -H "X-API-Key: your_simple_key" \
  -H "Content-Type: application/json" \
  -d '{
    "terminal_id": "ATM-001",
    "tickets_no": "TKT-2024-001",
    "status": "0.NEW",
    "dsp_flm": "FLM-001"
  }' \
  http://localhost:8080/api/v1/tickets
```

#### After (v2):

```bash
# Get all tickets with pagination (v2)
curl -H "X-API-Token: tok_live_abc123xyz456" \
  "http://localhost:8080/api/v2/tickets?page=1&per_page=20"

# Create ticket with updated field names (v2)
curl -X POST \
  -H "X-API-Token: tok_live_abc123xyz456" \
  -H "Content-Type: application/json" \
  -d '{
    "terminal_id": "ATM-001",
    "ticket_number": "TKT-2024-001",
    "status": "0.NEW",
    "flm_identifier": "FLM-001"
  }' \
  http://localhost:8080/api/v2/tickets
```

---

## Testing Your Migration

### Step 1: Get v2 Access

1. **Create a v2 token:**
   ```bash
   # Access admin dashboard
   URL: http://localhost:8080/admin/tokens
   Login with admin credentials
   Create new token with appropriate scopes
   ```

2. **Test authentication:**
   ```bash
   curl -H "X-API-Token: your_v2_token" \
     http://localhost:8080/api/v2/tickets/metadata
   ```

### Step 2: Test Parallel v1 and v2

Run both versions side-by-side to compare:

```javascript
// Test script to compare v1 vs v2
async function compareVersions() {
  // Fetch from v1
  const v1Response = await fetch('http://localhost:8080/api/v1/tickets', {
    headers: { 'X-API-Key': 'v1_key' }
  });
  const v1Data = await v1Response.json();

  // Fetch from v2
  const v2Response = await fetch('http://localhost:8080/api/v2/tickets?page=1&per_page=100', {
    headers: { 'X-API-Token': 'v2_token' }
  });
  const v2Data = await v2Response.json();

  // Compare data counts
  console.log('v1 count:', v1Data.total);
  console.log('v2 count:', v2Data.pagination.total_items);

  // Verify field mappings
  const v1Ticket = v1Data.data[0];
  const v2Ticket = v2Data.data[0];

  console.log('v1 tickets_no:', v1Ticket.tickets_no);
  console.log('v2 ticket_number:', v2Ticket.ticket_number);
  console.log('Match:', v1Ticket.tickets_no === v2Ticket.ticket_number);
}
```

### Step 3: Update Your Code Gradually

**Strategy:**

1. **Week 1-2:** Test v2 in development/staging
2. **Week 3-4:** Deploy v2 alongside v1 in production
3. **Month 2-3:** Migrate clients one by one
4. **Month 4-5:** Monitor for v1 usage decline
5. **Month 6:** Deprecate v1, return warnings
6. **Month 7:** Sunset v1, return 410 Gone

### Step 4: Monitor Migration Progress

Use the token dashboard to track:
- v1 vs v2 API usage
- Which tokens are still using v1
- Error rates during migration
- Performance comparisons

---

## FAQ

### Q: Can I use both v1 and v2 simultaneously?

**A:** Yes! Both versions will run side-by-side during the migration period. You can gradually migrate your applications one by one.

### Q: Do I need new database credentials?

**A:** No. Both v1 and v2 connect to the same databases. Only the API structure has changed.

### Q: What happens to my v1 API key?

**A:** v1 API keys will continue to work for v1 endpoints until the sunset date. However, you must create v2 tokens for v2 endpoints.

### Q: Can I automatically convert v1 responses to v2 format?

**A:** Yes, we provide conversion functions in our SDK. See the code examples above.

### Q: How do I handle pagination in v2?

**A:** v2 requires pagination for all list endpoints:
```javascript
// Get all tickets across pages
async function getAllTicketsPaginated() {
  let allTickets = [];
  let page = 1;
  let hasNext = true;

  while (hasNext) {
    const response = await getTickets(page, 100);
    allTickets = allTickets.concat(response.tickets);
    hasNext = response.pagination.has_next;
    page++;
  }

  return allTickets;
}
```

### Q: What if I find a bug during migration?

**A:** Report issues to:
- GitHub Issues: https://github.com/your-org/api-gateway/issues
- Email: support@example.com
- Include `request_id` from error responses for faster resolution

### Q: Can I test v2 without affecting production?

**A:** Yes! Create a separate token with "test" scope:
```
Scope: test
Environment: staging
Rate Limit: 10 req/sec
```

### Q: How do I rotate tokens securely?

**A:** Token rotation process:
1. Create new token (Token B)
2. Update half of your services to use Token B
3. Monitor for 24 hours
4. Update remaining services to use Token B
5. Disable old token (Token A)
6. Delete Token A after 7 days

### Q: Will v2 be faster?

**A:** Yes! v2 improvements include:
- Pagination reduces payload size
- Optimized database queries
- Better caching
- CDN support for metadata endpoints

Benchmarks show 30-50% faster response times on average.

### Q: What if I can't migrate by the deadline?

**A:** Contact our support team at least 30 days before the sunset date. We can:
- Extend v1 support for critical applications
- Provide migration assistance
- Offer custom transition plans

---

## Migration Checklist

Use this checklist to track your migration progress:

### Pre-Migration
- [ ] Read this migration guide completely
- [ ] Access v2 API documentation
- [ ] Create v2 token in admin dashboard
- [ ] Test v2 endpoints in development
- [ ] Update SDK/libraries to v2-compatible versions

### Code Changes
- [ ] Update authentication (X-API-Key → X-API-Token)
- [ ] Update base URL (/api/v1 → /api/v2)
- [ ] Update field names (see Field Mappings section)
- [ ] Add pagination handling
- [ ] Update error handling
- [ ] Update response parsing (handle new meta/pagination fields)

### Testing
- [ ] Test all read operations (GET endpoints)
- [ ] Test all write operations (POST/PUT/PATCH endpoints)
- [ ] Test error scenarios
- [ ] Compare v1 vs v2 responses
- [ ] Load testing with v2
- [ ] Test token permissions/scopes

### Deployment
- [ ] Deploy to staging environment
- [ ] Run integration tests
- [ ] Monitor error logs
- [ ] Gradual rollout to production (10% → 50% → 100%)
- [ ] Monitor performance metrics
- [ ] Verify token usage analytics

### Post-Migration
- [ ] Disable v1 token usage
- [ ] Remove v1 code from codebase
- [ ] Update internal documentation
- [ ] Archive v1 credentials securely
- [ ] Celebrate! 🎉

---

## Support & Resources

### Documentation
- **API v2 Documentation:** `API_DOCUMENTATION_V2.md`
- **Token Management Guide:** `TOKEN_MANAGEMENT_GUIDE.md`
- **Swagger UI v2:** http://localhost:8080/swagger/v2/index.html

### Tools
- **Token Dashboard:** http://localhost:8080/admin/tokens
- **Migration Validator:** http://localhost:8080/admin/migration-checker
- **API Playground:** http://localhost:8080/playground

### Support Channels
- **GitHub Issues:** https://github.com/your-org/api-gateway/issues
- **Email Support:** support@example.com
- **Slack Channel:** #api-migration
- **Office Hours:** Tuesdays 2-4 PM (for migration assistance)

### Example Code & SDKs
- **JavaScript SDK:** https://github.com/your-org/api-sdk-js
- **Python SDK:** https://github.com/your-org/api-sdk-python
- **Go SDK:** https://github.com/your-org/api-sdk-go
- **Migration Examples:** https://github.com/your-org/api-migration-examples

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2024-01-15 | Initial migration guide |
| 1.1 | 2024-01-20 | Added Python examples |
| 1.2 | 2024-01-25 | Added token rotation guide |
| 1.3 | 2024-02-01 | Added pagination FAQ |

---

*Last Updated: January 15, 2024*
*For questions or assistance, contact: support@example.com*
