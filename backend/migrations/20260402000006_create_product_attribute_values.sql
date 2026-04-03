-- +goose Up
CREATE TABLE product_attribute_values (
    id           UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    attribute_id UUID        NOT NULL REFERENCES product_attributes (id) ON DELETE CASCADE,
    value        TEXT        NOT NULL,
    slug         TEXT        NOT NULL,
    sort_order   INT         NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (attribute_id, slug)
);

-- +goose Down
DROP TABLE IF EXISTS product_attribute_values;
