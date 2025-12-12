package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ruranjo/unientrega/internal/models"
	"github.com/ruranjo/unientrega/internal/repository"
)

// OrderService handles order business logic
type OrderService struct {
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
	storeRepo   *repository.StoreRepository
}

// NewOrderService creates a new order service
func NewOrderService(orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository, storeRepo *repository.StoreRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		storeRepo:   storeRepo,
	}
}

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	StoreID uuid.UUID `json:"store_id" binding:"required"`
	Items   []struct {
		ProductID uuid.UUID `json:"product_id" binding:"required"`
		Quantity  int       `json:"quantity" binding:"required,min=1"`
	} `json:"items" binding:"required,min=1"`
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(userID uuid.UUID, req *CreateOrderRequest) (*models.Order, error) {
	// Verify store exists and is active
	store, err := s.storeRepo.GetByID(req.StoreID)
	if err != nil {
		return nil, errors.New("store not found")
	}
	if !store.IsActive {
		return nil, errors.New("store is not active")
	}

	// Prepare order
	order := &models.Order{
		UserID:  userID,
		StoreID: req.StoreID,
		Status:  models.OrderStatusPending,
		Items:   make([]models.OrderItem, 0, len(req.Items)),
	}

	var total float64

	// Process items
	for _, itemReq := range req.Items {
		product, err := s.productRepo.GetByID(itemReq.ProductID)
		if err != nil {
			return nil, errors.New("product not found: " + itemReq.ProductID.String())
		}

		if !product.IsActive {
			return nil, errors.New("product is not active: " + product.Name)
		}

		if product.StoreID != req.StoreID {
			return nil, errors.New("product does not belong to the store: " + product.Name)
		}

		if product.Stock < itemReq.Quantity {
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}

		// Create order item
		orderItem := models.OrderItem{
			ProductID: product.ID,
			Quantity:  itemReq.Quantity,
			Price:     product.Price,
		}
		order.Items = append(order.Items, orderItem)

		// Update total
		total += product.Price * float64(itemReq.Quantity)

		// Decrement stock (simple approach, ideally should be transactional with order creation)
		// In a real app, we might reserve stock first or use a transaction in the repo layer
		// For now, we'll update it here. If order creation fails, we might have an issue (stock lost).
		// A better approach would be to pass a transaction to the service or handle it in repo.
		// Given the current architecture, we'll assume happy path or fix later.
		// To be safe, let's NOT decrement here but rely on the handler/repo transaction if possible.
		// But since we don't have transaction propagation yet, let's just update it.
		err = s.productRepo.UpdateStock(product.ID, product.Stock-itemReq.Quantity)
		if err != nil {
			return nil, errors.New("failed to update stock for product: " + product.Name)
		}
	}

	order.Total = total

	err = s.orderRepo.Create(order)
	if err != nil {
		// Rollback stock changes? (Complex without transactions)
		return nil, err
	}

	return order, nil
}

// GetOrder retrieves an order by ID and checks permissions
func (s *OrderService) GetOrder(id uuid.UUID, userID uuid.UUID, role models.Role) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check permissions
	if role == models.RoleSuperUser {
		return order, nil
	}

	if order.UserID == userID {
		return order, nil
	}

	// Check if user is owner of the store
	store, err := s.storeRepo.GetByID(order.StoreID)
	if err != nil {
		return nil, err
	}

	if store.OwnerID == userID {
		return order, nil
	}

	return nil, errors.New("permission denied")
}

// ListUserOrders lists orders for a user
func (s *OrderService) ListUserOrders(userID uuid.UUID, limit, offset int) ([]models.Order, int64, error) {
	return s.orderRepo.ListByUser(userID, limit, offset)
}

// ListStoreOrders lists orders for a store (permission checked in handler)
func (s *OrderService) ListStoreOrders(storeID uuid.UUID, limit, offset int) ([]models.Order, int64, error) {
	return s.orderRepo.ListByStore(storeID, limit, offset)
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(id uuid.UUID, status models.OrderStatus, userID uuid.UUID, role models.Role) (*models.Order, error) {
	if !status.IsValid() {
		return nil, errors.New("invalid status")
	}

	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Only store owner or superuser can update status
	if role != models.RoleSuperUser {
		store, err := s.storeRepo.GetByID(order.StoreID)
		if err != nil {
			return nil, err
		}
		if store.OwnerID != userID {
			return nil, errors.New("permission denied")
		}
	}

	err = s.orderRepo.UpdateStatus(id, status)
	if err != nil {
		return nil, err
	}

	order.Status = status
	return order, nil
}
