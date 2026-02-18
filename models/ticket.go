package models

// OpenTicket represents a ticket record from ticket_master.dbo.open_ticket
// This model maps to the database table structure for open/active tickets
type OpenTicket struct {
	TerminalID         string     `json:"terminal_id" db:"Terminal ID" example:"ATM-001"`                                       // Terminal identifier
	TerminalName       string     `json:"terminal_name" db:"Terminal Name" example:"Main Branch ATM"`                           // Terminal name
	Priority           NullString `json:"priority" db:"Priority" swaggertype:"string" example:"1.High"`                        // Priority level (nullable): 1.High, 2.Middle, 3.Low, 4.Minimum
	Mode               NullString `json:"mode" db:"Mode" swaggertype:"string" example:"Off-line"`                              // Terminal mode (nullable): Closed, In Service, nan, Off-line, Supervisor
	InitialProblem     NullString `json:"initial_problem" db:"Initial Problem" swaggertype:"string" example:"Cash dispenser jam"` // Initial problem description (nullable)
	CurrentProblem     NullString `json:"current_problem" db:"Current Problem" swaggertype:"string" example:"Card reader error"` // Current problem description (nullable)
	PDuration          NullString `json:"p_duration" db:"P-Duration" swaggertype:"string" example:"2h 30m"`                     // Problem duration (nullable)
	IncidentStartTime  NullString `json:"incident_start_datetime" db:"Incident start datetime" swaggertype:"string" example:"2024-01-15 10:30:00"` // Incident start timestamp (nullable)
	Count              int        `json:"count" db:"Count" example:"5"`                                                         // Count value
	Status             NullString `json:"status" db:"Status" swaggertype:"string" example:"0.NEW"`                             // Ticket status (nullable): 0.NEW, 1.Req FD ke HD, 2.Kirim FLM, etc.
	Remarks            NullString `json:"remarks" db:"Remarks" swaggertype:"string" example:"Waiting for technician"`           // Remarks/notes (nullable)
	Balance            int        `json:"balance" db:"Balance" example:"1000000"`                                               // Balance amount
	Condition          NullString `json:"condition" db:"Condition" swaggertype:"string" example:"Critical"`                     // Condition status (nullable)
	TicketsNo          NullString `json:"tickets_no" db:"Tickets no" swaggertype:"string" example:"TKT-2024-001"`              // Ticket number (nullable)
	TicketsDuration    float64    `json:"tickets_duration" db:"Tickets duration" example:"150.5"`                               // Ticket duration in minutes (float)
	OpenTime           NullString `json:"open_time" db:"Open time" swaggertype:"string" example:"2024-01-15 08:00:00"`         // Ticket open time (nullable)
	CloseTime          NullString `json:"close_time" db:"Close time" swaggertype:"string" example:"2024-01-15 18:00:00"`       // Ticket close time (nullable)
	ProblemHistory     NullString `json:"problem_history" db:"Problem History" swaggertype:"string" example:"Card reader issue resolved"` // Problem history (nullable)
	ModeHistory        NullString `json:"mode_history" db:"Mode History" swaggertype:"string" example:"Online->Offline->Online"` // Mode change history (nullable)
	DSPFLM             NullString `json:"dsp_flm" db:"DSP FLM" swaggertype:"string" example:"FLM-001"`                         // DSP FLM identifier (nullable)
	DSPSLM             NullString `json:"dsp_slm" db:"DSP SLM" swaggertype:"string" example:"SLM-001"`                         // DSP SLM identifier (nullable)
	LastWithdrawal     NullTime   `json:"last_withdrawal" db:"Last Withdrawal" swaggertype:"string" example:"2024-01-15T09:30:00Z"` // Last withdrawal timestamp (nullable)
	ExportName         NullString `json:"export_name" db:"Export Name" swaggertype:"string" example:"ATM_Report_Jan2024"`      // Export name for reports (nullable)
}

