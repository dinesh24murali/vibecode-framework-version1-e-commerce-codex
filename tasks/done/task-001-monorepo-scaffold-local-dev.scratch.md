# Scratch: TASK-001 — Monorepo scaffold & local dev environment

## Plan

1. Create `.gitignore` at repo root
2. Create `.env.example` with all required env vars documented
3. Create `docker-compose.yml` with all 6 local dev services
4. Update root `Makefile` — add the missing targets (`dev`, `test`, `build`, `seed`, `migrate`, `migrate-status`, `sqlc-gen`, `codegen`) while keeping the existing ones (`scaffold`, `verify`, `adr`, `task`, `check-env`, `init`, `help`)
5. Update `README.md` with basic project setup instructions

## Files to Change

- `.gitignore` *(create)*
- `.env.example` *(create)*
- `docker-compose.yml` *(create)*
- `Makefile` *(update — add missing targets)*
- `README.md` *(update)*

## Uncertainties / Open Questions

- The `api` and `worker` Docker images don't exist yet (built in TASK-002 / TASK-040). The compose file will reference placeholder Dockerfiles with a simple `FROM scratch` or `FROM golang:1.26-alpine` base so `docker compose up` doesn't fail.
- The `frontend` image also doesn't exist. Same approach: placeholder Dockerfile.
- Port for the API: architecture doc says `/api/v1` is served from the backend — using `8080` internally is standard for Go services; expose as `8080:8080`.
- Frontend dev port: architecture doc mentions `3000`; expose as `3000:3000`.

## What Will Be Skipped

- Actual Go or Next.js source code — that belongs to TASK-002 and TASK-003
- Nginx or reverse proxy setup — not needed for local dev
- SOPS/age secrets encryption — that is a production concern (TASK-004+)
- GitHub Actions CI — listed as future work in the architecture doc

## Risks

- If placeholder Dockerfiles use `FROM scratch`, the containers won't start meaningfully. Using `FROM golang:1.26-alpine` with a `sleep infinity` CMD is safer for `docker compose up` to pass without errors while real code is absent.
- The `make verify` target calls `bash verify/scripts/run-e2e.sh` — if that script is missing or empty it may fail. Will ensure it exits 0 gracefully when no tests exist yet.
