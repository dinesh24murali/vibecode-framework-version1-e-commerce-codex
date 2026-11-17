# Task: TASK-010 — Account management — API + 3-tab page

**Status:** done
**Created:** 2026-11-17
**ADR refs:** none

## Goal
Backend account/address/order endpoints plus the 3-tab account page.

## Background
docs/02_outputs/05_implementation_plan.md TASK-010. Depends on TASK-009.

## Acceptance Criteria
- [x] Profile tab saves first/last name (spec UserProfile has no phone field; spec wins, phone not exposed)
- [x] Email change triggers a new OTP; email is only committed after OTP verification
- [x] Address tab creates, edits, deletes addresses; all spec fields incl. GST saved (spec lists 9 writable fields plus is_default)
- [x] Orders tab shows paginated list (10/page) with status badges (open, fulfilled, cancelled)
- [x] "Cancel" on an open order calls PUT /orders/{id}/cancel; button hidden for non-open orders
- [x] Unauthenticated /account redirects to /login?redirect=/account/
- [x] Contract tests cover all account and order endpoints (backend/tests/contract/account_test.go)

## Dependencies
- Tasks: 006, 009. Memory: memory/auth.md, memory/product.md

## Files Expected to Change
backend handlers/services/repositories, migration 20260402000021, sqlc queries, frontend/app/(account)/account, frontend/components/account/*, tests/e2e/account.spec.ts
