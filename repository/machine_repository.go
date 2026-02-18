package repository

import (
	"api-gateway/models"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// MachineRepository handles all database operations for machines/terminals
// Implements the data access layer for machine_master.dbo.atmi
type MachineRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewMachineRepository creates a new machine repository instance
func NewMachineRepository(db *sql.DB, logger *logrus.Logger) *MachineRepository {
	return &MachineRepository{
		db:     db,
		logger: logger,
	}
}

// GetAll retrieves machines with pagination support.
// If page <= 0, returns all machines (backwards compatible).
func (r *MachineRepository) GetAll(page, pageSize int) ([]*models.ATMI, int, error) {
	// Get total count
	var total int
	countErr := r.db.QueryRow("SELECT COUNT(*) FROM dbo.atmi").Scan(&total)
	if countErr != nil {
		r.logger.Errorf("Failed to count machines: %v", countErr)
		return nil, 0, fmt.Errorf("failed to count machines: %w", countErr)
	}

	var query string
	var rows *sql.Rows
	var err error

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = `
			SELECT
				terminal_id, store, store_code, store_name,
				date_of_activation, status, std,
				gps, lat, lon, province, [city/regency], district
			FROM dbo.atmi
			ORDER BY terminal_id ASC
			OFFSET @p1 ROWS FETCH NEXT @p2 ROWS ONLY
		`
		rows, err = r.db.Query(query, offset, pageSize)
	} else {
		query = `
			SELECT
				terminal_id, store, store_code, store_name,
				date_of_activation, status, std,
				gps, lat, lon, province, [city/regency], district
			FROM dbo.atmi
			ORDER BY terminal_id ASC
		`
		rows, err = r.db.Query(query)
	}

	if err != nil {
		r.logger.Errorf("Failed to query machines: %v", err)
		return nil, 0, fmt.Errorf("failed to query machines: %w", err)
	}
	defer rows.Close()

	machines := make([]*models.ATMI, 0, pageSize)
	for rows.Next() {
		machine := &models.ATMI{}
		err := rows.Scan(
			&machine.TerminalID,
			&machine.Store,
			&machine.StoreCode,
			&machine.StoreName,
			&machine.DateOfActivation,
			&machine.Status,
			&machine.Std,
			&machine.GPS,
			&machine.Lat,
			&machine.Lon,
			&machine.Province,
			&machine.CityRegency,
			&machine.District,
		)
		if err != nil {
			r.logger.Errorf("Failed to scan machine row: %v", err)
			continue
		}
		machines = append(machines, machine)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating machine rows: %w", err)
	}

	return machines, total, nil
}

// GetByTerminalID retrieves a single machine by its terminal ID
func (r *MachineRepository) GetByTerminalID(terminalID string) (*models.ATMI, error) {
	query := `
		SELECT
			terminal_id, store, store_code, store_name,
			date_of_activation, status, std,
			gps, lat, lon, province, [city/regency], district
		FROM dbo.atmi
		WHERE terminal_id = @p1
	`

	machine := &models.ATMI{}
	err := r.db.QueryRow(query, terminalID).Scan(
		&machine.TerminalID,
		&machine.Store,
		&machine.StoreCode,
		&machine.StoreName,
		&machine.DateOfActivation,
		&machine.Status,
		&machine.Std,
		&machine.GPS,
		&machine.Lat,
		&machine.Lon,
		&machine.Province,
		&machine.CityRegency,
		&machine.District,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("machine not found")
	}
	if err != nil {
		r.logger.Errorf("Failed to get machine by terminal ID: %v", err)
		return nil, fmt.Errorf("failed to get machine: %w", err)
	}

	return machine, nil
}

