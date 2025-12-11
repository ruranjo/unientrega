package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ruranjo/unientrega/internal/handlers"
)

// SetupHealthRoutes configures health and root routes
func SetupHealthRoutes(r *gin.Engine, healthHandler *handlers.HealthHandler) {
	r.GET("/", healthHandler.GetRoot)
	r.GET("/health", healthHandler.GetHealth)
}
