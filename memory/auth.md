# Auth Memory

## Token Lifetimes

| Token | TTL | Storage |
|-------|-----|---------|
| Access token (JWT) | 15 minutes | In-memory / Authorization header |
| Refresh token (opaque UUID) | 7 days | `refresh_tokens` Postgres table (hashed); raw value in HttpOnly cookie |
| OTP token (JWT) | 10 minutes | Stateless JWT — verified by signature only |
| OTP code (6-digit) | 10 minutes | Redis key `otp:<purpose>:<email>` with TTL |

## JWT Details

- **Algorithm:** EdDSA (Ed25519)
- **Private key env var:** `JWT_PRIVATE_KEY` (PEM PKCS#8 encoded)
- **Public key env var:** `JWT_PUBLIC_KEY` (PEM PKIX encoded)
- **Access token claims:** `sub` (user UUID), `role`, `iat`, `exp`
- **OTP token claims:** `sub` (email), `purpose` (`register_otp` | `password_reset_otp`), `iat`, `exp`

## Refresh Token

- Raw value: `uuid.New().String()` (UUID v4)
- Stored in DB as: SHA-256 hex digest
- Cookie: `refresh_token`; `HttpOnly=true`; `Secure=true`; `SameSite=Lax`; path `/`
- Also included in `AuthTokenResponse` JSON body to match API spec contract
- Token rotation: on every `/refresh` call, the old token is revoked and a new one issued in the same family
- Reuse detection: presenting a revoked token triggers `RevokeFamily` — all sessions in the family are invalidated

## Argon2id Parameters

| Parameter | Value |
|-----------|-------|
| Time | 1 |
| Memory | 64 MB (65536 KiB) |
| Threads | 4 |
| Key length | 32 bytes |
| Salt | 16 random bytes |
| Storage format | `hex(salt):hex(hash)` |

## OTP Delivery

- **Dev / Phase 1:** OTP is logged to stdout via `slog.Info` — no email sent.
- **Prod / Phase 2:** Replace log call with AWS SES send (see ADR-0005).
- OTP is single-use: Redis key is deleted immediately after successful verification.

## Schema Notes

- The `users` table stores `full_name` as a single column.
- The API spec `UserProfile` uses `first_name` + `last_name` separately.
- Mapping: write as `first_name + " " + last_name`; read by splitting on the first space.
- This is lossy for names with internal spaces — acceptable for the current schema.

## AuthMiddleware

- Extracts `Authorization: Bearer <token>` header
- Verifies Ed25519 signature and expiry
- Sets `userID` (uuid.UUID) and `userRole` (string) in Gin context
- Returns `401 Unauthorized` on any failure

## Frontend Store (TASK-008)

- `frontend/lib/store/` holds `authSlice` + `cartSlice` in one persisted store (`useStore`, `StoreState`).
- Access/refresh tokens are in-memory only; `partialize` persists only `user`. Cart is never persisted.
- `skipHydration: true` + `useHydrateStore` (mounted via `StoreHydrator` in root layout) rehydrates after mount (ADR-0003).
- `cartSlice.syncWithServer()` uses an injectable `CartApi` (`setCartApi`); TASK-007 must wire the generated client into it.

## Frontend Auth Pages (TASK-009)

- Routes: `/login`, `/signup`, `/signup/verify-otp` (OTP step then name+password step; email passed via `?email=`). Pages using `useSearchParams` are wrapped in Suspense (static export).
- `?redirect=` is sanitised by `lib/auth/helpers.ts#safeRedirect` (same-origin relative paths only).
- Facebook: loads the FB SDK script on click; needs `NEXT_PUBLIC_FACEBOOK_APP_ID`. Not covered by e2e (popup).
- e2e (`tests/e2e/auth.spec.ts`, config `tests/e2e/playwright.config.ts`) mocks `/api/v1/auth/*` with `page.route`; the static verify stack has no API.
