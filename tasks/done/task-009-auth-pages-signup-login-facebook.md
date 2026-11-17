# Task: TASK-009 — Auth pages — signup, login, Facebook

**Status:** done
**Created:** 2026-11-17
**ADR refs:** ADR-0003

## Goal
Build the 3-step signup flow, login page and Facebook login on the static frontend.

## Background
See docs/02_outputs/05_implementation_plan.md TASK-009. Depends on TASK-006 (backend auth), TASK-007 (generated client), TASK-008 (authSlice).

## Acceptance Criteria
- [x] Full 3-step signup flow (email -> OTP -> password) implemented against generated client
- [x] Invalid OTP shows an inline error message; does not clear the OTP input
- [x] Facebook login button initiates the OAuth popup and resolves to a logged-in session
- [x] Auth tokens stored in authSlice (in memory; user profile persisted per memory/auth.md)
- [x] `?redirect=` is honoured after successful login
- [x] E2e test in tests/e2e/auth.spec.ts covers signup and login happy paths
- [x] make verify passes

## Dependencies
- **Tasks:** 006, 007, 008
- **Memory files:** memory/auth.md
- **Docs:** docs/02_outputs/05_implementation_plan.md TASK-009

## Files Expected to Change
frontend/app/(auth)/login/page.tsx, signup/page.tsx, signup/verify-otp/page.tsx, frontend/components/auth/*, tests/e2e/auth.spec.ts

## Notes
Frontend tasks must pass make verify.
