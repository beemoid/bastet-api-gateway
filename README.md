# BASTET API Gateway

Middleware between on-premise SQL Server databases and cloud applications. Exposes a single unified REST endpoint that always JOINs `ticket_master` with `machine_master`, applies vendor-scoped access control from the API token, and includes a full token management admin dashboard.

## Architecture

```
[Cloud / Vendor App] <---> [API Gateway] <---> [On-Premise SQL Server]
                               |                  ├── ticket_master.dbo.open_ticket
                               |                  ├── machine_master.dbo.machine
                               |                  └── token_management (tokens, sessions, audit)
                               |
                               ├── /api/v1/data        ← unified data endpoint
                               ├── /admin              ← web dashboard
                               └── /api/v1/admin/*     ← token & analytics API
```

## Features

- **Single unified endpoint** — `/api/v1/data` always returns joined ticket + machine rows
- **Vendor-scoped tokens** — each token can be restricted to a specific vendor via `filter_column` / `filter_value` (e.g. `mm.[FLM name] = 'AVT'`)
- **Admin / Internal tokens** — `is_super_token=true` bypasses all filters using a customizable admin query
- **Full pagination** — `page`, `page_size`, `sort_by`, `sort_order`, `search`, `status`, `mode`, `priority`
- **Token management** — create, update, disable, delete tokens via dashboard or API
- **Rate limiting** — configurable per-token limits (per minute, hour, day)
- **IP whitelisting** — optional per-token IP restriction
- **Analytics & audit logs** — per-token usage tracking, endpoint stats, daily charts
- **Admin dashboard** — web UI at `/admin`
- **Gzip compression** — automatic response compression
- **Structured JSON logs** — via logrus
- **Graceful shutdown** — cleans up DB connections on SIGINT/SIGTERM
- **systemd service** — `service.sh` for bare-metal deployment
- **Docker support** — multi-stage image, `docker-compose.yml` ready

## Prerequisites

- Go 1.21+
- Microsoft SQL Server (on-premise)
- Access to `ticket_master`, `machine_master`, and `token_management` databases

---

## Quick Start

### 1. Configure environment

```bash
cp .env.example .env
# Edit .env with your database credentials and port
```

Key variables:

```env
SERVER_PORT=8080
GIN_MODE=release

TICKET_DB_HOST=localhost
TICKET_DB_PORT=1433
TICKET_DB_USER=sa
TICKET_DB_PASSWORD=your_password
TICKET_DB_NAME=ticket_master

MACHINE_DB_HOST=localhost
MACHINE_DB_PORT=1433
MACHINE_DB_USER=sa
MACHINE_DB_PASSWORD=your_password
MACHINE_DB_NAME=machine_master

TOKEN_DB_HOST=localhost
TOKEN_DB_PORT=1433
TOKEN_DB_USER=sa
TOKEN_DB_PASSWORD=your_password
TOKEN_DB_NAME=token_management

JWT_SECRET=change-this-to-a-long-random-secret
```

### 2. Install dependencies

```bash
make install
```

### 3. Run locally

```bash
make run
```

### 4. Access

| URL | Description |
|-----|-------------|
| `http://localhost:8080/health` | Health check |
| `http://localhost:8080/admin` | Admin dashboard |
| `http://localhost:8080/api/v1/data` | Data endpoint (requires token) |

---

## API Reference

### Authentication

All `/api/v1/data` endpoints require `X-API-Token`:

```
X-API-Token: tok_live_abc123xyz456
```

Two token types:

| Type | Behavior |
|------|----------|
| **Vendor token** (`filter_column` set) | Only sees rows matching `mm.[col] = value` |
| **Admin / Internal token** (`is_super_token=true`) | Full access, uses customizable admin query |

Admin dashboard endpoints use session auth (`X-Session-Token`).

---

### Data Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/data` | List all rows (paginated, filtered, sorted) |
| `GET` | `/api/v1/data/metadata` | Distinct status / mode / priority values |
| `GET` | `/api/v1/data/:terminal_id` | Single row by terminal ID |
| `PUT` | `/api/v1/data/:terminal_id` | Update ticket fields |

#### Query parameters for `GET /api/v1/data`

| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `page` | integer | Page number. Omit for all results. | — |
| `page_size` | integer | Items per page (max 500) | 100 |
| `sort_by` | string | Field to sort by (see below) | `incident_start_datetime` |
| `sort_order` | string | `asc` or `desc` | `desc` |
| `search` | string | Partial match on terminal_id or terminal_name | — |
| `status` | string | Exact match (e.g. `0.NEW`) | — |
| `mode` | string | Exact match (e.g. `Off-line`) | — |
| `priority` | string | Exact match (e.g. `1.High`) | — |

Sortable fields: `terminal_id`, `terminal_name`, `priority`, `mode`, `status`, `incident_start_datetime`, `count`, `balance`, `tickets_duration`, `open_time`, `close_time`, `flm_name`, `flm`, `slm`, `net`

#### Example requests

```bash
# All data for your vendor scope (no pagination)
curl -H "X-API-Token: tok_live_xxx" http://localhost:8080/api/v1/data

# Page 1, 50 rows, sorted by status ascending, only Off-line mode
curl -H "X-API-Token: tok_live_xxx" \
  "http://localhost:8080/api/v1/data?page=1&page_size=50&sort_by=status&sort_order=asc&mode=Off-line"

# Search by terminal name
curl -H "X-API-Token: tok_live_xxx" \
  "http://localhost:8080/api/v1/data?search=PULO+BAMBU"

# Update a ticket
curl -X PUT \
  -H "X-API-Token: tok_live_xxx" \
  -H "Content-Type: application/json" \
  -d '{"status": "2.Kirim FLM", "remarks": "Technician dispatched"}' \
  http://localhost:8080/api/v1/data/ATM-001
```

---

### Admin Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/admin/auth/login` | Admin login, returns session token |
| `POST` | `/api/v1/admin/auth/logout` | Invalidate session |
| `GET` | `/api/v1/admin/auth/me` | Current admin user |
| `GET` | `/api/v1/admin/tokens` | List all API tokens |
| `POST` | `/api/v1/admin/tokens` | Create API token |
| `GET` | `/api/v1/admin/tokens/:id` | Get token details |
| `PUT` | `/api/v1/admin/tokens/:id` | Update token |
| `DELETE` | `/api/v1/admin/tokens/:id` | Delete token |
| `PATCH` | `/api/v1/admin/tokens/:id/disable` | Disable token |
| `PATCH` | `/api/v1/admin/tokens/:id/enable` | Enable token |
| `GET` | `/api/v1/admin/tokens/:id/logs` | Token usage logs |
| `GET` | `/api/v1/admin/analytics/dashboard` | Dashboard statistics |
| `GET` | `/api/v1/admin/analytics/tokens/:id` | Per-token analytics |
| `GET` | `/api/v1/admin/analytics/endpoints` | Endpoint usage stats |
| `GET` | `/api/v1/admin/analytics/daily` | Daily usage chart |
| `GET` | `/api/v1/admin/audit-logs` | Audit log entries |

### Health Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | API + database health check |
| `GET` | `/ping` | Liveness check |

---

## Swagger / API Docs

Swagger is **not served from this application**. The JSON spec files are in `docs/` — open them in any external Swagger UI:

| File | For |
|------|-----|
| `docs/swagger.json` | Internal / full spec (all endpoints) |
| `docs/swagger_public.json` | External parties (data + health only, no admin) |

