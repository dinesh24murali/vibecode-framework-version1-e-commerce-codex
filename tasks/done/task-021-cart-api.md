# Task: TASK-021 — Cart API

**Status:** done
**Created:** 2026-11-17
**ADR refs:** none

## Goal
Authenticated cart endpoints with server-side totals.

## Acceptance Criteria
- [x] Unauthenticated requests return 401
- [x] Adding an existing product increments its quantity (no duplicate line)
- [x] DELETE /cart clears; GET returns an empty cart after
- [x] Cart total computed server-side with per-product discounts
- [x] PUT /cart/items/{id} with quantity 0 returns 400
- [x] Contract tests cover all 5 endpoints (backend/tests/contract/cart_test.go)
- [x] go build/vet/test/lint pass (docker)

## Dependencies
Tasks: 006, 005. Memory: memory/product.md.

## Files
cart_handler.go, cart_service.go, cart_repository.go, cart.sql (+sqlc), router.go, ProductService.SummariesByIDs, models Cart/CartItem, tests.
