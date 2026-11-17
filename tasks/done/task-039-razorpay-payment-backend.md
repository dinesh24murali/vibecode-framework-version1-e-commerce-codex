# Task: TASK-039 - Razorpay payment integration (backend)

**Status:** done
**Created:** 2026-11-17
**ADR refs:** ADR-0006

## Acceptance Criteria
- [x] POST /payments/initiate returns a Razorpay order id (fake gateway in tests)
- [x] Webhook with invalid signature returns 400 immediately
- [x] payment.captured transitions order to paid and updates payments (real-DB test)
- [x] payment.failed transitions order to payment_failed
- [x] Initiate is idempotent: replaying the Idempotency-Key returns the cached response (Redis test + contract test)
- [x] No secrets hardcoded; keys read from env (.env.example already lists them)
- [x] Contract tests cover both endpoints and signature validation (backend/tests/contract/payments_test.go)
- [x] go build/vet/test/lint, fe tsc + static build, no-raw-fetch pass (docker)

## Dependencies
Tasks: 023. Memory: memory/product.md, memory/infra.md.
