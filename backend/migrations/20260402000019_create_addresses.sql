-- +goose Up
CREATE TABLE addresses (
    id             UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id        UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    recipient_name TEXT        NOT NULL,
    line1          TEXT        NOT NULL,
    line2          TEXT,
    city           TEXT        NOT NULL,
    state          TEXT        NOT NULL,
    postal_code    TEXT        NOT NULL,
    country_code   TEXT        NOT NULL,
    phone_e164     TEXT,
    is_default     BOOLEAN     NOT NULL DEFAULT false,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_addresses_user ON addresses (user_id);

-- +goose Down
DROP TABLE IF EXISTS addresses;
