package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ruranjo/unientrega/internal/config"
	"github.com/ruranjo/unientrega/internal/database"
	"github.com/ruranjo/unientrega/internal/routes"
)

func main() {
	// Load .env file - search in multiple locations
	config.LoadEnv()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run database migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Set Gin mode based on environment
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	r := gin.Default()

	// Setup all routes
	routes.SetupRoutes(r, cfg)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port)
	log.Printf("Starting server on %s in %s mode", addr, cfg.App.Env)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
