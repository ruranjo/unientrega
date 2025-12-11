package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ruranjo/unientrega/internal/models"
)

// ProductRepository handles database operations for products
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create creates a new product
func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

// GetByID finds a product by ID
func (r *ProductRepository) GetByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("id = ?", id).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

// GetBySKU finds a product by SKU
func (r *ProductRepository) GetBySKU(sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("sku = ?", sku).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

// Update updates a product
func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

// Delete soft deletes a product
func (r *ProductRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Product{}, id).Error
}

// List returns a list of products with optional filters
func (r *ProductRepository) List(limit, offset int, category models.ProductCategory, storeID uuid.UUID, activeOnly bool) ([]*models.Product, error) {
	var products []*models.Product
	query := r.db.Limit(limit).Offset(offset)

	// Filter by category if provided
	if category != "" && category.IsValid() {
		query = query.Where("category = ?", category)
	}

	// Filter by store if provided
	if storeID != uuid.Nil {
		query = query.Where("store_id = ?", storeID)
	}

	// Filter by active status if requested
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	// Order by created_at descending (newest first)
	query = query.Order("created_at DESC")

	err := query.Find(&products).Error
	return products, err
}

// Count returns the total number of products with optional filters
func (r *ProductRepository) Count(category models.ProductCategory, storeID uuid.UUID, activeOnly bool) (int64, error) {
	var count int64
	query := r.db.Model(&models.Product{})

	// Filter by category if provided
	if category != "" && category.IsValid() {
		query = query.Where("category = ?", category)
	}

	// Filter by store if provided
	if storeID != uuid.Nil {
		query = query.Where("store_id = ?", storeID)
	}

	// Filter by active status if requested
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	err := query.Count(&count).Error
	return count, err
}

// ExistsBySKU checks if a product with the given SKU exists
func (r *ProductRepository) ExistsBySKU(sku string) (bool, error) {
	if sku == "" {
		return false, nil
	}
	var count int64
	err := r.db.Model(&models.Product{}).Where("sku = ?", sku).Count(&count).Error
	return count > 0, err
}

// UpdateStock updates the stock quantity of a product
func (r *ProductRepository) UpdateStock(id uuid.UUID, quantity int) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Update("stock", quantity).Error
}
