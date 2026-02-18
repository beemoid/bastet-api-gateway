package service

import (
	"api-gateway/models"
	"api-gateway/repository"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// TicketService handles business logic for ticket operations
// Acts as an intermediary between handlers and repository
type TicketService struct {
	repo   *repository.TicketRepository
	logger *logrus.Logger

	// Metadata caching
	metadataCache     *models.MetadataResponse
	metadataCacheMux  sync.RWMutex
	metadataLastFetch time.Time
	metadataCacheTTL  time.Duration
}

// NewTicketService creates a new ticket service instance
func NewTicketService(repo *repository.TicketRepository, logger *logrus.Logger) *TicketService {
	return &TicketService{
		repo:             repo,
		logger:           logger,
		metadataCacheTTL: 1 * time.Hour, // Cache metadata for 1 hour
	}
}

// GetAllTickets retrieves tickets with optional pagination
func (s *TicketService) GetAllTickets(page, pageSize int) ([]*models.OpenTicket, int, error) {
	s.logger.Info("Fetching all tickets")
	return s.repo.GetAll(page, pageSize)
}

// GetTicketByID retrieves a ticket by terminal ID
func (s *TicketService) GetTicketByID(terminalID string) (*models.OpenTicket, error) {
	s.logger.Infof("Fetching ticket with terminal ID: %s", terminalID)
	return s.repo.GetByTerminalID(terminalID)
}

// GetTicketByNumber retrieves a ticket by ticket number
func (s *TicketService) GetTicketByNumber(ticketNumber string) (*models.OpenTicket, error) {
	s.logger.Infof("Fetching ticket with number: %s", ticketNumber)
	return s.repo.GetByTicketNumber(ticketNumber)
}

// CreateTicket creates a new ticket
// Performs validation before creating
func (s *TicketService) CreateTicket(req *models.TicketCreateRequest) (*models.OpenTicket, error) {
	s.logger.Infof("Creating new ticket: %s", req.TicketsNo)

	// Check if ticket number already exists
	existing, _ := s.repo.GetByTicketNumber(req.TicketsNo)
	if existing != nil {
		s.logger.Warnf("Ticket number already exists: %s", req.TicketsNo)
		return nil, ErrTicketAlreadyExists
	}

	return s.repo.Create(req)
}

// UpdateTicket updates an existing ticket
func (s *TicketService) UpdateTicket(terminalID string, req *models.TicketUpdateRequest) (*models.OpenTicket, error) {
	s.logger.Infof("Updating ticket for terminal ID: %s", terminalID)

	// Verify ticket exists
	_, err := s.repo.GetByTerminalID(terminalID)
	if err != nil {
		return nil, err
	}

	return s.repo.Update(terminalID, req)
}

// GetTicketsByStatus retrieves tickets filtered by status
func (s *TicketService) GetTicketsByStatus(status string) ([]*models.OpenTicket, error) {
	s.logger.Infof("Fetching tickets with status: %s", status)
	return s.repo.GetByStatus(status)
}

// GetTicketsByTerminal retrieves tickets for a specific terminal
func (s *TicketService) GetTicketsByTerminal(terminalID string) ([]*models.OpenTicket, error) {
	s.logger.Infof("Fetching tickets for terminal: %s", terminalID)

	// Get single ticket by terminal ID
	ticket, err := s.repo.GetByTerminalID(terminalID)
	if err != nil {
		return nil, err
	}

	// Return as array
	return []*models.OpenTicket{ticket}, nil
}

// GetMetadata retrieves ticket metadata with intelligent caching
// Uses hybrid approach: queries database for actual values + adds descriptions from maps
func (s *TicketService) GetMetadata() (*models.MetadataResponse, error) {
	// Check cache first
	s.metadataCacheMux.RLock()
	if s.metadataCache != nil && time.Since(s.metadataLastFetch) < s.metadataCacheTTL {
		s.logger.Info("Returning cached ticket metadata")
		cached := s.metadataCache
		s.metadataCacheMux.RUnlock()
		return cached, nil
	}
	s.metadataCacheMux.RUnlock()

	// Cache miss or expired - query database
	s.logger.Info("Fetching fresh ticket metadata from database")

	// Query all distinct values from database (truly adaptive)
	statuses, err := s.repo.GetDistinctStatuses()
	if err != nil {
		return nil, err
	}

	modes, err := s.repo.GetDistinctModes()
	if err != nil {
		return nil, err
	}

	priorities, err := s.repo.GetDistinctPriorities()
	if err != nil {
		return nil, err
	}

	// Build response with optional descriptions
	statusInfos := make([]models.StatusInfo, 0, len(statuses))
	for _, status := range statuses {
		statusInfos = append(statusInfos, models.BuildStatusInfo(status))
	}

	modeInfos := make([]models.ModeInfo, 0, len(modes))
	for _, mode := range modes {
		modeInfos = append(modeInfos, models.BuildModeInfo(mode))
	}

	priorityInfos := make([]models.PriorityInfo, 0, len(priorities))
	for _, priority := range priorities {
		priorityInfos = append(priorityInfos, models.BuildPriorityInfo(priority))
	}

	// Create response
	response := &models.MetadataResponse{
		Success:     true,
		Message:     "Metadata retrieved successfully from database",
		Statuses:    statusInfos,
		Modes:       modeInfos,
		Priorities:  priorityInfos,
		LastUpdated: time.Now().Format(time.RFC3339),
	}

	// Update cache
	s.metadataCacheMux.Lock()
	s.metadataCache = response
	s.metadataLastFetch = time.Now()
	s.metadataCacheMux.Unlock()

	s.logger.Infof("Cached ticket metadata: %d statuses, %d modes, %d priorities", 
		len(statuses), len(modes), len(priorities))

	return response, nil
}

// RefreshMetadataCache forces a refresh of the metadata cache
// Useful when you know new values have been added to the database
func (s *TicketService) RefreshMetadataCache() error {
	s.metadataCacheMux.Lock()
	s.metadataCache = nil
	s.metadataCacheMux.Unlock()

	_, err := s.GetMetadata()
	return err
}
