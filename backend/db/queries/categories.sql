-- name: ListCategories :many
SELECT * FROM categories
ORDER BY name;

-- name: GetCategoryBySlug :one
SELECT * FROM categories
WHERE slug = $1;

-- name: GetCategoryByID :one
SELECT * FROM categories
WHERE id = $1;

-- name: CreateCategory :one
INSERT INTO categories (slug, name, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateCategory :one
UPDATE categories
SET name        = $2,
    description = $3
WHERE id = $1
RETURNING *;
