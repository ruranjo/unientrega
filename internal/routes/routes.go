package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ruranjo/unientrega/internal/config"
	"github.com/ruranjo/unientrega/internal/database"
	"github.com/ruranjo/unientrega/internal/handlers"
	"github.com/ruranjo/unientrega/internal/repository"
	"github.com/ruranjo/unientrega/internal/services"
)

// SetupRoutes configures all application routes
func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	// Get database instance
	db := database.GetDB()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	passwordResetRepo := repository.NewPasswordResetRepository(db)
	productRepo := repository.NewProductRepository(db)
	storeRepo := repository.NewStoreRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo, passwordResetRepo)
	authService := services.NewAuthService(userService)
	storeService := services.NewStoreService(storeRepo, userRepo)
	productService := services.NewProductService(productRepo)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(cfg)
	apiHandler := handlers.NewAPIHandler(cfg)
	authHandler := handlers.NewAuthHandler(authService, userService)
	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)
	storeHandler := handlers.NewStoreHandler(storeService)

	// Setup health and root routes
	SetupHealthRoutes(r, healthHandler)

	// API v1 routes
	v1 := r.Group("/api/v1")

	// Setup route groups
	SetupAPIRoutes(v1, apiHandler)
	SetupAuthRoutes(v1, authHandler)
	SetupUserRoutes(v1, userHandler)
	SetupStoreRoutes(v1, storeHandler)
	SetupProductRoutes(v1, productHandler)
}