// GetByStatus retrieves all machines with a specific status
func (r *MachineRepository) GetByStatus(status string) ([]*models.ATMI, error) {
	query := `
		SELECT
			terminal_id, store, store_code, store_name,
			date_of_activation, status, std,
			gps, lat, lon, province, [city/regency], district
		FROM dbo.atmi
		WHERE status = @p1
		ORDER BY terminal_id ASC
	`

	rows, err := r.db.Query(query, status)
	if err != nil {
		r.logger.Errorf("Failed to query machines by status: %v", err)
		return nil, fmt.Errorf("failed to query machines: %w", err)
	}
	defer rows.Close()

	machines := make([]*models.ATMI, 0)
	for rows.Next() {
		machine := &models.ATMI{}
		err := rows.Scan(
			&machine.TerminalID,
			&machine.Store,
			&machine.StoreCode,
			&machine.StoreName,
			&machine.DateOfActivation,
			&machine.Status,
			&machine.Std,
			&machine.GPS,
			&machine.Lat,
			&machine.Lon,
			&machine.Province,
			&machine.CityRegency,
			&machine.District,
		)
		if err != nil {
			r.logger.Errorf("Failed to scan machine row: %v", err)
			continue
		}
		machines = append(machines, machine)
	}

	return machines, nil
}

// GetByStoreCode retrieves all machines for a specific store code
func (r *MachineRepository) GetByStoreCode(storeCode string) ([]*models.ATMI, error) {
	query := `
		SELECT
			terminal_id, store, store_code, store_name,
			date_of_activation, status, std,
			gps, lat, lon, province, [city/regency], district
		FROM dbo.atmi
		WHERE store_code = @p1
		ORDER BY terminal_id ASC
	`

	rows, err := r.db.Query(query, storeCode)
	if err != nil {
		r.logger.Errorf("Failed to query machines by store code: %v", err)
		return nil, fmt.Errorf("failed to query machines: %w", err)
	}
	defer rows.Close()

	machines := make([]*models.ATMI, 0)
	for rows.Next() {
		machine := &models.ATMI{}
		err := rows.Scan(
			&machine.TerminalID,
			&machine.Store,
			&machine.StoreCode,
			&machine.StoreName,
			&machine.DateOfActivation,
			&machine.Status,
			&machine.Std,
			&machine.GPS,
			&machine.Lat,
			&machine.Lon,
			&machine.Province,
			&machine.CityRegency,
			&machine.District,
		)
		if err != nil {
			r.logger.Errorf("Failed to scan machine row: %v", err)
			continue
		}
		machines = append(machines, machine)
	}

	return machines, nil
}

// UpdateStatus updates the status and location of a machine
func (r *MachineRepository) UpdateStatus(req *models.MachineStatusUpdate) (*models.ATMI, error) {
	updates := []string{"status = @p1"}
	args := []interface{}{req.Status}
	paramCount := 2

	if req.GPS != "" {
		updates = append(updates, fmt.Sprintf("gps = @p%d", paramCount))
		args = append(args, req.GPS)
		paramCount++
	}

	if req.Lat != 0 {
		updates = append(updates, fmt.Sprintf("lat = @p%d", paramCount))
		args = append(args, req.Lat)
		paramCount++
	}

	if req.Lon != 0 {
		updates = append(updates, fmt.Sprintf("lon = @p%d", paramCount))
		args = append(args, req.Lon)
		paramCount++
	}

	// Add terminal_id as the last parameter
	args = append(args, req.TerminalID)

	query := fmt.Sprintf(
		"UPDATE dbo.atmi SET %s WHERE terminal_id = @p%d",
		strings.Join(updates, ", "),
		paramCount,
	)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		r.logger.Errorf("Failed to update machine status: %v", err)
		return nil, fmt.Errorf("failed to update machine status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("machine not found")
	}

	return r.GetByTerminalID(req.TerminalID)
}

