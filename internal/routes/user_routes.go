package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ruranjo/unientrega/internal/handlers"
	"github.com/ruranjo/unientrega/internal/middleware"
	"github.com/ruranjo/unientrega/internal/models"
)

// SetupUserRoutes configures user management routes
func SetupUserRoutes(v1 *gin.RouterGroup, userHandler *handlers.UserHandler) {
	users := v1.Group("/users")
	users.Use(middleware.AuthRequired())
	{
		// List users (admin only)
		users.GET("", middleware.RoleRequired(models.RoleSuperUser), userHandler.ListUsers)

		// Get user by ID (authenticated users can view)
		users.GET("/:id", userHandler.GetUser)

		// Update user (users can update themselves, admin can update anyone)
		users.PUT("/:id", userHandler.UpdateUser)

		// Delete user (admin only)
		users.DELETE("/:id", middleware.RoleRequired(models.RoleSuperUser), userHandler.DeleteUser)

		// Change password (users can change their own password)
		users.PUT("/:id/password", userHandler.ChangePassword)
	}
}
