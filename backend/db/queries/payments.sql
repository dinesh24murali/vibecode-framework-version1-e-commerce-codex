-- name: CreatePayment :one
INSERT INTO payments (order_id, provider, provider_order_id, status, amount_inr, gateway_payload)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPaymentByProviderOrderID :one
SELECT * FROM payments
WHERE provider_order_id = $1;

-- name: GetPaymentsByOrder :many
SELECT * FROM payments
WHERE order_id = $1
ORDER BY created_at DESC;

-- name: UpdatePaymentStatus :one
UPDATE payments
SET status              = $2,
    provider_payment_id = $3,
    gateway_payload     = $4,
    updated_at          = now()
WHERE id = $1
RETURNING *;
