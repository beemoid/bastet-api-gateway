package repository

import (
	"api-gateway/models"
	"api-gateway/repository/queries"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// ── Vendor filter infrastructure ─────────────────────────────────────────────

// vendorJoinSQL is the LEFT JOIN appended when a vendor filter is active.
const vendorJoinSQL = `
	LEFT JOIN machine_master.dbo.machine mm ON op.[Terminal ID] = mm.[Terminal ID]
`

// vendorFilterColumns maps the logical filter_column key (stored on a token)
// to the actual SQL column expression used in the WHERE clause.
// Add new entries here when new filterable dimensions are required.
var vendorFilterColumns = map[string]string{
	"flm_name":    "mm.[FLM name]",
	"flm":         "mm.[FLM]",
	"slm":         "mm.[SLM]",
	"net":         "mm.[Net]",
	"terminal_id": "op.[Terminal ID]",
	"status":      "op.[Status]",
	"priority":    "op.[Priority]",
}

// VendorFilter represents a parsed vendor scoping filter derived from the token.
type VendorFilter struct {
	IsSuperToken bool
	Column       string // resolved SQL column expression e.g. mm.[FLM name]
	Value        string
}

// ResolveVendorFilter translates a logical filter_column key to an SQL column expression.
func ResolveVendorFilter(filterColumn, filterValue string, isSuper bool) *VendorFilter {
	if isSuper {
		return &VendorFilter{IsSuperToken: true}
	}
	if filterColumn == "" || filterValue == "" {
		return nil
	}
	col, ok := vendorFilterColumns[strings.ToLower(filterColumn)]
	if !ok {
		col = filterColumn // fallback: treat raw value as column (admin-supplied)
	}
	return &VendorFilter{Column: col, Value: filterValue}
}

// ── Base SELECT shared by vendor queries ─────────────────────────────────────

// vendorDataSelect is the SELECT+FROM+JOIN block used for vendor-scoped queries.
// The JOIN is always included so machine columns are always populated.
const vendorDataSelect = `
	SELECT
		op.[Terminal ID],
		op.[Terminal Name],
		op.[Priority],
		op.[Mode],
		op.[Initial Problem],
		op.[Current Problem],
		op.[P-Duration],
		op.[Incident start datetime],
		op.[Count],
		op.[Status],
		op.[Remarks],
		op.[Balance],
		op.[Condition],
		op.[Tickets no],
		op.[Tickets duration],
		op.[Open time],
		op.[Close time],
		op.[Problem History],
		op.[Mode History],
		op.[DSP FLM],
		op.[DSP SLM],
		op.[Last Withdrawal],
		op.[Export Name],
		mm.[FLM name],
		mm.[FLM],
		mm.[SLM],
		mm.[Net]
	FROM ticket_master.dbo.open_ticket op
	LEFT JOIN machine_master.dbo.machine mm
		ON op.[Terminal ID] = mm.[Terminal ID]
`

// ── DataRepository ────────────────────────────────────────────────────────────

// DataRepository handles all database operations for the unified /api/v1/data endpoint.
// It joins ticket_master.dbo.open_ticket with machine_master.dbo.machine and applies
// vendor scoping based on the token's VendorFilter.
type DataRepository struct {
	// ticketDB is the connection to ticket_master (primary write target)
	ticketDB *sql.DB
	logger   *logrus.Logger
}

// NewDataRepository creates a new DataRepository.
// ticketDB must point to ticket_master — the machine_master JOIN is cross-database.
func NewDataRepository(ticketDB *sql.DB, logger *logrus.Logger) *DataRepository {
	return &DataRepository{
		ticketDB: ticketDB,
		logger:   logger,
	}
}

// scanDataRow scans a single result row into a DataRow.
// Column order must match vendorDataSelect / AdminDataQuery exactly (27 columns).
func scanDataRow(row interface {
	Scan(...interface{}) error
}) (*models.DataRow, error) {
	d := &models.DataRow{}
	return d, row.Scan(
		&d.TerminalID,
		&d.TerminalName,
		&d.Priority,
		&d.Mode,
		&d.InitialProblem,
		&d.CurrentProblem,
		&d.PDuration,
		&d.IncidentStartTime,
		&d.Count,
		&d.Status,
		&d.Remarks,
		&d.Balance,
		&d.Condition,
		&d.TicketsNo,
		&d.TicketsDuration,
		&d.OpenTime,
		&d.CloseTime,
		&d.ProblemHistory,
		&d.ModeHistory,
		&d.DSPFLM,
		&d.DSPSLM,
		&d.LastWithdrawal,
		&d.ExportName,
		&d.FLMName,
		&d.FLM,
		&d.SLM,
		&d.Net,
	)
}

// QueryParams holds all pagination, sorting, and filtering options for GetAll.
type QueryParams struct {
	Page      int
	PageSize  int
	SortBy    string // logical field name, e.g. "terminal_id", "status"
	SortOrder string // "asc" or "desc"
	Search    string // free-text search on terminal_id and terminal_name
	// Column filters
	Status   string
	Mode     string
	Priority string
}

// allowedSortColumns maps logical sort_by keys to safe SQL column expressions.
// Only keys present here are accepted — anything else falls back to the default.
var allowedSortColumns = map[string]string{
	"terminal_id":             "op.[Terminal ID]",
	"terminal_name":           "op.[Terminal Name]",
	"priority":                "op.[Priority]",
	"mode":                    "op.[Mode]",
	"status":                  "op.[Status]",
	"incident_start_datetime": "op.[Incident start datetime]",
	"count":                   "op.[Count]",
	"balance":                 "op.[Balance]",
	"tickets_duration":        "op.[Tickets duration]",
	"open_time":               "op.[Open time]",
	"close_time":              "op.[Close time]",
	"flm_name":                "mm.[FLM name]",
	"flm":                     "mm.[FLM]",
	"slm":                     "mm.[SLM]",
	"net":                     "mm.[Net]",
}

// buildOrderBy returns a safe ORDER BY clause from QueryParams.
func buildOrderBy(p QueryParams) string {
	col, ok := allowedSortColumns[strings.ToLower(p.SortBy)]
	if !ok {
		col = "op.[Incident start datetime]"
	}
	dir := "DESC"
	if strings.ToLower(p.SortOrder) == "asc" {
		dir = "ASC"
	}
	return fmt.Sprintf("ORDER BY %s %s", col, dir)
}

// GetAll retrieves rows with optional vendor scoping, pagination, sorting, and filtering.
// - filter == nil            → no vendor restriction (legacy / unrestricted token)
// - filter.IsSuperToken=true → uses AdminDataQuery from repository/queries package
// - filter has Column+Value  → vendor-scoped query with WHERE clause
// If page <= 0 all rows are returned (no pagination).
func (r *DataRepository) GetAll(filter *VendorFilter, p QueryParams) ([]*models.DataRow, int, error) {
	var baseSelect string
	var conditions []string
	var args []interface{}
	paramIdx := 1

	if filter != nil && filter.IsSuperToken {
		baseSelect = queries.AdminDataQuery
	} else {
		baseSelect = vendorDataSelect
		if filter != nil && filter.Column != "" && filter.Value != "" {
			conditions = append(conditions, fmt.Sprintf("%s = @p%d", filter.Column, paramIdx))
			args = append(args, filter.Value)
			paramIdx++
		}
	}

	// Column filters
	if p.Status != "" {
		conditions = append(conditions, fmt.Sprintf("op.[Status] = @p%d", paramIdx))
		args = append(args, p.Status)
		paramIdx++
	}
	if p.Mode != "" {
		conditions = append(conditions, fmt.Sprintf("op.[Mode] = @p%d", paramIdx))
		args = append(args, p.Mode)
		paramIdx++
	}
	if p.Priority != "" {
		conditions = append(conditions, fmt.Sprintf("op.[Priority] = @p%d", paramIdx))
		args = append(args, p.Priority)
		paramIdx++
	}

	// Free-text search on terminal_id and terminal_name
	if p.Search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(op.[Terminal ID] LIKE @p%d OR op.[Terminal Name] LIKE @p%d)",
			paramIdx, paramIdx,
		))
		args = append(args, "%"+p.Search+"%")
		paramIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	orderBy := buildOrderBy(p)

	// Count query
	countQuery := "SELECT COUNT(*) FROM ticket_master.dbo.open_ticket op LEFT JOIN machine_master.dbo.machine mm ON op.[Terminal ID] = mm.[Terminal ID]"
	if whereClause != "" {
		countQuery += " " + whereClause
	}
	var total int
	if err := r.ticketDB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		r.logger.Errorf("Failed to count data rows: %v", err)
		return nil, 0, fmt.Errorf("failed to count rows: %w", err)
	}

	// Build data query
	query := baseSelect
	if whereClause != "" {
		query += "\n" + whereClause
	}

	var rows *sql.Rows
	var err error

	if p.Page > 0 && p.PageSize > 0 {
		offset := (p.Page - 1) * p.PageSize
		query += fmt.Sprintf("\n%s\nOFFSET @p%d ROWS FETCH NEXT @p%d ROWS ONLY", orderBy, paramIdx, paramIdx+1)
		rows, err = r.ticketDB.Query(query, append(args, offset, p.PageSize)...)
	} else {
		query += "\n" + orderBy
		rows, err = r.ticketDB.Query(query, args...)
	}

	if err != nil {
		r.logger.Errorf("Failed to query data: %v", err)
		return nil, 0, fmt.Errorf("failed to query data: %w", err)
	}
	defer rows.Close()

	result := make([]*models.DataRow, 0, p.PageSize)
	for rows.Next() {
		d, err := scanDataRow(rows)
		if err != nil {
			r.logger.Errorf("Failed to scan data row: %v", err)
			continue
		}
		result = append(result, d)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return result, total, nil
}

