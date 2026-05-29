-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE products (
    id         UUID PRIMARY KEY NOT NULL,
    name       TEXT NOT NULL,
    price      NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted_at TIMESTAMPTZ
);

-- +migrate Down
DROP TABLE IF EXISTS products;
DROP EXTENSION IF EXISTS "uuid-ossp" CASCADE;