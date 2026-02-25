# Migration Guide: Tickets/Machines Endpoints → Unified Data Endpoint

## Overview

This guide covers two migration paths:

1. **Legacy API Key → Token-based auth** — if you were using `X-API-Key`
2. **Old endpoints → `/api/v1/data`** — if you were calling `/api/v1/tickets` or `/api/v1/machines`

Both changes happened together. The new system uses a single unified endpoint that always JOINs ticket and machine data, with access automatically scoped to your vendor via the token.

---

## What Changed

### Removed

| Old | Replaced by |
|-----|-------------|
| `GET /api/v1/tickets` | `GET /api/v1/data` |
| `GET /api/v1/tickets/:id` | `GET /api/v1/data/:terminal_id` |
| `GET /api/v1/tickets/metadata` | `GET /api/v1/data/metadata` |
| `PUT /api/v1/tickets/:id` | `PUT /api/v1/data/:terminal_id` |
| `GET /api/v1/machines` | merged into `GET /api/v1/data` |
| `GET /api/v1/machines/:terminal_id` | merged into `GET /api/v1/data/:terminal_id` |
| `X-API-Key` header | `X-API-Token` header |
| Swagger at `/swagger/index.html` | Static JSON files in `docs/` |

### Added

| Feature | Details |
|---------|---------|
| Unified JOIN response | Every row includes both ticket fields and machine fields (`flm_name`, `flm`, `slm`, `net`) |
| Vendor-scoped tokens | Token carries `filter_column` + `filter_value` — DB filters automatically applied |
| Admin / Internal tokens | `is_super_token=true` — sees all data via customizable admin query |
| Full pagination | `page`, `page_size`, `sort_by`, `sort_order`, `search`, `status`, `mode`, `priority` |
| Metadata endpoint | `GET /api/v1/data/metadata` — distinct status/mode/priority from live DB (1-hour cache) |
| Token management dashboard | Web UI at `/admin` |
| systemd service | `service.sh install/start/stop/restart/status/log` |

---

## Authentication Migration

### Before (Legacy API Key)

```bash
curl -H "X-API-Key: your_simple_api_key" \
  http://localhost:8080/api/v1/tickets
```

### After (API Token)

```bash
curl -H "X-API-Token: tok_live_abc123xyz456" \
  http://localhost:8080/api/v1/data
```

**Only two things change:** the header name and the endpoint URL. Response format has a new shape (see below).

---

## Endpoint Migration

### Get all data

```diff
- GET /api/v1/tickets
- GET /api/v1/machines
+ GET /api/v1/data
```

The new response combines both. Vendor tokens automatically see only their rows — you don't need to add any filter yourself.

```bash
# Old — two separate calls needed
curl -H "X-API-Key: key" http://localhost:8080/api/v1/tickets
curl -H "X-API-Key: key" http://localhost:8080/api/v1/machines

# New — one call, joined result
curl -H "X-API-Token: tok_live_xxx" http://localhost:8080/api/v1/data
```

### Get single row

```diff
- GET /api/v1/tickets/:id
+ GET /api/v1/data/:terminal_id
```

```bash
# Old
curl -H "X-API-Key: key" http://localhost:8080/api/v1/tickets/ATM-001

# New
curl -H "X-API-Token: tok_live_xxx" http://localhost:8080/api/v1/data/ATM-001
```

### Update ticket fields

```diff
- PUT /api/v1/tickets/:id
+ PUT /api/v1/data/:terminal_id
```

Request body fields are the same:

```bash
curl -X PUT \
  -H "X-API-Token: tok_live_xxx" \
  -H "Content-Type: application/json" \
  -d '{"status": "2.Kirim FLM", "remarks": "Technician dispatched"}' \
  http://localhost:8080/api/v1/data/ATM-001
```

Vendor tokens return `403` if the terminal is outside their scope. Admin tokens can update any terminal.

### Metadata

```diff
- GET /api/v1/tickets/metadata
- GET /api/v1/machines/metadata
+ GET /api/v1/data/metadata
```

```bash
curl -H "X-API-Token: tok_live_xxx" http://localhost:8080/api/v1/data/metadata
```

---

## Response Format Changes

### List response

Old `/api/v1/tickets` returned a flat array or simple wrapper. New `/api/v1/data` always returns:

