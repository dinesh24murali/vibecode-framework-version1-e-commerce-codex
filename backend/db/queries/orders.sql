-- name: CreateOrder :one
INSERT INTO orders (user_id, order_number, status, subtotal_inr, shipping_inr, tax_inr, total_inr, currency_code, coupon_code, coupon_discount_inr, placed_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, now())
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

-- name: ListOrdersByUserFiltered :many
SELECT * FROM orders
WHERE user_id = $1
  AND (cardinality($2::order_status[]) = 0 OR status = ANY($2::order_status[]))
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountOrdersByUserFiltered :one
SELECT count(*) FROM orders
WHERE user_id = $1
  AND (cardinality($2::order_status[]) = 0 OR status = ANY($2::order_status[]));

-- name: AdminListOrders :many
SELECT * FROM orders
WHERE (cardinality(sqlc.arg(statuses)::order_status[]) = 0 OR status = ANY(sqlc.arg(statuses)::order_status[]))
ORDER BY created_at DESC, id
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: AdminCountOrders :one
SELECT count(*) FROM orders
WHERE (cardinality(sqlc.arg(statuses)::order_status[]) = 0 OR status = ANY(sqlc.arg(statuses)::order_status[]));
