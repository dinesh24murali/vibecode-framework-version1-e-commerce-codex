-- Admin coupon management.

-- name: CountCoupons :one
SELECT count(*) FROM coupons;

-- name: AdminUpdateCoupon :one
UPDATE coupons
SET name = $2, description = $3, is_enable = $4, type = $5, percent = $6, flat = $7,
    coupon_expiry = $8, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteCoupon :execrows
DELETE FROM coupons WHERE id = $1;

-- name: CountCouponRedemptions :one
SELECT count(*) FROM coupon_users WHERE coupon_id = $1;
