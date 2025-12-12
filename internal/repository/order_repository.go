package repository

import (
	"github.com/google/uuid"
	"github.com/ruranjo/unientrega/internal/models"
	"gorm.io/gorm"
)

// OrderRepository handles database operations for orders
type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create creates a new order with items
func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

// GetByID retrieves an order by ID with its items
func (r *OrderRepository) GetByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Items").First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// ListByUser retrieves orders for a specific user
func (r *OrderRepository) ListByUser(userID uuid.UUID, limit, offset int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.db.Model(&models.Order{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Items").Limit(limit).Offset(offset).Order("created_at desc").Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// ListByStore retrieves orders for a specific store
func (r *OrderRepository) ListByStore(storeID uuid.UUID, limit, offset int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.db.Model(&models.Order{}).Where("store_id = ?", storeID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Items").Limit(limit).Offset(offset).Order("created_at desc").Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// UpdateStatus updates the status of an order
func (r *OrderRepository) UpdateStatus(id uuid.UUID, status models.OrderStatus) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}
