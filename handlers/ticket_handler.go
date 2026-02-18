package handlers

import (
	"api-gateway/models"
	"api-gateway/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// TicketHandler handles HTTP requests for ticket operations
type TicketHandler struct {
	service *service.TicketService
	logger  *logrus.Logger
}

// NewTicketHandler creates a new ticket handler instance
func NewTicketHandler(service *service.TicketService, logger *logrus.Logger) *TicketHandler {
	return &TicketHandler{
		service: service,
		logger:  logger,
	}
}

// GetAll handles GET /api/tickets - retrieves all tickets
// @Summary Get all tickets
// @Description Retrieve all tickets from the system. Supports pagination via query params.
// @Tags Tickets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number (default: all results)" minimum(1)
// @Param page_size query int false "Items per page (default: 100, max: 500)" minimum(1) maximum(500)
// @Success 200 {object} models.TicketListResponse "List of tickets retrieved successfully"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tickets [get]
func (h *TicketHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "100"))

	if pageSize > 500 {
		pageSize = 500
	}
	if pageSize < 1 {
		pageSize = 100
	}

	tickets, total, err := h.service.GetAllTickets(page, pageSize)
	if err != nil {
		h.logger.Errorf("Error fetching tickets: %v", err)
		c.JSON(http.StatusInternalServerError, models.TicketListResponse{
			Success: false,
			Message: "Failed to fetch tickets",
			Data:    nil,
			Total:   0,
		})
		return
	}

	resp := models.TicketListResponse{
		Success: true,
		Message: "Tickets retrieved successfully",
		Data:    tickets,
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

// GetByID handles GET /api/tickets/:id - retrieves a ticket by terminal ID
// @Summary Get ticket by terminal ID
// @Description Retrieve a specific ticket by its terminal ID
// @Tags Tickets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Terminal ID"
// @Success 200 {object} models.TicketResponse "Ticket retrieved successfully"
// @Failure 404 {object} models.ErrorResponse "Ticket not found"
// @Router /tickets/{id} [get]
func (h *TicketHandler) GetByID(c *gin.Context) {
	// Parse terminal ID from URL parameter
	terminalID := c.Param("id")

	ticket, err := h.service.GetTicketByID(terminalID)
	if err != nil {
		h.logger.Errorf("Error fetching ticket: %v", err)
		c.JSON(http.StatusNotFound, models.TicketResponse{
			Success: false,
			Message: "Ticket not found",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.TicketResponse{
		Success: true,
		Message: "Ticket retrieved successfully",
		Data:    ticket,
	})
}

// GetByNumber handles GET /api/tickets/number/:number - retrieves a ticket by ticket number
// @Summary Get ticket by number
// @Description Retrieve a ticket by its unique ticket number
// @Tags Tickets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param number path string true "Ticket Number"
// @Success 200 {object} models.TicketResponse "Ticket retrieved successfully"
// @Failure 404 {object} models.ErrorResponse "Ticket not found"
// @Router /tickets/number/{number} [get]
func (h *TicketHandler) GetByNumber(c *gin.Context) {
	ticketNumber := c.Param("number")

	ticket, err := h.service.GetTicketByNumber(ticketNumber)
	if err != nil {
		h.logger.Errorf("Error fetching ticket: %v", err)
		c.JSON(http.StatusNotFound, models.TicketResponse{
			Success: false,
			Message: "Ticket not found",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.TicketResponse{
		Success: true,
		Message: "Ticket retrieved successfully",
		Data:    ticket,
	})
}

// Create handles POST /api/tickets - creates a new ticket
// @Summary Create a new ticket
// @Description Create a new ticket in the system
// @Tags Tickets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param ticket body models.TicketCreateRequest true "Ticket creation data"
// @Success 201 {object} models.TicketResponse "Ticket created successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request data"
// @Failure 409 {object} models.ErrorResponse "Ticket number already exists"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tickets [post]
func (h *TicketHandler) Create(c *gin.Context) {
	var req models.TicketCreateRequest

	// Bind and validate JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, models.TicketResponse{
			Success: false,
			Message: "Invalid request data: " + err.Error(),
			Data:    nil,
		})
		return
	}

	ticket, err := h.service.CreateTicket(&req)
	if err != nil {
		h.logger.Errorf("Error creating ticket: %v", err)

		// Check for duplicate ticket error
		if err == service.ErrTicketAlreadyExists {
			c.JSON(http.StatusConflict, models.TicketResponse{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.TicketResponse{
			Success: false,
			Message: "Failed to create ticket",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, models.TicketResponse{
		Success: true,
		Message: "Ticket created successfully",
		Data:    ticket,
	})
}

// Update handles PUT /api/tickets/:id - updates an existing ticket
// @Summary Update a ticket
// @Description Update an existing ticket by terminal ID
// @Tags Tickets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Terminal ID"
// @Param ticket body models.TicketUpdateRequest true "Ticket update data"
// @Success 200 {object} models.TicketResponse "Ticket updated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request data"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tickets/{id} [put]
func (h *TicketHandler) Update(c *gin.Context) {
	// Parse terminal ID from URL parameter
	terminalID := c.Param("id")

	var req models.TicketUpdateRequest

	// Bind and validate JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, models.TicketResponse{
			Success: false,
			Message: "Invalid request data: " + err.Error(),
			Data:    nil,
		})
		return
	}

	ticket, err := h.service.UpdateTicket(terminalID, &req)
	if err != nil {
		h.logger.Errorf("Error updating ticket: %v", err)
		c.JSON(http.StatusInternalServerError, models.TicketResponse{
			Success: false,
			Message: "Failed to update ticket",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.TicketResponse{
		Success: true,
		Message: "Ticket updated successfully",
		Data:    ticket,
	})
}

// GetByStatus handles GET /api/tickets/status/:status - retrieves tickets by status
// @Summary Get tickets by status
// @Description Retrieve all tickets with a specific status
// @Tags Tickets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param status path string true "Ticket Status" Enums(0.NEW, 1.Req FD ke HD, 2.Kirim FLM, 21.Req Replenish, 3.SLM, 4.SLM-Net, 5.Menunggu Update, 6.Follow-up Sales team, 8.Wait transaction)
// @Success 200 {object} models.TicketListResponse "Tickets retrieved successfully"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tickets/status/{status} [get]
func (h *TicketHandler) GetByStatus(c *gin.Context) {
	status := c.Param("status")

	tickets, err := h.service.GetTicketsByStatus(status)
	if err != nil {
		h.logger.Errorf("Error fetching tickets: %v", err)
		c.JSON(http.StatusInternalServerError, models.TicketListResponse{
			Success: false,
			Message: "Failed to fetch tickets",
			Data:    nil,
			Total:   0,
		})
		return
	}

	c.JSON(http.StatusOK, models.TicketListResponse{
		Success: true,
		Message: "Tickets retrieved successfully",
		Data:    tickets,
		Total:   len(tickets),
	})
}

// GetByTerminal handles GET /api/tickets/terminal/:terminal_id - retrieves tickets by terminal
// @Summary Get tickets by terminal
// @Description Retrieve all tickets associated with a specific terminal
// @Tags Tickets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param terminal_id path string true "Terminal ID"
// @Success 200 {object} models.TicketListResponse "Tickets retrieved successfully"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tickets/terminal/{terminal_id} [get]
func (h *TicketHandler) GetByTerminal(c *gin.Context) {
	terminalID := c.Param("terminal_id")

	tickets, err := h.service.GetTicketsByTerminal(terminalID)
	if err != nil {
		h.logger.Errorf("Error fetching tickets: %v", err)
		c.JSON(http.StatusInternalServerError, models.TicketListResponse{
			Success: false,
			Message: "Failed to fetch tickets",
			Data:    nil,
			Total:   0,
		})
		return
	}

	c.JSON(http.StatusOK, models.TicketListResponse{
		Success: true,
		Message: "Tickets retrieved successfully",
		Data:    tickets,
		Total:   len(tickets),
	})
}

// GetMetadata handles GET /api/tickets/metadata - retrieves valid values for ticket fields
// @Summary Get ticket metadata
// @Description Retrieve all valid values for ticket status, mode, and priority fields from the database
// @Description Values are cached for 1 hour for performance. New values automatically appear after cache refresh.
// @Tags Tickets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.MetadataResponse "Metadata retrieved successfully"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tickets/metadata [get]
func (h *TicketHandler) GetMetadata(c *gin.Context) {
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
