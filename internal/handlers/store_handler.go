package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ruranjo/unientrega/internal/models"
	"github.com/ruranjo/unientrega/internal/services"
)

// StoreHandler handles store management requests
type StoreHandler struct {
	storeService *services.StoreService
}

// NewStoreHandler creates a new store handler
func NewStoreHandler(storeService *services.StoreService) *StoreHandler {
	return &StoreHandler{
		storeService: storeService,
	}
}

// CreateStore creates a new store
// @Summary Create store
// @Tags stores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Store true "Store data"
// @Success 201 {object} models.Store
// @Router /api/v1/stores [post]
func (h *StoreHandler) CreateStore(c *gin.Context) {
	var store models.Store

	if err := c.ShouldBindJSON(&store); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user ID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDStr.(uuid.UUID)

	// Set owner ID to current user if not provided (or override it if user is not superuser)
	// For now, let's assume the creator is the owner
	store.OwnerID = userID

	err := h.storeService.CreateStore(&store)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, store)
}

// GetStore returns a store by ID
// @Summary Get store by ID
// @Tags stores
// @Produce json
// @Security BearerAuth
// @Param id path string true "Store ID"
// @Success 200 {object} models.Store
// @Router /api/v1/stores/{id} [get]
func (h *StoreHandler) GetStore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	store, err := h.storeService.GetStoreByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	c.JSON(http.StatusOK, store)
}

// ListStores returns a list of stores
// @Summary List stores
// @Tags stores
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Param active_only query bool false "Show only active stores" default(false)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/stores [get]
func (h *StoreHandler) ListStores(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	activeOnlyStr := c.DefaultQuery("active_only", "false")

	activeOnly := activeOnlyStr == "true" || activeOnlyStr == "1"

	stores, err := h.storeService.ListStores(limit, offset, activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get total count for pagination
	total, err := h.storeService.CountStores(activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stores": stores,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// UpdateStore updates a store
// @Summary Update store
// @Tags stores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Store ID"
// @Param request body models.Store true "Store data"
// @Success 200 {object} models.Store
// @Router /api/v1/stores/{id} [put]
func (h *StoreHandler) UpdateStore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	var updateData models.Store
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing store
	store, err := h.storeService.GetStoreByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	// Check permissions: only owner or superuser can update
	userIDStr, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")
	userID := userIDStr.(uuid.UUID)

	if store.OwnerID != userID && userRole != models.RoleSuperUser {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this store"})
		return
	}

	// Update fields
	store.Name = updateData.Name
	store.Description = updateData.Description
	store.Location = updateData.Location
	store.IsActive = updateData.IsActive

	err = h.storeService.UpdateStore(store)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, store)
}

// DeleteStore deletes a store
// @Summary Delete store
// @Tags stores
// @Produce json
// @Security BearerAuth
// @Param id path string true "Store ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/stores/{id} [delete]
func (h *StoreHandler) DeleteStore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		return
	}

	// Check permissions: only owner or superuser can delete
	// For safety, maybe only superuser should delete stores?
	// Let's allow owner too for now, but check ownership
	store, err := h.storeService.GetStoreByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")
	userID := userIDStr.(uuid.UUID)

	if store.OwnerID != userID && userRole != models.RoleSuperUser {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this store"})
		return
	}

	err = h.storeService.DeleteStore(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Store deleted successfully"})
}
