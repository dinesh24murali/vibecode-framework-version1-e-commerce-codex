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

-- name: AdminCreateCategory :one
INSERT INTO categories (slug, name, image_url)
VALUES ($1, $2, $3)
RETURNING *;

-- name: AdminUpdateCategory :one
UPDATE categories
SET slug = $2, name = $3, image_url = $4, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: AdminDeleteCategory :execrows
DELETE FROM categories WHERE id = $1;

-- name: CountProductsInCategory :one
SELECT count(*) FROM products WHERE category_id = $1;

-- name: CountCategories :one
SELECT count(*) FROM categories;

-- name: ListCategoriesPaged :many
SELECT * FROM categories
ORDER BY name
LIMIT $1 OFFSET $2;
