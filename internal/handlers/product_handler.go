package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ruranjo/unientrega/internal/models"
	"github.com/ruranjo/unientrega/internal/services"
)

// ProductHandler handles product management requests
type ProductHandler struct {
	productService *services.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct creates a new product
// @Summary Create product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Product true "Product data"
// @Success 201 {object} models.Product
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.productService.CreateProduct(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProduct returns a product by ID
// @Summary Get product by ID
// @Tags products
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.productService.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// ListProducts returns a list of products
// @Summary List products
// @Tags products
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Param category query string false "Filter by category"
// @Param store_id query string false "Filter by store ID"
// @Param active_only query bool false "Show only active products" default(false)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	categoryStr := c.Query("category")
	storeIDStr := c.Query("store_id")
	activeOnlyStr := c.DefaultQuery("active_only", "false")

	var category models.ProductCategory
	if categoryStr != "" {
		category = models.ProductCategory(categoryStr)
	}

	var storeID uuid.UUID
	if storeIDStr != "" {
		var err error
		storeID, err = uuid.Parse(storeIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
			return
		}
	}

	activeOnly := activeOnlyStr == "true" || activeOnlyStr == "1"

	products, err := h.productService.ListProducts(limit, offset, category, storeID, activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get total count for pagination
	total, err := h.productService.CountProducts(category, storeID, activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

// UpdateProduct updates a product
// @Summary Update product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param request body models.Product true "Product data"
// @Success 200 {object} models.Product
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var updateData models.Product
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing product
	product, err := h.productService.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Update fields
	product.Name = updateData.Name
	product.Description = updateData.Description
	product.Category = updateData.Category
	product.Price = updateData.Price
	product.Stock = updateData.Stock
	product.StoreID = updateData.StoreID
	product.SKU = updateData.SKU
	product.ImageURL = updateData.ImageURL
	product.IsActive = updateData.IsActive

	err = h.productService.UpdateProduct(product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct deletes a product
// @Summary Delete product
// @Tags products
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	err = h.productService.DeleteProduct(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// UpdateStock updates the stock quantity of a product
// @Summary Update product stock
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param request body map[string]int true "Stock quantity"
// @Success 200 {object} map[string]string
// @Router /api/v1/products/{id}/stock [patch]
func (h *ProductHandler) UpdateStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req struct {
		Stock int `json:"stock" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.productService.UpdateStock(id, req.Stock)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock updated successfully"})
}
