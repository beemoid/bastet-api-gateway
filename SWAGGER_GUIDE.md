# Swagger/OpenAPI Documentation Guide

This guide explains how to use and maintain the Swagger documentation for the API Gateway.

## Overview

The API Gateway uses **Swagger/OpenAPI 2.0** for interactive API documentation. Swagger provides:

- **Interactive UI** - Test API endpoints directly from the browser
- **Auto-generated docs** - Based on code annotations
- **Type definitions** - Request/response schemas
- **Authentication** - Built-in API key and token testing

## Quick Start

### 1. Install Swagger CLI Tool

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Or use the Makefile:

```bash
make install
```

### 2. Generate Swagger Documentation

```bash
swag init -g main.go --output docs
```

Or use the Makefile:

```bash
make swagger
```

This will generate the following files in the `docs/` directory:
- `docs.go` - Go package with embedded docs
- `swagger.json` - OpenAPI JSON specification
- `swagger.yaml` - OpenAPI YAML specification

### 3. Run the Application

```bash
make run
```

Or manually:

```bash
go run main.go
```

### 4. Access Swagger UI

Open your browser and navigate to:

```
http://localhost:8080/swagger/index.html
```

## Using Swagger UI

### Testing Endpoints

1. **Click on an endpoint** to expand it
2. **Click "Try it out"** button
3. **Fill in parameters** (path params, query params, request body)
4. **Add API Key** (click "Authorize" button at the top)
   - Enter your API key in the `X-API-Key` field
   - Click "Authorize"
5. **Click "Execute"** to send the request
6. **View the response** below

### Endpoint Tags

The API endpoints are organized into the following tag groups:

| Tag | Description |
|-----|-------------|
| **Tickets** | Ticket CRUD operations and metadata |
| **Machines** | Machine/terminal operations and metadata |
| **Health** | Health check and ping endpoints |
| **Admin Auth** | Admin login, logout, session management |
| **Token Management** | API token CRUD, disable/enable |
| **Analytics** | Dashboard stats, token analytics, endpoint stats, daily usage |

### Example: Creating a Ticket

1. Navigate to `POST /api/v1/tickets`
2. Click "Try it out"
3. Click "Authorize" and enter your API token
4. Edit the request body:
   ```json
   {
     "terminal_id": "ATM-001",
     "terminal_name": "Main Branch ATM",
     "priority": "1.High",
     "mode": "Off-line",
     "initial_problem": "Card reader malfunction",
     "status": "0.NEW",
     "tickets_no": "TKT-2024-001"
   }
   ```
5. Click "Execute"
6. Check the response

### Example: Managing API Tokens

1. Navigate to `POST /api/v1/admin/auth/login` under **Admin Auth**
2. Login with admin credentials to get a session token
3. Use the session token in `X-Session-Token` header
4. Navigate to `POST /api/v1/admin/tokens` under **Token Management**
5. Create a new token with scopes and rate limits

## Adding Swagger Annotations

### Handler Function Example

```go
// GetAll handles GET /api/v1/tickets - retrieves all tickets
// @Summary Get all tickets
// @Description Retrieve all tickets from the system
// @Tags Tickets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.TicketListResponse "List of tickets retrieved successfully"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tickets [get]
func (h *TicketHandler) GetAll(c *gin.Context) {
    // Implementation
}
```

### Token Management Handler Example

```go
// CreateToken handles POST /api/v1/admin/tokens
// @Summary Create API Token
// @Description Create a new API token with scopes and rate limits
// @Tags Token Management
// @Accept json
// @Produce json
// @Param token body models.CreateTokenRequest true "Token Details"
// @Success 201 {object} models.CreateTokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/tokens [post]
func (h *TokenHandler) CreateToken(c *gin.Context) {
    // Implementation
}
```

### Analytics Handler Example

```go
// GetDashboardStats handles GET /api/v1/admin/analytics/dashboard
// @Summary Get Dashboard Stats
// @Description Get overview statistics for the admin dashboard
// @Tags Analytics
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/analytics/dashboard [get]
func (h *TokenHandler) GetDashboardStats(c *gin.Context) {
    // Implementation
}
```

### Annotation Breakdown

- `@Summary` - Short description (appears in the endpoint list)
- `@Description` - Detailed description
- `@Tags` - Group endpoints by category (Tickets, Machines, Health, Admin Auth, Token Management, Analytics)
- `@Accept` - Request content type
- `@Produce` - Response content type
- `@Security` - Security scheme (ApiKeyAuth)
- `@Param` - Parameter definition (path, query, body)
- `@Success` - Successful response definition
- `@Failure` - Error response definition
- `@Router` - Route path and HTTP method

### Parameter Annotations

**Path Parameter:**
```go
// @Param id path int true "Token ID"
```

**Query Parameter:**
```go
// @Param status query string false "Filter by status"
// @Param days query int false "Number of days (default 7)"
// @Param limit query int false "Limit results (default 20)"
```

**Body Parameter:**
```go
// @Param token body models.CreateTokenRequest true "Token creation data"
// @Param login body models.LoginRequest true "Login credentials"
```

**Enum Parameter:**
```go
// @Param environment path string true "Environment" Enums(production, staging, development, test)
```

### Response Annotations

```go
// @Success 200 {object} models.TicketResponse "Ticket retrieved successfully"
// @Success 201 {object} models.CreateTokenResponse "Token created successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Not found"
// @Failure 429 {object} models.ErrorResponse "Rate limit exceeded"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
```

