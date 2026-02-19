# Migration Guide: API Key to Token-Based Authentication

## Overview

This guide helps you migrate from the legacy API key authentication (`X-API-Key`) to the new token-based authentication system (`X-API-Token`). The token system was introduced alongside the admin dashboard, analytics, and audit logging features.

**Important:** Both authentication methods currently work side-by-side on all `/api/v1` endpoints. There is no separate `/api/v2` prefix -- all new features are available under the existing `/api/v1` path.

---

## Table of Contents

1. [What Changed](#what-changed)
2. [Why Migrate](#why-migrate)
3. [Migration Steps](#migration-steps)
4. [Authentication Comparison](#authentication-comparison)
5. [Code Examples](#code-examples)
6. [New Features Available After Migration](#new-features-available-after-migration)
7. [FAQ](#faq)

---

## What Changed

### New Systems Added (v1.1.0)

| Feature | Description |
|---------|-------------|
| **Token Management** | Scoped API tokens with rate limiting, IP whitelisting, and expiration |
| **Admin Dashboard** | Web UI at `/admin` for managing tokens and viewing analytics |
| **Admin Auth** | Session-based admin authentication (`POST /api/v1/admin/auth/login`) |
| **Analytics** | Dashboard stats, endpoint stats, daily usage tracking |
| **Audit Logging** | Complete history of all administrative actions |
| **Rate Limiting** | Per-token configurable limits (per minute, hour, day) |
| **Usage Tracking** | Every API call made with a token is logged with response time, status, IP, etc. |
| **Gzip Compression** | Automatic response compression |
| **Third Database** | New `token_management` database for tokens, sessions, and audit logs |

### What Did NOT Change

- All existing `/api/v1` endpoints remain identical
- Request/response formats are unchanged
- Field names are unchanged
- No mandatory pagination was added
- Database schemas for `ticket_master` and `machine_master` are untouched
- `X-API-Key` authentication still works

---

## Why Migrate

### Legacy API Key Limitations

| Feature | API Key (`X-API-Key`) | API Token (`X-API-Token`) |
|---------|----------------------|--------------------------|
| Multiple keys per app | No (single shared key) | Yes (unlimited tokens) |
| Scoped permissions | No | Yes |
| Rate limiting | No | Yes (per minute/hour/day) |
| Usage analytics | No | Yes (per-token tracking) |
| IP whitelisting | No | Yes (optional) |
| Token expiration | No | Yes (configurable) |
| Disable/enable | No (must change key) | Yes (instant toggle) |
| Audit trail | No | Yes (full history) |
| Environment separation | No | Yes (production/staging/dev/test) |

---

## Migration Steps

### Step 1: Access the Admin Dashboard

Navigate to the admin dashboard:
```
http://localhost:8080/admin
```

Login with admin credentials configured in the `token_management` database.

### Step 2: Create an API Token

**Via Dashboard:**
1. Click "Create Token"
2. Fill in token details:
   - **Name:** e.g., "Production App"
   - **Environment:** production, staging, development, or test
   - **Scopes:** e.g., tickets:read, tickets:write, machines:read
   - **Rate Limits:** requests per minute/hour/day
   - **IP Whitelist:** (optional) restrict to specific IPs
   - **Expiration:** (optional) set an expiry date
3. Click "Create"
4. **Save the token immediately** -- it won't be shown again

**Via API:**
```bash
# Login to get session token
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "your_password"}' \
  http://localhost:8080/api/v1/admin/auth/login

# Create API token
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

### Step 3: Update Your Application

Replace the `X-API-Key` header with `X-API-Token`:

```diff
- headers: { 'X-API-Key': 'your_simple_api_key' }
+ headers: { 'X-API-Token': 'tok_live_abc123xyz456' }
```

**That's it.** No URL changes, no field name changes, no response format changes.

### Step 4: Add Rate Limit Handling

Since tokens have rate limits, add handling for `429 Too Many Requests`:

```javascript
const response = await fetch('/api/v1/tickets', {
  headers: { 'X-API-Token': token }
});

if (response.status === 429) {
  // Back off and retry
  const delay = 5000; // 5 seconds
  await new Promise(resolve => setTimeout(resolve, delay));
  // Retry request...
}
```

### Step 5: Monitor Usage

After migrating, monitor your token's usage via:
- **Admin Dashboard:** `http://localhost:8080/admin`
- **Token Analytics API:** `GET /api/v1/admin/analytics/tokens/:id`
- **Usage Logs API:** `GET /api/v1/admin/tokens/:id/logs`

---

## Authentication Comparison

### Before (Legacy API Key)

```bash
curl -H "X-API-Key: your_simple_api_key" \
  http://localhost:8080/api/v1/tickets
```

### After (API Token)

```bash
curl -H "X-API-Token: tok_live_abc123xyz456" \
  http://localhost:8080/api/v1/tickets
```

**Note:** The endpoint URL, request format, and response format are identical. Only the authentication header changes.

---

## Code Examples

### JavaScript/Node.js

#### Before:
```javascript
const API_KEY = 'your_simple_api_key';

async function getTickets() {
  const response = await fetch('http://localhost:8080/api/v1/tickets', {
    headers: { 'X-API-Key': API_KEY }
  });
  return response.json();
}
```

#### After:
```javascript
const API_TOKEN = 'tok_live_abc123xyz456';

async function getTickets() {
  const response = await fetch('http://localhost:8080/api/v1/tickets', {
    headers: { 'X-API-Token': API_TOKEN }
  });

  // Handle rate limiting
  if (response.status === 429) {
    await new Promise(r => setTimeout(r, 5000));
    return getTickets(); // Retry
  }

  return response.json();
}
```

### Python

#### Before:
```python
import requests

API_KEY = 'your_simple_api_key'
BASE_URL = 'http://localhost:8080/api/v1'

def get_tickets():
    headers = {'X-API-Key': API_KEY}
    response = requests.get(f'{BASE_URL}/tickets', headers=headers)
    return response.json()
```

#### After:
```python
import requests
import time

API_TOKEN = 'tok_live_abc123xyz456'
BASE_URL = 'http://localhost:8080/api/v1'

def get_tickets():
    headers = {'X-API-Token': API_TOKEN}
    response = requests.get(f'{BASE_URL}/tickets', headers=headers)

    # Handle rate limiting
    if response.status_code == 429:
        time.sleep(5)
        return get_tickets()  # Retry

    return response.json()
```

### cURL

#### Before:
```bash
# Get tickets
curl -H "X-API-Key: your_key" http://localhost:8080/api/v1/tickets

# Create ticket
curl -X POST \
  -H "X-API-Key: your_key" \
  -H "Content-Type: application/json" \
  -d '{"terminal_id": "ATM-001", ...}' \
  http://localhost:8080/api/v1/tickets
```

#### After:
```bash
# Get tickets (same URL, different auth header)
curl -H "X-API-Token: tok_live_abc123" http://localhost:8080/api/v1/tickets

# Create ticket (same URL, different auth header)
curl -X POST \
  -H "X-API-Token: tok_live_abc123" \
  -H "Content-Type: application/json" \
  -d '{"terminal_id": "ATM-001", ...}' \
  http://localhost:8080/api/v1/tickets
```

---

## New Features Available After Migration

Once you've migrated to token-based auth, you gain access to these capabilities:

### 1. Per-Token Analytics

View detailed usage data for each token:
```bash
curl -H "X-Session-Token: sess_abc123" \
  http://localhost:8080/api/v1/admin/analytics/tokens/1?days=7
```

### 2. Endpoint Statistics

See which endpoints are used most:
```bash
curl -H "X-Session-Token: sess_abc123" \
  http://localhost:8080/api/v1/admin/analytics/endpoints?days=7&limit=20
```

### 3. Daily Usage Trends

Track daily request volumes:
```bash
curl -H "X-Session-Token: sess_abc123" \
  http://localhost:8080/api/v1/admin/analytics/daily?days=30
```

### 4. Audit Logs

Review all administrative actions:
```bash
curl -H "X-Session-Token: sess_abc123" \
  http://localhost:8080/api/v1/admin/audit-logs?limit=100
```

### 5. Token Lifecycle Management

```bash
# Disable a token temporarily
curl -X PATCH \
  -H "X-Session-Token: sess_abc123" \
  http://localhost:8080/api/v1/admin/tokens/1/disable

# Re-enable it
curl -X PATCH \
  -H "X-Session-Token: sess_abc123" \
  http://localhost:8080/api/v1/admin/tokens/1/enable
```

### 6. Token Rotation

Rotate tokens without downtime:

1. Create a new token (Token B)
2. Update half of your services to use Token B
3. Monitor for 24 hours
4. Update remaining services to use Token B
5. Disable old token (Token A)
6. Delete Token A after 7 days

---

## FAQ

### Q: Can I use both `X-API-Key` and `X-API-Token` simultaneously?

**A:** Yes. Both authentication methods work on all `/api/v1` endpoints. You can migrate your applications one at a time.

### Q: Do the URLs change?

**A:** No. All endpoints remain at `/api/v1/...`. Only the authentication header changes.

### Q: Do response formats change?

**A:** No. Request and response formats are identical regardless of which auth method you use.

### Q: What happens if I don't migrate?

**A:** The `X-API-Key` auth method continues to work. However, you won't benefit from rate limiting, usage analytics, scoped permissions, or any of the token management features.

### Q: Do I need new database credentials?

**A:** No. Both auth methods connect to the same `ticket_master` and `machine_master` databases. The token system uses a separate `token_management` database that the gateway manages internally.

### Q: How do I handle rate limits?

**A:** When rate limits are exceeded, the API returns `429 Too Many Requests`. Implement exponential backoff in your client:

```javascript
if (response.status === 429) {
  await new Promise(r => setTimeout(r, 5000));
  // Retry...
}
```

### Q: Can I create tokens without the web dashboard?

**A:** Yes. Use the admin API endpoints:
1. Login: `POST /api/v1/admin/auth/login`
2. Create token: `POST /api/v1/admin/tokens`

### Q: How many tokens can I create?

**A:** There is no hard limit. Create separate tokens per application, environment, or team as needed.

### Q: What scopes are available?

**A:** Scopes are flexible strings. Common patterns include:
- `tickets:read` - Read ticket data
- `tickets:write` - Create/update tickets
- `machines:read` - Read machine data
- `machines:write` - Update machine status

If no scopes are defined on a token, all endpoints are accessible (backward compatibility).

### Q: What if a token expires?

**A:** Expired tokens return `401 Unauthorized`. Create a new token via the admin dashboard or API before the old one expires.

---

## Migration Checklist

### Pre-Migration
- [ ] Read this migration guide
- [ ] Access the admin dashboard at `/admin`
- [ ] Create an API token with appropriate scopes and rate limits
- [ ] Save the token securely

### Code Changes
- [ ] Replace `X-API-Key` header with `X-API-Token` header
- [ ] Add `429 Too Many Requests` (rate limit) handling
- [ ] Store the token securely (environment variable, secrets manager)

### Testing
- [ ] Test all endpoints with the new token
- [ ] Verify rate limits work as expected
- [ ] Check usage logs in admin dashboard

### Post-Migration
- [ ] Monitor token analytics for errors
- [ ] Disable old API key if all applications have migrated
- [ ] Set up token rotation schedule

---

## Support

- **Admin Dashboard:** `http://localhost:8080/admin`
- **Swagger UI:** `http://localhost:8080/swagger/index.html`
- **API Documentation:** See `API_DOCUMENTATION.md`
- **Swagger Guide:** See `SWAGGER_GUIDE.md`

---

*Last Updated: February 2025*
