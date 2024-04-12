// Package database provides functionality for establishing and managing connections
// to a PostgreSQL database, as well as handling tasks related to database operations.
package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/emso-c/konzek-go-assignment/src/modules/logger"
	_ "github.com/lib/pq"
)

// DatabaseConnection represents a configuration for establishing a database connection.
type DatabaseConnection struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
	SSLMode  string
}

// _NewDatabaseConnection creates a new DatabaseConnection instance based on environment variables.
func _NewDatabaseConnection() *DatabaseConnection {
	return &DatabaseConnection{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Dbname:   os.Getenv("POSTGRES_DB"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}
}

// Connect establishes a connection to the PostgreSQL database using the configured parameters.
func (d *DatabaseConnection) Connect() *sql.DB {
	logger := logger.GetLogger()
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Dbname,
		d.SSLMode,
	)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		errStr := fmt.Sprintf("Error connecting to database: %s", err)
		logger.Fatal(errStr)
	}

	// Create tasks table if it does not exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		errStr := fmt.Sprintf("Error creating tasks table: %s", err)
		logger.Fatal(errStr)
	}

	logger.Info("Successfully connected to database")
	return db
}

var dbInstance *sql.DB = nil

// GetDatabase returns a singleton instance of the database connection.
// It establishes the connection based on the environment variables if it has not been established yet.
func GetDatabase() *sql.DB {
	if dbInstance == nil {
		dbConn := _NewDatabaseConnection()
		dbInstance = dbConn.Connect()
	}
	return dbInstance
}

// Close closes the database connection.
func Close() {
	if dbInstance != nil {
		err := dbInstance.Close()
		if err != nil {
			logger := logger.GetLogger()
			errStr := fmt.Sprintf("Error closing the database connection: %s", err)
			logger.Fatal(errStr)
		}
		dbInstance = nil
	}
}
