package models

// ATMI represents a terminal/machine record from machine_master.dbo.atmi
// This model contains detailed information about ATM terminals
type ATMI struct {
	TerminalID        string   `json:"terminal_id" db:"terminal_id" example:"ATM-001"`                             // Unique terminal identifier
	Store             string   `json:"store" db:"store" example:"Main Branch"`                                     // Store name
	StoreCode         string   `json:"store_code" db:"store_code" example:"STR-001"`                               // Store code
	StoreName         string   `json:"store_name" db:"store_name" example:"Jakarta Main Branch"`                   // Full store name
	DateOfActivation  NullTime `json:"date_of_activation" db:"date_of_activation" swaggertype:"string" example:"2023-01-15T00:00:00Z"` // Activation date (nullable)
	Status            string   `json:"status" db:"status" example:"Active"`                                        // Operational status
	Std               int      `json:"std" db:"std" example:"1"`                                                   // Standard value
	GPS               string   `json:"gps" db:"gps" example:"-6.200000,106.816666"`                                // GPS coordinates
	Lat               float64  `json:"lat" db:"lat" example:"-6.200000"`                                           // Latitude
	Lon               float64  `json:"lon" db:"lon" example:"106.816666"`                                          // Longitude
	Province          string   `json:"province" db:"province" example:"DKI Jakarta"`                               // Province name
	CityRegency       string   `json:"city_regency" db:"city/regency" example:"Jakarta Pusat"`                     // City or regency name
	District          string   `json:"district" db:"district" example:"Menteng"`                                   // District name
}

// MachineStatusUpdate represents a status update for a terminal
// Used when receiving status updates from monitoring systems or cloud app
type MachineStatusUpdate struct {
	TerminalID string  `json:"terminal_id" binding:"required" example:"ATM-001"` // Required: terminal to update
	Status     string  `json:"status" binding:"required" example:"Active"`       // Required: new status
	GPS        string  `json:"gps" example:"-6.200000,106.816666"`               // Optional: GPS coordinates
	Lat        float64 `json:"lat" example:"-6.200000"`                          // Optional: latitude
	Lon        float64 `json:"lon" example:"106.816666"`                         // Optional: longitude
}

// MachineResponse is the standardized response format for machine operations
type MachineResponse struct {
	Success bool   `json:"success"`        // Indicates if operation was successful
	Message string `json:"message"`        // Human-readable message
	Data    *ATMI  `json:"data,omitempty"` // Machine data (if applicable)
}

// MachineListResponse is the response format for listing multiple machines
type MachineListResponse struct {
	Success    bool    `json:"success"`                  // Indicates if operation was successful
	Message    string  `json:"message"`                  // Human-readable message
	Data       []*ATMI `json:"data,omitempty"`           // Array of machines
	Total      int     `json:"total"`                    // Total count of machines
	Page       int     `json:"page,omitempty"`           // Current page number
	PageSize   int     `json:"page_size,omitempty"`      // Items per page
	TotalPages int     `json:"total_pages,omitempty"`    // Total number of pages
}

// MachineFilter represents query parameters for filtering machines
// Used in list/search operations
type MachineFilter struct {
	Status      string `form:"status"`       // Filter by status
	StoreCode   string `form:"store_code"`   // Filter by store code
	Province    string `form:"province"`     // Filter by province
	CityRegency string `form:"city_regency"` // Filter by city/regency
	District    string `form:"district"`     // Filter by district
}
