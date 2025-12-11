package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ruranjo/unientrega/internal/handlers"
	"github.com/ruranjo/unientrega/internal/middleware"
	"github.com/ruranjo/unientrega/internal/models"
)

// SetupStoreRoutes configures store management routes
func SetupStoreRoutes(v1 *gin.RouterGroup, storeHandler *handlers.StoreHandler) {
	stores := v1.Group("/stores")
	stores.Use(middleware.AuthRequired())
	{
		// List stores (all authenticated users can view)
		stores.GET("", storeHandler.ListStores)

		// Get store by ID (all authenticated users can view)
		stores.GET("/:id", storeHandler.GetStore)

		// Create store (superuser only)
		stores.POST("", middleware.RoleRequired(models.RoleSuperUser), storeHandler.CreateStore)

		// Update store (owner or superuser - logic in handler)
		stores.PUT("/:id", middleware.RoleRequired(models.RoleSuperUser, models.RoleStore), storeHandler.UpdateStore)

		// Delete store (owner or superuser - logic in handler)
		stores.DELETE("/:id", middleware.RoleRequired(models.RoleSuperUser, models.RoleStore), storeHandler.DeleteStore)
	}
}
