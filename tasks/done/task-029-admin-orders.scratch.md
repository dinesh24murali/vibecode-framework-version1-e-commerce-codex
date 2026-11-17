# Scratch: task-029-admin-orders

## Plan
1. sqlc AdminListOrders/AdminCountOrders; AdminOrderRepository extends OrderRepository.
2. AdminOrderService (reuses SpecStatus/dbStatuses, joins customer), handler (reuses orderJSON + customer fields), routes.
3. Spec: add optional customer_name/customer_email to Order (needed for the modal); client regenerated.
4. Orders page with status tabs + pagination, OrderTable, OrderDetailModal (fetches by id).
5. Real-DB contract test + mocked e2e.

## Uncertainties / decisions
- Transition rules: only open-class orders can change; fulfilled requires paid/packed/shipped (unpaid -> 422); `open` is never a valid target; invalid -> 422.
- fulfilled = DB `delivered`; cancelled = DB `cancelled`.
- Cancelling does not release reserved stock (same as the customer cancel path); the stock_release job does not exist yet.

## What Will Be Skipped
Refund handling on cancel, status history/audit, order search.

## Risks
Reserved stock stays reserved on admin cancel until a release job exists.
