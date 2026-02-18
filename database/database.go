package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/microsoft/go-mssqldb" // SQL Server driver
	"github.com/sirupsen/logrus"
)

// DBManager manages multiple database connections
// Holds separate connections for ticket, machine, and token databases
type DBManager struct {
	TicketDB  *sql.DB // Connection to ticket_master database
	MachineDB *sql.DB // Connection to machine_master database
	TokenDB   *sql.DB // Connection to token_management database
	logger    *logrus.Logger
}

// NewDBManager creates a new database manager with connections to all databases
// Connections are non-fatal: if a database is unavailable at startup, the app
// keeps running and the connection will succeed automatically once the database
// becomes available. Health endpoint reports real-time status.
func NewDBManager(ticketDSN, machineDSN, tokenDSN string, logger *logrus.Logger) *DBManager {
	manager := &DBManager{
		logger: logger,
	}

	manager.TicketDB = openDB(ticketDSN, "ticket_master", logger)
	manager.MachineDB = openDB(machineDSN, "machine_master", logger)

	if tokenDSN != "" {
		manager.TokenDB = openDB(tokenDSN, "token_management", logger)
	} else {
		logger.Warn("Token database DSN not configured, token management will be unavailable")
	}

	return manager
}

// openDB opens a database connection, configures the pool, and pings.
// Always returns the *sql.DB even if ping fails — Go's database/sql
// will automatically reconnect when the database becomes available.
func openDB(dsn, name string, logger *logrus.Logger) *sql.DB {
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		logger.Warnf("Failed to open %s database: %v", name, err)
		return nil
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		logger.Warnf("%s database not available at startup: %v (will retry automatically)", name, err)
	} else {
		logger.Infof("Successfully connected to %s database", name)
	}

	return db
}

// Close gracefully closes all database connections
// Should be called when the application shuts down
func (dm *DBManager) Close() error {
	var firstErr error

	if dm.TicketDB != nil {
		if err := dm.TicketDB.Close(); err != nil {
			dm.logger.Errorf("Error closing ticket database: %v", err)
			if firstErr == nil {
				firstErr = err
			}
		} else {
			dm.logger.Info("Ticket database connection closed")
		}
	}

	if dm.MachineDB != nil {
		if err := dm.MachineDB.Close(); err != nil {
			dm.logger.Errorf("Error closing machine database: %v", err)
			if firstErr == nil {
				firstErr = err
			}
		} else {
			dm.logger.Info("Machine database connection closed")
		}
	}

	if dm.TokenDB != nil {
		if err := dm.TokenDB.Close(); err != nil {
			dm.logger.Errorf("Error closing token database: %v", err)
			if firstErr == nil {
				firstErr = err
			}
		} else {
			dm.logger.Info("Token database connection closed")
		}
	}

	return firstErr
}

// DatabaseHealth holds the real-time health status of each database connection
type DatabaseHealth struct {
	TicketDB  string `json:"ticket_database"`
	MachineDB string `json:"machine_database"`
	TokenDB   string `json:"token_database"`
}

// HealthCheck pings all databases in real-time and returns their current status.
// If a database was unavailable at startup but has since come online, this will
// reflect the updated status.
func (dm *DBManager) HealthCheck() *DatabaseHealth {
	return &DatabaseHealth{
		TicketDB:  checkDB(dm.TicketDB),
		MachineDB: checkDB(dm.MachineDB),
		TokenDB:   checkDB(dm.TokenDB),
	}
}

// checkDB pings a single database and returns its status string
func checkDB(db *sql.DB) string {
	if db == nil {
		return "not configured"
	}
	if err := db.Ping(); err != nil {
		return fmt.Sprintf("disconnected: %v", err)
	}
	return "connected"
}
