# Memory: Product

## Scope

- `e-commerce-site` is a single-store direct-to-consumer bookstore, not a marketplace.
- v1 includes storefront browsing, customer accounts, cart, checkout without payment gateway, and an admin portal.
- The storefront taxonomy currently includes 12 categories: Colouring Books, Drawing Books, Journals, Calendars, Sketch Pads, Planners, Stained Art, Fabrii Art, Note Book, Sticker Book, Sticker Colouring Book, and Activity Cards.

## Core Product Rules

- Product listing filters must support category, level, theme, and size.
- Product search and filter state must persist in URL query parameters.
- Customer account area contains My profile, My address, and My orders.
- Order statuses are limited to `Open`, `Fulfilled`, and `Cancelled`.
- Customers may cancel only `Open` orders.
- Checkout supports one coupon code per order and no online payment in v1.

## Seed Data

- Default coupons: `SUPERHIT` for flat 150 discount and `FIRSTTIME` for 10% discount.
- Initial seed data must include one default admin user and 5 dummy products across 3 categories.

## Open Product Decisions

- Tax regime details are not finalized.
- Footer link target for `Library` is not yet explicitly defined.
- `Printing supplies` behavior and destination are still unspecified.
