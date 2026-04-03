-- name: AssignAttributeToProduct :one
INSERT INTO product_attribute_assignments (product_id, attribute_value_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetAttributesByProduct :many
SELECT pav.*, pa.code AS attribute_code, pa.name AS attribute_name
FROM product_attribute_assignments paa
JOIN product_attribute_values pav ON pav.id = paa.attribute_value_id
JOIN product_attributes pa ON pa.id = pav.attribute_id
WHERE paa.product_id = $1
ORDER BY pa.sort_order, pav.sort_order;

-- name: RemoveAttributeFromProduct :exec
DELETE FROM product_attribute_assignments
WHERE product_id = $1 AND attribute_value_id = $2;
