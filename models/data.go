package models

// DataRow is the unified response row returned by GET /api/v1/data.
// It combines all ticket fields from ticket_master.dbo.open_ticket
// with machine dimension columns from machine_master.dbo.machine via a LEFT JOIN.
type DataRow struct {
	// ── Ticket fields ────────────────────────────────────────────
	TerminalID        string     `json:"terminal_id" example:"ATM-001"`
	TerminalName      string     `json:"terminal_name" example:"Main Branch ATM"`
	Priority          NullString `json:"priority" swaggertype:"string" example:"1.High"`
	Mode              NullString `json:"mode" swaggertype:"string" example:"Off-line"`
	InitialProblem    NullString `json:"initial_problem" swaggertype:"string" example:"Cash dispenser jam"`
	CurrentProblem    NullString `json:"current_problem" swaggertype:"string" example:"Card reader error"`
	PDuration         NullString `json:"p_duration" swaggertype:"string" example:"2h 30m"`
	IncidentStartTime NullString `json:"incident_start_datetime" swaggertype:"string" example:"2024-01-15 10:30:00"`
	Count             int        `json:"count" example:"5"`
	Status            NullString `json:"status" swaggertype:"string" example:"0.NEW"`
	Remarks           NullString `json:"remarks" swaggertype:"string" example:"Waiting for technician"`
	Balance           int        `json:"balance" example:"1000000"`
	Condition         NullString `json:"condition" swaggertype:"string" example:"Critical"`
	TicketsNo         NullString `json:"tickets_no" swaggertype:"string" example:"TKT-2024-001"`
	TicketsDuration   float64    `json:"tickets_duration" example:"150.5"`
	OpenTime          NullString `json:"open_time" swaggertype:"string" example:"2024-01-15 08:00:00"`
	CloseTime         NullString `json:"close_time" swaggertype:"string" example:"2024-01-15 18:00:00"`
	ProblemHistory    NullString `json:"problem_history" swaggertype:"string" example:"Card reader issue resolved"`
	ModeHistory       NullString `json:"mode_history" swaggertype:"string" example:"Online->Offline->Online"`
	DSPFLM            NullString `json:"dsp_flm" swaggertype:"string" example:"FLM-001"`
	DSPSLM            NullString `json:"dsp_slm" swaggertype:"string" example:"SLM-001"`
	LastWithdrawal    NullTime   `json:"last_withdrawal" swaggertype:"string" example:"2024-01-15T09:30:00Z"`
	ExportName        NullString `json:"export_name" swaggertype:"string" example:"ATM_Report_Jan2024"`

	// ── Machine dimension fields (LEFT JOIN machine_master.dbo.machine) ──
	FLMName NullString `json:"flm_name" swaggertype:"string" example:"AVT"`        // mm.[FLM name]
	FLM     NullString `json:"flm" swaggertype:"string" example:"AVT - BANDUNG"`   // mm.[FLM]
	SLM     NullString `json:"slm" swaggertype:"string" example:"KGP - WINCOR DW"` // mm.[SLM]
	Net     NullString `json:"net" swaggertype:"string" example:"NOSAIRIS"`        // mm.[Net]
}

// DataUpdateRequest represents updatable ticket fields sent in PUT /api/v1/data/:terminal_id
type DataUpdateRequest struct {
	Priority       string `json:"priority" example:"1.High"`
	Mode           string `json:"mode" example:"Off-line"`
	CurrentProblem string `json:"current_problem" example:"Card reader fixed"`
	Status         string `json:"status" example:"2.Kirim FLM"`
	Remarks        string `json:"remarks" example:"Technician dispatched"`
	Condition      string `json:"condition" example:"Normal"`
	CloseTime      string `json:"close_time" example:"2024-01-15 18:00:00"`
	ProblemHistory string `json:"problem_history" example:"Card reader issue resolved"`
	ModeHistory    string `json:"mode_history" example:"Online->Offline->Online"`
}

// DataResponse is the standardized single-row response
type DataResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    *DataRow `json:"data,omitempty"`
}

// DataListResponse is the standardized list response
type DataListResponse struct {
	Success    bool       `json:"success"`
	Message    string     `json:"message"`
	Data       []*DataRow `json:"data,omitempty"`
	Total      int        `json:"total"`
	Page       int        `json:"page,omitempty"`
	PageSize   int        `json:"page_size,omitempty"`
	TotalPages int        `json:"total_pages,omitempty"`

	SortBy    string `json:"sort_by,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
	Search    string `json:"search,omitempty"`
	Status    string `json:"status,omitempty"`
	Mode      string `json:"mode,omitempty"`
	Priority  string `json:"priority,omitempty"`
}
