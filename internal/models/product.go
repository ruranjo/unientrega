package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductCategory represents a product category in the system
type ProductCategory string

// ProductCategory constants
const (
	CategoryPhotocopies ProductCategory = "photocopies"
	CategoryStationery  ProductCategory = "stationery"
	CategoryFood        ProductCategory = "food"
	CategoryBeverages   ProductCategory = "beverages"
	CategorySnacks      ProductCategory = "snacks"
	CategoryOther       ProductCategory = "other"
)

// IsValid checks if the product category is valid
func (pc ProductCategory) IsValid() bool {
	switch pc {
	case CategoryPhotocopies, CategoryStationery, CategoryFood, CategoryBeverages, CategorySnacks, CategoryOther:
		return true
	}
	return false
}

// String returns the string representation of the product category
func (pc ProductCategory) String() string {
	return string(pc)
}

// Product represents a product in the store
type Product struct {
	ID          uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string          `gorm:"size:200;not null" json:"name"`
	Description string          `gorm:"type:text" json:"description"`
	Category    ProductCategory `gorm:"type:varchar(50);not null" json:"category"`
	Price       float64         `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock       int             `gorm:"not null;default:0" json:"stock"`
	StoreID     uuid.UUID       `gorm:"type:uuid;not null" json:"store_id"` // Foreign key to Store
	SKU         string          `gorm:"size:100;uniqueIndex" json:"sku,omitempty"`
	ImageURL    string          `gorm:"size:500" json:"image_url,omitempty"`
	IsActive    bool            `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"` // Soft delete
}

// TableName specifies the table name for Product model
func (Product) TableName() string {
	return "products"
}

// BeforeCreate is a GORM hook that runs before creating a product
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID if not set
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
