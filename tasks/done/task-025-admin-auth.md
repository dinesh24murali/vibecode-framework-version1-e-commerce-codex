# Task: task-025-admin-auth

**Status:** done
**Created:** 2026-11-17
**ADR refs:** none

## Goal
Admin-only authentication: AdminMiddleware, `POST /api/v1/admin/auth/login`, and a guarded admin area with its own login page.

## Background
TASK-025 in docs/02_outputs/05_implementation_plan.md. Gate for admin tasks 026-029.

## Acceptance Criteria
- [x] Admin login returns a JWT with `role: admin` claim
- [x] Non-admin JWTs receive `403` on all `/api/v1/admin/*` routes
- [x] Admin login page uses a separate layout with no public site chrome
- [x] Unauthenticated access to any `/admin/*` page redirects to `/admin/login`
- [x] Contract test validates the `POST /api/v1/admin/auth/login` response
- [x] E2e test covers admin login happy path and redirect behaviour
- [x] `make verify` passes

## Dependencies
- **Tasks:** 006
- **Memory files:** memory/auth.md

## Files Expected to Change
backend/internal/auth/{admin_middleware,middleware}.go, services/auth_service.go, handlers/auth_handler.go, server/router.go, backend/tests/contract/admin_auth_test.go, frontend admin layout/login/guard/slice, lib/api/client.ts, tests/e2e/admin-login.spec.ts

## Notes
Contract tests live in backend/tests/contract (existing convention).
