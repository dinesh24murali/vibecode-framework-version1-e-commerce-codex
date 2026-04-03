-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, product_title, product_sku, quantity, unit_price_inr, line_total_inr)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: ListOrderItems :many
SELECT * FROM order_items
WHERE order_id = $1;
