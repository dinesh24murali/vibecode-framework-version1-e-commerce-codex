# Task: TASK-006 — JWT auth system — backend

**Status:** active
**Created:** 2026-03-31
**ADR refs:** ADR needed — OTP delivery mechanism (stdout in dev vs SMTP/SES in production)

---

## Goal

Implement the full backend authentication pipeline: 3-step OTP registration, email/password login, Facebook OAuth login, JWT access tokens (Ed25519), rotating refresh tokens (Redis-backed), and all password-reset endpoints.

## Background

Auth is the gating dependency for every protected API endpoint. The architecture mandates short-lived access JWTs signed with EdDSA (Ed25519), opaque rotating refresh tokens hashed before persistence, and Redis-backed revocation. JWTs must never be stored in `localStorage` on the client; the refresh token is delivered in an `HttpOnly` cookie.

The password-reset flow (`POST /api/v1/auth/password-reset/*`) was designed in the done task `tasks/done/add-forgot-password-flow.md`. This task wires those three endpoints using the same OTP service — no re-design required.

See `docs/02_outputs/03_tech_architecture.md` §5 and `docs/02_outputs/04_api_spec.yaml` auth endpoints.

## Acceptance Criteria

- [ ] `POST /api/v1/auth/register/init` — sends OTP (logged to stdout in dev); returns `200` even for unknown email to prevent enumeration
- [ ] `POST /api/v1/auth/register/verify-otp` — returns a short-lived `otp_token` (JWT or signed token); expires after 10 minutes
- [ ] `POST /api/v1/auth/register/complete` — requires valid `otp_token`; creates user record and returns `AuthTokenResponse`
- [ ] `POST /api/v1/auth/login` — verifies Argon2id password hash; returns `access_token` + sets `refresh_token` cookie
- [ ] `POST /api/v1/auth/login/facebook` — exchanges Facebook access token for a site JWT
- [ ] `POST /api/v1/auth/refresh` — issues a new access token; invalidates the old refresh token (token rotation)
- [ ] `POST /api/v1/auth/logout` — revokes the refresh token family and all active jtis in Redis
- [ ] All three `/api/v1/auth/password-reset/*` routes are wired and return spec-correct responses
- [ ] `AuthMiddleware` rejects requests without a valid JWT with `401 Unauthorized`
- [ ] Refresh token stored as a hash in Postgres (`refresh_tokens` table); family metadata stored in Redis
- [ ] Contract test in `tests/contract/auth_test.go` covers all 10 auth endpoints against `docs/02_outputs/04_api_spec.yaml`
- [ ] ADR created for OTP delivery mechanism

## Dependencies

- **Tasks:** TASK-005 (schema must exist — `users`, `refresh_tokens`, `consents` tables required)
- **Memory files:** `memory/auth.md` *(read before starting — create if it does not exist)*
- **Docs:** `docs/02_outputs/03_tech_architecture.md` §5 (Auth & Authorization Flow), `docs/02_outputs/04_api_spec.yaml` (auth endpoint schemas), `tasks/done/add-forgot-password-flow.md` (password-reset design)

## Files Expected to Change

- `backend/internal/auth/jwt.go` *(create — Ed25519 sign/verify)*
- `backend/internal/auth/middleware.go` *(create — Gin AuthMiddleware)*
- `backend/internal/auth/otp.go` *(create — OTP generation and verification)*
- `backend/internal/auth/facebook.go` *(create — Facebook token exchange)*
- `backend/internal/handlers/auth_handler.go` *(create)*
- `backend/internal/services/auth_service.go` *(create)*
- `backend/internal/repositories/user_repository.go` *(create)*
- `backend/db/queries/users.sql` *(create or update)*
- `backend/db/queries/refresh_tokens.sql` *(create)*
- `tests/contract/auth_test.go` *(create)*
- `memory/auth.md` *(create — document token lifetimes, cookie settings, revocation strategy)*
- `adr/NNNN-otp-delivery-mechanism.md` *(create)*

## Notes

- JWT signing algorithm: `EdDSA (Ed25519)` — keys loaded from `JWT_PRIVATE_KEY` / `JWT_PUBLIC_KEY` env vars
- Access token lifetime: 15 minutes (document in `memory/auth.md`)
- Refresh token: opaque random value, hashed with SHA-256 before persistence; stored in `HttpOnly; Secure; SameSite=Lax` cookie
- Argon2id parameters: document in `memory/auth.md` (suggested: time=1, memory=64MB, threads=4)
- OTP: 6-digit numeric; valid for 10 minutes; stored in Redis with TTL; single-use
- Follow Repository → Service → Handler layering with interfaces at each boundary (Bridge pattern between repository interface and concrete Postgres implementation)
- Token rotation: on refresh, old refresh token is revoked in Redis before issuing a new one; reuse of a revoked token invalidates the entire family
- Facebook login: validate the Facebook access token with the Graph API before creating/linking the user account

---

## Scratch File

AI: before implementing, create `tasks/active/task-006-jwt-auth-system-backend.scratch.md` with your plan.
See `tasks/AGENTS.md` for the required scratch file format.
