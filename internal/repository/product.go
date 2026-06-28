package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sirovia/bardista/internal/domain"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll(ctx context.Context) ([]domain.Product, error) {
	query := `
		SELECT id, name, description, price, is_available, created_at, updated_at
		FROM products
		WHERE is_available = true AND deleted_at IS NULL
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Price, &p.IsAvailable,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, err
}

func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	query := `
		SELECT id, name, description, price, is_available, created_at, updated_at
		FROM products
		WHERE id = $1 AND deleted_at IS NULL
	`

	p := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.IsAvailable, &p.CreatedAt, &p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *ProductRepository) Create(ctx context.Context, p *domain.Product) error {
	query := `
		INSERT INTO products (id, name, description, price, is_available, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.Name, p.Description, p.Price,
		p.IsAvailable, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *ProductRepository) Update(ctx context.Context, p *domain.Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, is_available = $4, updated_at = $5
		WHERE id = $6 AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		p.Name, p.Description, p.Price, p.IsAvailable, time.Now(), p.ID,
	)
	return err
}

func (r *ProductRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE products SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}
