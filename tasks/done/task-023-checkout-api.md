# Task: TASK-023 - Checkout API

**Status:** done
**Created:** 2026-11-17
**ADR refs:** none

## Acceptance Criteria
- [x] POST /checkout/calculate returns CheckoutSummary with GST line
- [x] POST /checkout/coupons/validate: 200 with discount for valid codes; unknown/expired codes return 200 is_valid=false (spec has no 404; spec wins)
- [x] POST /checkout/orders creates the order and atomically reserves stock (real-DB concurrency test)
- [x] Empty cart returns 400
- [x] Coupon applied before tax
- [x] Contract tests cover all three endpoints, flat (SUPERHIT 150) and percentage (FIRSTTIME 10%) coupons and edge cases (backend/tests/contract/checkout_test.go)
- [x] go build/vet/test/lint pass (docker); migration 24 applied

## Dependencies
Tasks: 021. Memory: memory/product.md.
