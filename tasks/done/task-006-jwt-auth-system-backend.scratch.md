# Scratch: TASK-006 — JWT auth system — backend

## Plan (step-by-step)

1. Add `github.com/golang-jwt/jwt/v5` to `go.mod` (Ed25519/EdDSA support)
2. Create `backend/internal/auth/jwt.go` — load Ed25519 keys from PEM env vars, sign/verify access tokens (15 min) and OTP tokens (10 min)
3. Create `backend/internal/auth/otp.go` — generate 6-digit crypto-random OTP, store in Redis with TTL, single-use verify+delete
4. Create `backend/internal/auth/facebook.go` — exchange Facebook access_token via Graph API (`/me?fields=id,email,first_name,last_name`)
5. Create `backend/internal/auth/middleware.go` — Gin `AuthMiddleware`: extract Bearer token, verify Ed25519 sig, set `userID` in context; return 401 on failure
6. Create `backend/internal/repositories/user_repository.go` — `UserRepository` interface + `pgUserRepository` backed by sqlc `Queries`
7. Create `backend/internal/repositories/refresh_token_repository.go` — `RefreshTokenRepository` interface + pgImpl
8. Create `backend/internal/services/auth_service.go` — all business logic: register flow, login, Facebook login, refresh, logout, password-reset flow
9. Create `backend/internal/handlers/auth_handler.go` — bind JSON, delegate to service, return spec-correct JSON responses
10. Update `backend/internal/server/router.go` — add `RouterDeps` struct, wire `/api/v1/auth/*` routes, apply `AuthMiddleware` to logout endpoint
11. Update `backend/internal/server/server.go` — pass deps through to `newRouter`
12. Create `memory/auth.md` — document token lifetimes, Argon2id params, cookie settings, revocation strategy
13. Create `adr/0005-otp-delivery-mechanism.md` — document stdout (dev) vs SMTP/SES (prod) decision
14. Create `tests/contract/auth_test.go` — 10 endpoint contract tests against API spec

## Files to Change

**Create:**
- `backend/internal/auth/jwt.go`
- `backend/internal/auth/otp.go`
- `backend/internal/auth/middleware.go`
- `backend/internal/auth/facebook.go`
- `backend/internal/repositories/user_repository.go`
- `backend/internal/repositories/refresh_token_repository.go`
- `backend/internal/services/auth_service.go`
- `backend/internal/handlers/auth_handler.go`
- `memory/auth.md`
- `adr/0005-otp-delivery-mechanism.md`
- `tests/contract/auth_test.go`

**Modify:**
- `backend/internal/server/router.go` — add RouterDeps, wire auth group
- `backend/internal/server/server.go` — thread deps into newRouter
- `backend/go.mod` / `backend/go.sum` — add golang-jwt/jwt/v5
- `memory/_index.md` — add auth.md entry

## Uncertainties / Open Questions

1. **full_name vs first_name/last_name mismatch** — the DB schema has a single `full_name` column but the API spec defines `first_name` + `last_name` separately in `UserProfile` and `RegisterCompleteRequest`. Decision: concatenate on write (`first_name + " " + last_name`), split on first space on read. This is lossy for names with multiple spaces but acceptable given the existing schema. Note this in `memory/auth.md`.

2. **Facebook env vars** — need `FACEBOOK_APP_ID` and `FACEBOOK_APP_SECRET` to verify the token with the `/debug_token` endpoint. Simpler alternative (MVP): just call `/me` with the user-provided token and trust Facebook's validation (the token is signed by Facebook and expires). Using `/me` without `/debug_token` is weaker (token stolen from another app would work). Decision for this task: use `/me` only for MVP simplicity; flag as a known limitation in the ADR.

3. **OTP token type** — using a standard JWT with a `purpose` claim (`"register_otp"` or `"password_reset_otp"`) and `sub` = email. Signed with the same Ed25519 key as access tokens but with 10 min expiry. Verified in `register/complete` and `password-reset/complete` handlers before proceeding.

4. **Refresh token in cookie vs body** — the API spec `AuthTokenResponse` includes `refresh_token` in the JSON body AND the task notes say HttpOnly cookie. Decision: set the cookie AND include in the response body to match the spec schema exactly. The cookie is the secure path; the body field satisfies the OpenAPI contract.

5. **`updatePasswordHash` sqlc query** — the password-reset/complete endpoint needs to update a user's password_hash. Checking if this query exists in sqlc or needs to be added. Will add `UpdatePasswordHash` query if missing.

## What Will Be Skipped

- SMTP/SES email delivery — OTP is logged to stdout in dev (this task). Production email delivery is deferred to the infra/email task.
- Facebook `/debug_token` endpoint validation — using `/me` only for MVP (noted as limitation).
- Admin auth endpoints (`/api/v1/admin/auth/login`) — that is a separate task scope.
- Rate limiting on auth endpoints — deferred to a future hardening task.

## Risks

1. Ed25519 PEM parsing — `golang-jwt/jwt/v5` requires `crypto/ed25519` private key (not `edwards25519` package directly). Must parse PEM with `x509.ParsePKCS8PrivateKey` or `encoding/pem` + `ed25519` stdlib. Test this carefully.
2. `pqtype.Inet` for IP address — the `CreateRefreshTokenParams` requires a `pqtype.Inet`. Need to parse the client IP from `c.ClientIP()` correctly using `net.ParseIP`.
3. Token rotation reuse detection — if a revoked refresh token is re-used, the entire family must be revoked. This requires the `RevokeTokenFamily` sqlc query (already generated) AND checking `revoked_at` before issuing a new token.
4. sqlc `UpdatePasswordHash` may be missing — need to check `backend/db/sqlc/users.sql.go` and add the query if absent.
