package handlers

import (
	"api-gateway/models"
	"api-gateway/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// MachineHandler handles HTTP requests for machine operations
type MachineHandler struct {
	service *service.MachineService
	logger  *logrus.Logger
}

// NewMachineHandler creates a new machine handler instance
func NewMachineHandler(service *service.MachineService, logger *logrus.Logger) *MachineHandler {
	return &MachineHandler{
		service: service,
		logger:  logger,
	}
}

// GetAll handles GET /api/machines - retrieves all machines
// @Summary Get all machines
// @Description Retrieve all machines/terminals from the system. Supports pagination via query params.
// @Tags Machines
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number (default: all results)" minimum(1)
// @Param page_size query int false "Items per page (default: 100, max: 500)" minimum(1) maximum(500)
// @Success 200 {object} models.MachineListResponse "List of machines retrieved successfully"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /machines [get]
func (h *MachineHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "100"))

	if pageSize > 500 {
		pageSize = 500
	}
	if pageSize < 1 {
		pageSize = 100
	}

	machines, total, err := h.service.GetAllMachines(page, pageSize)
	if err != nil {
		h.logger.Errorf("Error fetching machines: %v", err)
		c.JSON(http.StatusInternalServerError, models.MachineListResponse{
			Success: false,
			Message: "Failed to fetch machines",
			Data:    nil,
			Total:   0,
		})
		return
	}

	resp := models.MachineListResponse{
		Success: true,
		Message: "Machines retrieved successfully",
		Data:    machines,
		Total:   total,
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

// GetByTerminalID handles GET /api/machines/:terminal_id - retrieves a machine by terminal ID
// @Summary Get machine by terminal ID
// @Description Retrieve a specific machine by its terminal ID
// @Tags Machines
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param terminal_id path string true "Terminal ID"
// @Success 200 {object} models.MachineResponse "Machine retrieved successfully"
// @Failure 404 {object} models.ErrorResponse "Machine not found"
// @Router /machines/{terminal_id} [get]
func (h *MachineHandler) GetByTerminalID(c *gin.Context) {
	terminalID := c.Param("terminal_id")

	machine, err := h.service.GetMachineByTerminalID(terminalID)
	if err != nil {
		h.logger.Errorf("Error fetching machine: %v", err)
		c.JSON(http.StatusNotFound, models.MachineResponse{
			Success: false,
			Message: "Machine not found",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.MachineResponse{
		Success: true,
		Message: "Machine retrieved successfully",
		Data:    machine,
	})
}

// GetByStatus handles GET /api/machines/status/:status - retrieves machines by status
// @Summary Get machines by status
// @Description Retrieve all machines with a specific operational status
// @Tags Machines
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param status path string true "Machine Status" Enums(Active, Inactive, Maintenance, Offline)
// @Success 200 {object} models.MachineListResponse "Machines retrieved successfully"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /machines/status/{status} [get]
func (h *MachineHandler) GetByStatus(c *gin.Context) {
	status := c.Param("status")

	machines, err := h.service.GetMachinesByStatus(status)
	if err != nil {
		h.logger.Errorf("Error fetching machines: %v", err)
		c.JSON(http.StatusInternalServerError, models.MachineListResponse{
			Success: false,
			Message: "Failed to fetch machines",
			Data:    nil,
			Total:   0,
		})
		return
	}

	c.JSON(http.StatusOK, models.MachineListResponse{
		Success: true,
		Message: "Machines retrieved successfully",
		Data:    machines,
		Total:   len(machines),
	})
}

// GetByBranch handles GET /api/machines/branch/:branch_code - retrieves machines by store code
// @Summary Get machines by store code
// @Description Retrieve all machines for a specific store code
// @Tags Machines
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param branch_code path string true "Store Code"
// @Success 200 {object} models.MachineListResponse "Machines retrieved successfully"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /machines/branch/{branch_code} [get]
func (h *MachineHandler) GetByBranch(c *gin.Context) {
	storeCode := c.Param("branch_code")

	machines, err := h.service.GetMachinesByBranch(storeCode)
	if err != nil {
		h.logger.Errorf("Error fetching machines: %v", err)
		c.JSON(http.StatusInternalServerError, models.MachineListResponse{
			Success: false,
			Message: "Failed to fetch machines",
			Data:    nil,
			Total:   0,
		})
		return
	}

	c.JSON(http.StatusOK, models.MachineListResponse{
		Success: true,
		Message: "Machines retrieved successfully",
		Data:    machines,
		Total:   len(machines),
	})
}

// UpdateStatus handles PATCH /api/machines/status - updates machine status
// @Summary Update machine status
// @Description Update the operational status of a machine
// @Tags Machines
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param machine body models.MachineStatusUpdate true "Machine status update data"
// @Success 200 {object} models.MachineResponse "Machine status updated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request data"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /machines/status [patch]
func (h *MachineHandler) UpdateStatus(c *gin.Context) {
	var req models.MachineStatusUpdate

	// Bind and validate JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, models.MachineResponse{
			Success: false,
			Message: "Invalid request data: " + err.Error(),
			Data:    nil,
		})
		return
	}

	machine, err := h.service.UpdateMachineStatus(&req)
	if err != nil {
		h.logger.Errorf("Error updating machine status: %v", err)
		c.JSON(http.StatusInternalServerError, models.MachineResponse{
			Success: false,
			Message: "Failed to update machine status",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.MachineResponse{
		Success: true,
		Message: "Machine status updated successfully",
		Data:    machine,
	})
}

// Search handles GET /api/machines/search - searches machines with filters
// @Summary Search machines
// @Description Search machines using multiple filter criteria
// @Tags Machines
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param status query string false "Filter by status"
// @Param store_code query string false "Filter by store code"
// @Param province query string false "Filter by province"
// @Param city_regency query string false "Filter by city/regency"
// @Param district query string false "Search by district (partial match)"
// @Success 200 {object} models.MachineListResponse "Search completed successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /machines/search [get]
func (h *MachineHandler) Search(c *gin.Context) {
	var filter models.MachineFilter

	// Bind query parameters to filter struct
	if err := c.ShouldBindQuery(&filter); err != nil {
		h.logger.Errorf("Invalid query parameters: %v", err)
		c.JSON(http.StatusBadRequest, models.MachineListResponse{
			Success: false,
			Message: "Invalid query parameters: " + err.Error(),
			Data:    nil,
			Total:   0,
		})
		return
	}

	machines, err := h.service.SearchMachines(&filter)
	if err != nil {
		h.logger.Errorf("Error searching machines: %v", err)
		c.JSON(http.StatusInternalServerError, models.MachineListResponse{
			Success: false,
			Message: "Failed to search machines",
			Data:    nil,
			Total:   0,
		})
		return
	}

	c.JSON(http.StatusOK, models.MachineListResponse{
		Success: true,
		Message: "Search completed successfully",
		Data:    machines,
		Total:   len(machines),
	})
}

// GetMetadata handles GET /api/machines/metadata - retrieves valid values for machine fields
// @Summary Get machine metadata
// @Description Retrieve all valid values for machine SLM, FLM, NET, and FLM Name fields from the database
// @Description Values are cached for 1 hour for performance. New values automatically appear after cache refresh.
// @Tags Machines
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.MachineMetadataResponse "Metadata retrieved successfully"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /machines/metadata [get]
func (h *MachineHandler) GetMetadata(c *gin.Context) {
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
