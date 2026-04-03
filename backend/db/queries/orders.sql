-- name: CreateOrder :one
INSERT INTO orders (user_id, order_number, status, subtotal_inr, shipping_inr, tax_inr, total_inr, currency_code, placed_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now())
RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1;

-- name: GetOrderByNumber :one
SELECT * FROM orders
WHERE order_number = $1;

-- name: ListOrdersByUser :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status     = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;
