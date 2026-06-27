CREATE TABLE IF NOT EXISTS products (
    id              UUID                PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255)        NOT NULL,
    description     TEXT,
    price           NUMERIC(10,2)       NOT NULL CHECK (price >= 0),
    is_available    BOOLEAN             NOT NULL DEFAULT TRUE,
    deleted_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ          NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ          NOT NULL DEFAULT NOW()
);