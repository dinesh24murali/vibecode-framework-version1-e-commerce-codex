-- +goose Up
CREATE TABLE coupon_users (
    id         UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    coupon_id  UUID        NOT NULL REFERENCES coupons (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, coupon_id)
);

CREATE INDEX idx_coupon_users_coupon_created ON coupon_users (coupon_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS coupon_users;
