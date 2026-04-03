-- +goose Up
CREATE TABLE product_attributes (
    id            UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    code          TEXT        NOT NULL UNIQUE,
    name          TEXT        NOT NULL,
    is_filterable BOOLEAN     NOT NULL DEFAULT false,
    sort_order    INT         NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS product_attributes;
