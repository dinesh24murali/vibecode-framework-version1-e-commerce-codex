-- name: GetCartByUser :many
SELECT * FROM cart
WHERE user_id = $1
ORDER BY updated_at DESC;

-- name: GetCartItem :one
SELECT * FROM cart
WHERE user_id = $1 AND product_id = $2;

-- name: UpsertCartItem :one
INSERT INTO cart (user_id, product_id, quantity)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, product_id)
DO UPDATE SET quantity   = EXCLUDED.quantity,
              updated_at = now()
RETURNING *;

-- name: DeleteCartItem :exec
DELETE FROM cart
WHERE user_id = $1 AND product_id = $2;

-- name: ClearCart :exec
DELETE FROM cart
WHERE user_id = $1;
