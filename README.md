# API Gateway for On-Premise to Cloud Communication

This API Gateway serves as the middleware between your on-premise databases and cloud applications. It provides RESTful APIs to manage tickets and terminal/machine information, with a complete token management system and admin dashboard.

## Architecture

```
[Cloud App] <---> [API Gateway] <---> [On-Premise Databases]
                   |                   ├── ticket_master.dbo.open_ticket
                   |                   ├── machine_master.dbo.atmi
                   |                   └── token_management (tokens, sessions, audit)
                   |
                   ├── Admin Dashboard (Web UI)
                   ├── Token Management API
                   ├── Analytics & Audit Logs
                   └── Swagger Documentation
```

## Features

- **RESTful API** for ticket and machine management (v1)
- **Token Management System** - Scoped API tokens with rate limiting, IP whitelisting, and expiration
- **Admin Dashboard** - Web UI for managing tokens, viewing analytics, and audit logs
- **Dual Authentication** - Supports both `X-API-Key` (legacy) and `X-API-Token` (token-based)
- **Rate Limiting** - Per-token configurable limits (per minute, hour, day)
- **Triple Database Support** - Connects to ticket_master, machine_master, and token_management databases
- **Hybrid Adaptive Metadata** - Auto-discovers field values from database with 1-hour caching
- **Analytics & Monitoring** - Dashboard stats, endpoint analytics, daily usage tracking
- **Audit Logging** - Complete history of all administrative actions
- **Gzip Compression** - Automatic response compression for reduced bandwidth
- **CORS Support** - Configurable cross-origin requests
- **Health Checks** - Monitor API and database status
- **Structured Logging** - JSON formatted logs with logrus
- **Graceful Shutdown** - Properly closes database connections
- **Swagger/OpenAPI** - Interactive API documentation

## Prerequisites

- Go 1.21 or higher
- SQL Server (on-premise databases)
- Access to `ticket_master`, `machine_master`, and `token_management` databases

## Installation

1. **Clone or navigate to the project directory**
   ```bash
   cd /path/to/api-gateway
   ```

2. **Install dependencies**
   ```bash
   make install
   # Or manually:
   # go mod download
   # go install github.com/swaggo/swag/cmd/swag@latest
   ```

3. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

4. **Generate Swagger documentation**
   ```bash
   make swagger
   ```

5. **Run the application**
   ```bash
   make run
   # Or manually:
   # go run main.go
   ```

## Configuration

Create a `.env` file based on `.env.example`:

```env
# Server Configuration
SERVER_PORT=8080
GIN_MODE=debug

# Ticket Database (ticket_master)
TICKET_DB_HOST=localhost
TICKET_DB_PORT=1433
TICKET_DB_USER=your_username
TICKET_DB_PASSWORD=your_password
TICKET_DB_NAME=ticket_master

# Machine Database (machine_master)
MACHINE_DB_HOST=localhost
MACHINE_DB_PORT=1433
MACHINE_DB_USER=your_username
MACHINE_DB_PASSWORD=your_password
MACHINE_DB_NAME=machine_master

# Token Management Database (token_management)
# Defaults to same host/credentials as Ticket DB if not set
TOKEN_DB_HOST=localhost
TOKEN_DB_PORT=1433
TOKEN_DB_USER=your_username
TOKEN_DB_PASSWORD=your_password
TOKEN_DB_NAME=token_management

# Cloud App Configuration
CLOUD_APP_URL=https://your-cloud-app.com
CLOUD_APP_API_KEY=your_api_key

# Security
JWT_SECRET=your_jwt_secret_key
API_KEY=your_internal_api_key
```

## API Documentation

### Interactive Swagger Documentation

Access the interactive API documentation at:
```
http://localhost:8080/swagger/index.html
```

### Admin Dashboard

Access the token management dashboard at:
```
http://localhost:8080/admin
```

The dashboard provides:
- Token CRUD management (create, view, update, disable, delete)
- Usage analytics and endpoint statistics
- Audit log viewer
- Daily usage charts

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication

The API supports two authentication methods:

**1. API Key (legacy)** - Simple key-based authentication:
```
X-API-Key: your_api_key
```

**2. API Token (recommended)** - Token-based authentication with scopes, rate limiting, and analytics:
```
X-API-Token: tok_live_abc123xyz456
```

Both methods are accepted on all `/api/v1` endpoints. Health checks (`/health`, `/ping`) and Swagger docs require no authentication.

### Additional Documentation

