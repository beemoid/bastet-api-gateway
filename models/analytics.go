package models

// DashboardStats provides comprehensive statistics for the monitoring dashboard
type DashboardStats struct {
	Success bool              `json:"success" example:"true"`
	Message string            `json:"message" example:"Dashboard statistics retrieved successfully"`
	Data    DashboardStatsData `json:"data"`
}

// DashboardStatsData contains all dashboard statistics
type DashboardStatsData struct {
	Overview        OverviewStats       `json:"overview"`
	TicketStats     TicketStatistics    `json:"ticket_stats"`
	MachineStats    MachineStatistics   `json:"machine_stats"`
	MaintenanceStats MaintenanceStats   `json:"maintenance_stats"`
	GeographicStats GeographicStats     `json:"geographic_stats"`
}

// OverviewStats provides high-level overview numbers
type OverviewStats struct {
	TotalMachines       int     `json:"total_machines" example:"1500"`
	ActiveMachines      int     `json:"active_machines" example:"1350"`
	TotalOpenTickets    int     `json:"total_open_tickets" example:"45"`
	CriticalTickets     int     `json:"critical_tickets" example:"12"`
	MachineAvailability float64 `json:"machine_availability_percent" example:"90.5"`
}

// TicketStatistics provides ticket-related statistics
type TicketStatistics struct {
	ByStatus   []StatusCount   `json:"by_status"`
	ByPriority []PriorityCount `json:"by_priority"`
	ByMode     []ModeCount     `json:"by_mode"`
	AvgDuration float64        `json:"avg_duration_minutes" example:"125.5"`
	TotalCount  int            `json:"total_count" example:"45"`
}

// MachineStatistics provides machine-related statistics
type MachineStatistics struct {
	ByStatus     []MachineStatusCount `json:"by_status"`
	ByProvince   []ProvinceCount      `json:"by_province"`
	BySLM        []SLMCount           `json:"by_slm"`
	ByFLMName    []FLMNameCount       `json:"by_flm_name"`
	ByNetwork    []NetworkCount       `json:"by_network"`
	TotalCount   int                  `json:"total_count" example:"1500"`
}

// MaintenanceStats provides maintenance workload statistics
type MaintenanceStats struct {
	ByFLMProvider []FLMWorkloadCount `json:"by_flm_provider"`
	BySLMProvider []SLMWorkloadCount `json:"by_slm_provider"`
	TopBusyAreas  []AreaWorkloadCount `json:"top_busy_areas"`
}

// GeographicStats provides location-based statistics
type GeographicStats struct {
	ByProvince []ProvinceStats `json:"by_province"`
	ByCity     []CityStats     `json:"by_city"`
}

// StatusCount represents count by ticket status
type StatusCount struct {
	Status string `json:"status" example:"0.NEW"`
	Count  int    `json:"count" example:"15"`
}

// PriorityCount represents count by ticket priority
type PriorityCount struct {
	Priority string `json:"priority" example:"1.High"`
	Count    int    `json:"count" example:"12"`
}

// ModeCount represents count by terminal mode
type ModeCount struct {
	Mode  string `json:"mode" example:"Off-line"`
	Count int    `json:"count" example:"8"`
}

// MachineStatusCount represents count by machine status
type MachineStatusCount struct {
	Status string `json:"status" example:"Active"`
	Count  int    `json:"count" example:"1350"`
}

// ProvinceCount represents count by province
type ProvinceCount struct {
	Province string `json:"province" example:"DKI Jakarta"`
	Count    int    `json:"count" example:"250"`
}

// SLMCount represents count by SLM provider
type SLMCount struct {
	SLM   string `json:"slm" example:"KGP - WINCOR DW"`
	Count int    `json:"count" example:"120"`
}

// FLMNameCount represents count by FLM name
type FLMNameCount struct {
	FLMName string `json:"flm_name" example:"AVT"`
	Count   int    `json:"count" example:"800"`
}

// NetworkCount represents count by network provider
type NetworkCount struct {
	Network string `json:"network" example:"NOSAIRIS"`
	Count   int    `json:"count" example:"500"`
}

// FLMWorkloadCount represents FLM workload (machines + tickets)
type FLMWorkloadCount struct {
	FLM            string `json:"flm" example:"AVT - BANDUNG"`
	Area           string `json:"area" example:"BANDUNG"`
	MachineCount   int    `json:"machine_count" example:"45"`
	OpenTickets    int    `json:"open_tickets" example:"5"`
	WorkloadScore  int    `json:"workload_score" example:"50"` // MachineCount + (OpenTickets * 2)
}

// SLMWorkloadCount represents SLM workload
type SLMWorkloadCount struct {
	SLM          string `json:"slm" example:"KGP - WINCOR DW"`
	MachineCount int    `json:"machine_count" example:"120"`
	OpenTickets  int    `json:"open_tickets" example:"8"`
}

