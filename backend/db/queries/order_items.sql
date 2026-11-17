-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, product_title, product_sku, quantity, unit_price_inr, line_total_inr)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: ListOrderItems :many
SELECT * FROM order_items
WHERE order_id = $1;

-- name: ListOrderItemsWithProduct :many
SELECT oi.id, oi.order_id, oi.product_id, oi.product_title, oi.quantity,
       oi.unit_price_inr, oi.line_total_inr,
       COALESCE(p.slug, '') AS product_slug
FROM order_items oi
LEFT JOIN products p ON p.id = oi.product_id
WHERE oi.order_id = $1
ORDER BY oi.id;