To view: paste the file contents into [https://editor.swagger.io](https://editor.swagger.io), or serve it with any static Swagger UI.

---

## Project Structure

```
api-gateway/
├── config/
│   └── config.go                        # Env-driven configuration
├── database/
│   ├── database.go                      # DB connection manager
│   └── migrations/
│       ├── 001_create_token_management_schema.sql
│       └── 002_add_vendor_filter_to_tokens.sql
├── docs/
│   ├── swagger.json                     # Full private API spec
│   └── swagger_public.json             # Public API spec (data + health only)
├── handlers/
│   ├── data_handler.go                  # GET/PUT /api/v1/data
│   ├── health_handler.go                # /health, /ping
│   └── token_handler.go                 # Admin, token management, analytics
├── middleware/
│   ├── auth.go                          # CombinedAuth (token validation + vendor filter)
│   ├── token_auth.go                    # Admin session auth, rate limit, scope check
│   ├── cors.go                          # CORS
│   └── logger.go                        # Request logging
├── models/
│   ├── data.go                          # DataRow, DataListResponse, DataUpdateRequest
│   ├── token.go                         # APIToken, AdminUser, session, audit models
│   ├── analytics.go                     # Analytics response types
│   ├── nullable.go                      # NullString, NullTime helpers
│   ├── ticket_constants.go              # Ticket status/mode/priority metadata
│   └── machine_constants.go             # Machine metadata
├── repository/
│   ├── data_repository.go               # GetAll, GetByTerminalID, Update + VendorFilter
│   ├── queries/
│   │   └── admin_data_query.go          # Customizable admin SELECT query
│   └── token_repository.go             # Token CRUD, sessions, audit, analytics
├── routes/
│   └── routes.go                        # All route definitions
├── service/
│   ├── data_service.go                  # Data business logic + metadata cache
│   ├── token_service.go                 # Token validation, rate limiting, analytics
│   └── errors.go                        # Custom error types
├── templates/
│   ├── login.html                       # Admin login page
│   ├── dashboard.html                   # Token management dashboard
│   └── assets/                          # Static assets
├── main.go                              # Entry point
├── go.mod / go.sum
├── .env.example                         # Config template
├── Makefile                             # Build, dist, docker commands
├── Dockerfile                           # Multi-stage build
├── docker-compose.yml                   # Docker deployment
├── service.sh                           # systemd service manager
├── README.md
├── API_DOCUMENTATION.md                 # Full API reference
└── MIGRATION_GUIDE_V1_TO_V2.md         # Legacy → token auth migration guide
```

---

## Deployment

### Bare-metal (systemd)

```bash
# 1. Build release package
make dist
# → dist/api-gateway/ contains: binary, templates/, docs/, .env.example, service.sh

# 2. Copy to server
scp -r dist/api-gateway/ user@server:/opt/bastet-api-gateway/

# 3. On the server
cd /opt/bastet-api-gateway
cp .env.example .env      # fill in your values
sudo ./service.sh install # install + enable (auto-start on reboot)
sudo ./service.sh start
```

**Service commands:**

```bash
./service.sh start      # start the service
./service.sh stop       # stop the service
./service.sh restart    # restart the service
./service.sh status     # show current status
./service.sh log        # tail live logs (Ctrl+C to exit)
./service.sh log-all    # full log history
sudo ./service.sh install    # install as systemd service
sudo ./service.sh uninstall  # remove service
```

### Docker

```bash
# 1. Build Docker release package
make docker-dist
# → dist/docker/ contains: Dockerfile, docker-compose.yml, .env.example

# 2. Copy to server and deploy
cp .env.example .env   # fill in your values
docker compose up -d
```

Or build and run directly:

```bash
make docker-build
# Image: api-gateway:latest
# Then on the server: docker compose up -d
```

### Makefile reference

```bash
make install       # download Go dependencies
make build         # compile → bin/api-gateway
make run           # build + run locally
make test          # run tests
make clean         # remove bin/ and dist/
make dist          # bare-metal release package → dist/api-gateway/
make docker-build  # build Docker image
make docker-dist   # Docker release package → dist/docker/
```

---

## Customizing the Admin Query

Admin / Internal tokens (`is_super_token=true`) use a dedicated SELECT defined in:

```
repository/queries/admin_data_query.go
```

Edit the `AdminDataQuery` constant there to change columns, ordering, or add filters for the admin view — without touching any other code.

---

## Security

- **Vendor-scoped tokens** — DB-level filtering, not application-level
- **Admin session management** — session tokens with expiration
- **Rate limiting** — per-token, per minute/hour/day
- **IP whitelisting** — optional per-token
- **Token expiration** — configurable
- **Audit logging** — all admin actions recorded
- **Parameterized queries** — no SQL injection risk; `sort_by` uses an allowlist
- **Non-root Docker user** — container runs as unprivileged `appuser`
- **CORS** — configurable in `middleware/cors.go`

---

## Production Checklist

- [ ] `GIN_MODE=release` in `.env`
- [ ] Strong `JWT_SECRET` (random, 32+ characters)
- [ ] Run DB migration `002_add_vendor_filter_to_tokens.sql`
- [ ] Create at least one admin user in `token_management`
- [ ] Configure rate limits on all tokens
- [ ] Put a reverse proxy (nginx) with TLS in front
- [ ] Restrict CORS origins in `middleware/cors.go`
- [ ] Set `restart: unless-stopped` (Docker) or use `service.sh install` (systemd)

---

## License

Proprietary — for internal use only.

## Support

For issues or questions, contact the development team.

---

*Built with Go · Gin · SQL Server · systemd*
