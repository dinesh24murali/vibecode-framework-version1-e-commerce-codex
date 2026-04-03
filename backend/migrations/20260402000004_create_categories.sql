-- +goose Up
CREATE TABLE categories (
    id          UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    slug        TEXT        NOT NULL UNIQUE,
    name        TEXT        NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS categories;
