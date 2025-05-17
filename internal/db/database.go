package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"services-api/internal/models"
)

// Initialize creates a database connection with the specified database URL.
// It configures the connection with appropriate settings for production use.
// Returns a GORM database instance or an error if the connection fails.
func Initialize(databaseURL string) (*gorm.DB, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("database URL is required")
	}

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(databaseURL), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// Migrate runs database migrations for all models.
// It creates or updates the necessary tables, indexes, and constraints.
// Returns an error if the migration fails.
func Migrate(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	return db.AutoMigrate(
		&models.Service{},
		&models.Version{},
	)
}

// ConfigureConnectionPool sets up the connection pool settings for the database.
// It's important to call this after initializing the database to ensure proper resource usage.
func ConfigureConnectionPool(db *gorm.DB, maxIdleConns, maxOpenConns int, connMaxLifetime time.Duration) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database SQL instance: %w", err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(maxIdleConns)           // Maximum number of idle connections
	sqlDB.SetMaxOpenConns(maxOpenConns)           // Maximum number of open connections
	sqlDB.SetConnMaxLifetime(connMaxLifetime)     // Maximum connection lifetime

	// Verify connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// Close cleanly shuts down the database connection pool.
// Returns an error if there's a problem closing the connection.
func Close(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database SQL instance: %w", err)
	}

	return sqlDB.Close()
}