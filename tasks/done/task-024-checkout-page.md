# Task: TASK-024 - Checkout page (frontend)

**Status:** done
**Created:** 2026-11-17
**ADR refs:** none

## Acceptance Criteria
- [x] Empty-cart redirect to /cart (server cart decides, since the slice is empty right after login)
- [x] Unauthenticated redirect to /login/?redirect=/checkout/
- [x] Coupon shows live valid/invalid feedback on blur; total updates when a valid coupon is applied
- [x] Successful order placement navigates to /account/?tab=orders (account page honours ?tab=)
- [x] E2e covers add to cart -> checkout -> coupon -> place order, plus guards, no-address and failed-order cases (tests/e2e/checkout.spec.ts)
- [x] tsc, static build, no-raw-fetch, e2e (39) and DOM checks pass (docker)

## Dependencies
Tasks: 023, 022. Memory: memory/product.md, memory/auth.md.
