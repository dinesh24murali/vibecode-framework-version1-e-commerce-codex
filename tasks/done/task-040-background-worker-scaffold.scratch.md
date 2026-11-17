# Scratch: task-040 background worker

## Plan
1. dispatcher.go: Handler interface + HandlerFunc + Dispatcher registry (Register/Dispatch); core never changes when handlers are added.
2. worker.go: Redis Streams consumer (stream `jobs`, group `workers`, XREADGROUP). Success -> XACK. Failure -> if attempts < 3, XADD a copy with attempts+1 after backoff (1s doubling, cap 60s) via time.AfterFunc, then XACK original; else XADD to `jobs:dead-letter` and XACK. Idempotency: Redis key `job:done:<id>` checked before and set after success. Reclaims stale pending via XAUTOCLAIM.
3. handlers/stubs.go: register four Phase 2 job types as stub handlers.
4. cmd/worker/main.go: reads DATABASE_URL/REDIS_URL directly (config.Load requires JWT keys the worker does not need), fast-fails if Postgres/Redis unreachable (log.Fatalf = non-zero), graceful shutdown.
5. Tests: dispatcher unit tests; worker integration tests against compose Redis (skip if REDIS_URL unset).
6. docker-compose worker service already exists; verify it builds. Contract test criterion N/A (no API endpoints).

## Files
backend/cmd/worker/main.go, internal/worker/{dispatcher,worker}.go + tests, internal/worker/handlers/stubs.go

## Uncertainties (resolved)
- Backoff via AfterFunc: crash before re-add leaves original pending, recovered by XAUTOCLAIM. Acceptable.
- Retry count: 3 retries after the first attempt, then dead-letter.

## Skipped
- Real handler logic (Phase 2).

## Risks
- Lint strictness (errcheck).
