-- +goose Up
-- Order.coupon_code / Order.coupon_discount from the API spec. tax_inr already holds the GST amount.
ALTER TABLE orders
    ADD COLUMN coupon_code         TEXT,
    ADD COLUMN coupon_discount_inr NUMERIC(12, 2) NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE orders
    DROP COLUMN coupon_code,
    DROP COLUMN coupon_discount_inr;
