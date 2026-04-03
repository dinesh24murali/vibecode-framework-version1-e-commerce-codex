-- name: CreateConsent :one
INSERT INTO consents (user_id, consent_type, policy_version, granted_at)
VALUES ($1, $2, $3, now())
RETURNING *;

-- name: GetConsentsByUser :many
SELECT * FROM consents
WHERE user_id = $1
ORDER BY granted_at DESC;

-- name: WithdrawConsent :exec
UPDATE consents
SET withdrawn_at = now()
WHERE user_id = $1 AND consent_type = $2 AND withdrawn_at IS NULL;
