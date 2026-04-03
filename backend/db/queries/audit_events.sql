-- name: CreateAuditEvent :one
INSERT INTO audit_events (order_id, actor_user_id, event_type, payload)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListAuditEventsByOrder :many
SELECT * FROM audit_events
WHERE order_id = $1
ORDER BY created_at DESC;
