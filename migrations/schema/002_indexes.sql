-- +migrate Up
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_products_name_trgm ON products USING gin (name gin_trgm_ops);

CREATE INDEX idx_products_active_created_at ON products (created_at)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_products_active_price ON products (price)
    WHERE deleted_at IS NULL;

-- +migrate Down
DROP INDEX IF EXISTS idx_products_active_price;
DROP INDEX IF EXISTS idx_products_active_created_at;
DROP INDEX IF EXISTS idx_products_name_trgm;
DROP EXTENSION IF EXISTS pg_trgm;
