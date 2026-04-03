-- +goose Up
CREATE TYPE order_status AS ENUM (
    'pending_payment',
    'paid',
    'packed',
    'shipped',
    'delivered',
    'cancelled',
    'payment_failed',
    'refunded'
);

CREATE TABLE orders (
    id            UUID         NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id       UUID         NOT NULL REFERENCES users (id),
    order_number  TEXT         NOT NULL UNIQUE,
    status        order_status NOT NULL DEFAULT 'pending_payment',
    subtotal_inr  NUMERIC(12, 2) NOT NULL,
    shipping_inr  NUMERIC(12, 2) NOT NULL DEFAULT 0,
    tax_inr       NUMERIC(12, 2) NOT NULL DEFAULT 0,
    total_inr     NUMERIC(12, 2) NOT NULL,
    currency_code CHAR(3)      NOT NULL DEFAULT 'INR',
    placed_at     TIMESTAMPTZ,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_orders_user_created ON orders (user_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS order_status;
