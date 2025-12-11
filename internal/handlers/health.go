package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ruranjo/unientrega/internal/config"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	config *config.Config
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{
		config: cfg,
	}
}

// GetHealth returns the health status of the API
func (h *HealthHandler) GetHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
		"env":    h.config.App.Env,
	})
}

// GetRoot returns a simple hello world message
func (h *HealthHandler) GetRoot(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World",
	})
}
