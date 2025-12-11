package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ruranjo/unientrega/internal/handlers"
	"github.com/ruranjo/unientrega/internal/middleware"
	"github.com/ruranjo/unientrega/internal/models"
)

// SetupProductRoutes configures product management routes
func SetupProductRoutes(v1 *gin.RouterGroup, productHandler *handlers.ProductHandler) {
	products := v1.Group("/products")
	products.Use(middleware.AuthRequired())
	{
		// List products (all authenticated users can view)
		products.GET("", productHandler.ListProducts)

		// Get product by ID (all authenticated users can view)
		products.GET("/:id", productHandler.GetProduct)

		// Create product (store and superuser only)
		products.POST("", middleware.RoleRequired(models.RoleSuperUser, models.RoleStore), productHandler.CreateProduct)

		// Update product (store and superuser only)
		products.PUT("/:id", middleware.RoleRequired(models.RoleSuperUser, models.RoleStore), productHandler.UpdateProduct)

		// Delete product (store and superuser only)
		products.DELETE("/:id", middleware.RoleRequired(models.RoleSuperUser, models.RoleStore), productHandler.DeleteProduct)

		// Update stock (store and superuser only)
		products.PATCH("/:id/stock", middleware.RoleRequired(models.RoleSuperUser, models.RoleStore), productHandler.UpdateStock)
	}
}
