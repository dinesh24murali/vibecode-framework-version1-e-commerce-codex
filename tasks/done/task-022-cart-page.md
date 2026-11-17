# Task: TASK-022 - Cart page (frontend)

**Status:** done
**Created:** 2026-11-17
**ADR refs:** none

## Acceptance Criteria
- [x] Quantity changes are optimistic in cartSlice, then confirmed by the API (rollback on failure)
- [x] Empty cart state with Continue Shopping link
- [x] Cart persists across refresh for authenticated users (server cart; re-login needed as tokens are memory-only)
- [x] Proceed to Checkout disabled with login prompt when unauthenticated
- [x] E2e: add from product page, cart view, quantity update, remove, failure rollback (tests/e2e/cart.spec.ts)
- [x] fe-build, tsc, check-no-fetch, e2e and DOM checks pass (docker)

## Dependencies
Tasks: 021, 007, 008. Memory: memory/product.md, memory/auth.md.
