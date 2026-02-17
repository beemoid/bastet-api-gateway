package models

import (
	"database/sql"
	"time"
)

// ATMI represents a terminal/machine record from machine_master.dbo.atmi
// This model contains detailed information about ATM terminals
type ATMI struct {
	ID            int            `json:"id" db:"id"`                         // Primary key
	TerminalID    string         `json:"terminal_id" db:"terminal_id"`       // Unique terminal identifier
	TerminalName  string         `json:"terminal_name" db:"terminal_name"`   // Friendly name of the terminal
	Location      string         `json:"location" db:"location"`             // Physical location/address
	BranchCode    string         `json:"branch_code" db:"branch_code"`       // Branch/site code
	IPAddress     string         `json:"ip_address" db:"ip_address"`         // Network IP address
	Model         string         `json:"model" db:"model"`                   // Terminal model/type
	Manufacturer  string         `json:"manufacturer" db:"manufacturer"`     // Equipment manufacturer
	SerialNumber  string         `json:"serial_number" db:"serial_number"`   // Hardware serial number
	Status        string         `json:"status" db:"status"`                 // Operational status: Active, Inactive, Maintenance, Offline
	LastPingTime  sql.NullTime   `json:"last_ping_time" db:"last_ping_time"` // Last communication timestamp (nullable)
	InstallDate   time.Time      `json:"install_date" db:"install_date"`     // Installation date
	WarrantyExp   sql.NullTime   `json:"warranty_exp" db:"warranty_exp"`     // Warranty expiration date (nullable)
	Notes         sql.NullString `json:"notes" db:"notes"`                   // Additional notes (nullable)
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`         // Record creation timestamp
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`         // Last update timestamp
}

// MachineStatusUpdate represents a status update for a terminal
// Used when receiving status updates from monitoring systems or cloud app
type MachineStatusUpdate struct {
	TerminalID   string `json:"terminal_id" binding:"required"`                               // Required: terminal to update
	Status       string `json:"status" binding:"required,oneof=Active Inactive Maintenance Offline"` // Required: new status
	LastPingTime string `json:"last_ping_time"`                                               // Optional: timestamp of last ping
	Notes        string `json:"notes"`                                                        // Optional: additional notes
}

// MachineResponse is the standardized response format for machine operations
type MachineResponse struct {
	Success bool   `json:"success"`        // Indicates if operation was successful
	Message string `json:"message"`        // Human-readable message
	Data    *ATMI  `json:"data,omitempty"` // Machine data (if applicable)
}

// MachineListResponse is the response format for listing multiple machines
type MachineListResponse struct {
	Success bool    `json:"success"`        // Indicates if operation was successful
	Message string  `json:"message"`        // Human-readable message
	Data    []*ATMI `json:"data,omitempty"` // Array of machines
	Total   int     `json:"total"`          // Total count of machines
}

// MachineFilter represents query parameters for filtering machines
// Used in list/search operations
type MachineFilter struct {
	Status     string `form:"status"`      // Filter by status
	BranchCode string `form:"branch_code"` // Filter by branch
	Location   string `form:"location"`    // Search by location
}