## Main Application Annotations

In `main.go`, the general API information is configured:

```go
// @title API Gateway for On-Premise to Cloud Communication
// @version 1.0
// @description This API Gateway serves as middleware between on-premise databases and cloud applications, providing RESTful APIs for ticket and machine management.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Proprietary
// @license.url http://www.example.com/license

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @description API key for authentication. Required for all endpoints except /health and /ping.
```

## Common Tasks

### Regenerate Documentation After Code Changes

Every time you modify handler annotations:

```bash
make swagger
```

### Change API Base URL

Edit `main.go`:

```go
// @host your-domain.com:8080
```

Then regenerate:

```bash
make swagger
```

### Add New Endpoint

1. Add handler function with annotations
2. Register route in `routes/routes.go`
3. Regenerate swagger docs: `make swagger`
4. Restart application

### Model Documentation

Swagger automatically reads struct tags from your models:

```go
type CreateTokenRequest struct {
    Name               string   `json:"name" binding:"required,min=3,max=200" example:"Production App"`
    Description        string   `json:"description" example:"API token for production"`
    Environment        string   `json:"environment" binding:"required,oneof=production staging development test" example:"production"`
    Scopes             []string `json:"scopes" example:"tickets:read,tickets:write"`
    RateLimitPerMinute int      `json:"rate_limit_per_minute" example:"60"`
}
```

### Nullable Type Documentation

For nullable fields, use the `swaggertype` tag to override the type:

```go
type OpenTicket struct {
    Priority   NullString `json:"priority" db:"Priority" swaggertype:"string" example:"1.High"`
    ExpiresAt  NullTime   `json:"expires_at" db:"expires_at" swaggertype:"string" example:"2024-12-31T23:59:59Z"`
}
```

## Customization

### Change Swagger UI Theme

Modify the Swagger endpoint in `routes/routes.go`:

```go
// Dark theme
router.GET("/swagger/*any", ginSwagger.WrapHandler(
    swaggerFiles.Handler,
    ginSwagger.DefaultModelsExpandDepth(-1),
    ginSwagger.DocExpansion("none"),
    ginSwagger.DarkMode(),
))
```

### Custom Swagger Configuration

```go
router.GET("/swagger/*any", ginSwagger.WrapHandler(
    swaggerFiles.Handler,
    ginSwagger.DefaultModelsExpandDepth(-1), // Hide models section
    ginSwagger.DocExpansion("list"),         // Expand endpoints list
    ginSwagger.DeepLinking(true),           // Enable deep linking
    ginSwagger.PersistAuthorization(true),  // Remember authorization
))
```

## Deployment

### Include Swagger in Production

The `docs/` directory is automatically included in the Docker build.

To disable Swagger in production, use build tags:

```bash
go build -tags=!swagger -o api-gateway
```

Or conditionally register the route:

```go
if os.Getenv("GIN_MODE") != "release" {
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
```

## Troubleshooting

### Swagger Not Loading

1. Check if `docs/` directory exists
2. Verify `import _ "api-gateway/docs"` in main.go
3. Regenerate docs: `make swagger`
4. Restart application

### Models Not Showing

1. Ensure models are referenced in handler annotations
2. Models must be in a package that's imported
3. Check `swaggertype` tags for custom types (NullString, NullTime)
4. Regenerate docs: `make swagger`

### Authentication Not Working

1. Click "Authorize" button in Swagger UI
2. Enter API key in the value field
3. Click "Authorize"
4. Try the request again

**Note:** Admin endpoints use session-based auth (`X-Session-Token`), not API key auth. Login via `POST /api/v1/admin/auth/login` first, then use the session token.

### Changes Not Reflecting

1. Stop the application
2. Run `make swagger`
3. Start the application again
4. Hard refresh browser (Ctrl+F5 or Cmd+Shift+R)

### Rate Limit Errors in Testing

If you get `429 Too Many Requests` while testing:
1. Check the token's rate limits in the admin dashboard
2. Increase the rate limits or wait for the time window to reset
3. Use a different token for testing

## Resources

- **Swag GitHub**: https://github.com/swaggo/swag
- **Swag Annotations**: https://github.com/swaggo/swag#declarative-comments-format
- **OpenAPI Specification**: https://swagger.io/specification/v2/
- **Gin Swagger**: https://github.com/swaggo/gin-swagger

## Best Practices

1. **Always regenerate** docs after changing annotations
2. **Use descriptive summaries** for each endpoint
3. **Document all parameters** with clear descriptions
4. **Include examples** in model definitions using `example` tags
5. **Group endpoints** logically using tags (Tickets, Machines, Token Management, Analytics, etc.)
6. **Document error responses** for all possible cases (400, 401, 404, 429, 500, 503)
7. **Keep descriptions concise** but informative
8. **Use enums** for fields with limited values (environment, role, status)
9. **Version your API** properly in the main annotations
10. **Test endpoints** using Swagger UI before releasing
11. **Use `swaggertype` tag** for custom types (NullString, NullTime) so they render correctly

## Security Notes

- Swagger UI exposes your API structure
- Consider disabling in production or protecting with authentication
- Never commit API keys or tokens in annotations
- Use environment variables for sensitive data
- Review Swagger JSON before deploying
- Admin endpoints are protected by session authentication, not shown in Swagger's "Authorize" dialog

---

**Happy Documenting!**
