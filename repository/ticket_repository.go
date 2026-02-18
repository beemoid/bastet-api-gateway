package repository

import (
	"api-gateway/models"
	"database/sql"
	"fmt"
	"strings"

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

// GetAll retrieves tickets with pagination support.
// If page <= 0, returns all tickets (backwards compatible).
func (r *TicketRepository) GetAll(page, pageSize int) ([]*models.OpenTicket, int, error) {
	// Get total count
	var total int
	countErr := r.db.QueryRow("SELECT COUNT(*) FROM dbo.open_ticket").Scan(&total)
	if countErr != nil {
		r.logger.Errorf("Failed to count tickets: %v", countErr)
		return nil, 0, fmt.Errorf("failed to count tickets: %w", countErr)
	}

	var query string
	var rows *sql.Rows
	var err error

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = `
			SELECT
				[Terminal ID], [Terminal Name], [Priority], [Mode],
				[Initial Problem], [Current Problem], [P-Duration],
				[Incident start datetime], [Count], [Status], [Remarks],
				[Balance], [Condition], [Tickets no], [Tickets duration],
				[Open time], [Close time], [Problem History], [Mode History],
				[DSP FLM], [DSP SLM], [Last Withdrawal], [Export Name]
			FROM dbo.open_ticket
			ORDER BY [Incident start datetime] DESC
			OFFSET @p1 ROWS FETCH NEXT @p2 ROWS ONLY
		`
		rows, err = r.db.Query(query, offset, pageSize)
	} else {
		query = `
			SELECT
				[Terminal ID], [Terminal Name], [Priority], [Mode],
				[Initial Problem], [Current Problem], [P-Duration],
				[Incident start datetime], [Count], [Status], [Remarks],
				[Balance], [Condition], [Tickets no], [Tickets duration],
				[Open time], [Close time], [Problem History], [Mode History],
				[DSP FLM], [DSP SLM], [Last Withdrawal], [Export Name]
			FROM dbo.open_ticket
			ORDER BY [Incident start datetime] DESC
		`
		rows, err = r.db.Query(query)
	}

	if err != nil {
		r.logger.Errorf("Failed to query tickets: %v", err)
		return nil, 0, fmt.Errorf("failed to query tickets: %w", err)
	}
	defer rows.Close()

	tickets := make([]*models.OpenTicket, 0, pageSize)
	for rows.Next() {
		ticket := &models.OpenTicket{}
		err := rows.Scan(
			&ticket.TerminalID,
			&ticket.TerminalName,
			&ticket.Priority,
			&ticket.Mode,
			&ticket.InitialProblem,
			&ticket.CurrentProblem,
			&ticket.PDuration,
			&ticket.IncidentStartTime,
			&ticket.Count,
			&ticket.Status,
			&ticket.Remarks,
			&ticket.Balance,
			&ticket.Condition,
			&ticket.TicketsNo,
			&ticket.TicketsDuration,
			&ticket.OpenTime,
			&ticket.CloseTime,
			&ticket.ProblemHistory,
			&ticket.ModeHistory,
			&ticket.DSPFLM,
			&ticket.DSPSLM,
			&ticket.LastWithdrawal,
			&ticket.ExportName,
		)
		if err != nil {
			r.logger.Errorf("Failed to scan ticket row: %v", err)
			continue
		}
		tickets = append(tickets, ticket)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating ticket rows: %w", err)
	}

	return tickets, total, nil
}

