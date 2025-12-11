package services

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/ruranjo/unientrega/internal/models"
	"github.com/ruranjo/unientrega/internal/repository"
)

// ProductService handles business logic for products
type ProductService struct {
	productRepo *repository.ProductRepository
}

// NewProductService creates a new product service
func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

// CreateProduct creates a new product with validation
func (s *ProductService) CreateProduct(product *models.Product) error {
	// Validate product name
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("product name is required")
	}

	// Validate price
	if product.Price < 0 {
		return errors.New("product price must be non-negative")
	}

	// Validate stock
	if product.Stock < 0 {
		return errors.New("product stock must be non-negative")
	}

	// Validate category
	if !product.Category.IsValid() {
		return errors.New("invalid product category")
	}

	// Validate store ID
	if product.StoreID == uuid.Nil {
		return errors.New("store ID is required")
	}

	// Check SKU uniqueness if provided
	if product.SKU != "" {
		exists, err := s.productRepo.ExistsBySKU(product.SKU)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("product with this SKU already exists")
		}
	}

	return s.productRepo.Create(product)
}

// GetProductByID retrieves a product by ID
func (s *ProductService) GetProductByID(id uuid.UUID) (*models.Product, error) {
	return s.productRepo.GetByID(id)
}

// GetProductBySKU retrieves a product by SKU
func (s *ProductService) GetProductBySKU(sku string) (*models.Product, error) {
	if sku == "" {
		return nil, errors.New("SKU is required")
	}
	return s.productRepo.GetBySKU(sku)
}

// UpdateProduct updates a product with validation
func (s *ProductService) UpdateProduct(product *models.Product) error {
	// Validate product name
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("product name is required")
	}

	// Validate price
	if product.Price < 0 {
		return errors.New("product price must be non-negative")
	}

	// Validate stock
	if product.Stock < 0 {
		return errors.New("product stock must be non-negative")
	}

	// Validate category
	if !product.Category.IsValid() {
		return errors.New("invalid product category")
	}

	// Validate store ID
	if product.StoreID == uuid.Nil {
		return errors.New("store ID is required")
	}

	// Check if product exists
	existingProduct, err := s.productRepo.GetByID(product.ID)
	if err != nil {
		return err
	}

	// Check SKU uniqueness if changed
	if product.SKU != "" && product.SKU != existingProduct.SKU {
		exists, err := s.productRepo.ExistsBySKU(product.SKU)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("product with this SKU already exists")
		}
	}

	return s.productRepo.Update(product)
}

// DeleteProduct soft deletes a product
func (s *ProductService) DeleteProduct(id uuid.UUID) error {
	// Check if product exists
	_, err := s.productRepo.GetByID(id)
	if err != nil {
		return err
	}

	return s.productRepo.Delete(id)
}

// ListProducts returns a list of products with optional filters
func (s *ProductService) ListProducts(limit, offset int, category models.ProductCategory, storeID uuid.UUID, activeOnly bool) ([]*models.Product, error) {
	// Set default limit if not provided or invalid
	if limit <= 0 {
		limit = 10
	}
	// Cap maximum limit to prevent excessive queries
	if limit > 100 {
		limit = 100
	}

	// Ensure offset is non-negative
	if offset < 0 {
		offset = 0
	}

	return s.productRepo.List(limit, offset, category, storeID, activeOnly)
}

// CountProducts returns the total number of products with optional filters
func (s *ProductService) CountProducts(category models.ProductCategory, storeID uuid.UUID, activeOnly bool) (int64, error) {
	return s.productRepo.Count(category, storeID, activeOnly)
}

// UpdateStock updates the stock quantity of a product with validation
func (s *ProductService) UpdateStock(id uuid.UUID, quantity int) error {
	// Validate stock quantity
	if quantity < 0 {
		return errors.New("stock quantity must be non-negative")
	}

	// Check if product exists
	_, err := s.productRepo.GetByID(id)
	if err != nil {
		return err
	}

	return s.productRepo.UpdateStock(id, quantity)
}

// IsAvailable checks if a product is available (active and in stock)
func (s *ProductService) IsAvailable(product *models.Product) bool {
	return product.IsActive && product.Stock > 0
}