```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": [ { ...DataRow... } ],
  "total": 671,
  "page": 1,
  "page_size": 100,
  "total_pages": 7,
  "sort_by": "incident_start_datetime",
  "sort_order": "desc",
  "search": "",
  "status": "",
  "mode": "",
  "priority": ""
}
```

Pagination fields (`page`, `page_size`, `total_pages`) are omitted when `page` is not requested.

### DataRow now includes machine fields

Every row now includes four additional fields from the `machine_master` JOIN:

```json
{
  "terminal_id": "ATM-001",
  "terminal_name": "Main Branch ATM",
  "status": "0.NEW",
  "mode": "Off-line",
  "...all ticket fields...",
  "flm_name": "AVT",
  "flm": "AVT - BANDUNG",
  "slm": "KGP - WINCOR DW",
  "net": "SMS"
}
```

---

## Pagination

### Get all results (no pagination)

Simply omit `page`. All matching rows are returned:

```bash
curl -H "X-API-Token: tok_live_xxx" http://localhost:8080/api/v1/data
```

### Paginated

```bash
# Page 1, 50 rows per page
curl -H "X-API-Token: tok_live_xxx" \
  "http://localhost:8080/api/v1/data?page=1&page_size=50"
```

### Sorted

```bash
# Sort by status ascending
curl -H "X-API-Token: tok_live_xxx" \
  "http://localhost:8080/api/v1/data?sort_by=status&sort_order=asc"
```

### Filtered

```bash
# Only Off-line terminals, high priority
curl -H "X-API-Token: tok_live_xxx" \
  "http://localhost:8080/api/v1/data?mode=Off-line&priority=1.High"

# Use /metadata to discover valid filter values
curl -H "X-API-Token: tok_live_xxx" http://localhost:8080/api/v1/data/metadata
```

---

## Token Creation

### Via Admin Dashboard

1. Go to `http://localhost:8080/admin`
2. Click **Create Token**
3. Fill in:
   - **Name** — e.g. `AVT Production`
   - **Environment** — `production`
   - **Vendor Name** — e.g. `AVT`
   - **Filter Column** — e.g. `flm_name`
   - **Filter Value** — e.g. `AVT`
   - **Admin / Internal Token** — toggle ON to bypass all vendor filters
   - **Rate Limits** — requests per minute/hour/day
4. Click **Create** — **save the token immediately, it won't be shown again**

### Filter Column Options

| `filter_column` | SQL expression | Example `filter_value` |
|-----------------|----------------|------------------------|
| `flm_name` | `mm.[FLM name]` | `AVT` |
| `flm` | `mm.[FLM]` | `AVT - BANDUNG` |
| `slm` | `mm.[SLM]` | `KGP - WINCOR DW` |
| `net` | `mm.[Net]` | `SMS` |
| `terminal_id` | `op.[Terminal ID]` | `ATM-001` |
| `status` | `op.[Status]` | `0.NEW` |
| `priority` | `op.[Priority]` | `1.High` |

### Via API

```bash
# 1. Login
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "your_password"}' \
  http://localhost:8080/api/v1/admin/auth/login

# 2. Create a vendor token
curl -X POST \
  -H "X-Session-Token: sess_abc123" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "AVT Production",
    "environment": "production",
    "vendor_name": "AVT",
    "filter_column": "flm_name",
    "filter_value": "AVT",
    "is_super_token": false,
    "rate_limit_per_minute": 60,
    "rate_limit_per_hour": 1000,
    "rate_limit_per_day": 10000
  }' \
  http://localhost:8080/api/v1/admin/tokens
```

---

## Code Examples

### JavaScript / Fetch

```javascript
const API_TOKEN = 'tok_live_abc123xyz456';
const BASE = 'http://localhost:8080/api/v1';

async function getData(params = {}) {
  const qs = new URLSearchParams(params).toString();
  const url = `${BASE}/data${qs ? '?' + qs : ''}`;

  const res = await fetch(url, {
    headers: { 'X-API-Token': API_TOKEN }
  });

  if (res.status === 429) {
    // Rate limited — wait and retry
    await new Promise(r => setTimeout(r, 5000));
    return getData(params);
  }

  return res.json();
}

// Get all data
const all = await getData();

// Paginated, sorted, filtered
const page1 = await getData({
  page: 1,
  page_size: 50,
  sort_by: 'status',
  sort_order: 'asc',
  mode: 'Off-line'
});
```

### Python

