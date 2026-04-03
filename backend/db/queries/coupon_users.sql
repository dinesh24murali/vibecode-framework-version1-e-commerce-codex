-- name: CreateCouponUser :one
INSERT INTO coupon_users (user_id, coupon_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetCouponUserByUserAndCoupon :one
SELECT * FROM coupon_users
WHERE user_id = $1 AND coupon_id = $2;

-- name: ListCouponUsersByCoupon :many
SELECT * FROM coupon_users
WHERE coupon_id = $1
ORDER BY created_at DESC;
