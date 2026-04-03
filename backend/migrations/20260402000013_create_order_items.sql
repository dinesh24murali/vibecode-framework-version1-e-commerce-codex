-- +goose Up
CREATE TABLE order_items (
    id              UUID           NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    order_id        UUID           NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    product_id      UUID           REFERENCES products (id) ON DELETE SET NULL,
    product_title   TEXT           NOT NULL,
    product_sku     TEXT           NOT NULL,
    quantity        INT            NOT NULL,
    unit_price_inr  NUMERIC(12, 2) NOT NULL,
    line_total_inr  NUMERIC(12, 2) NOT NULL
);

CREATE INDEX idx_order_items_order ON order_items (order_id);

-- +goose Down
DROP TABLE IF EXISTS order_items;
