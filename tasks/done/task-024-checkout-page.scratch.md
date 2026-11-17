# Scratch: task-024-checkout-page
## Plan
1. AddressSelector, CouponInput, CheckoutSummary components. 2. CheckoutPageContent: guards, load cart+addresses, calculate on address/coupon change, place order. 3. /checkout page. 4. Account page reads ?tab=. 5. e2e with API mocks.
## Decisions
- Empty-cart guard uses GET /cart (slice is empty after login because tokens are memory-only).
- ?tab= read in an effect from window.location (static export; avoids a Suspense boundary).
- Payment is not started here: the order stays pending_payment; Razorpay Checkout JS is a later task.
## Skipped
Razorpay checkout, adding addresses inline.
## Risks
Reload on /checkout/ loses the token and bounces to login (known memory-only token behaviour).
