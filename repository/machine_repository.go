package repository

import (
	"api-gateway/models"
	"database/sql"
	"fmt"
	"strings"
	"time"

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

// GetAll retrieves all machines from the database
// Returns a slice of ATMI and any error encountered
func (r *MachineRepository) GetAll() ([]*models.ATMI, error) {
	query := `
		SELECT
			id, terminal_id, terminal_name, location, branch_code,
			ip_address, model, manufacturer, serial_number, status,
			last_ping_time, install_date, warranty_exp, notes,
			created_at, updated_at
		FROM dbo.atmi
		ORDER BY terminal_id ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Errorf("Failed to query machines: %v", err)
		return nil, fmt.Errorf("failed to query machines: %w", err)
	}
	defer rows.Close()

	machines := make([]*models.ATMI, 0)
	for rows.Next() {
		machine := &models.ATMI{}
		err := rows.Scan(
			&machine.ID,
			&machine.TerminalID,
			&machine.TerminalName,
			&machine.Location,
			&machine.BranchCode,
			&machine.IPAddress,
			&machine.Model,
			&machine.Manufacturer,
			&machine.SerialNumber,
			&machine.Status,
			&machine.LastPingTime,
			&machine.InstallDate,
			&machine.WarrantyExp,
			&machine.Notes,
			&machine.CreatedAt,
			&machine.UpdatedAt,
		)
		if err != nil {
			r.logger.Errorf("Failed to scan machine row: %v", err)
			continue
		}
		machines = append(machines, machine)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating machine rows: %w", err)
	}

	return machines, nil
}

// GetByTerminalID retrieves a single machine by its terminal ID
func (r *MachineRepository) GetByTerminalID(terminalID string) (*models.ATMI, error) {
	query := `
		SELECT
			id, terminal_id, terminal_name, location, branch_code,
			ip_address, model, manufacturer, serial_number, status,
			last_ping_time, install_date, warranty_exp, notes,
			created_at, updated_at
		FROM dbo.atmi
		WHERE terminal_id = @p1
	`

	machine := &models.ATMI{}
	err := r.db.QueryRow(query, terminalID).Scan(
		&machine.ID,
		&machine.TerminalID,
		&machine.TerminalName,
		&machine.Location,
		&machine.BranchCode,
		&machine.IPAddress,
		&machine.Model,
		&machine.Manufacturer,
		&machine.SerialNumber,
		&machine.Status,
		&machine.LastPingTime,
		&machine.InstallDate,
		&machine.WarrantyExp,
		&machine.Notes,
		&machine.CreatedAt,
		&machine.UpdatedAt,
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
			id, terminal_id, terminal_name, location, branch_code,
			ip_address, model, manufacturer, serial_number, status,
			last_ping_time, install_date, warranty_exp, notes,
			created_at, updated_at
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
			&machine.ID,
			&machine.TerminalID,
			&machine.TerminalName,
			&machine.Location,
			&machine.BranchCode,
			&machine.IPAddress,
			&machine.Model,
			&machine.Manufacturer,
			&machine.SerialNumber,
			&machine.Status,
			&machine.LastPingTime,
			&machine.InstallDate,
			&machine.WarrantyExp,
			&machine.Notes,
			&machine.CreatedAt,
			&machine.UpdatedAt,
		)
		if err != nil {
			r.logger.Errorf("Failed to scan machine row: %v", err)
			continue
		}
		machines = append(machines, machine)
	}

	return machines, nil
}

// GetByBranchCode retrieves all machines for a specific branch
func (r *MachineRepository) GetByBranchCode(branchCode string) ([]*models.ATMI, error) {
	query := `
		SELECT
			id, terminal_id, terminal_name, location, branch_code,
			ip_address, model, manufacturer, serial_number, status,
			last_ping_time, install_date, warranty_exp, notes,
			created_at, updated_at
		FROM dbo.atmi
		WHERE branch_code = @p1
		ORDER BY terminal_id ASC
	`

	rows, err := r.db.Query(query, branchCode)
	if err != nil {
		r.logger.Errorf("Failed to query machines by branch: %v", err)
		return nil, fmt.Errorf("failed to query machines: %w", err)
	}
	defer rows.Close()

	machines := make([]*models.ATMI, 0)
	for rows.Next() {
		machine := &models.ATMI{}
		err := rows.Scan(
			&machine.ID,
			&machine.TerminalID,
			&machine.TerminalName,
			&machine.Location,
			&machine.BranchCode,
			&machine.IPAddress,
			&machine.Model,
			&machine.Manufacturer,
			&machine.SerialNumber,
			&machine.Status,
			&machine.LastPingTime,
			&machine.InstallDate,
			&machine.WarrantyExp,
			&machine.Notes,
			&machine.CreatedAt,
			&machine.UpdatedAt,
		)
		if err != nil {
			r.logger.Errorf("Failed to scan machine row: %v", err)
			continue
		}
		machines = append(machines, machine)
	}

	return machines, nil
}

// UpdateStatus updates the status and related fields of a machine
func (r *MachineRepository) UpdateStatus(req *models.MachineStatusUpdate) (*models.ATMI, error) {
	updates := []string{"status = @p1", "updated_at = @p2"}
	args := []interface{}{req.Status, time.Now()}
	paramCount := 3

	// Add last_ping_time if provided
	if req.LastPingTime != "" {
		pingTime, err := time.Parse(time.RFC3339, req.LastPingTime)
		if err != nil {
			// Try alternative format
			pingTime, err = time.Parse("2006-01-02 15:04:05", req.LastPingTime)
			if err != nil {
				return nil, fmt.Errorf("invalid last_ping_time format: %w", err)
			}
		}
		updates = append(updates, fmt.Sprintf("last_ping_time = @p%d", paramCount))
		args = append(args, pingTime)
		paramCount++
	}

	// Add notes if provided
	if req.Notes != "" {
		updates = append(updates, fmt.Sprintf("notes = @p%d", paramCount))
		args = append(args, req.Notes)
		paramCount++
	}

	// Add terminal_id as the last parameter for WHERE clause
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

	// Return the updated machine
	return r.GetByTerminalID(req.TerminalID)
}

// Search performs a flexible search across multiple fields
func (r *MachineRepository) Search(filter *models.MachineFilter) ([]*models.ATMI, error) {
	query := `
		SELECT
			id, terminal_id, terminal_name, location, branch_code,
			ip_address, model, manufacturer, serial_number, status,
			last_ping_time, install_date, warranty_exp, notes,
			created_at, updated_at
		FROM dbo.atmi
		WHERE 1=1
	`

	args := []interface{}{}
	paramCount := 1

	// Add filters dynamically
	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = @p%d", paramCount)
		args = append(args, filter.Status)
		paramCount++
	}

	if filter.BranchCode != "" {
		query += fmt.Sprintf(" AND branch_code = @p%d", paramCount)
		args = append(args, filter.BranchCode)
		paramCount++
	}

	if filter.Location != "" {
		query += fmt.Sprintf(" AND location LIKE @p%d", paramCount)
		args = append(args, "%"+filter.Location+"%")
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
			&machine.ID,
			&machine.TerminalID,
			&machine.TerminalName,
			&machine.Location,
			&machine.BranchCode,
			&machine.IPAddress,
			&machine.Model,
			&machine.Manufacturer,
			&machine.SerialNumber,
			&machine.Status,
			&machine.LastPingTime,
			&machine.InstallDate,
			&machine.WarrantyExp,
			&machine.Notes,
			&machine.CreatedAt,
			&machine.UpdatedAt,
		)
		if err != nil {
			r.logger.Errorf("Failed to scan machine row: %v", err)
			continue
		}
		machines = append(machines, machine)
	}

	return machines, nil
}
