-- name: CreateContactQuery :one
INSERT INTO contact_queries (name, email, phone, subject, message)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListContactQueries :many
SELECT * FROM contact_queries
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
