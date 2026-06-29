package service

import (
	"context"
	"errors"

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
	return nil, nil
}
