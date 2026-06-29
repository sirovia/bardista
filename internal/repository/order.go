package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/sirovia/bardista/internal/domain"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	orderQuery := `
		INSERT INTO orders (id, user_id, status, total_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = tx.ExecContext(ctx, orderQuery,
		order.ID, order.UserID, order.Status,
		order.TotalPrice, order.CreatedAt, order.UpdatedAt,
	)
	if err != nil {
		return err
	}

	itemQuery := `
		INSERT INTO order_items (id, order_id, product_id, quantity, unit_price)
		VALUES ($1, $2, $3, $4, $5)
	`

	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, itemQuery,
			item.ID, order.ID, item.ProductID, item.Quantity, item.UnitPrice,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	orderQuery := `
		SELECT id, user_id, status, total_price, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	order := &domain.Order{}
	err := r.db.QueryRowContext(ctx, orderQuery, id).Scan(
		&order.ID, &order.UserID, &order.Status,
		&order.TotalPrice, &order.CreatedAt, &order.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	itemQuery := `
		SELECT id, order_id, product_id, quantity, unit_price
		FROM order_items WHERE order_id = $1
	`
	rows, err := r.db.QueryContext(ctx, itemQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID,
			&item.Quantity, &item.UnitPrice,
		); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}
	return order, nil
}

func (r *OrderRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Order, error) {
	orderQuery := `
		SELECT id, user_id, status, total_price, created_at, updated_at
		FROM orders WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, orderQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(
			&o.ID, &o.UserID, &o.Status,
			&o.TotalPrice, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, err
		}

		orders = append(orders, o)
	}

	return orders, nil
}

func (r *OrderRepository) GetAll(ctx context.Context) ([]domain.Order, error) {
	orderQuery := `
		SELECT id, user_id, status, total_price, created_at, updated_at
		FROM orders ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, orderQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(
			&o.ID, &o.UserID, &o.Status,
			&o.TotalPrice, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, err
		}

		orders = append(orders, o)
	}
	return orders, err
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, id, status)
	return err
}
