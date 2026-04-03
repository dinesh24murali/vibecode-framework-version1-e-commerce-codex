-- name: ListAddressesByUser :many
SELECT * FROM addresses
WHERE user_id = $1
ORDER BY is_default DESC, created_at DESC;

-- name: GetAddressByID :one
SELECT * FROM addresses
WHERE id = $1;

-- name: CreateAddress :one
INSERT INTO addresses (user_id, recipient_name, line1, line2, city, state, postal_code, country_code, phone_e164, is_default)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: UpdateAddress :one
UPDATE addresses
SET recipient_name = $2,
    line1          = $3,
    line2          = $4,
    city           = $5,
    state          = $6,
    postal_code    = $7,
    country_code   = $8,
    phone_e164     = $9
WHERE id = $1
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
