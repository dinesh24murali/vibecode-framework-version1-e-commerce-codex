-- name: GetCouponByCode :one
SELECT * FROM coupons
WHERE code = $1 AND is_enable = true;

-- name: GetCouponByID :one
SELECT * FROM coupons
WHERE id = $1;

-- name: ListCoupons :many
SELECT * FROM coupons
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateCoupon :one
INSERT INTO coupons (name, description, code, is_enable, type, percent, flat, coupon_expiry)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateCoupon :one
UPDATE coupons
SET name          = $2,
    description   = $3,
    is_enable     = $4,
    coupon_expiry = $5,
    updated_at    = now()
WHERE id = $1
RETURNING *;
