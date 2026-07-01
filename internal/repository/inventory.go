package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sirovia/bardista/internal/domain"
)

type InventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func scanInventoryItem(scanner interface {
	Scan(dest ...any) error
}) (*domain.InventoryItem, error) {
	item := &domain.InventoryItem{}
	var threshold sql.NullFloat64

	err := scanner.Scan(
		&item.ID,
		&item.Name,
		&item.Unit,
		&item.Quantity,
		&threshold,
		&item.DeletedAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if threshold.Valid {
		item.LowStockThreshold = &threshold.Float64
	}

	return item, nil
}

const inventorySelectColumns = `
	id, name, unit, quantity, low_stock_threshold, deleted_at, created_at, updated_at
`

func (r *InventoryRepository) GetAll(ctx context.Context) ([]domain.InventoryItem, error) {
	query := `
		SELECT ` + inventorySelectColumns + `
		FROM inventory_items
		WHERE deleted_at IS NULL
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.InventoryItem
	for rows.Next() {
		item, err := scanInventoryItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, rows.Err()
}

func (r *InventoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.InventoryItem, error) {
	query := `
		SELECT ` + inventorySelectColumns + `
		FROM inventory_items
		WHERE id = $1 AND deleted_at IS NULL
	`

	item, err := scanInventoryItem(r.db.QueryRowContext(ctx, query, id))
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *InventoryRepository) Create(ctx context.Context, item *domain.InventoryItem) error {
	query := `
		INSERT INTO inventory_items (
			id, name, unit, quantity, low_stock_threshold, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		item.ID,
		item.Name,
		item.Unit,
		item.Quantity,
		nullableFloat(item.LowStockThreshold),
		item.CreatedAt,
		item.UpdatedAt,
	)
	return err
}

func (r *InventoryRepository) Update(ctx context.Context, item *domain.InventoryItem) error {
	query := `
		UPDATE inventory_items
		SET name = $1, unit = $2, low_stock_threshold = $3, updated_at = $4
		WHERE id = $5 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		item.Name,
		item.Unit,
		nullableFloat(item.LowStockThreshold),
		item.UpdatedAt,
		item.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *InventoryRepository) AdjustStockByDelta(ctx context.Context, id uuid.UUID, delta float64) (*domain.InventoryItem, error) {
	query := `
		UPDATE inventory_items
		SET quantity = quantity + $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL AND quantity + $1 >= 0
		RETURNING ` + inventorySelectColumns

	item, err := scanInventoryItem(r.db.QueryRowContext(ctx, query, delta, time.Now(), id))
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *InventoryRepository) SetStock(ctx context.Context, id uuid.UUID, quantity float64) (*domain.InventoryItem, error) {
	query := `
		UPDATE inventory_items
		SET quantity = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
		RETURNING ` + inventorySelectColumns

	item, err := scanInventoryItem(r.db.QueryRowContext(ctx, query, quantity, time.Now(), id))
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *InventoryRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE inventory_items
		SET deleted_at = $1, updated_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *InventoryRepository) CountRecipeReferences(ctx context.Context, id uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM product_ingredients
		WHERE inventory_item_id = $1
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, id).Scan(&count)
	return count, err
}

func nullableFloat(value *float64) sql.NullFloat64 {
	if value == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *value, Valid: true}
}
