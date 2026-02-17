package service

import (
	"api-gateway/models"
	"api-gateway/repository"

	"github.com/sirupsen/logrus"
)

// MachineService handles business logic for machine operations
// Acts as an intermediary between handlers and repository
type MachineService struct {
	repo   *repository.MachineRepository
	logger *logrus.Logger
}

// NewMachineService creates a new machine service instance
func NewMachineService(repo *repository.MachineRepository, logger *logrus.Logger) *MachineService {
	return &MachineService{
		repo:   repo,
		logger: logger,
	}
}

// GetAllMachines retrieves all machines
func (s *MachineService) GetAllMachines() ([]*models.ATMI, error) {
	s.logger.Info("Fetching all machines")
	return s.repo.GetAll()
}

// GetMachineByTerminalID retrieves a machine by terminal ID
func (s *MachineService) GetMachineByTerminalID(terminalID string) (*models.ATMI, error) {
	s.logger.Infof("Fetching machine with terminal ID: %s", terminalID)
	return s.repo.GetByTerminalID(terminalID)
}

// GetMachinesByStatus retrieves machines filtered by status
func (s *MachineService) GetMachinesByStatus(status string) ([]*models.ATMI, error) {
	s.logger.Infof("Fetching machines with status: %s", status)
	return s.repo.GetByStatus(status)
}

// GetMachinesByBranch retrieves machines for a specific branch
func (s *MachineService) GetMachinesByBranch(branchCode string) ([]*models.ATMI, error) {
	s.logger.Infof("Fetching machines for branch: %s", branchCode)
	return s.repo.GetByBranchCode(branchCode)
}

// UpdateMachineStatus updates the status of a machine
func (s *MachineService) UpdateMachineStatus(req *models.MachineStatusUpdate) (*models.ATMI, error) {
	s.logger.Infof("Updating status for terminal: %s", req.TerminalID)

	// Verify machine exists
	_, err := s.repo.GetByTerminalID(req.TerminalID)
	if err != nil {
		return nil, err
	}

	return s.repo.UpdateStatus(req)
}

// SearchMachines performs a flexible search based on filters
func (s *MachineService) SearchMachines(filter *models.MachineFilter) ([]*models.ATMI, error) {
	s.logger.Info("Searching machines with filters")
	return s.repo.Search(filter)
}
