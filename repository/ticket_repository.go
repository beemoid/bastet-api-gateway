package repository

import (
	"api-gateway/models"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// TicketRepository handles all database operations for tickets
// Implements the data access layer for ticket_master.dbo.open_ticket
type TicketRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewTicketRepository creates a new ticket repository instance
func NewTicketRepository(db *sql.DB, logger *logrus.Logger) *TicketRepository {
	return &TicketRepository{
		db:     db,
		logger: logger,
	}
}

// GetAll retrieves all open tickets from the database
// Returns a slice of OpenTicket and any error encountered
func (r *TicketRepository) GetAll() ([]*models.OpenTicket, error) {
	query := `
		SELECT
			id, ticket_number, terminal_id, description, priority,
			status, category, reported_by, assigned_to, created_at,
			updated_at, resolved_at, resolution_notes
		FROM dbo.open_ticket
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Errorf("Failed to query tickets: %v", err)
		return nil, fmt.Errorf("failed to query tickets: %w", err)
	}
	defer rows.Close()

	tickets := make([]*models.OpenTicket, 0)
	for rows.Next() {
		ticket := &models.OpenTicket{}
		err := rows.Scan(
			&ticket.ID,
			&ticket.TicketNumber,
			&ticket.TerminalID,
			&ticket.Description,
			&ticket.Priority,
			&ticket.Status,
			&ticket.Category,
			&ticket.ReportedBy,
			&ticket.AssignedTo,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
			&ticket.ResolvedAt,
			&ticket.ResolutionNotes,
		)
		if err != nil {
			r.logger.Errorf("Failed to scan ticket row: %v", err)
			continue
		}
		tickets = append(tickets, ticket)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating ticket rows: %w", err)
	}

	return tickets, nil
}

// GetByID retrieves a single ticket by its ID
func (r *TicketRepository) GetByID(id int) (*models.OpenTicket, error) {
	query := `
		SELECT
			id, ticket_number, terminal_id, description, priority,
			status, category, reported_by, assigned_to, created_at,
			updated_at, resolved_at, resolution_notes
		FROM dbo.open_ticket
		WHERE id = @p1
	`

	ticket := &models.OpenTicket{}
	err := r.db.QueryRow(query, id).Scan(
		&ticket.ID,
		&ticket.TicketNumber,
		&ticket.TerminalID,
		&ticket.Description,
		&ticket.Priority,
		&ticket.Status,
		&ticket.Category,
		&ticket.ReportedBy,
		&ticket.AssignedTo,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
		&ticket.ResolvedAt,
		&ticket.ResolutionNotes,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ticket not found")
	}
	if err != nil {
		r.logger.Errorf("Failed to get ticket by ID: %v", err)
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	return ticket, nil
}

// GetByTicketNumber retrieves a ticket by its unique ticket number
func (r *TicketRepository) GetByTicketNumber(ticketNumber string) (*models.OpenTicket, error) {
	query := `
		SELECT
			id, ticket_number, terminal_id, description, priority,
			status, category, reported_by, assigned_to, created_at,
			updated_at, resolved_at, resolution_notes
		FROM dbo.open_ticket
		WHERE ticket_number = @p1
	`

	ticket := &models.OpenTicket{}
	err := r.db.QueryRow(query, ticketNumber).Scan(
		&ticket.ID,
		&ticket.TicketNumber,
		&ticket.TerminalID,
		&ticket.Description,
		&ticket.Priority,
		&ticket.Status,
		&ticket.Category,
		&ticket.ReportedBy,
		&ticket.AssignedTo,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
		&ticket.ResolvedAt,
		&ticket.ResolutionNotes,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ticket not found")
	}
	if err != nil {
		r.logger.Errorf("Failed to get ticket by number: %v", err)
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	return ticket, nil
}

// Create inserts a new ticket into the database
// Returns the created ticket with populated ID and timestamps
func (r *TicketRepository) Create(req *models.TicketCreateRequest) (*models.OpenTicket, error) {
	query := `
		INSERT INTO dbo.open_ticket
		(ticket_number, terminal_id, description, priority, status, category, reported_by, assigned_to, created_at, updated_at)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10);
		SELECT SCOPE_IDENTITY();
	`

	now := time.Now()
	status := "Open" // Default status for new tickets

	var assignedTo sql.NullString
	if req.AssignedTo != "" {
		assignedTo = sql.NullString{String: req.AssignedTo, Valid: true}
	}

	// Execute insert and get the new ID
	var newID int64
	err := r.db.QueryRow(
		query,
		req.TicketNumber,
		req.TerminalID,
		req.Description,
		req.Priority,
		status,
		req.Category,
		req.ReportedBy,
		assignedTo,
		now,
		now,
	).Scan(&newID)

	if err != nil {
		r.logger.Errorf("Failed to create ticket: %v", err)
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	// Retrieve and return the created ticket
	return r.GetByID(int(newID))
}

// Update modifies an existing ticket
// Only updates fields that are provided (non-empty)
func (r *TicketRepository) Update(id int, req *models.TicketUpdateRequest) (*models.OpenTicket, error) {
	// Build dynamic update query based on provided fields
	updates := []string{}
	args := []interface{}{}
	paramCount := 1

	if req.Status != "" {
		updates = append(updates, fmt.Sprintf("status = @p%d", paramCount))
		args = append(args, req.Status)
		paramCount++

		// If status is Resolved, set resolved_at timestamp
		if req.Status == "Resolved" {
			updates = append(updates, fmt.Sprintf("resolved_at = @p%d", paramCount))
			args = append(args, time.Now())
			paramCount++
		}
	}

	if req.Priority != "" {
		updates = append(updates, fmt.Sprintf("priority = @p%d", paramCount))
		args = append(args, req.Priority)
		paramCount++
	}

	if req.AssignedTo != "" {
		updates = append(updates, fmt.Sprintf("assigned_to = @p%d", paramCount))
		args = append(args, req.AssignedTo)
		paramCount++
	}

	if req.Description != "" {
		updates = append(updates, fmt.Sprintf("description = @p%d", paramCount))
		args = append(args, req.Description)
		paramCount++
	}

	if req.ResolutionNotes != "" {
		updates = append(updates, fmt.Sprintf("resolution_notes = @p%d", paramCount))
		args = append(args, req.ResolutionNotes)
		paramCount++
	}

	// Always update the updated_at timestamp
	updates = append(updates, fmt.Sprintf("updated_at = @p%d", paramCount))
	args = append(args, time.Now())
	paramCount++

	// Add ID as the last parameter for WHERE clause
	args = append(args, id)

	if len(updates) == 1 { // Only updated_at was added
		return nil, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(
		"UPDATE dbo.open_ticket SET %s WHERE id = @p%d",
		strings.Join(updates, ", "),
		paramCount,
	)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		r.logger.Errorf("Failed to update ticket: %v", err)
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("ticket not found")
	}

	// Return the updated ticket
	return r.GetByID(id)
}

// GetByStatus retrieves all tickets with a specific status
func (r *TicketRepository) GetByStatus(status string) ([]*models.OpenTicket, error) {
	query := `
		SELECT
			id, ticket_number, terminal_id, description, priority,
			status, category, reported_by, assigned_to, created_at,
			updated_at, resolved_at, resolution_notes
		FROM dbo.open_ticket
		WHERE status = @p1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, status)
	if err != nil {
		r.logger.Errorf("Failed to query tickets by status: %v", err)
		return nil, fmt.Errorf("failed to query tickets: %w", err)
	}
	defer rows.Close()

	tickets := make([]*models.OpenTicket, 0)
	for rows.Next() {
		ticket := &models.OpenTicket{}
		err := rows.Scan(
			&ticket.ID,
			&ticket.TicketNumber,
			&ticket.TerminalID,
			&ticket.Description,
			&ticket.Priority,
			&ticket.Status,
			&ticket.Category,
			&ticket.ReportedBy,
			&ticket.AssignedTo,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
			&ticket.ResolvedAt,
			&ticket.ResolutionNotes,
		)
		if err != nil {
			r.logger.Errorf("Failed to scan ticket row: %v", err)
			continue
		}
		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

// GetByTerminalID retrieves all tickets for a specific terminal
func (r *TicketRepository) GetByTerminalID(terminalID string) ([]*models.OpenTicket, error) {
	query := `
		SELECT
			id, ticket_number, terminal_id, description, priority,
			status, category, reported_by, assigned_to, created_at,
			updated_at, resolved_at, resolution_notes
		FROM dbo.open_ticket
		WHERE terminal_id = @p1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, terminalID)
	if err != nil {
		r.logger.Errorf("Failed to query tickets by terminal: %v", err)
		return nil, fmt.Errorf("failed to query tickets: %w", err)
	}
	defer rows.Close()

	tickets := make([]*models.OpenTicket, 0)
	for rows.Next() {
		ticket := &models.OpenTicket{}
		err := rows.Scan(
			&ticket.ID,
			&ticket.TicketNumber,
			&ticket.TerminalID,
			&ticket.Description,
			&ticket.Priority,
			&ticket.Status,
			&ticket.Category,
			&ticket.ReportedBy,
			&ticket.AssignedTo,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
			&ticket.ResolvedAt,
			&ticket.ResolutionNotes,
		)
		if err != nil {
			r.logger.Errorf("Failed to scan ticket row: %v", err)
			continue
		}
		tickets = append(tickets, ticket)
	}

	return tickets, nil
}
