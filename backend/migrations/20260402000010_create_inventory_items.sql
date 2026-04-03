-- +goose Up
CREATE TABLE inventory_items (
    product_id        UUID        NOT NULL PRIMARY KEY REFERENCES products (id) ON DELETE CASCADE,
    on_hand           INT         NOT NULL DEFAULT 0,
    reserved          INT         NOT NULL DEFAULT 0,
    reorder_threshold INT         NOT NULL DEFAULT 0,
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS inventory_items;
