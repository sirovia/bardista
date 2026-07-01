CREATE TABLE IF NOT EXISTS product_ingredients (
    id                  UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id          UUID            NOT NULL REFERENCES products(id),
    inventory_item_id   UUID            NOT NULL REFERENCES inventory_items(id),
    quantity            NUMERIC(12, 3)  NOT NULL CHECK (quantity > 0),
    created_at          TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    UNIQUE (product_id, inventory_item_id)
);

CREATE INDEX IF NOT EXISTS idx_product_ingredients_product_id
    ON product_ingredients(product_id);

CREATE INDEX IF NOT EXISTS idx_product_ingredients_inventory_item_id
    ON product_ingredients(inventory_item_id);
