package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ruranjo/unientrega/internal/config"
)

// APIHandler handles API-related requests
type APIHandler struct {
	config *config.Config
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(cfg *config.Config) *APIHandler {
	return &APIHandler{
		config: cfg,
	}
}

// GetWelcome returns the API welcome message
func (h *APIHandler) GetWelcome(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to UniEntrega API",
		"version": "1.0.0",
	})
}

// GetConfig returns the current configuration (for testing)
func (h *APIHandler) GetConfig(c *gin.Context) {
	c.JSON(200, gin.H{
		"database": gin.H{
			"host": h.config.Database.Host,
			"port": h.config.Database.Port,
			"name": h.config.Database.Name,
			"user": h.config.Database.User,
		},
		"app": gin.H{
			"env":  h.config.App.Env,
			"port": h.config.App.Port,
		},
	})
}
