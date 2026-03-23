# Memory: Product

## Scope

- `e-commerce-site` is a single-store direct-to-consumer bookstore, not a marketplace.
- v1 includes storefront browsing, customer accounts, cart, checkout without payment gateway, and an admin portal.
- The storefront taxonomy currently includes 12 categories: Colouring Books, Drawing Books, Journals, Calendars, Sketch Pads, Planners, Stained Art, Fabrii Art, Note Book, Sticker Book, Sticker Colouring Book, and Activity Cards.

## Core Product Rules

- Product listing filters must support category, level, theme, and size.
- Product search and filter state must persist in URL query parameters.
- Customer account area contains My profile, My address, and My orders.
- Customers can maintain multiple addresses and choose one default address.
- Order statuses are limited to `Open`, `Fulfilled`, and `Cancelled`.
- Customers may cancel only `Open` orders and may provide an optional cancellation reason.
- Checkout supports one coupon code per order and no online payment in v1.
- Tax is calculated at individual product level using the tax rate stored on each product.
- Contact inquiries must both persist in the database and send a notification email.
- Facebook login is controlled by feature flag and may be hidden until platform approval is complete.
- All admin product fields are mandatory in v1.

## Seed Data

- Default coupons: `SUPERHIT` for flat 150 discount and `FIRSTTIME` for 10% discount.
- Initial seed data must include one default admin user and 5 dummy products across 3 categories.

## Operational Expectations

- Order placement triggers customer confirmation email and manual phone follow-up by the business.
- `Printing supplies` exists as a blank public placeholder page in v1.

## Open Product Decisions