// Search performs a flexible search across multiple fields
func (r *MachineRepository) Search(filter *models.MachineFilter) ([]*models.ATMI, error) {
	query := `
		SELECT
			terminal_id, store, store_code, store_name,
			date_of_activation, status, std,
			gps, lat, lon, province, [city/regency], district
		FROM dbo.atmi
		WHERE 1=1
	`

	args := []interface{}{}
	paramCount := 1

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = @p%d", paramCount)
		args = append(args, filter.Status)
		paramCount++
	}

	if filter.StoreCode != "" {
		query += fmt.Sprintf(" AND store_code = @p%d", paramCount)
		args = append(args, filter.StoreCode)
		paramCount++
	}

	if filter.Province != "" {
		query += fmt.Sprintf(" AND province = @p%d", paramCount)
		args = append(args, filter.Province)
		paramCount++
	}

	if filter.CityRegency != "" {
		query += fmt.Sprintf(" AND [city/regency] = @p%d", paramCount)
		args = append(args, filter.CityRegency)
		paramCount++
	}

	if filter.District != "" {
		query += fmt.Sprintf(" AND district LIKE @p%d", paramCount)
		args = append(args, "%"+filter.District+"%")
		paramCount++
	}

	query += " ORDER BY terminal_id ASC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		r.logger.Errorf("Failed to search machines: %v", err)
		return nil, fmt.Errorf("failed to search machines: %w", err)
	}
	defer rows.Close()

	machines := make([]*models.ATMI, 0)
	for rows.Next() {
		machine := &models.ATMI{}
		err := rows.Scan(
			&machine.TerminalID,
			&machine.Store,
			&machine.StoreCode,
			&machine.StoreName,
			&machine.DateOfActivation,
			&machine.Status,
			&machine.Std,
			&machine.GPS,
			&machine.Lat,
			&machine.Lon,
			&machine.Province,
			&machine.CityRegency,
			&machine.District,
		)
		if err != nil {
			r.logger.Errorf("Failed to scan machine row: %v", err)
			continue
		}
		machines = append(machines, machine)
	}

	return machines, nil
}

// GetDistinctSLMs retrieves all unique SLM values from the database
// This provides a truly adaptive list of what SLM types are actually in use
func (r *MachineRepository) GetDistinctSLMs() ([]string, error) {
	query := `
		SELECT DISTINCT [slm]
		FROM dbo.atmi
		WHERE [slm] IS NOT NULL
		ORDER BY [slm]
	`

	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Errorf("Failed to query distinct SLMs: %v", err)
		return nil, fmt.Errorf("failed to query SLMs: %w", err)
	}
	defer rows.Close()

	slms := []string{}
	for rows.Next() {
		var slm string
		if err := rows.Scan(&slm); err != nil {
			r.logger.Errorf("Failed to scan SLM: %v", err)
			continue
		}
		slms = append(slms, slm)
	}

	return slms, nil
}

// GetDistinctFLMs retrieves all unique FLM values from the database
func (r *MachineRepository) GetDistinctFLMs() ([]string, error) {
	query := `
		SELECT DISTINCT [flm]
		FROM dbo.atmi
		WHERE [flm] IS NOT NULL
		ORDER BY [flm]
	`

	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Errorf("Failed to query distinct FLMs: %v", err)
		return nil, fmt.Errorf("failed to query FLMs: %w", err)
	}
	defer rows.Close()

	flms := []string{}
	for rows.Next() {
		var flm string
		if err := rows.Scan(&flm); err != nil {
			r.logger.Errorf("Failed to scan FLM: %v", err)
			continue
		}
		flms = append(flms, flm)
	}

	return flms, nil
}

// GetDistinctNETs retrieves all unique network provider values from the database
func (r *MachineRepository) GetDistinctNETs() ([]string, error) {
	query := `
		SELECT DISTINCT [net]
		FROM dbo.atmi
		WHERE [net] IS NOT NULL
		ORDER BY [net]
	`

	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Errorf("Failed to query distinct NETs: %v", err)
		return nil, fmt.Errorf("failed to query NETs: %w", err)
	}
	defer rows.Close()

	nets := []string{}
	for rows.Next() {
		var net string
		if err := rows.Scan(&net); err != nil {
			r.logger.Errorf("Failed to scan NET: %v", err)
			continue
		}
		nets = append(nets, net)
	}

	return nets, nil
}

// GetDistinctFLMNames retrieves all unique FLM name values from the database
func (r *MachineRepository) GetDistinctFLMNames() ([]string, error) {
	query := `
		SELECT DISTINCT [flm_name]
		FROM dbo.atmi
		WHERE [flm_name] IS NOT NULL
		ORDER BY [flm_name]
	`

	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Errorf("Failed to query distinct FLM names: %v", err)
		return nil, fmt.Errorf("failed to query FLM names: %w", err)
	}
	defer rows.Close()

	flmNames := []string{}
	for rows.Next() {
		var flmName string
		if err := rows.Scan(&flmName); err != nil {
			r.logger.Errorf("Failed to scan FLM name: %v", err)
			continue
		}
		flmNames = append(flmNames, flmName)
	}

	return flmNames, nil
}
