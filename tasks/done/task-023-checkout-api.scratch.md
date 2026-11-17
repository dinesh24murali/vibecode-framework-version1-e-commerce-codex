# Scratch: task-023-checkout-api
## Plan
1. Migration: orders.coupon_code, coupon_discount_inr. 2. sqlc: ReserveStockIfAvailable (conditional UPDATE), CreateOrder with coupon cols. 3. Repos: coupon, checkout (single tx: reserve, order, items, address, coupon use, clear cart). 4. Services: CouponService, OrderBuilder, CheckoutService. 5. Handler + routes. 6. Contract tests with fakes + real-DB concurrency test.
## Decisions
- Spec wins: validate returns 200 is_valid=false for bad codes; invalid coupon on calculate/order is 422; empty cart 400; insufficient stock and bad address 422.
- A coupon can be redeemed once per user (coupon_users unique), which also gives FIRSTTIME its semantics.
- Tax: 18% on (subtotal - discount); all items taxable (no taxable flag in schema). Discount capped at subtotal.
- Stock: reserved += qty at order time (available = on_hand - reserved); TASK-039 deducts on payment.
- product_sku uses product slug (products.sku not exposed in summaries).
## Skipped
Shipping cost (0), seeding SUPERHIT/FIRSTTIME (TASK-030), frontend checkout page.
## Risks
sqlc version drift: headers restored to v1.30.0.
