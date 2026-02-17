package models

import (
	"database/sql"
	"time"
)

// OpenTicket represents a ticket record from ticket_master.dbo.open_ticket
// This model maps to the database table structure for open/active tickets
type OpenTicket struct {
	ID              int            `json:"id" db:"id"`                           // Primary key
	TicketNumber    string         `json:"ticket_number" db:"ticket_number"`     // Unique ticket identifier
	TerminalID      string         `json:"terminal_id" db:"terminal_id"`         // Reference to ATMI terminal
	Description     string         `json:"description" db:"description"`         // Issue description
	Priority        string         `json:"priority" db:"priority"`               // Priority level: Low, Medium, High, Critical
	Status          string         `json:"status" db:"status"`                   // Current status: Open, In Progress, Pending, Resolved
	Category        string         `json:"category" db:"category"`               // Ticket category
	ReportedBy      string         `json:"reported_by" db:"reported_by"`         // User who reported the issue
	AssignedTo      sql.NullString `json:"assigned_to" db:"assigned_to"`         // Technician assigned to ticket (nullable)
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`           // Ticket creation timestamp
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`           // Last update timestamp
	ResolvedAt      sql.NullTime   `json:"resolved_at" db:"resolved_at"`         // Resolution timestamp (nullable)
	ResolutionNotes sql.NullString `json:"resolution_notes" db:"resolution_notes"` // Notes about resolution (nullable)
}

// TicketCreateRequest represents the payload for creating a new ticket
// Used when receiving ticket creation requests from the cloud app
type TicketCreateRequest struct {
	TicketNumber string `json:"ticket_number" binding:"required"`     // Required: unique ticket number
	TerminalID   string `json:"terminal_id" binding:"required"`       // Required: associated terminal
	Description  string `json:"description" binding:"required"`       // Required: issue description
	Priority     string `json:"priority" binding:"required,oneof=Low Medium High Critical"` // Required: must be one of the valid priorities
	Category     string `json:"category" binding:"required"`          // Required: ticket category
	ReportedBy   string `json:"reported_by" binding:"required"`       // Required: reporter name
	AssignedTo   string `json:"assigned_to"`                          // Optional: assigned technician
}

// TicketUpdateRequest represents the payload for updating an existing ticket
// Used when the cloud app sends ticket updates
type TicketUpdateRequest struct {
	Status          string `json:"status" binding:"omitempty,oneof=Open InProgress Pending Resolved"` // Optional: new status
	Priority        string `json:"priority" binding:"omitempty,oneof=Low Medium High Critical"`       // Optional: new priority
	AssignedTo      string `json:"assigned_to"`                                                       // Optional: reassign ticket
	Description     string `json:"description"`                                                       // Optional: update description
	ResolutionNotes string `json:"resolution_notes"`                                                  // Optional: add resolution notes
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
	Success bool          `json:"success"`        // Indicates if operation was successful
	Message string        `json:"message"`        // Human-readable message
	Data    []*OpenTicket `json:"data,omitempty"` // Array of tickets
	Total   int           `json:"total"`          // Total count of tickets
}