// GetByTerminalID retrieves a single ticket by terminal ID
func (r *TicketRepository) GetByTerminalID(terminalID string) (*models.OpenTicket, error) {
	query := `
		SELECT
			[Terminal ID], [Terminal Name], [Priority], [Mode],
			[Initial Problem], [Current Problem], [P-Duration],
			[Incident start datetime], [Count], [Status], [Remarks],
			[Balance], [Condition], [Tickets no], [Tickets duration],
			[Open time], [Close time], [Problem History], [Mode History],
			[DSP FLM], [DSP SLM], [Last Withdrawal], [Export Name]
		FROM dbo.open_ticket
		WHERE [Terminal ID] = @p1
	`

	ticket := &models.OpenTicket{}
	err := r.db.QueryRow(query, terminalID).Scan(
		&ticket.TerminalID,
		&ticket.TerminalName,
		&ticket.Priority,
		&ticket.Mode,
		&ticket.InitialProblem,
		&ticket.CurrentProblem,
		&ticket.PDuration,
		&ticket.IncidentStartTime,
		&ticket.Count,
		&ticket.Status,
		&ticket.Remarks,
		&ticket.Balance,
		&ticket.Condition,
		&ticket.TicketsNo,
		&ticket.TicketsDuration,
		&ticket.OpenTime,
		&ticket.CloseTime,
		&ticket.ProblemHistory,
		&ticket.ModeHistory,
		&ticket.DSPFLM,
		&ticket.DSPSLM,
		&ticket.LastWithdrawal,
		&ticket.ExportName,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ticket not found")
	}
	if err != nil {
		r.logger.Errorf("Failed to get ticket by terminal ID: %v", err)
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	return ticket, nil
}

// GetByTicketNumber retrieves a ticket by its unique ticket number
func (r *TicketRepository) GetByTicketNumber(ticketNumber string) (*models.OpenTicket, error) {
	query := `
		SELECT
			[Terminal ID], [Terminal Name], [Priority], [Mode],
			[Initial Problem], [Current Problem], [P-Duration],
			[Incident start datetime], [Count], [Status], [Remarks],
			[Balance], [Condition], [Tickets no], [Tickets duration],
			[Open time], [Close time], [Problem History], [Mode History],
			[DSP FLM], [DSP SLM], [Last Withdrawal], [Export Name]
		FROM dbo.open_ticket
		WHERE [Tickets no] = @p1
	`

	ticket := &models.OpenTicket{}
	err := r.db.QueryRow(query, ticketNumber).Scan(
		&ticket.TerminalID,
		&ticket.TerminalName,
		&ticket.Priority,
		&ticket.Mode,
		&ticket.InitialProblem,
		&ticket.CurrentProblem,
		&ticket.PDuration,
		&ticket.IncidentStartTime,
		&ticket.Count,
		&ticket.Status,
		&ticket.Remarks,
		&ticket.Balance,
		&ticket.Condition,
		&ticket.TicketsNo,
		&ticket.TicketsDuration,
		&ticket.OpenTime,
		&ticket.CloseTime,
		&ticket.ProblemHistory,
		&ticket.ModeHistory,
		&ticket.DSPFLM,
		&ticket.DSPSLM,
		&ticket.LastWithdrawal,
		&ticket.ExportName,
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
func (r *TicketRepository) Create(req *models.TicketCreateRequest) (*models.OpenTicket, error) {
	query := `
		INSERT INTO dbo.open_ticket
		([Terminal ID], [Terminal Name], [Priority], [Mode], [Initial Problem],
		 [Current Problem], [P-Duration], [Incident start datetime], [Status],
		 [Remarks], [Condition], [Tickets no], [Export Name])
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13)
	`

	_, err := r.db.Exec(
		query,
		req.TerminalID,
		req.TerminalName,
		req.Priority,
		req.Mode,
		req.InitialProblem,
		req.CurrentProblem,
		req.PDuration,
		req.IncidentStartTime,
		req.Status,
		req.Remarks,
		req.Condition,
		req.TicketsNo,
		req.ExportName,
	)

	if err != nil {
		r.logger.Errorf("Failed to create ticket: %v", err)
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	// Retrieve and return the created ticket
	return r.GetByTicketNumber(req.TicketsNo)
}

// Update modifies an existing ticket
func (r *TicketRepository) Update(terminalID string, req *models.TicketUpdateRequest) (*models.OpenTicket, error) {
	updates := []string{}
	args := []interface{}{}
	paramCount := 1

	if req.Priority != "" {
		updates = append(updates, fmt.Sprintf("[Priority] = @p%d", paramCount))
		args = append(args, req.Priority)
		paramCount++
	}

	if req.Mode != "" {
		updates = append(updates, fmt.Sprintf("[Mode] = @p%d", paramCount))
		args = append(args, req.Mode)
		paramCount++
	}

	if req.CurrentProblem != "" {
		updates = append(updates, fmt.Sprintf("[Current Problem] = @p%d", paramCount))
		args = append(args, req.CurrentProblem)
		paramCount++
	}

	if req.Status != "" {
		updates = append(updates, fmt.Sprintf("[Status] = @p%d", paramCount))
		args = append(args, req.Status)
		paramCount++
	}

	if req.Remarks != "" {
		updates = append(updates, fmt.Sprintf("[Remarks] = @p%d", paramCount))
		args = append(args, req.Remarks)
		paramCount++
	}

	if req.Condition != "" {
		updates = append(updates, fmt.Sprintf("[Condition] = @p%d", paramCount))
		args = append(args, req.Condition)
		paramCount++
	}

	if req.CloseTime != "" {
		updates = append(updates, fmt.Sprintf("[Close time] = @p%d", paramCount))
		args = append(args, req.CloseTime)
		paramCount++
	}

	if req.ProblemHistory != "" {
		updates = append(updates, fmt.Sprintf("[Problem History] = @p%d", paramCount))
		args = append(args, req.ProblemHistory)
		paramCount++
	}

	if req.ModeHistory != "" {
		updates = append(updates, fmt.Sprintf("[Mode History] = @p%d", paramCount))
		args = append(args, req.ModeHistory)
		paramCount++
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// Add terminal ID as the last parameter
	args = append(args, terminalID)

	query := fmt.Sprintf(
		"UPDATE dbo.open_ticket SET %s WHERE [Terminal ID] = @p%d",
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

	return r.GetByTerminalID(terminalID)
}

// GetByStatus retrieves all tickets with a specific status
func (r *TicketRepository) GetByStatus(status string) ([]*models.OpenTicket, error) {
	query := `
		SELECT
			[Terminal ID], [Terminal Name], [Priority], [Mode],
			[Initial Problem], [Current Problem], [P-Duration],
			[Incident start datetime], [Count], [Status], [Remarks],
			[Balance], [Condition], [Tickets no], [Tickets duration],
			[Open time], [Close time], [Problem History], [Mode History],
			[DSP FLM], [DSP SLM], [Last Withdrawal], [Export Name]
		FROM dbo.open_ticket
		WHERE [Status] = @p1
		ORDER BY [Incident start datetime] DESC
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
			&ticket.TerminalID,
			&ticket.TerminalName,
			&ticket.Priority,
			&ticket.Mode,
			&ticket.InitialProblem,
			&ticket.CurrentProblem,
			&ticket.PDuration,
			&ticket.IncidentStartTime,
			&ticket.Count,
			&ticket.Status,
			&ticket.Remarks,
			&ticket.Balance,
			&ticket.Condition,
			&ticket.TicketsNo,
			&ticket.TicketsDuration,
			&ticket.OpenTime,
			&ticket.CloseTime,
			&ticket.ProblemHistory,
			&ticket.ModeHistory,
			&ticket.DSPFLM,
			&ticket.DSPSLM,
			&ticket.LastWithdrawal,
			&ticket.ExportName,
		)
		if err != nil {
			r.logger.Errorf("Failed to scan ticket row: %v", err)
			continue
		}
		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

// GetDistinctStatuses retrieves all unique status values from the database
// This provides a truly adaptive list of what statuses are actually in use
func (r *TicketRepository) GetDistinctStatuses() ([]string, error) {
	query := `
		SELECT DISTINCT [Status]
		FROM dbo.open_ticket
		WHERE [Status] IS NOT NULL AND [Status] != ''
		ORDER BY [Status]
	`

	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Errorf("Failed to query distinct statuses: %v", err)
		return nil, fmt.Errorf("failed to query statuses: %w", err)
	}
	defer rows.Close()

	statuses := []string{}
	for rows.Next() {
		var status string
		if err := rows.Scan(&status); err != nil {
			r.logger.Errorf("Failed to scan status: %v", err)
			continue
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

// GetDistinctModes retrieves all unique mode values from the database
func (r *TicketRepository) GetDistinctModes() ([]string, error) {
	query := `
		SELECT DISTINCT [Mode]
		FROM dbo.open_ticket
		WHERE [Mode] IS NOT NULL AND [Mode] != ''
		ORDER BY [Mode]
	`

	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Errorf("Failed to query distinct modes: %v", err)
		return nil, fmt.Errorf("failed to query modes: %w", err)
	}
	defer rows.Close()

	modes := []string{}
	for rows.Next() {
		var mode string
		if err := rows.Scan(&mode); err != nil {
			r.logger.Errorf("Failed to scan mode: %v", err)
			continue
		}
		modes = append(modes, mode)
	}

	return modes, nil
}

// GetDistinctPriorities retrieves all unique priority values from the database
func (r *TicketRepository) GetDistinctPriorities() ([]string, error) {
	query := `
		SELECT DISTINCT [Priority]
		FROM dbo.open_ticket
		WHERE [Priority] IS NOT NULL AND [Priority] != ''
		ORDER BY [Priority]
	`

	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Errorf("Failed to query distinct priorities: %v", err)
		return nil, fmt.Errorf("failed to query priorities: %w", err)
	}
	defer rows.Close()

	priorities := []string{}
	for rows.Next() {
		var priority string
		if err := rows.Scan(&priority); err != nil {
			r.logger.Errorf("Failed to scan priority: %v", err)
			continue
		}
		priorities = append(priorities, priority)
	}

	return priorities, nil
}