// GetByTerminalID retrieves a single row by terminal ID with optional vendor scoping.
func (r *DataRepository) GetByTerminalID(terminalID string, filter *VendorFilter) (*models.DataRow, error) {
	var query string
	var args []interface{}

	if filter != nil && filter.IsSuperToken {
		// Admin path: use customizable query + simple WHERE
		query = queries.AdminDataQuery + "\nWHERE op.[Terminal ID] = @p1"
		args = []interface{}{terminalID}
	} else if filter != nil && filter.Column != "" && filter.Value != "" {
		// Vendor path: vendor filter + terminal filter
		query = vendorDataSelect + fmt.Sprintf(
			"WHERE op.[Terminal ID] = @p1 AND %s = @p2", filter.Column,
		)
		args = []interface{}{terminalID, filter.Value}
	} else {
		// Unrestricted token (legacy or no filter set)
		query = vendorDataSelect + "WHERE op.[Terminal ID] = @p1"
		args = []interface{}{terminalID}
	}

	d, err := scanDataRow(r.ticketDB.QueryRow(query, args...))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("not found")
	}
	if err != nil {
		r.logger.Errorf("Failed to get row by terminal ID: %v", err)
		return nil, fmt.Errorf("failed to get row: %w", err)
	}
	return d, nil
}

