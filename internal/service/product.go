package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sirovia/bardista/internal/domain"
	"github.com/sirovia/bardista/internal/repository"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) GetAll(ctx context.Context) ([]domain.Product, error) {
	return s.productRepo.GetAll(ctx)
}

func (s *ProductService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

func (s *ProductService) Create(ctx context.Context, name, description string, price float64, isAvailable bool) (*domain.Product, error) {
	now := time.Now()
	p := &domain.Product{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Price:       price,
		IsAvailable: isAvailable,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.productRepo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProductService) Update(ctx context.Context, id uuid.UUID, name, description string, price float64, isAvailable bool) (*domain.Product, error) {
	p, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	p.Name = name
	p.Description = description
	p.Price = price
	p.IsAvailable = isAvailable
	p.UpdatedAt = time.Now()

	if err := s.productRepo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProductService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.productRepo.SoftDelete(ctx, id)
}
