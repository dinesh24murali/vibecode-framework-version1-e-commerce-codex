# Memory: Product

## Scope

- `e-commerce-site` is a single-store direct-to-consumer bookstore, not a marketplace.
- v1 includes storefront browsing, customer accounts, cart, checkout without payment gateway, and an admin portal.
- The storefront taxonomy currently includes 12 categories: Colouring Books, Drawing Books, Journals, Calendars, Sketch Pads, Planners, Stained Art, Fabrii Art, Note Book, Sticker Book, Sticker Colouring Book, and Activity Cards.

## Core Product Rules

- User accounts are role-based with `admin` and `customer` as the only valid roles.
- Product listing filters must support category plus extensible product attributes such as level, theme, and size.
- Product search and filter state must persist in URL query parameters.
- Product records store a default thumbnail in `products.image_url` and additional gallery images in `product_image`.
- Products may carry an optional `discount_percent` value capped at 100.
- Customer account area contains My profile, My address, and My orders.
- Customers can maintain multiple addresses and choose one default address.
- Order shipment details must be copied into `order_address` at checkout so historical orders are not affected by later profile-address edits.
- Cart state is stored as one row per user-product pair in a single `cart` table.
- Order statuses are limited to `Open`, `Fulfilled`, and `Cancelled`.
- Customers may cancel only `Open` orders and may provide an optional cancellation reason.
- Checkout supports one coupon code per order and no online payment in v1.
- Coupon usage history is tracked in `coupon_user`.
- Tax is calculated at individual product level using the tax rate stored on each product.
- Contact inquiries must both persist in the database and send a notification email.
- Facebook login is controlled by feature flag and may be hidden until platform approval is complete.
- All admin product fields are mandatory in v1.

## Seed Data

- Default coupons: `SUPERHIT` for flat 150 discount and `FIRSTTIME` for 10% discount.
- Initial seed data must include one default admin user and 5 dummy products across 3 categories.
- Seed extensible product attributes for at least `level`, `theme`, and `size` so catalog filters work without hard-coded schema columns.

## Operational Expectations

- Order placement triggers customer confirmation email and manual phone follow-up by the business.
- `Printing supplies` exists as a blank public placeholder page in v1.

## Open Product Decisions

- **Variant model deferred (v1 decision):** The current generic product-attribute model (level, theme, size as labels) is intentionally kept for v1. If any future product attribute changes SKU, inventory, or price semantics (e.g. same book sold in A4 and A3 at different prices with separate stock), it should move from the generic attribute model into a true variant model. Migration would involve: new `product_variants` table, moving `sku`/`price_inr` from `products` to variants, re-pointing `inventory_items` and `cart` to variants, and adding `product_variant_attribute_values`.
