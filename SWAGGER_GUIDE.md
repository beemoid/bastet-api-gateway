# Swagger/OpenAPI Documentation Guide

This guide explains how to use and maintain the Swagger documentation for the API Gateway.

## 📚 Overview

The API Gateway uses **Swagger/OpenAPI 2.0** for interactive API documentation. Swagger provides:

- **Interactive UI** - Test API endpoints directly from the browser
- **Auto-generated docs** - Based on code annotations
- **Type definitions** - Request/response schemas
- **Authentication** - Built-in API key testing

## 🚀 Quick Start

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

## 🎯 Using Swagger UI

### Testing Endpoints

1. **Click on an endpoint** to expand it
2. **Click "Try it out"** button
3. **Fill in parameters** (path params, query params, request body)
4. **Add API Key** (click "Authorize" button at the top)
   - Enter your API key in the `X-API-Key` field
   - Click "Authorize"
5. **Click "Execute"** to send the request
6. **View the response** below

### Example: Creating a Ticket

1. Navigate to `POST /api/v1/tickets`
2. Click "Try it out"
3. Click "Authorize" and enter your API key
4. Edit the request body:
   ```json
   {
     "ticket_number": "TKT-2024-001",
     "terminal_id": "ATM-001",
     "description": "Card reader malfunction",
     "priority": "High",
     "category": "Hardware",
     "reported_by": "John Doe",
     "assigned_to": "Tech Team A"
   }
   ```
5. Click "Execute"
6. Check the response

## 📝 Adding Swagger Annotations

### Handler Function Example

```go
// GetAll handles GET /api/tickets - retrieves all tickets
// @Summary Get all tickets
// @Description Retrieve all tickets from the system
// @Tags Tickets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.TicketListResponse "List of tickets retrieved successfully"
// @Failure 500 {object} models.TicketListResponse "Internal server error"
// @Router /tickets [get]
func (h *TicketHandler) GetAll(c *gin.Context) {
    // Implementation
}
```

### Annotation Breakdown

- `@Summary` - Short description (appears in the endpoint list)
- `@Description` - Detailed description
- `@Tags` - Group endpoints by category (Tickets, Machines, Health)
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
// @Param id path int true "Ticket ID"
```

**Query Parameter:**
```go
// @Param status query string false "Filter by status"
```

**Body Parameter:**
```go
// @Param ticket body models.TicketCreateRequest true "Ticket creation data"
```

**Enum Parameter:**
```go
// @Param status path string true "Ticket Status" Enums(Open, InProgress, Pending, Resolved)
```

### Response Annotations

```go
// @Success 200 {object} models.TicketResponse "Ticket retrieved successfully"
// @Failure 404 {object} models.TicketResponse "Ticket not found"
// @Failure 500 {object} models.TicketResponse "Internal server error"
```

## 🏗️ Main Application Annotations

In `main.go`, add general API information:

```go
// @title API Gateway for On-Premise to Cloud Communication
// @version 1.0
// @description This API Gateway serves as middleware between on-premise databases and cloud applications
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
// @description API key for authentication
```

## 🔧 Common Tasks

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
type TicketCreateRequest struct {
    TicketNumber string `json:"ticket_number" binding:"required" example:"TKT-2024-001"`
    TerminalID   string `json:"terminal_id" binding:"required" example:"ATM-001"`
    Description  string `json:"description" binding:"required" example:"Card reader malfunction"`
    Priority     string `json:"priority" binding:"required" example:"High"`
}
```

## 🎨 Customization

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

## 📦 Deployment

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

## 🔍 Troubleshooting

### Swagger Not Loading

1. Check if `docs/` directory exists
2. Verify `import _ "api-gateway/docs"` in main.go
3. Regenerate docs: `make swagger`
4. Restart application

### Models Not Showing

1. Ensure models are referenced in handler annotations
2. Models must be in a package that's imported
3. Regenerate docs: `make swagger`

### Authentication Not Working

1. Click "Authorize" button in Swagger UI
2. Enter API key in the value field
3. Click "Authorize"
4. Try the request again

### Changes Not Reflecting

1. Stop the application
2. Run `make swagger`
3. Start the application again
4. Hard refresh browser (Ctrl+F5 or Cmd+Shift+R)

## 📚 Resources

- **Swag GitHub**: https://github.com/swaggo/swag
- **Swag Annotations**: https://github.com/swaggo/swag#declarative-comments-format
- **OpenAPI Specification**: https://swagger.io/specification/v2/
- **Gin Swagger**: https://github.com/swaggo/gin-swagger

## 🎯 Best Practices

1. **Always regenerate** docs after changing annotations
2. **Use descriptive summaries** for each endpoint
3. **Document all parameters** with clear descriptions
4. **Include examples** in model definitions
5. **Group endpoints** logically using tags
6. **Document error responses** for all possible cases
7. **Keep descriptions concise** but informative
8. **Use enums** for fields with limited values
9. **Version your API** properly in the main annotations
10. **Test endpoints** using Swagger UI before releasing

## 🔐 Security Notes

- Swagger UI exposes your API structure
- Consider disabling in production or protecting with authentication
- Never commit API keys in annotations
- Use environment variables for sensitive data
- Review Swagger JSON before deploying

---

**Happy Documenting!** 📖✨
