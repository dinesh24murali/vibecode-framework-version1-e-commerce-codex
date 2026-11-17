# Scratch: TASK-017

## Plan
1. lib/constants/navigation.ts: 12 categories (name+slug), 11 products offered, 13 inhouse services (names from intake questionnaire), nav links.
2. components/layout: Navbar (client; desktop hover/focus menus, mobile drawer with accordions, cart badge), BookStoreMegaMenu, PrintingServicesMegaMenu, SearchBar.
3. Mount Navbar in app/(shop)/layout.tsx.
4. e2e tests/e2e/navbar.spec.ts (mocked API); build + verify.

## Files
frontend/lib/constants/navigation.ts, frontend/components/layout/*.tsx, frontend/app/(shop)/layout.tsx, tests/e2e/navbar.spec.ts, CHANGELOG, memory/product.md, PROGRESS.

## Decisions
- Existing /products page (TASK-015) uses `category=<slug>` and `q`, so the plan's `category=<name>` and `search=` become slug and `q`.
- Categories are a static constant (fixed taxonomy) so the menu renders without the API; slug = kebab-case of the name.
- Menus open on hover, focus and click (keyboard users). Destination pages (wall-of-fame, contact, printing) arrive in 018-020.

## Skipped
Destination pages; admin-editable menus.

## Risks
Category slugs must match seeded slugs (TASK-030 seed should use kebab-case of the names).
