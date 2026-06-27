CREATE TABLE IF NOT EXISTS order_items (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id        UUID            NOT NULL REFERENCES orders(id),
    product_id      UUID            NOT NULL REFERENCES products(id),
    quantity        INTEGER         NOT NULL CHECK (quantity >= 1),
    unit_price      NUMERIC(10,2)   NOT NULL CHECK (unit_price >= 0)
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id           ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id     ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_products_available       ON products(is_available) WHERE is_available IS NULL;