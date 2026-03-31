# Task: TASK-001 — Monorepo scaffold & local dev environment

**Status:** done
**Created:** 2026-03-31
**ADR refs:** none

---

## Goal

Create the top-level monorepo structure with `backend/` and `frontend/` directories, a full `docker-compose.yml` for all local dev services, and a root `Makefile` with all standard targets.

## Background

This is the foundation task. Nothing in Phase 0 or beyond can begin until this is done. The local dev environment must mirror the production topology (Postgres, Redis, API, Worker, Frontend, MailHog) so that all subsequent tasks can be developed and tested consistently.

See `docs/02_outputs/03_tech_architecture.md` §6 for the full container layout and `docs/02_outputs/05_implementation_plan.md` TASK-001.

## Acceptance Criteria

- [ ] `docker compose up` starts Postgres (5432), Redis (6379), MailHog (1025/8025), api, worker, and frontend with no errors
- [ ] `make dev` starts both backend and frontend dev servers concurrently
- [ ] `.env.example` documents every required environment variable with a one-line description, including: `REDIS_URL`, `DATABASE_URL`, `MAILHOG_HOST`, `RAZORPAY_KEY_ID`, `RAZORPAY_KEY_SECRET`, `JWT_PRIVATE_KEY`, `JWT_PUBLIC_KEY`, `NEXT_PUBLIC_API_URL`
- [ ] `.gitignore` excludes `.env`, build artefacts, `vendor/`, `node_modules/`, `out/`
- [ ] Root `Makefile` has targets: `dev`, `test`, `build`, `seed`, `migrate`, `migrate-status`, `sqlc-gen`, `codegen`, `verify`, `adr`, `task`

## Dependencies

- **Tasks:** none
- **Memory files:** none
- **Docs:** `docs/02_outputs/03_tech_architecture.md` §6 (Infrastructure & Deployment), `docs/02_outputs/05_implementation_plan.md` TASK-001

## Files Expected to Change

- `Makefile` *(create)*
- `.gitignore` *(create)*
- `.env.example` *(create)*
- `docker-compose.yml` *(create)*
- `README.md` *(create)*
- `backend/` *(create directory)*
- `frontend/` *(create directory)*

## Notes

- Docker Compose service names must match the architecture doc exactly: `frontend`, `api`, `worker`, `postgres`, `redis`, `mailhog`
- The `api` and `worker` services can reference placeholder Dockerfiles at this stage — they will be built out in TASK-002 and TASK-040
- Production uses a single EC2 host with Docker Compose; keep the local compose file aligned with that topology
- MailHog replaces SES in local dev (SMTP on 1025, web UI on 8025)

---

## Scratch File

AI: before implementing, create `tasks/active/task-001-monorepo-scaffold-local-dev.scratch.md` with your plan.
See `tasks/AGENTS.md` for the required scratch file format.
