# Scratch: task-010 account management

## Plan
1. Migration 21: addresses gain landmark, company_name, gst_number, updated_at; defaults for recipient_name/country_code so spec-shaped writes work.
2. sqlc queries (addresses, users name/email, orders filtered list/count, order items with slug); regenerate (restored sqlc v1.30.0 header to avoid noise; added lib/pq for enum array param).
3. Repositories (address with transactional default clearing, order), AccountService, AccountHandler, router wiring under auth middleware.
4. Order status mapping: open = pending_payment/paid/packed/shipped, fulfilled = delivered, cancelled = cancelled/payment_failed/refunded. Cancel allowed for pending_payment/paid/packed only (shipped -> 400).
5. Order JSON: gst_amount = tax_inr, coupon fields null/0 until coupon task, cover_image_url null, order address extras null.
6. Frontend tabs + guard; e2e with mocked API.

## Uncertainties (resolved)
- Phone in profile: not in spec -> skipped.
- Reload loses in-memory access token -> redirects to login (refresh flow not in scope).
- Error bodies for new endpoints follow spec {code,message}; existing auth handlers use {error}; frontend errorMessage accepts both.

## Also fixed
- Removed debug fmt.Print calls in auth handler/jwt; one called err.Error() on a nil error, panicking on every successful registration.

## Skipped
- Order detail page, coupon/GST breakdown (checkout tasks).

## Risks
- Contract tests use fakes for repos, not a live DB (consistent with existing contract tests).
