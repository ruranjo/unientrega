package database

import (
	"log"

	"github.com/ruranjo/unientrega/internal/models"
)

// Migrate runs database migrations
func Migrate() error {
	log.Println("Running database migrations...")

	// Auto-migrate models
	err := db.AutoMigrate(
		&models.User{},
		&models.PasswordReset{},
		&models.Store{},
		&models.Product{},
		// Add more models here as you create them
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}
