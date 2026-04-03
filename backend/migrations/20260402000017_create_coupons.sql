-- +goose Up
CREATE TABLE coupons (
    id             UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    name           TEXT        NOT NULL,
    description    TEXT,
    code           TEXT        NOT NULL UNIQUE,
    is_enable      BOOLEAN     NOT NULL DEFAULT true,
    type           TEXT        NOT NULL CHECK (type IN ('flat', 'percentage')),
    percent        INT,
    flat           INT,
    coupon_expiry  TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS coupons;