```python
import requests
import time

API_TOKEN = 'tok_live_abc123xyz456'
BASE = 'http://localhost:8080/api/v1'
HEADERS = {'X-API-Token': API_TOKEN}

def get_data(**params):
    res = requests.get(f'{BASE}/data', headers=HEADERS, params=params)

    if res.status_code == 429:
        time.sleep(5)
        return get_data(**params)

    res.raise_for_status()
    return res.json()

# All data
all_data = get_data()

# Paginated + filtered
page = get_data(page=1, page_size=50, sort_by='status', mode='Off-line')

# Update a terminal
def update_terminal(terminal_id, fields):
    res = requests.put(
        f'{BASE}/data/{terminal_id}',
        headers=HEADERS,
        json=fields
    )
    res.raise_for_status()
    return res.json()

update_terminal('ATM-001', {'status': '2.Kirim FLM', 'remarks': 'Dispatched'})
```

### cURL

```bash
TOKEN="tok_live_abc123"

# All data
curl -H "X-API-Token: $TOKEN" http://localhost:8080/api/v1/data

# Page 1, sort by status ascending, filter by mode
curl -H "X-API-Token: $TOKEN" \
  "http://localhost:8080/api/v1/data?page=1&page_size=50&sort_by=status&sort_order=asc&mode=Off-line"

# Search
curl -H "X-API-Token: $TOKEN" \
  "http://localhost:8080/api/v1/data?search=PULO+BAMBU"

# Single terminal
curl -H "X-API-Token: $TOKEN" http://localhost:8080/api/v1/data/1001200

# Update
curl -X PUT \
  -H "X-API-Token: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "2.Kirim FLM"}' \
  http://localhost:8080/api/v1/data/1001200

# Valid filter values
curl -H "X-API-Token: $TOKEN" http://localhost:8080/api/v1/data/metadata
```

---

## Swagger / API Docs

Swagger is no longer served from the application itself. Use the static JSON files:

| File | Use |
|------|-----|
| `docs/swagger.json` | Full spec — internal use |
| `docs/swagger_public.json` | Public spec — share with vendors (data + health only) |

Open either file at [https://editor.swagger.io](https://editor.swagger.io) by pasting the contents.

---

## Migration Checklist

### Authentication
- [ ] Replace `X-API-Key` header with `X-API-Token`
- [ ] Create tokens via admin dashboard at `/admin` or via `POST /api/v1/admin/tokens`
- [ ] Add `429 Too Many Requests` handling with retry/backoff

### Endpoint URLs
- [ ] Change `/api/v1/tickets` → `/api/v1/data`
- [ ] Change `/api/v1/tickets/:id` → `/api/v1/data/:terminal_id`
- [ ] Change `/api/v1/machines` → `/api/v1/data` (machine fields included automatically)
- [ ] Change metadata endpoint → `/api/v1/data/metadata`
- [ ] Remove any `X-API-Key` references

### Response parsing
- [ ] Update to read new `DataListResponse` shape (`data`, `total`, `page`, etc.)
- [ ] Map new field names if needed (`flm_name`, `flm`, `slm`, `net` now always present)

### Post-migration
- [ ] Monitor token analytics via admin dashboard or `GET /api/v1/admin/analytics/tokens/:id`
- [ ] Set token expiry and rate limits appropriate for each consumer

---

## FAQ

**Q: Do I still get machine data separately?**
A: No — it's always included in every row. `flm_name`, `flm`, `slm`, `net` are always present in the response (may be `null` if the terminal has no machine record).

**Q: Can I get all data without pagination?**
A: Yes — omit the `page` parameter entirely. All rows matching your vendor scope (and any filters) are returned.

**Q: My token is a vendor token. Do I need to pass filter params?**
A: No. The filter is read from the token itself by the middleware. You never need to pass `filter_column` or `filter_value` in requests.

**Q: What if I need to see all vendors' data?**
A: Create an Admin / Internal token (`is_super_token=true`) via the dashboard. That token bypasses all vendor filters.

**Q: How do I know which status/mode/priority values are valid?**
A: Call `GET /api/v1/data/metadata` — it returns distinct values from the live database, cached for 1 hour.

**Q: Where is Swagger UI now?**
A: Not hosted by the app anymore. Use the JSON files in `docs/` with any external Swagger UI (e.g. editor.swagger.io). Share `docs/swagger_public.json` with external parties.

---

*Last Updated: February 2026*
