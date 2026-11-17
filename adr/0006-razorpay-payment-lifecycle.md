# ADR-0006: Razorpay payment lifecycle

**Status:** Accepted
**Date:** 2026-11-17
**Deciders:** AI agent (auto-accepted)

---

## Context

TASK-039 adds payments. Orders are created by `POST /checkout/orders` in `pending_payment` with stock already reserved (`reserved += qty`). Razorpay confirms payment asynchronously by webhook, and webhooks are delivered at least once and possibly out of order. The frontend can also retry `initiate` after a flaky network.

## Decision

- **Initiate:** `POST /payments/initiate` (JWT, `Idempotency-Key` required) creates a Razorpay order for the user's own `pending_payment` order and stores a `payments` row (`created`). The response is cached in Redis (24h, keyed by user, order and key); a replay returns it without calling Razorpay. A short `pending` claim (30s) rejects concurrent duplicates with 409.
- **Webhook:** `POST /payments/webhooks/razorpay` is public; the HMAC-SHA256 signature of the raw body (`X-Razorpay-Signature`) is the authentication. Bad signature returns 400 before any parsing or writes. Unknown events and unknown payments are acknowledged with 200 so Razorpay stops retrying.
- **State changes** happen in one transaction that locks the payment row: `payment.captured` sets payment `captured`, order `paid` and deducts stock (`on_hand` and `reserved`); `payment.failed` sets payment `failed` and order `payment_failed`. A captured payment is final, so replays and late failures are no-ops.
- **Failed payments keep their reservation.** Razorpay lets the customer retry on the same order, and the existing `stock_release` job type is the place to expire abandoned reservations.
- **Side effects are asynchronous:** after commit the service enqueues `order_confirmation_email` on the Redis Streams queue with the job id `order_confirmation:<orderID>`, which the worker already deduplicates. The job is enqueued on every captured webhook, including replays, so a crash between commit and publish is healed by Razorpay's retry.
- Keys come from `RAZORPAY_KEY_ID`, `RAZORPAY_KEY_SECRET`, `RAZORPAY_WEBHOOK_SECRET`. Without the key pair initiate answers 503; without the webhook secret every webhook is rejected. The gateway sits behind `PaymentGateway`, so tests use a fake.

## Consequences

### Positive
- Safe under retries and duplicate or reordered webhooks; no money-taken-but-order-unpaid window beyond Razorpay's own retries.
- No real credentials needed to run the test suite.

### Negative
- Abandoned or failed orders hold stock until the `stock_release` job (still a stub) is implemented.
- The idempotency store is payment-specific; TASK-038 may generalise it into middleware.
