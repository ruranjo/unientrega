package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ruranjo/unientrega/internal/config"
)

var db *gorm.DB

// Connect initializes the database connection
func Connect(cfg *config.Config) error {
	var err error

	// Build DSN (Data Source Name)
	dsn := cfg.Database.GetDSN()

	// Configure GORM logger
	gormLogger := logger.Default.LogMode(logger.Info)
	if cfg.App.Env == "production" {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	// Open database connection
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return db
}

// Close closes the database connection
func Close() error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
