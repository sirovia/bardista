package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sirovia/bardista/internal/domain"
	"github.com/sirovia/bardista/internal/repository"
)

var ErrInvalidTransition = errors.New("invalid status transition")
var ErrProductUnavailable = errors.New("product unavailable")

var validTransitions = map[string][]string{
	"pending":   {"confirmed", "cancelled"},
	"confirmed": {"completed"},
	"completed": {},
	"cancelled": {},
}

type OrderService struct {
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
}

func NewOrderService(orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository) *OrderService {
	return &OrderService{orderRepo: orderRepo, productRepo: productRepo}
}

type OrderItemInput struct {
	ProductID uuid.UUID
	Quantity  int
}

func (s *OrderService) CreateOrder(ctx context.Context, userID uuid.UUID, items []OrderItemInput) (*domain.Order, error) {
	if len(items) == 0 {
		return nil, errors.New("Order must contain at least one item")
	}

	orderID := uuid.New()
	var orderItems []domain.OrderItem
	var totalPrice float64

	for _, item := range items {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err == repository.ErrNotFound {
			return nil, errors.New("product not found")
		}
		if err != nil {
			return nil, err
		}
		if !product.IsAvailable {
			return nil, ErrProductUnavailable
		}

		orderItems = append(orderItems, domain.OrderItem{
			ID:        uuid.New(),
			OrderID:   orderID,
			ProductID: product.ID,
			Quantity:  item.Quantity,
			UnitPrice: product.Price,
		})

		totalPrice += product.Price * float64(item.Quantity)
	}

	now := time.Now()
	order := &domain.Order{
		ID:         uuid.New(),
		UserID:     userID,
		Status:     "pending",
		TotalPrice: totalPrice,
		Items:      orderItems,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID, role string) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if role != "admin" && order.UserID != userID {
		return nil, repository.ErrNotFound
	}

	return order, nil
}

func (s *OrderService) GetOrders(ctx context.Context, userID uuid.UUID, role string) ([]domain.Order, error) {
	if role == "admin" {
		return s.orderRepo.GetAll(ctx)
	}
	return s.orderRepo.GetByUserID(ctx, userID)
}

func (s *OrderService) UpdateStatus(ctx context.Context, id uuid.UUID, newStatus string) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	allowed := validTransitions[order.Status]
	valid := false
	for _, s := range allowed {
		if s == newStatus {
			valid = true
			break
		}
	}

	if !valid {
		return nil, ErrInvalidTransition
	}

	if err := s.orderRepo.UpdateStatus(ctx, id, newStatus); err != nil {
		return nil, err
	}

	order.Status = newStatus
	return order, nil
}
