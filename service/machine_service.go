package service

import (
	"api-gateway/models"
	"api-gateway/repository"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// MachineService handles business logic for machine operations
// Acts as an intermediary between handlers and repository
type MachineService struct {
	repo   *repository.MachineRepository
	logger *logrus.Logger

	// Metadata caching
	metadataCache     *models.MachineMetadataResponse
	metadataCacheMux  sync.RWMutex
	metadataLastFetch time.Time
	metadataCacheTTL  time.Duration
}

// NewMachineService creates a new machine service instance
func NewMachineService(repo *repository.MachineRepository, logger *logrus.Logger) *MachineService {
	return &MachineService{
		repo:             repo,
		logger:           logger,
		metadataCacheTTL: 1 * time.Hour, // Cache metadata for 1 hour
	}
}

// GetAllMachines retrieves machines with optional pagination
func (s *MachineService) GetAllMachines(page, pageSize int) ([]*models.ATMI, int, error) {
	s.logger.Info("Fetching all machines")
	return s.repo.GetAll(page, pageSize)
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

// GetMachinesByBranch retrieves machines for a specific store code
func (s *MachineService) GetMachinesByBranch(storeCode string) ([]*models.ATMI, error) {
	s.logger.Infof("Fetching machines for store code: %s", storeCode)
	return s.repo.GetByStoreCode(storeCode)
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

// GetMetadata retrieves machine metadata with intelligent caching
// Uses hybrid approach: queries database for actual values + adds descriptions from maps
func (s *MachineService) GetMetadata() (*models.MachineMetadataResponse, error) {
	// Check cache first
	s.metadataCacheMux.RLock()
	if s.metadataCache != nil && time.Since(s.metadataLastFetch) < s.metadataCacheTTL {
		s.logger.Info("Returning cached machine metadata")
		cached := s.metadataCache
		s.metadataCacheMux.RUnlock()
		return cached, nil
	}
	s.metadataCacheMux.RUnlock()

	// Cache miss or expired - query database
	s.logger.Info("Fetching fresh machine metadata from database")

	// Query all distinct values from database (truly adaptive)
	slms, err := s.repo.GetDistinctSLMs()
	if err != nil {
		return nil, err
	}

	flms, err := s.repo.GetDistinctFLMs()
	if err != nil {
		return nil, err
	}

	nets, err := s.repo.GetDistinctNETs()
	if err != nil {
		return nil, err
	}

	flmNames, err := s.repo.GetDistinctFLMNames()
	if err != nil {
		return nil, err
	}

	// Build response with optional descriptions
	slmInfos := make([]models.SLMInfo, 0, len(slms))
	for _, slm := range slms {
		slmInfos = append(slmInfos, models.BuildSLMInfo(slm))
	}

	flmInfos := make([]models.FLMInfo, 0, len(flms))
	for _, flm := range flms {
		flmInfos = append(flmInfos, models.BuildFLMInfo(flm))
	}

	netInfos := make([]models.NETInfo, 0, len(nets))
	for _, net := range nets {
		netInfos = append(netInfos, models.BuildNETInfo(net))
	}

	flmNameInfos := make([]models.FLMNameInfo, 0, len(flmNames))
	for _, flmName := range flmNames {
		flmNameInfos = append(flmNameInfos, models.BuildFLMNameInfo(flmName))
	}

	// Create response
	response := &models.MachineMetadataResponse{
		Success:     true,
		Message:     "Machine metadata retrieved successfully from database",
		SLMs:        slmInfos,
		FLMs:        flmInfos,
		NETs:        netInfos,
		FLMNames:    flmNameInfos,
		LastUpdated: time.Now().Format(time.RFC3339),
	}

	// Update cache
	s.metadataCacheMux.Lock()
	s.metadataCache = response
	s.metadataLastFetch = time.Now()
	s.metadataCacheMux.Unlock()

	s.logger.Infof("Cached machine metadata: %d SLMs, %d FLMs, %d NETs, %d FLM names",
		len(slms), len(flms), len(nets), len(flmNames))

	return response, nil
}

// RefreshMetadataCache forces a refresh of the metadata cache
// Useful when you know new values have been added to the database
func (s *MachineService) RefreshMetadataCache() error {
	s.metadataCacheMux.Lock()
	s.metadataCache = nil
	s.metadataCacheMux.Unlock()

	_, err := s.GetMetadata()
	return err
}
