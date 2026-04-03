-- +goose Up
-- order_address is an immutable snapshot captured at checkout time.
-- It has no FK to the addresses table — profile edits must not rewrite historical shipment data.
CREATE TABLE order_address (
    id             UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    order_id       UUID        NOT NULL UNIQUE REFERENCES orders (id) ON DELETE CASCADE,
    recipient_name TEXT        NOT NULL,
    line1          TEXT        NOT NULL,
    line2          TEXT,
    city           TEXT        NOT NULL,
    state          TEXT        NOT NULL,
    postal_code    TEXT        NOT NULL,
    country_code   TEXT        NOT NULL,
    phone_e164     TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS order_address;
