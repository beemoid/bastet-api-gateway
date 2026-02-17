package service

import (
	"api-gateway/models"
	"api-gateway/repository"

	"github.com/sirupsen/logrus"
)

// TicketService handles business logic for ticket operations
// Acts as an intermediary between handlers and repository
type TicketService struct {
	repo   *repository.TicketRepository
	logger *logrus.Logger
}

// NewTicketService creates a new ticket service instance
func NewTicketService(repo *repository.TicketRepository, logger *logrus.Logger) *TicketService {
	return &TicketService{
		repo:   repo,
		logger: logger,
	}
}

// GetAllTickets retrieves all tickets
func (s *TicketService) GetAllTickets() ([]*models.OpenTicket, error) {
	s.logger.Info("Fetching all tickets")
	return s.repo.GetAll()
}

// GetTicketByID retrieves a ticket by ID
func (s *TicketService) GetTicketByID(id int) (*models.OpenTicket, error) {
	s.logger.Infof("Fetching ticket with ID: %d", id)
	return s.repo.GetByID(id)
}

// GetTicketByNumber retrieves a ticket by ticket number
func (s *TicketService) GetTicketByNumber(ticketNumber string) (*models.OpenTicket, error) {
	s.logger.Infof("Fetching ticket with number: %s", ticketNumber)
	return s.repo.GetByTicketNumber(ticketNumber)
}

// CreateTicket creates a new ticket
// Performs validation before creating
func (s *TicketService) CreateTicket(req *models.TicketCreateRequest) (*models.OpenTicket, error) {
	s.logger.Infof("Creating new ticket: %s", req.TicketNumber)

	// Check if ticket number already exists
	existing, _ := s.repo.GetByTicketNumber(req.TicketNumber)
	if existing != nil {
		s.logger.Warnf("Ticket number already exists: %s", req.TicketNumber)
		return nil, ErrTicketAlreadyExists
	}

	return s.repo.Create(req)
}

// UpdateTicket updates an existing ticket
func (s *TicketService) UpdateTicket(id int, req *models.TicketUpdateRequest) (*models.OpenTicket, error) {
	s.logger.Infof("Updating ticket ID: %d", id)

	// Verify ticket exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.Update(id, req)
}

// GetTicketsByStatus retrieves tickets filtered by status
func (s *TicketService) GetTicketsByStatus(status string) ([]*models.OpenTicket, error) {
	s.logger.Infof("Fetching tickets with status: %s", status)
	return s.repo.GetByStatus(status)
}

// GetTicketsByTerminal retrieves tickets for a specific terminal
func (s *TicketService) GetTicketsByTerminal(terminalID string) ([]*models.OpenTicket, error) {
	s.logger.Infof("Fetching tickets for terminal: %s", terminalID)
	return s.repo.GetByTerminalID(terminalID)
}
