# API Gateway for On-Premise to Cloud Communication

This API Gateway serves as the middleware between your on-premise databases and cloud applications. It provides RESTful APIs to manage tickets and terminal/machine information.

## 🏗️ Architecture

```
[Cloud App] <---> [API Gateway] <---> [On-Premise Databases]
                                       ├── ticket_master.dbo.open_ticket
                                       └── machine_master.dbo.atmi
```

## 🚀 Features

- **RESTful API** for ticket and machine management
- **Dual Database Support** - Connects to both ticket_master and machine_master databases
- **Security** - API key authentication middleware
- **CORS Support** - Enables cross-origin requests from cloud app
- **Health Checks** - Monitor API and database status
- **Structured Logging** - JSON formatted logs with logrus
- **Graceful Shutdown** - Properly closes database connections

## 📋 Prerequisites

- Go 1.21 or higher
- SQL Server (On-premise databases)
- Access to `ticket_master` and `machine_master` databases

## 🛠️ Installation

1. **Clone or navigate to the project directory**
   ```bash
   cd /Users/bee/proj/bastet-cloud/api-gateway
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

## ⚙️ Configuration

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

# Cloud App Configuration
CLOUD_APP_URL=https://your-cloud-app.com
CLOUD_APP_API_KEY=your_api_key

# Security
API_KEY=your_internal_api_key
```

## 📚 API Documentation

### Interactive Swagger Documentation

Access the interactive API documentation at:
```
http://localhost:8080/swagger/index.html
```

The Swagger UI allows you to:
- Browse all API endpoints
- View request/response schemas
- Test endpoints directly from the browser
- Authenticate with API key

**Quick Swagger Guide:**
1. Open `http://localhost:8080/swagger/index.html`
2. Click "Authorize" button
3. Enter your API key
4. Test any endpoint by clicking "Try it out"

See [SWAGGER_GUIDE.md](SWAGGER_GUIDE.md) for detailed Swagger documentation.

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication
All API endpoints (except `/health` and `/ping`) require an API key in the request header:
```
X-API-Key: your_api_key
```

### Additional Documentation
- **[API_DOCUMENTATION.md](API_DOCUMENTATION.md)** - Comprehensive API reference with cURL examples
- **[SWAGGER_GUIDE.md](SWAGGER_GUIDE.md)** - Complete Swagger/OpenAPI documentation guide

---

### 🎫 Ticket Endpoints

#### Get All Tickets
```http
GET /api/v1/tickets
```

**Response:**
```json
{
  "success": true,
  "message": "Tickets retrieved successfully",
  "data": [...],
  "total": 10
}
```

#### Get Ticket by ID
```http
GET /api/v1/tickets/:id
```

#### Get Ticket by Number
```http
GET /api/v1/tickets/number/:number
```

#### Get Tickets by Status
```http
GET /api/v1/tickets/status/:status
```
Valid statuses: `Open`, `InProgress`, `Pending`, `Resolved`

#### Get Tickets by Terminal
```http
GET /api/v1/tickets/terminal/:terminal_id
```

#### Create Ticket
```http
POST /api/v1/tickets
Content-Type: application/json

{
  "ticket_number": "TKT-2024-001",
  "terminal_id": "ATM-001",
  "description": "Card reader not working",
  "priority": "High",
  "category": "Hardware",
  "reported_by": "John Doe",
  "assigned_to": "Tech Team"
}
```

#### Update Ticket
```http
PUT /api/v1/tickets/:id
Content-Type: application/json

{
  "status": "Resolved",
  "priority": "Medium",
  "resolution_notes": "Replaced card reader"
}
```

---

### 🖥️ Machine Endpoints

#### Get All Machines
```http
GET /api/v1/machines
```

#### Get Machine by Terminal ID
```http
GET /api/v1/machines/:terminal_id
```

#### Get Machines by Status
```http
GET /api/v1/machines/status/:status
```
Valid statuses: `Active`, `Inactive`, `Maintenance`, `Offline`

#### Get Machines by Branch
```http
GET /api/v1/machines/branch/:branch_code
```

#### Search Machines
```http
GET /api/v1/machines/search?status=Active&branch_code=BR001&location=Jakarta
```