// TicketCreateRequest represents the payload for creating a new ticket
// Used when receiving ticket creation requests from the cloud app
type TicketCreateRequest struct {
	TerminalID        string `json:"terminal_id" binding:"required" example:"ATM-001"`                    // Required: terminal identifier
	TerminalName      string `json:"terminal_name" binding:"required" example:"Main Branch ATM"`          // Required: terminal name
	Priority          string `json:"priority" binding:"required" example:"1.High"`                        // Required: priority level (1.High, 2.Middle, 3.Low, 4.Minimum)
	Mode              string `json:"mode" binding:"required" example:"Off-line"`                          // Required: terminal mode (Closed, In Service, nan, Off-line, Supervisor)
	InitialProblem    string `json:"initial_problem" binding:"required" example:"Cash dispenser jam"`    // Required: initial problem
	CurrentProblem    string `json:"current_problem" example:"Card reader error"`                        // Optional: current problem
	PDuration         string `json:"p_duration" example:"2h 30m"`                                         // Optional: problem duration
	IncidentStartTime string `json:"incident_start_datetime" example:"2024-01-15 10:30:00"`              // Optional: incident start time
	Status            string `json:"status" example:"0.NEW"`                                              // Optional: status (0.NEW, 1.Req FD ke HD, etc.)
	Remarks           string `json:"remarks" example:"Waiting for technician"`                            // Optional: remarks
	Condition         string `json:"condition" example:"Critical"`                                        // Optional: condition
	TicketsNo         string `json:"tickets_no" example:"TKT-2024-001"`                                   // Optional: ticket number
	ExportName        string `json:"export_name" example:"ATM_Report_Jan2024"`                            // Optional: export name
}

// TicketUpdateRequest represents the payload for updating an existing ticket
// Used when the cloud app sends ticket updates
type TicketUpdateRequest struct {
	Priority       string `json:"priority" example:"1.High"`                            // Optional: update priority (1.High, 2.Middle, 3.Low, 4.Minimum)
	Mode           string `json:"mode" example:"Off-line"`                              // Optional: update mode (Closed, In Service, nan, Off-line, Supervisor)
	CurrentProblem string `json:"current_problem" example:"Card reader fixed"`          // Optional: update current problem
	Status         string `json:"status" example:"0.NEW"`                               // Optional: update status (0.NEW, 1.Req FD ke HD, etc.)
	Remarks        string `json:"remarks" example:"Issue resolved"`                     // Optional: update remarks
	Condition      string `json:"condition" example:"Normal"`                           // Optional: update condition
	CloseTime      string `json:"close_time" example:"2024-01-15 18:00:00"`             // Optional: set close time
	ProblemHistory string `json:"problem_history" example:"Card reader issue resolved"` // Optional: update problem history
	ModeHistory    string `json:"mode_history" example:"Online->Offline->Online"`       // Optional: update mode history
}

// TicketResponse is the standardized response format for ticket operations
// Includes both the ticket data and metadata about the response
type TicketResponse struct {
	Success bool        `json:"success"`         // Indicates if operation was successful
	Message string      `json:"message"`         // Human-readable message
	Data    *OpenTicket `json:"data,omitempty"`  // Ticket data (if applicable)
}

// TicketListResponse is the response format for listing multiple tickets
type TicketListResponse struct {
	Success    bool          `json:"success"`                  // Indicates if operation was successful
	Message    string        `json:"message"`                  // Human-readable message
	Data       []*OpenTicket `json:"data,omitempty"`           // Array of tickets
	Total      int           `json:"total"`                    // Total count of tickets
	Page       int           `json:"page,omitempty"`           // Current page number
	PageSize   int           `json:"page_size,omitempty"`      // Items per page
	TotalPages int           `json:"total_pages,omitempty"`    // Total number of pages
}

// ErrorResponse is the standardized error response format
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"` // Always false for errors
	Message string `json:"message" example:"Error message describing what went wrong"` // Error message
	Error   string `json:"error,omitempty" example:"detailed error information"` // Optional detailed error
}
