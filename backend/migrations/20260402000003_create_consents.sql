-- +goose Up
CREATE TABLE consents (
    id              UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id         UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    consent_type    TEXT        NOT NULL,
    policy_version  TEXT        NOT NULL,
    granted_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    withdrawn_at    TIMESTAMPTZ
);

CREATE INDEX idx_consents_user_type_granted ON consents (user_id, consent_type, granted_at DESC);

-- +goose Down
DROP TABLE IF EXISTS consents;
