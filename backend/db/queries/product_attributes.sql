-- name: ListAttributes :many
SELECT * FROM product_attributes
ORDER BY sort_order, name;

-- name: ListFilterableAttributes :many
SELECT * FROM product_attributes
WHERE is_filterable = true
ORDER BY sort_order, name;

-- name: GetAttributeByCode :one
SELECT * FROM product_attributes
WHERE code = $1;

-- name: CreateAttribute :one
INSERT INTO product_attributes (code, name, is_filterable, sort_order)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAttributeByID :one
SELECT * FROM product_attributes
WHERE id = $1;

-- name: ListAttributesPaged :many
SELECT * FROM product_attributes
ORDER BY sort_order, name
LIMIT $1 OFFSET $2;

-- name: CountAttributes :one
SELECT count(*) FROM product_attributes;

-- name: UpdateAttribute :one
UPDATE product_attributes
SET name = $2, is_filterable = $3, sort_order = $4, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteAttribute :execrows
DELETE FROM product_attributes WHERE id = $1;
