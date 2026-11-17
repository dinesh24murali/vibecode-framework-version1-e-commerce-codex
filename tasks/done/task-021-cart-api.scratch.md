# Scratch: task-021-cart-api
## Plan
1. sqlc queries by item id. 2. Repo/service/handler layered. 3. Wire protected routes. 4. Contract test with fakes.
## Uncertainties / decisions
- Spec says unit_price is price at add time; schema has no price column, so current discounted price is used.
- Spec lists 422 for validation; project convention (apiError 400) kept; qty 0 -> 400 as required by the plan.
- Lines for inactive products are omitted from the cart.
## Skipped
Tax/coupons (TASK-023). Stock checks.
## Risks
sqlc version drift (newer sqlc rewrote headers) - only changed files kept, headers restored.
