package service

import "errors"

// Common service errors
// These errors provide standardized error messages across the service layer
var (
	ErrTicketAlreadyExists = errors.New("ticket with this number already exists")
	ErrTicketNotFound      = errors.New("ticket not found")
	ErrMachineNotFound     = errors.New("machine not found")
	ErrInvalidInput        = errors.New("invalid input data")
)
