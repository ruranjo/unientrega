package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ruranjo/unientrega/internal/handlers"
	"github.com/ruranjo/unientrega/internal/middleware"
	"github.com/ruranjo/unientrega/internal/models"
)

// SetupOrderRoutes configures order management routes
func SetupOrderRoutes(v1 *gin.RouterGroup, orderHandler *handlers.OrderHandler) {
	orders := v1.Group("/orders")
	orders.Use(middleware.AuthRequired())
	{
		// Create order (authenticated users)
		orders.POST("", orderHandler.CreateOrder)

		// List orders (authenticated users - logic in handler)
		orders.GET("", orderHandler.ListOrders)

		// Get order by ID (authenticated users - logic in handler)
		orders.GET("/:id", orderHandler.GetOrder)

		// Update order status (store owner or superuser)
		orders.PATCH("/:id/status", middleware.RoleRequired(models.RoleSuperUser, models.RoleStore), orderHandler.UpdateOrderStatus)
	}
}
