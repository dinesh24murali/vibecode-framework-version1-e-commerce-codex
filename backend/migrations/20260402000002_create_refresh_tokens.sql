-- +goose Up
CREATE TABLE refresh_tokens (
    id              UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id         UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_family_id TEXT        NOT NULL,
    token_hash      TEXT        NOT NULL,
    user_agent      TEXT,
    ip_address      INET,
    expires_at      TIMESTAMPTZ NOT NULL,
    revoked_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_refresh_tokens_user_expires ON refresh_tokens (user_id, expires_at);

-- +goose Down
DROP TABLE IF EXISTS refresh_tokens;
