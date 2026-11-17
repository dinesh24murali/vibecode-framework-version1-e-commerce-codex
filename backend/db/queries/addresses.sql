-- name: ListAddressesByUser :many
SELECT * FROM addresses
WHERE user_id = $1
ORDER BY is_default DESC, created_at DESC;

-- name: GetAddressByID :one
SELECT * FROM addresses
WHERE id = $1;

-- name: GetAddressForUser :one
SELECT * FROM addresses
WHERE id = $1 AND user_id = $2;

-- name: CreateAddress :one
INSERT INTO addresses (user_id, line1, line2, landmark, city, state, postal_code, company_name, gst_number, phone_e164, is_default)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: UpdateAddress :one
UPDATE addresses
SET line1        = $3,
    line2        = $4,
    landmark     = $5,
    city         = $6,
    state        = $7,
    postal_code  = $8,
    company_name = $9,
    gst_number   = $10,
    phone_e164   = $11,
    is_default   = $12,
    updated_at   = now()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteAddress :exec
DELETE FROM addresses
WHERE id = $1 AND user_id = $2;

-- name: ClearDefaultAddresses :exec
UPDATE addresses
SET is_default = false
WHERE user_id = $1;

-- name: SetDefaultAddress :one
UPDATE addresses
SET is_default = true
WHERE id = $1 AND user_id = $2
RETURNING *;
