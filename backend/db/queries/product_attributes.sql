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
