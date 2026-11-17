# Scratch: task-025-admin-auth

## Plan
1. AdminMiddleware (authenticate + role check; refactor Middleware to share `authenticate`).
2. AuthService.AdminLogin / AuthHandler.AdminLogin; route public `/admin/auth/login`, other `/admin/*` behind AdminMiddleware.
3. adminAuthSlice (in-memory tokens), apiFetch uses admin token for `/api/v1/admin/*` (except auth).
4. AdminGuard in (admin) layout, login page + form, dashboard placeholder.
5. Contract + middleware test, e2e spec.

## Files to Create or Modify
See task file.

## Uncertainties
- Non-admin credentials on admin login: decided 401 (same as bad password) to avoid account enumeration.
- Admin tokens memory-only: reload requires re-login (same as customer, no refresh flow).

## What Will Be Skipped
Admin seed user (TASK-030), refresh flow, admin nav/dashboard content.

## Risks
Calling c.Next() inside a shared middleware would run handlers early; avoided via `authenticate`.
