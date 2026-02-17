package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
// It includes server settings, database connections, and external service configurations
type Config struct {
	Server      ServerConfig
	TicketDB    DatabaseConfig
	MachineDB   DatabaseConfig
	CloudApp    CloudAppConfig
	Security    SecurityConfig
}

// ServerConfig contains server-related configuration
type ServerConfig struct {
	Port    string // Port number for the API server
	GinMode string // Gin framework mode: debug, release, or test
}

// DatabaseConfig holds database connection parameters
type DatabaseConfig struct {
	Host     string // Database server hostname or IP
	Port     string // Database server port
	User     string // Database username
	Password string // Database password
	Database string // Database name
}

// CloudAppConfig contains configuration for the cloud application
type CloudAppConfig struct {
	URL    string // Base URL of the cloud application
	APIKey string // API key for authentication with cloud app
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	JWTSecret string // Secret key for JWT token generation/validation
	APIKey    string // Internal API key for securing endpoints
}

// Load reads configuration from environment variables
// It first loads the .env file, then populates the Config struct
// Returns error if required environment variables are missing
func Load() (*Config, error) {
	// Load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port:    getEnv("SERVER_PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		TicketDB: DatabaseConfig{
			Host:     getEnv("TICKET_DB_HOST", "localhost"),
			Port:     getEnv("TICKET_DB_PORT", "1433"),
			User:     getEnv("TICKET_DB_USER", ""),
			Password: getEnv("TICKET_DB_PASSWORD", ""),
			Database: getEnv("TICKET_DB_NAME", "ticket_master"),
		},
		MachineDB: DatabaseConfig{
			Host:     getEnv("MACHINE_DB_HOST", "localhost"),
			Port:     getEnv("MACHINE_DB_PORT", "1433"),
			User:     getEnv("MACHINE_DB_USER", ""),
			Password: getEnv("MACHINE_DB_PASSWORD", ""),
			Database: getEnv("MACHINE_DB_NAME", "machine_master"),
		},
		CloudApp: CloudAppConfig{
			URL:    getEnv("CLOUD_APP_URL", ""),
			APIKey: getEnv("CLOUD_APP_API_KEY", ""),
		},
		Security: SecurityConfig{
			JWTSecret: getEnv("JWT_SECRET", "default-secret-change-in-production"),
			APIKey:    getEnv("API_KEY", ""),
		},
	}

	return config, nil
}

// GetDSN generates a connection string for SQL Server
// Format: sqlserver://username:password@host:port?database=dbname
func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Database,
	)
}

// getEnv retrieves an environment variable value or returns a default value
// This helper function prevents nil pointer errors when env vars are not set
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
