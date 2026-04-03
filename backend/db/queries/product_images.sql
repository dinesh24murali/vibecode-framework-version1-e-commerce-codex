-- name: ListImagesByProduct :many
SELECT * FROM product_images
WHERE product_id = $1
ORDER BY created_at DESC;

-- name: CreateProductImage :one
INSERT INTO product_images (product_id, image_url)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteProductImage :exec
DELETE FROM product_images
WHERE id = $1;
