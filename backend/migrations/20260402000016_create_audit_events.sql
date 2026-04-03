-- +goose Up
-- actor_user_id has no FK so audit records survive user deletion.
CREATE TABLE audit_events (
    id            UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    order_id      UUID        REFERENCES orders (id) ON DELETE SET NULL,
    actor_user_id UUID,
    event_type    TEXT        NOT NULL,
    payload       JSONB,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_events_order_created ON audit_events (order_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS audit_events;