- **[API_DOCUMENTATION.md](API_DOCUMENTATION.md)** - Comprehensive API reference with cURL examples
- **[SWAGGER_GUIDE.md](SWAGGER_GUIDE.md)** - Complete Swagger/OpenAPI documentation guide
- **[MIGRATION_GUIDE_V1_TO_V2.md](MIGRATION_GUIDE_V1_TO_V2.md)** - Migration guide from v1 to v2

---

### Ticket Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/tickets` | Get all tickets |
| `GET` | `/api/v1/tickets/:id` | Get ticket by terminal ID |
| `GET` | `/api/v1/tickets/number/:number` | Get ticket by ticket number |
| `GET` | `/api/v1/tickets/status/:status` | Get tickets by status |
| `GET` | `/api/v1/tickets/terminal/:terminal_id` | Get tickets by terminal |
| `GET` | `/api/v1/tickets/metadata` | Get adaptive ticket metadata |
| `POST` | `/api/v1/tickets` | Create a new ticket |
| `PUT` | `/api/v1/tickets/:id` | Update a ticket |

### Machine Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/machines` | Get all machines |
| `GET` | `/api/v1/machines/:terminal_id` | Get machine by terminal ID |
| `GET` | `/api/v1/machines/status/:status` | Get machines by status |
| `GET` | `/api/v1/machines/branch/:branch_code` | Get machines by branch |
| `GET` | `/api/v1/machines/search` | Search machines with filters |
| `GET` | `/api/v1/machines/metadata` | Get adaptive machine metadata |
| `PATCH` | `/api/v1/machines/status` | Update machine status |

### Admin / Token Management Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/admin/auth/login` | Admin login |
| `POST` | `/api/v1/admin/auth/logout` | Admin logout |
| `GET` | `/api/v1/admin/auth/me` | Get current admin user |
| `GET` | `/api/v1/admin/tokens` | List all tokens |
| `POST` | `/api/v1/admin/tokens` | Create a new token |
| `GET` | `/api/v1/admin/tokens/:id` | Get token details |
| `PUT` | `/api/v1/admin/tokens/:id` | Update a token |
| `DELETE` | `/api/v1/admin/tokens/:id` | Delete a token |
| `PATCH` | `/api/v1/admin/tokens/:id/disable` | Disable a token |
| `PATCH` | `/api/v1/admin/tokens/:id/enable` | Enable a token |
| `GET` | `/api/v1/admin/tokens/:id/logs` | Get token usage logs |

### Analytics Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/admin/analytics/dashboard` | Overall dashboard statistics |
| `GET` | `/api/v1/admin/analytics/tokens/:id` | Token-specific analytics |
| `GET` | `/api/v1/admin/analytics/endpoints` | Endpoint usage statistics |
| `GET` | `/api/v1/admin/analytics/daily` | Daily usage breakdown |
| `GET` | `/api/v1/admin/audit-logs` | Administrative audit logs |

### Health Check Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check (API + databases) |
| `GET` | `/ping` | Simple ping/pong |

---

## Project Structure

```
api-gateway/
├── config/                  # Configuration management
│   └── config.go
├── database/                # Database connection management
│   └── database.go
├── docs/                    # Generated Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── handlers/                # HTTP request handlers
│   ├── ticket_handler.go
│   ├── machine_handler.go
│   ├── health_handler.go
│   └── token_handler.go     # Token management & analytics handlers
├── middleware/              # HTTP middleware
│   ├── auth.go              # API key authentication (legacy)
│   ├── token_auth.go        # Token auth, admin sessions, rate limiting, scope checking
│   ├── cors.go              # CORS configuration
│   └── logger.go            # Request logging
├── models/                  # Data models
│   ├── ticket.go            # Ticket model & request/response types
│   ├── machine.go           # Machine model & request/response types
│   ├── token.go             # Token, admin user, session, audit log models
│   ├── analytics.go         # Dashboard stats, workload, geographic analytics
│   ├── ticket_constants.go  # Ticket metadata constants
│   ├── machine_constants.go # Machine metadata constants
│   └── nullable.go          # Nullable type helpers for SQL
├── repository/              # Data access layer
│   ├── ticket_repository.go
│   ├── machine_repository.go
│   └── token_repository.go  # Token, session, audit log data access
├── routes/                  # Route definitions
│   └── routes.go
├── service/                 # Business logic layer
│   ├── ticket_service.go
│   ├── machine_service.go
│   ├── token_service.go     # Token validation, rate limiting, analytics
│   └── errors.go            # Custom error types
├── templates/               # Admin dashboard HTML templates
│   ├── login.html           # Admin login page
│   ├── dashboard.html       # Token management dashboard
│   └── assets/              # Static assets (CSS, JS)
├── main.go                  # Application entry point
├── go.mod                   # Go module dependencies
├── Makefile                 # Development commands
├── Dockerfile               # Multi-stage Docker build
├── .env.example             # Environment variables template
├── .gitignore               # Git ignore rules
├── README.md                # This file
├── API_DOCUMENTATION.md     # Comprehensive API reference
├── SWAGGER_GUIDE.md         # Swagger documentation guide
└── MIGRATION_GUIDE_V1_TO_V2.md  # v1 to v2 migration guide
```

