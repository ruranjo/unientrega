package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ruranjo/unientrega/internal/handlers"
)

// SetupAPIRoutes configures public API routes
func SetupAPIRoutes(v1 *gin.RouterGroup, apiHandler *handlers.APIHandler) {
	v1.GET("/", apiHandler.GetWelcome)
	v1.GET("/config", apiHandler.GetConfig)
}
