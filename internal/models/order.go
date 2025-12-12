package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusReady     OrderStatus = "ready"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// IsValid checks if the order status is valid
func (os OrderStatus) IsValid() bool {
	switch os {
	case OrderStatusPending, OrderStatusConfirmed, OrderStatusPreparing, OrderStatusReady, OrderStatusCompleted, OrderStatusCancelled:
		return true
	}
	return false
}

// String returns the string representation of the order status
func (os OrderStatus) String() string {
	return string(os)
}

// Order represents a customer order
type Order struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	StoreID          uuid.UUID      `gorm:"type:uuid;not null" json:"store_id"`
	DeliveryPersonID *uuid.UUID     `gorm:"type:uuid" json:"delivery_person_id"` // Nullable if not assigned
	Status           OrderStatus    `gorm:"type:varchar(50);not null;default:'pending'" json:"status"`
	Total            float64        `gorm:"type:decimal(10,2);not null" json:"total"`
	Items            []OrderItem    `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Order model
func (Order) TableName() string {
	return "orders"
}

// BeforeCreate is a GORM hook that runs before creating an order
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	if o.Status == "" {
		o.Status = OrderStatusPending
	}
	return nil
}

// OrderItem represents an item within an order
type OrderItem struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID   uuid.UUID      `gorm:"type:uuid;not null" json:"order_id"`
	ProductID uuid.UUID      `gorm:"type:uuid;not null" json:"product_id"`
	Quantity  int            `gorm:"not null" json:"quantity"`
	Price     float64        `gorm:"type:decimal(10,2);not null" json:"price"` // Snapshot price at time of order
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for OrderItem model
func (OrderItem) TableName() string {
	return "order_items"
}

// BeforeCreate is a GORM hook that runs before creating an order item
func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == uuid.Nil {
		oi.ID = uuid.New()
	}
	return nil
}