// AreaWorkloadCount represents workload by area
type AreaWorkloadCount struct {
	Area         string `json:"area" example:"BANDUNG"`
	MachineCount int    `json:"machine_count" example:"85"`
	OpenTickets  int    `json:"open_tickets" example:"7"`
}

// ProvinceStats represents detailed province statistics
type ProvinceStats struct {
	Province       string  `json:"province" example:"DKI Jakarta"`
	MachineCount   int     `json:"machine_count" example:"250"`
	ActiveMachines int     `json:"active_machines" example:"230"`
	OpenTickets    int     `json:"open_tickets" example:"12"`
	Availability   float64 `json:"availability_percent" example:"92.0"`
}

// CityStats represents detailed city statistics
type CityStats struct {
	City           string  `json:"city" example:"Jakarta Pusat"`
	MachineCount   int     `json:"machine_count" example:"85"`
	ActiveMachines int     `json:"active_machines" example:"78"`
	OpenTickets    int     `json:"open_tickets" example:"4"`
	Availability   float64 `json:"availability_percent" example:"91.8"`
}

// TerminalWithTicket combines machine and ticket information
type TerminalWithTicket struct {
	Machine *ATMI        `json:"machine"`          // Machine/terminal details
	Ticket  *OpenTicket  `json:"ticket,omitempty"` // Associated ticket (if any)
	HasTicket bool       `json:"has_ticket"`       // Whether terminal has an open ticket
}

// TerminalWithTicketResponse is the response format
type TerminalWithTicketResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Data    []*TerminalWithTicket `json:"data,omitempty"`
	Total   int                   `json:"total"`
}

// FLMWorkloadResponse provides FLM workload analysis
type FLMWorkloadResponse struct {
	Success bool               `json:"success" example:"true"`
	Message string             `json:"message" example:"FLM workload retrieved successfully"`
	Data    []FLMWorkloadCount `json:"data"`
	Total   int                `json:"total" example:"58"`
}

// AreaAnalysisResponse provides area-based analysis
type AreaAnalysisResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Area analysis retrieved successfully"`
	Data    []AreaStats `json:"data"`
	Total   int         `json:"total" example:"50"`
}

// AreaStats represents comprehensive area statistics
type AreaStats struct {
	Area           string   `json:"area" example:"BANDUNG"`
	Province       string   `json:"province" example:"Jawa Barat"`
	MachineCount   int      `json:"machine_count" example:"85"`
	ActiveMachines int      `json:"active_machines" example:"78"`
	OpenTickets    int      `json:"open_tickets" example:"7"`
	Availability   float64  `json:"availability_percent" example:"91.8"`
	FLMProviders   []string `json:"flm_providers" example:"AVT - BANDUNG,BRS - BANDUNG"`
	TopIssues      []string `json:"top_issues,omitempty" example:"Card reader error,Cash dispenser jam"`
}

// TicketsByFLMResponse provides tickets grouped by FLM
type TicketsByFLMResponse struct {
	Success bool              `json:"success" example:"true"`
	Message string            `json:"message" example:"Tickets by FLM retrieved successfully"`
	Data    []FLMTicketsGroup `json:"data"`
	Total   int               `json:"total" example:"58"`
}

// FLMTicketsGroup groups tickets by FLM provider
type FLMTicketsGroup struct {
	FLM     string         `json:"flm" example:"AVT - BANDUNG"`
	Area    string         `json:"area" example:"BANDUNG"`
	Tickets []*OpenTicket  `json:"tickets"`
	Count   int            `json:"count" example:"5"`
}

// CriticalTerminalsResponse provides list of critical terminals
type CriticalTerminalsResponse struct {
	Success bool                `json:"success" example:"true"`
	Message string              `json:"message" example:"Critical terminals retrieved successfully"`
	Data    []CriticalTerminal  `json:"data"`
	Total   int                 `json:"total" example:"12"`
}

// CriticalTerminal represents a terminal with critical issues
type CriticalTerminal struct {
	TerminalID      string  `json:"terminal_id" example:"ATM-001"`
	TerminalName    string  `json:"terminal_name" example:"Main Branch ATM"`
	Location        string  `json:"location" example:"DKI Jakarta - Jakarta Pusat"`
	Status          string  `json:"status" example:"Off-line"`
	TicketStatus    string  `json:"ticket_status" example:"0.NEW"`
	Priority        string  `json:"priority" example:"1.High"`
	Duration        float64 `json:"duration_minutes" example:"758.8"`
	Problem         string  `json:"problem" example:"Card reader error"`
	FLM             string  `json:"flm" example:"AVT - CIDENG"`
	SLM             string  `json:"slm" example:"KGP - WINCOR DW"`
	GPS             string  `json:"gps" example:"-6.200000,106.816666"`
}
