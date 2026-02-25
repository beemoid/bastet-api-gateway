package handlers

import (
	"api-gateway/models"
	"api-gateway/repository"
	"api-gateway/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// DataHandler handles HTTP requests for the unified /api/v1/data endpoint.
type DataHandler struct {
	service *service.DataService
	logger  *logrus.Logger
}

// NewDataHandler creates a new DataHandler instance.
func NewDataHandler(service *service.DataService, logger *logrus.Logger) *DataHandler {
	return &DataHandler{
		service: service,
		logger:  logger,
	}
}

// vendorFilterFromContext extracts the vendor filter set by TokenAuthMiddleware.
// Returns a super-token filter for admin/internal tokens, a scoped filter for vendor
// tokens, or nil for unrestricted (legacy) tokens.
func vendorFilterFromContext(c *gin.Context) *repository.VendorFilter {
	isSuper, _ := c.Get("token_is_super")
	if isSuperBool, ok := isSuper.(bool); ok && isSuperBool {
		return &repository.VendorFilter{IsSuperToken: true}
	}

	filterColumn, _ := c.Get("token_filter_column")
	filterValue, _ := c.Get("token_filter_value")

	col, _ := filterColumn.(string)
	val, _ := filterValue.(string)

	if col == "" || val == "" {
		return nil
	}
	return repository.ResolveVendorFilter(col, val, false)
}

// GetAll handles GET /api/v1/data
// @Summary Get all data
// @Description Retrieve joined ticket+machine rows with pagination, sorting, and filtering. Vendor-scoped tokens only see rows matching their filter. Admin/Internal tokens see all rows.
// @Tags Data
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number (default: all results)" minimum(1)
// @Param page_size query int false "Items per page (default: 100, max: 500)" minimum(1) maximum(500)
// @Param sort_by query string false "Sort field: terminal_id, terminal_name, priority, mode, status, incident_start_datetime, count, balance, tickets_duration, open_time, close_time, flm_name, flm, slm, net"
// @Param sort_order query string false "Sort direction: asc or desc (default: desc)"
// @Param search query string false "Search by terminal_id or terminal_name (partial match)"
// @Param status query string false "Filter by exact status value (e.g. 0.NEW)"
// @Param mode query string false "Filter by exact mode value (e.g. Off-line)"
// @Param priority query string false "Filter by exact priority value (e.g. 1.High)"
// @Success 200 {object} models.DataListResponse "Data retrieved successfully"
// @Failure 401 {object} models.ErrorResponse "Missing or invalid API token"
// @Failure 429 {object} models.ErrorResponse "Rate limit exceeded"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /data [get]
func (h *DataHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "100"))

	if pageSize > 500 {
		pageSize = 500
	}
	if pageSize < 1 {
		pageSize = 100
	}

	sortBy := c.DefaultQuery("sort_by", "incident_start_datetime")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	params := repository.QueryParams{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Search:    strings.TrimSpace(c.Query("search")),
		Status:    strings.TrimSpace(c.Query("status")),
		Mode:      strings.TrimSpace(c.Query("mode")),
		Priority:  strings.TrimSpace(c.Query("priority")),
	}

	filter := vendorFilterFromContext(c)
	rows, total, err := h.service.GetAll(filter, params)
	if err != nil {
		h.logger.Errorf("Error fetching data: %v", err)
		c.JSON(http.StatusInternalServerError, models.DataListResponse{
			Success: false,
			Message: "Failed to fetch data",
		})
		return
	}

	resp := models.DataListResponse{
		Success:   true,
		Message:   "Data retrieved successfully",
		Data:      rows,
		Total:     total,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Search:    params.Search,
		Status:    params.Status,
		Mode:      params.Mode,
		Priority:  params.Priority,
	}

	if page > 0 {
		resp.Page = page
		resp.PageSize = pageSize
		totalPages := total / pageSize
		if total%pageSize > 0 {
			totalPages++
		}
		resp.TotalPages = totalPages
	}

	c.JSON(http.StatusOK, resp)
}

// GetByID handles GET /api/v1/data/:terminal_id
// @Summary Get data by terminal ID
// @Description Retrieve a single joined row by terminal ID. Vendor tokens return 404 if the terminal is outside their scope.
// @Tags Data
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param terminal_id path string true "Terminal ID"
// @Success 200 {object} models.DataResponse "Data retrieved successfully"
// @Failure 404 {object} models.ErrorResponse "Not found"
// @Router /data/{terminal_id} [get]
func (h *DataHandler) GetByID(c *gin.Context) {
	terminalID := c.Param("terminal_id")
	filter := vendorFilterFromContext(c)

	row, err := h.service.GetByTerminalID(terminalID, filter)
	if err != nil {
		h.logger.Errorf("Error fetching data row: %v", err)
		c.JSON(http.StatusNotFound, models.DataResponse{
			Success: false,
			Message: "Not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{
		Success: true,
		Message: "Data retrieved successfully",
		Data:    row,
	})
}

// Update handles PUT /api/v1/data/:terminal_id
// @Summary Update ticket fields
// @Description Update ticket fields for a terminal. Vendor tokens can only update terminals within their scope (returns 403 otherwise). Admin/Internal tokens can update any terminal.
// @Tags Data
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param terminal_id path string true "Terminal ID"
// @Param body body models.DataUpdateRequest true "Fields to update"
// @Success 200 {object} models.DataResponse "Updated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 403 {object} models.ErrorResponse "Outside vendor scope"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /data/{terminal_id} [put]
func (h *DataHandler) Update(c *gin.Context) {
	terminalID := c.Param("terminal_id")
	filter := vendorFilterFromContext(c)

	var req models.DataUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, models.DataResponse{
			Success: false,
			Message: "Invalid request data: " + err.Error(),
		})
		return
	}

	row, err := h.service.Update(terminalID, &req, filter)
	if err != nil {
		h.logger.Errorf("Error updating data row: %v", err)
		statusCode := http.StatusInternalServerError
		msg := "Failed to update"
		errMsg := err.Error()
		if errMsg == "not found or not accessible for this vendor" {
			statusCode = http.StatusForbidden
			msg = errMsg
		} else if errMsg == "not found" {
			statusCode = http.StatusNotFound
			msg = errMsg
		} else if errMsg == "no fields to update" {
			statusCode = http.StatusBadRequest
			msg = errMsg
		}
		c.JSON(statusCode, models.DataResponse{
			Success: false,
			Message: msg,
		})
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{
		Success: true,
		Message: "Updated successfully",
		Data:    row,
	})
}

// GetMetadata handles GET /api/v1/data/metadata
// @Summary Get field metadata
// @Description Retrieve all valid values for status, mode, and priority fields. Cached for 1 hour.
// @Tags Data
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.MetadataResponse "Metadata retrieved successfully"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /data/metadata [get]
func (h *DataHandler) GetMetadata(c *gin.Context) {
	metadata, err := h.service.GetMetadata()
	if err != nil {
		h.logger.Errorf("Error fetching metadata: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to fetch metadata",
			Error:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, metadata)
}
