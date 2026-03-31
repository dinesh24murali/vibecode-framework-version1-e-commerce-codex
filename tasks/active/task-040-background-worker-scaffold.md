# Task: TASK-040 — Background worker scaffold

**Status:** active
**Created:** 2026-03-31
**ADR refs:** none

---

## Goal

Scaffold the Go background worker as a separate binary (`cmd/worker`) that connects to Postgres and Redis, consumes jobs from Redis Streams, and dispatches to registered handlers — ready for Phase 2 job handlers to be plugged in.

## Background

The worker must be a separate process from the API even on a single host, so that synchronous traffic spikes do not starve background work (email, payment reconciliation, stock release, audit fan-out). Redis Streams with consumer groups (`XREADGROUP`) is the chosen queue mechanism for the single-EC2 phase. All job handlers must be idempotent — replaying a job must not create duplicate side effects.

See `docs/02_outputs/03_tech_architecture.md` §2 (Background Worker), §2 (Queue) and `docs/02_outputs/05_implementation_plan.md` TASK-040.

## Acceptance Criteria

- [ ] `go build ./cmd/worker` succeeds with zero errors
- [ ] Worker connects to Redis and Postgres on startup; exits non-zero if either is unreachable
- [ ] Worker consumes from a Redis Stream using `XREADGROUP` and dispatches to registered handlers by job type
- [ ] Failed jobs are requeued with exponential backoff; after 3 retries the job moves to a dead-letter stream
- [ ] Successful jobs are acknowledged with `XACK`
- [ ] Worker runs as a separate `worker` service in `docker-compose.yml`
- [ ] `golangci-lint run` passes
- [ ] Job handler interface is defined so Phase 2 handlers (`order_confirmation_email`, `payment_reconciliation`, `stock_release`, `audit_fan_out`) can be registered without modifying dispatcher core

## Dependencies

- **Tasks:** TASK-004 (Postgres pool and Redis client must exist)
- **Memory files:** `memory/infra.md` (check Redis AOF and stream configuration)
- **Docs:** `docs/02_outputs/03_tech_architecture.md` §2 (Background Worker, Queue), `docs/02_outputs/05_implementation_plan.md` TASK-040

## Files Expected to Change

- `backend/cmd/worker/main.go` *(create)*
- `backend/internal/worker/worker.go` *(create — main run loop)*
- `backend/internal/worker/dispatcher.go` *(create — job type → handler registry)*
- `backend/internal/worker/handlers/` *(create directory with `.gitkeep` — Phase 2 handlers go here)*
- `docker-compose.yml` *(update — add `worker` service)*

## Notes

- The worker binary is at `backend/cmd/worker/main.go`; the API binary is at `backend/cmd/api/main.go` (or `backend/main.go` — align with TASK-002's decision)
- Redis Streams consumer group name: use a constant like `"workers"` so all worker replicas share the same group
- Exponential backoff: start at 1s, double each retry, cap at 60s; after 3 retries write to `dead-letter` stream
- Idempotency: each handler receives a job ID; handlers must check for prior completion before executing side effects (use Redis or Postgres idempotency key table)
- Phase 2 job types to pre-register as stubs: `order_confirmation_email`, `payment_reconciliation`, `stock_release`, `audit_fan_out`

---

## Scratch File

AI: before implementing, create `tasks/active/task-040-background-worker-scaffold.scratch.md` with your plan.
See `tasks/AGENTS.md` for the required scratch file format.
