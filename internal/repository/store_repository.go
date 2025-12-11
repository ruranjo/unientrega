package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ruranjo/unientrega/internal/models"
)

// StoreRepository handles database operations for stores
type StoreRepository struct {
	db *gorm.DB
}

// NewStoreRepository creates a new store repository
func NewStoreRepository(db *gorm.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

// Create creates a new store
func (r *StoreRepository) Create(store *models.Store) error {
	return r.db.Create(store).Error
}

// GetByID finds a store by ID
func (r *StoreRepository) GetByID(id uuid.UUID) (*models.Store, error) {
	var store models.Store
	err := r.db.Where("id = ?", id).First(&store).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("store not found")
		}
		return nil, err
	}
	return &store, nil
}

// GetByOwnerID finds stores managed by a specific user
func (r *StoreRepository) GetByOwnerID(ownerID uuid.UUID) ([]*models.Store, error) {
	var stores []*models.Store
	err := r.db.Where("owner_id = ?", ownerID).Find(&stores).Error
	return stores, err
}

// Update updates a store
func (r *StoreRepository) Update(store *models.Store) error {
	return r.db.Save(store).Error
}

// Delete soft deletes a store
func (r *StoreRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Store{}, id).Error
}

// List returns a list of stores with optional filters
func (r *StoreRepository) List(limit, offset int, activeOnly bool) ([]*models.Store, error) {
	var stores []*models.Store
	query := r.db.Limit(limit).Offset(offset)

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	// Order by name
	query = query.Order("name ASC")

	err := query.Find(&stores).Error
	return stores, err
}

// Count returns the total number of stores
func (r *StoreRepository) Count(activeOnly bool) (int64, error) {
	var count int64
	query := r.db.Model(&models.Store{})

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	err := query.Count(&count).Error
	return count, err
}
