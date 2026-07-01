CREATE TABLE IF NOT EXISTS inventory_items (
    id                      UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    name                    VARCHAR(255)    NOT NULL,
    unit                    VARCHAR(20)     NOT NULL,
    quantity                NUMERIC(12, 3)  NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    low_stock_threshold     NUMERIC(12, 3)  CHECK (low_stock_threshold IS NULL OR low_stock_threshold >= 0),
    deleted_at              TIMESTAMPTZ,
    created_at              TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ      NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_inventory_items_name_active
    ON inventory_items(name) WHERE deleted_at IS NULL;
