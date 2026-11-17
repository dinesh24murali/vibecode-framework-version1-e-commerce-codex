# Scratch: task-039-razorpay-payment-backend
## Plan
1. Add payment paths/schemas to the spec, run orval. 2. Config for Razorpay env. 3. Repos: payment (tx, idempotent Apply). 4. Services: gateway interface + Razorpay HTTP adapter, idempotency store (Redis), PaymentService. 5. Handler + routes (webhook public). 6. worker.Queue publisher. 7. Contract + real DB + Redis tests. 8. ADR-0006.
## Decisions
- Spec had no payment endpoints; spec is authoritative, so it was extended (not silently diverged from).
- 409 for non-pending order / in-flight key, 503 when unconfigured, 404 for other users' orders.
- Failed payment keeps reservation (see ADR-0006).
- Confirmation job enqueued on every captured webhook (idempotent job id).
## Skipped
Frontend Razorpay checkout (later task), real Razorpay calls, stock_release job.
## Risks
Webhook payload shape taken from Razorpay docs, not verified against a live account.
