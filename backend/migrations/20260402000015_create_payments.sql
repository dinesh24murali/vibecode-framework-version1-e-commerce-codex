-- +goose Up
CREATE TABLE payments (
    id                  UUID           NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    order_id            UUID           NOT NULL REFERENCES orders (id),
    provider            TEXT           NOT NULL,
    provider_order_id   TEXT           NOT NULL UNIQUE,
    provider_payment_id TEXT,
    status              TEXT           NOT NULL,
    amount_inr          NUMERIC(12, 2) NOT NULL,
    gateway_payload     JSONB,
    created_at          TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE INDEX idx_payments_order_status ON payments (order_id, status);

-- +goose Down
DROP TABLE IF EXISTS payments;
