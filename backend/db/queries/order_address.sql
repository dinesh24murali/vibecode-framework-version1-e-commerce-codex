-- name: CreateOrderAddress :one
INSERT INTO order_address (order_id, recipient_name, line1, line2, city, state, postal_code, country_code, phone_e164)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetOrderAddress :one
SELECT * FROM order_address
WHERE order_id = $1;
