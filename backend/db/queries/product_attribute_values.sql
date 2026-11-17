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

-- name: UpdateAttributeValue :one
UPDATE product_attribute_values
SET value = $3, slug = $4, sort_order = $5, updated_at = now()
WHERE id = $1 AND attribute_id = $2
RETURNING *;

-- name: DeleteAttributeValue :execrows
DELETE FROM product_attribute_values WHERE id = $1 AND attribute_id = $2;

-- name: CountAssignmentsForValue :one
SELECT count(*) FROM product_attribute_assignments WHERE attribute_value_id = $1;
