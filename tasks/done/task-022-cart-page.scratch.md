# Scratch: task-022-cart-page
## Plan
1. CartPageContent fetches the server cart via generated client (needs product name/image, which cartSlice lines lack). 2. Optimistic update on page state + slice, rollback on error. 3. CartLineItem, CartSummary. 4. AddToCartButton writes to server when signed in. 5. e2e with stateful API fake.
## Decisions
- Unauthenticated users get a login prompt (no product data for local-only lines).
- Summary: subtotal = list price, discount = difference, total = server subtotal; tax at checkout.
## Skipped
Checkout page; local-cart merge on login beyond existing syncWithServer.
## Risks
Memory-only tokens: reload loses auth until re-login.
