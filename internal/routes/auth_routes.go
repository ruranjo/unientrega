package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ruranjo/unientrega/internal/handlers"
	"github.com/ruranjo/unientrega/internal/middleware"
)

// SetupAuthRoutes configures authentication routes
func SetupAuthRoutes(v1 *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	auth := v1.Group("/auth")
	{
		// Public auth routes
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)

		// Password reset (public)
		auth.POST("/password-reset/request", authHandler.RequestPasswordReset)
		auth.POST("/password-reset/validate", authHandler.ValidateResetToken)
		auth.POST("/password-reset/confirm", authHandler.ResetPassword)

		// Protected auth routes
		authProtected := auth.Group("")
		authProtected.Use(middleware.AuthRequired())
		{
			authProtected.GET("/me", authHandler.GetMe)
			authProtected.POST("/logout", authHandler.Logout)
		}
	}
}