## Architecture Layers

1. **Handlers** - Handle HTTP requests/responses
2. **Services** - Contain business logic (token validation, rate limiting, caching)
3. **Repositories** - Data access and database operations
4. **Models** - Data structures and validation
5. **Middleware** - Cross-cutting concerns (auth, token auth, logging, CORS, compression)

## Security

- **Dual authentication** - API key (legacy) and token-based authentication
- **Admin session management** - JWT-based sessions with expiration for dashboard access
- **Scoped permissions** - Fine-grained token scopes for access control
- **Rate limiting** - Configurable per-token rate limits (minute/hour/day)
- **IP whitelisting** - Optional IP restriction per token
- **Token expiration** - Configurable token expiry dates
- **Token revocation** - Ability to revoke tokens with audit trail
- **Audit logging** - All administrative actions logged
- **CORS configuration** - Configurable cross-origin access
- **SQL injection prevention** - Parameterized queries throughout
- **Gzip compression** - Reduces response payload sizes
- **Non-root Docker user** - Container runs as unprivileged user

## Testing

### Using Swagger UI (Recommended)

1. Open `http://localhost:8080/swagger/index.html`
2. Click "Authorize" and enter your API key
3. Test any endpoint interactively

### Using cURL

```bash
# Health check
curl http://localhost:8080/health

# Using API key (legacy)
curl -H "X-API-Key: your_api_key" http://localhost:8080/api/v1/tickets

# Using API token
curl -H "X-API-Token: tok_live_abc123" http://localhost:8080/api/v1/tickets
```

### Running Tests

```bash
make test
```

## Deployment

### Using Docker

**Build and run with Docker:**
```bash
make docker-build
make docker-run
```

**Or manually:**
```bash
docker build -t api-gateway:latest .
docker run -p 8080:8080 --env-file .env api-gateway:latest
```

The Docker image uses a multi-stage build with:
- Alpine-based images for minimal size
- Non-root user for security
- Built-in health check
- Swagger docs included

### Production Checklist

1. **Set production environment**
   ```env
   GIN_MODE=release
   ```

2. **Configure JWT secret** - Use a strong, random secret:
   ```env
   JWT_SECRET=your_strong_random_secret
   ```

3. **Update CORS settings** - Restrict origins in `middleware/cors.go`

4. **Secure API keys and tokens** - Use strong, random values

5. **Disable/Protect Swagger UI** - Remove or add auth to Swagger in production

6. **Enable HTTPS** - Use reverse proxy (nginx) with SSL certificates

7. **Database optimization** - Adjust connection pool settings in `database/database.go`

8. **Configure rate limits** - Set appropriate per-token rate limits

9. **Monitor audit logs** - Regularly review admin actions

## Monitoring

- **Health endpoint**: `/health` - Check API and database status
- **Admin dashboard**: `/admin` - View token analytics, endpoint stats, daily usage
- **Analytics API**: `/api/v1/admin/analytics/*` - Programmatic access to metrics
- **Audit logs**: `/api/v1/admin/audit-logs` - Track administrative actions
- **Logs**: JSON formatted logs for easy parsing

## Development

### Makefile Commands

```bash
make help          # Show all available commands
make install       # Install dependencies including Swagger
make swagger       # Generate Swagger documentation
make run           # Generate docs and run the application
make build         # Build the application binary
make test          # Run tests
make clean         # Clean build artifacts
make docker-build  # Build Docker image
make docker-run    # Run Docker container
```

### Adding New Endpoints

1. Create model in `models/`
2. Create repository in `repository/`
3. Create service in `service/`
4. Create handler in `handlers/` with Swagger annotations
5. Register routes in `routes/routes.go`
6. Run `make swagger` to regenerate documentation
7. Test using Swagger UI

## License

This project is proprietary software for internal use.

## Support

For issues or questions, contact the development team.

---

**Built with Go, Gin Framework, and SQL Server**
