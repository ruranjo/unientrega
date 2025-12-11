package services

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/ruranjo/unientrega/internal/models"
	"github.com/ruranjo/unientrega/internal/repository"
)

// StoreService handles business logic for stores
type StoreService struct {
	storeRepo *repository.StoreRepository
	userRepo  *repository.UserRepository
}

// NewStoreService creates a new store service
func NewStoreService(storeRepo *repository.StoreRepository, userRepo *repository.UserRepository) *StoreService {
	return &StoreService{
		storeRepo: storeRepo,
		userRepo:  userRepo,
	}
}

// CreateStore creates a new store with validation
func (s *StoreService) CreateStore(store *models.Store) error {
	// Validate store name
	if strings.TrimSpace(store.Name) == "" {
		return errors.New("store name is required")
	}

	// Validate owner exists
	if store.OwnerID != uuid.Nil {
		_, err := s.userRepo.GetByID(store.OwnerID)
		if err != nil {
			return errors.New("invalid owner ID")
		}
	} else {
		return errors.New("owner ID is required")
	}

	return s.storeRepo.Create(store)
}

// GetStoreByID retrieves a store by ID
func (s *StoreService) GetStoreByID(id uuid.UUID) (*models.Store, error) {
	return s.storeRepo.GetByID(id)
}

// GetStoresByOwner retrieves stores managed by a specific user
func (s *StoreService) GetStoresByOwner(ownerID uuid.UUID) ([]*models.Store, error) {
	return s.storeRepo.GetByOwnerID(ownerID)
}

// UpdateStore updates a store with validation
func (s *StoreService) UpdateStore(store *models.Store) error {
	// Validate store name
	if strings.TrimSpace(store.Name) == "" {
		return errors.New("store name is required")
	}

	// Check if store exists
	_, err := s.storeRepo.GetByID(store.ID)
	if err != nil {
		return err
	}

	return s.storeRepo.Update(store)
}

// DeleteStore soft deletes a store
func (s *StoreService) DeleteStore(id uuid.UUID) error {
	// Check if store exists
	_, err := s.storeRepo.GetByID(id)
	if err != nil {
		return err
	}

	return s.storeRepo.Delete(id)
}

// ListStores returns a list of stores with optional filters
func (s *StoreService) ListStores(limit, offset int, activeOnly bool) ([]*models.Store, error) {
	// Set default limit if not provided or invalid
	if limit <= 0 {
		limit = 10
	}
	// Cap maximum limit
	if limit > 100 {
		limit = 100
	}

	// Ensure offset is non-negative
	if offset < 0 {
		offset = 0
	}

	return s.storeRepo.List(limit, offset, activeOnly)
}

// CountStores returns the total number of stores
func (s *StoreService) CountStores(activeOnly bool) (int64, error) {
	return s.storeRepo.Count(activeOnly)
}

// IsStoreOwner checks if a user is the owner of a store
func (s *StoreService) IsStoreOwner(userID, storeID uuid.UUID) (bool, error) {
	store, err := s.storeRepo.GetByID(storeID)
	if err != nil {
		return false, err
	}
	return store.OwnerID == userID, nil
}
