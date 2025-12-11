package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Store represents a physical or logical store unit (e.g., Copy Center, Cafeteria)
type Store struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"size:200;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Location    string         `gorm:"size:200" json:"location"`
	OwnerID     uuid.UUID      `gorm:"type:uuid;not null" json:"owner_id"` // User who manages this store
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}

// TableName specifies the table name for Store model
func (Store) TableName() string {
	return "stores"
}

// BeforeCreate is a GORM hook that runs before creating a store
func (s *Store) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID if not set
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