// Update modifies ticket fields for a given terminal ID with vendor filter enforcement.
// For vendor-scoped tokens the UPDATE+JOIN pattern ensures 0 rows → 403 at handler level.
func (r *DataRepository) Update(terminalID string, req *models.DataUpdateRequest, filter *VendorFilter) (*models.DataRow, error) {
	updates := []string{}
	args := []interface{}{}
	p := 1

	add := func(col, val string) {
		if val != "" {
			updates = append(updates, fmt.Sprintf("[%s] = @p%d", col, p))
			args = append(args, val)
			p++
		}
	}

	add("Priority", req.Priority)
	add("Mode", req.Mode)
	add("Current Problem", req.CurrentProblem)
	add("Status", req.Status)
	add("Remarks", req.Remarks)
	add("Condition", req.Condition)
	add("Close time", req.CloseTime)
	add("Problem History", req.ProblemHistory)
	add("Mode History", req.ModeHistory)

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	var query string
	if filter != nil && !filter.IsSuperToken && filter.Column != "" {
		// Vendor-scoped: UPDATE via FROM+JOIN so vendor check is enforced at DB level
		args = append(args, terminalID, filter.Value)
		query = fmt.Sprintf(
			`UPDATE op SET %s
			 FROM ticket_master.dbo.open_ticket op
			 LEFT JOIN machine_master.dbo.machine mm ON op.[Terminal ID] = mm.[Terminal ID]
			 WHERE op.[Terminal ID] = @p%d AND %s = @p%d`,
			strings.Join(updates, ", "),
			p, filter.Column, p+1,
		)
	} else {
		// Admin / unrestricted token: simple UPDATE
		args = append(args, terminalID)
		query = fmt.Sprintf(
			"UPDATE ticket_master.dbo.open_ticket SET %s WHERE [Terminal ID] = @p%d",
			strings.Join(updates, ", "),
			p,
		)
	}

	result, err := r.ticketDB.Exec(query, args...)
	if err != nil {
		r.logger.Errorf("Failed to update: %v", err)
		return nil, fmt.Errorf("failed to update: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		if filter != nil && !filter.IsSuperToken {
			return nil, fmt.Errorf("not found or not accessible for this vendor")
		}
		return nil, fmt.Errorf("not found")
	}

	return r.GetByTerminalID(terminalID, filter)
}

// GetDistinctStatuses returns distinct Status values from open_ticket.
func (r *DataRepository) GetDistinctStatuses() ([]string, error) {
	rows, err := r.ticketDB.Query(`
		SELECT DISTINCT [Status] FROM ticket_master.dbo.open_ticket
		WHERE [Status] IS NOT NULL AND [Status] != '' ORDER BY [Status]
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err == nil {
			out = append(out, v)
		}
	}
	return out, nil
}

// GetDistinctModes returns distinct Mode values from open_ticket.
func (r *DataRepository) GetDistinctModes() ([]string, error) {
	rows, err := r.ticketDB.Query(`
		SELECT DISTINCT [Mode] FROM ticket_master.dbo.open_ticket
		WHERE [Mode] IS NOT NULL AND [Mode] != '' ORDER BY [Mode]
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err == nil {
			out = append(out, v)
		}
	}
	return out, nil
}

// GetDistinctPriorities returns distinct Priority values from open_ticket.
func (r *DataRepository) GetDistinctPriorities() ([]string, error) {
	rows, err := r.ticketDB.Query(`
		SELECT DISTINCT [Priority] FROM ticket_master.dbo.open_ticket
		WHERE [Priority] IS NOT NULL AND [Priority] != '' ORDER BY [Priority]
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err == nil {
			out = append(out, v)
		}
	}
	return out, nil
}
