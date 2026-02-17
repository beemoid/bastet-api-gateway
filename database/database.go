package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/microsoft/go-mssqldb" // SQL Server driver
	"github.com/sirupsen/logrus"
)

// DBManager manages multiple database connections
// Holds separate connections for ticket and machine databases
type DBManager struct {
	TicketDB  *sql.DB // Connection to ticket_master database
	MachineDB *sql.DB // Connection to machine_master database
	logger    *logrus.Logger
}

// NewDBManager creates a new database manager with connections to both databases
// Parameters:
//   - ticketDSN: Connection string for ticket database
//   - machineDSN: Connection string for machine database
//   - logger: Logger instance for database operations
// Returns error if any database connection fails
func NewDBManager(ticketDSN, machineDSN string, logger *logrus.Logger) (*DBManager, error) {
	manager := &DBManager{
		logger: logger,
	}

	// Connect to Ticket Database
	ticketDB, err := sql.Open("sqlserver", ticketDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ticket database: %w", err)
	}

	// Configure ticket database connection pool
	ticketDB.SetMaxOpenConns(25)                 // Maximum number of open connections
	ticketDB.SetMaxIdleConns(5)                  // Maximum number of idle connections
	ticketDB.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of a connection

	// Verify ticket database connection
	if err := ticketDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping ticket database: %w", err)
	}

	manager.TicketDB = ticketDB
	logger.Info("Successfully connected to ticket_master database")

	// Connect to Machine Database
	machineDB, err := sql.Open("sqlserver", machineDSN)
	if err != nil {
		ticketDB.Close() // Clean up ticket connection
		return nil, fmt.Errorf("failed to connect to machine database: %w", err)
	}

	// Configure machine database connection pool
	machineDB.SetMaxOpenConns(25)
	machineDB.SetMaxIdleConns(5)
	machineDB.SetConnMaxLifetime(5 * time.Minute)

	// Verify machine database connection
	if err := machineDB.Ping(); err != nil {
		ticketDB.Close()  // Clean up ticket connection
		machineDB.Close() // Clean up machine connection
		return nil, fmt.Errorf("failed to ping machine database: %w", err)
	}

	manager.MachineDB = machineDB
	logger.Info("Successfully connected to machine_master database")

	return manager, nil
}

// Close gracefully closes all database connections
// Should be called when the application shuts down
func (dm *DBManager) Close() error {
	var ticketErr, machineErr error

	// Close ticket database connection
	if dm.TicketDB != nil {
		ticketErr = dm.TicketDB.Close()
		if ticketErr != nil {
			dm.logger.Errorf("Error closing ticket database: %v", ticketErr)
		} else {
			dm.logger.Info("Ticket database connection closed")
		}
	}

	// Close machine database connection
	if dm.MachineDB != nil {
		machineErr = dm.MachineDB.Close()
		if machineErr != nil {
			dm.logger.Errorf("Error closing machine database: %v", machineErr)
		} else {
			dm.logger.Info("Machine database connection closed")
		}
	}

	// Return first error encountered, if any
	if ticketErr != nil {
		return ticketErr
	}
	return machineErr
}

// HealthCheck verifies that both database connections are alive
// Returns error if either database is unreachable
// Used by health check endpoint to monitor database status
func (dm *DBManager) HealthCheck() error {
	// Check ticket database
	if err := dm.TicketDB.Ping(); err != nil {
		return fmt.Errorf("ticket database health check failed: %w", err)
	}

	// Check machine database
	if err := dm.MachineDB.Ping(); err != nil {
		return fmt.Errorf("machine database health check failed: %w", err)
	}

	return nil
}
