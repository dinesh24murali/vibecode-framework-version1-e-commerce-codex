-- name: ListValuesByAttribute :many
SELECT * FROM product_attribute_values
WHERE attribute_id = $1
ORDER BY sort_order, value;

-- name: GetAttributeValueByID :one
SELECT * FROM product_attribute_values
WHERE id = $1;

-- name: CreateAttributeValue :one
INSERT INTO product_attribute_values (attribute_id, value, slug, sort_order)
VALUES ($1, $2, $3, $4)
RETURNING *;