#### Update Machine Status
```http
PATCH /api/v1/machines/status
Content-Type: application/json

{
  "terminal_id": "ATM-001",
  "status": "Maintenance",
  "last_ping_time": "2024-01-15T10:30:00Z",
  "notes": "Scheduled maintenance"
}
```

---

### 🏥 Health Check Endpoints

#### Health Check
```http
GET /health
```

**Response:**
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

#### Ping
```http
GET /ping
```

## 📁 Project Structure

```
api-gateway/
├── config/              # Configuration management
│   └── config.go
├── database/            # Database connection management
│   └── database.go
├── handlers/            # HTTP request handlers
│   ├── ticket_handler.go
│   ├── machine_handler.go
│   └── health_handler.go
├── middleware/          # HTTP middleware
│   ├── auth.go
│   ├── cors.go
│   └── logger.go
├── models/              # Data models
│   ├── ticket.go
│   └── machine.go
├── repository/          # Data access layer
│   ├── ticket_repository.go
│   └── machine_repository.go
├── routes/              # Route definitions
│   └── routes.go
├── service/             # Business logic layer
│   ├── ticket_service.go
│   ├── machine_service.go
│   └── errors.go
├── main.go              # Application entry point
├── go.mod               # Go module dependencies
├── .env.example         # Environment variables template
├── .gitignore          # Git ignore rules
└── README.md           # This file
```

## 🏛️ Architecture Layers

1. **Handlers** - Handle HTTP requests/responses
2. **Services** - Contain business logic
3. **Repositories** - Data access and database operations
4. **Models** - Data structures and validation
5. **Middleware** - Cross-cutting concerns (auth, logging, CORS)

## 🔒 Security

- API key authentication for all endpoints
- CORS configuration for cloud app access
- SQL injection prevention through parameterized queries
- Structured error handling (no sensitive data exposure)

## 🧪 Testing

### Using Swagger UI (Recommended)

1. Open `http://localhost:8080/swagger/index.html`
2. Click "Authorize" and enter your API key
3. Test any endpoint interactively

### Using cURL

Run the health check to verify the setup:
```bash
curl http://localhost:8080/health
```

Test with API key:
```bash
curl -H "X-API-Key: your_api_key" http://localhost:8080/api/v1/tickets
```

### Running Tests

```bash
make test
```

## 🚀 Deployment

### Using Docker

**Build and run with Docker:**
```bash
make docker-build
make docker-run
```

**Or using Docker Compose:**
```bash
docker-compose up -d
```

**Or manually:**
```bash
docker build -t api-gateway:latest .
docker run -p 8080:8080 --env-file .env api-gateway:latest
```

### Production Checklist

1. **Set production environment**
   ```env
   GIN_MODE=release
   ```

2. **Update CORS settings**
   - Edit `middleware/cors.go`
   - Replace `AllowOrigins: []string{"*"}` with your cloud app URL

3. **Secure API keys**
   - Use strong, random API keys
   - Store in environment variables or secret management system

4. **Disable/Protect Swagger UI in production**
   - Option 1: Remove Swagger route in production
   - Option 2: Add authentication to Swagger endpoint
   - See [SWAGGER_GUIDE.md](SWAGGER_GUIDE.md) for details

5. **Enable HTTPS**
   - Use reverse proxy (nginx, Apache) or load balancer
   - Configure SSL certificates

6. **Database optimization**
   - Adjust connection pool settings in `database/database.go`
   - Monitor connection usage

7. **Logging**
   - Configure log rotation
   - Set appropriate log levels

## 📊 Monitoring

- **Health endpoint**: `/health` - Check API and database status
- **Logs**: JSON formatted logs for easy parsing
- **Metrics**: Monitor latency, status codes, and error rates

## 🛠️ Development

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

### Code Style

- Each function/struct has documentation comments
- Add Swagger annotations to all handlers
- Use structured logging with context
- Follow Go naming conventions
- Error handling at every layer

## 📝 License

This project is proprietary software for internal use.

## 👥 Support

For issues or questions, contact the development team.

---

**Built with** ❤️ **using Go and Gin Framework**
