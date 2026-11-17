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

## Account and Orders API (TASK-010)

- Addresses (spec fields): address_line1/2, landmark, city, state, pincode, company_name, gst_number, phone, is_default. DB columns: line1, line2, landmark, city, state, postal_code, company_name, gst_number, phone_e164. Setting is_default clears other defaults in one transaction.
- Order statuses in the API are open/fulfilled/cancelled, mapped from the DB enum (see `services.SpecStatus`). Only pending_payment, paid and packed orders can be cancelled.
- Order JSON `gst_amount` comes from `tax_inr`; `coupon_code`/`coupon_discount` are null/0 until coupons are stored on orders.
- Email change: OTP purpose `email_change_otp`, keyed by the new email; email only changes after confirm.

## Attributes API (TASK-011)

- `GET /api/v1/attributes` lists only `is_filterable` attributes, each with ordered values; `GET /attributes/{id}/values` returns 404 for unknown attribute ids (empty list when it has no values). Both are public. Wire names follow the spec: `code`, `value_slug`.

## Catalogue API (TASK-012)

- `GET /products` query params follow the spec: `q`, `category[]` (slugs, OR), `attr_value_id[]` (OR within an attribute, AND across attributes), `sort_by` (price_asc|price_desc|newest|popular; popular = newest for now), `page`, `limit` (default 20, max 100). Response: `data`, `meta` (total, page, limit, total_pages), `filters`.
- Only `is_active = true` products are public (products has no deleted_at). Detail `stock` = on_hand - reserved. Category per product is single in the DB, returned as a one-element `categories` array.
- Filter SQL lives in `db/queries/catalog.sql`; verified by the real-DB test in `internal/repositories/product_repository_test.go`.

## Landing page (TASK-014)

- `GET /api/v1/landing`: `carousel_slides` (table `carousel_slides`, active only, max 5, ordered by sort_order) and `featured_products` (8 most recent active products; no featured flag exists).
- Frontend `app/(shop)/page.tsx` fetches the payload at build time; if the API is unreachable then, `LandingContent` fetches it on mount. Empty slides show a fallback banner; no products shows "Products are temporarily unavailable."

## Product list page (TASK-015)

- `/products/` URL params: repeated `category` (slugs), repeated `attr` (attribute value UUIDs), `q`, `sort` (omitted when newest), `page` (omitted when 1). `filterSlice` in the root store mirrors them; `useProductFilterUrlSync` keeps URL and store in sync and gates rendering until hydrated. Filters are not persisted.

## Product detail page (TASK-016)

- `/products/[slug]/` is pre-rendered by `generateStaticParams` (pages through `GET /products`). The spec only has `GET /products/{id}`, so slug is resolved to id via the list API (build time, and client-side via `?q=<slug>` when the page was built without data). If the API is down at build, a single `sample-product` shell is generated and resolves from the URL on the client. Zoom is custom (hover scale 2x, pinch up to 4x); no image-zoom dependency.

## Navbar (TASK-017)

- `components/layout/Navbar` is mounted in `app/(shop)/layout.tsx`. Menu content is static in `lib/constants/navigation.ts` (12 categories, 11 printing products, 13 inhouse services). Category links use `/products/?category=<kebab-case name>`, so seeded category slugs must be kebab-case of the names. Search goes to `/products/?q=`.

## Footer (TASK-018)

- `components/layout/Footer` is mounted in the (shop) layout. Social URLs come from `NEXT_PUBLIC_FACEBOOK_URL`, `NEXT_PUBLIC_INSTAGRAM_URL`, `NEXT_PUBLIC_WHATSAPP_URL` via `lib/constants/social.ts` (reuse on the contact page).

## Contact (TASK-020)

- `POST /api/v1/contact` returns 200 + MessageResponse (spec), requires name/email/subject/message, optional phone (422 on bad email/phone, 400 on missing fields). Enquiry stored in `contact_queries`; email goes to `CONTACT_NOTIFY_EMAIL` over `SMTP_HOST:SMTP_PORT` when set, failures are only logged. Page address comes from `NEXT_PUBLIC_BUSINESS_ADDRESS`.

## Cart API (TASK-021)
Cart rows are per (user, product). Cart.id is the user id. Totals use the current discounted price (no price snapshot column). Adding an existing product sums quantities; quantity < 1 is a 400. Tax and coupons belong to checkout (TASK-023).

## Cart page (TASK-022)
`/cart/` loads the server cart with the generated client (needs product details) and mirrors it into cartSlice. Quantity edits are optimistic and roll back on API error. Add to Cart calls the server when signed in, otherwise only the local slice. Tokens are memory-only, so tests must reach /cart/ by client-side navigation after login.

## Checkout API (TASK-023)
Pricing: subtotal (server-discounted cart) minus coupon = taxable; GST 18% on taxable; total = taxable + GST. Orders persist coupon_code/coupon_discount_inr; tax_inr holds GST. Coupons: flat (INR off) or percentage, capped at subtotal, one redemption per user (coupon_users). Order placement runs in one transaction and reserves stock with a conditional UPDATE (no oversell); status starts pending_payment. validate returns 200 with is_valid=false for unusable codes.

## Payments (TASK-039, ADR-0006)
Order lifecycle: pending_payment -> (webhook payment.captured) paid, or payment_failed. Capture deducts on_hand and reserved; failure keeps the reservation. Webhook is public and authenticated only by X-Razorpay-Signature. initiate needs an Idempotency-Key header (8-128 chars). Env: RAZORPAY_KEY_ID, RAZORPAY_KEY_SECRET, RAZORPAY_WEBHOOK_SECRET.
