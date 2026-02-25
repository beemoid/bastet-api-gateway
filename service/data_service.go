package service

import (
	"api-gateway/models"
	"api-gateway/repository"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// DataService handles business logic for the unified /api/v1/data endpoint.
type DataService struct {
	repo   *repository.DataRepository
	logger *logrus.Logger

	// Metadata caching
	metadataCache     *models.MetadataResponse
	metadataCacheMux  sync.RWMutex
	metadataLastFetch time.Time
	metadataCacheTTL  time.Duration
}

// NewDataService creates a new DataService instance.
func NewDataService(repo *repository.DataRepository, logger *logrus.Logger) *DataService {
	return &DataService{
		repo:             repo,
		logger:           logger,
		metadataCacheTTL: 1 * time.Hour,
	}
}

// GetAll retrieves data rows with optional vendor scoping, pagination, sorting, and filtering.
func (s *DataService) GetAll(filter *repository.VendorFilter, p repository.QueryParams) ([]*models.DataRow, int, error) {
	s.logger.Info("Fetching data rows")
	return s.repo.GetAll(filter, p)
}

// GetByTerminalID retrieves a single row by terminal ID with vendor scoping.
func (s *DataService) GetByTerminalID(terminalID string, filter *repository.VendorFilter) (*models.DataRow, error) {
	s.logger.Infof("Fetching data row for terminal: %s", terminalID)
	return s.repo.GetByTerminalID(terminalID, filter)
}

// Update modifies ticket fields with vendor filter enforcement.
func (s *DataService) Update(terminalID string, req *models.DataUpdateRequest, filter *repository.VendorFilter) (*models.DataRow, error) {
	s.logger.Infof("Updating data row for terminal: %s", terminalID)
	return s.repo.Update(terminalID, req, filter)
}

// GetMetadata returns distinct status/mode/priority values with 1-hour caching.
func (s *DataService) GetMetadata() (*models.MetadataResponse, error) {
	s.metadataCacheMux.RLock()
	if s.metadataCache != nil && time.Since(s.metadataLastFetch) < s.metadataCacheTTL {
		s.logger.Info("Returning cached metadata")
		cached := s.metadataCache
		s.metadataCacheMux.RUnlock()
		return cached, nil
	}
	s.metadataCacheMux.RUnlock()

	s.logger.Info("Fetching fresh metadata from database")

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

	statusInfos := make([]models.StatusInfo, 0, len(statuses))
	for _, v := range statuses {
		statusInfos = append(statusInfos, models.BuildStatusInfo(v))
	}
	modeInfos := make([]models.ModeInfo, 0, len(modes))
	for _, v := range modes {
		modeInfos = append(modeInfos, models.BuildModeInfo(v))
	}
	priorityInfos := make([]models.PriorityInfo, 0, len(priorities))
	for _, v := range priorities {
		priorityInfos = append(priorityInfos, models.BuildPriorityInfo(v))
	}

	resp := &models.MetadataResponse{
		Success:     true,
		Message:     "Metadata retrieved successfully",
		Statuses:    statusInfos,
		Modes:       modeInfos,
		Priorities:  priorityInfos,
		LastUpdated: time.Now().Format(time.RFC3339),
	}

	s.metadataCacheMux.Lock()
	s.metadataCache = resp
	s.metadataLastFetch = time.Now()
	s.metadataCacheMux.Unlock()

	return resp, nil
}
