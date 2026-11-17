# Scratch: task-014 landing page

## Plan
1. Migration 23: carousel_slides table (no table existed). Query: active slides ordered, max 5.
2. LandingService: slides + Featured (8 newest active products; schema has no featured flag, functional spec says "featured or recent"). Refactor ProductService summaries into a shared helper.
3. Frontend: (shop)/page.tsx server component fetches landing at build; LandingContent client component takes initial data and fetches on mount when the build-time fetch failed (the build container has no API); HeroCarousel, ProductCard, ProductGrid; fallbacks per functional spec.
4. Removed app/page.tsx placeholder (route conflict with (shop)/page.tsx).

## Uncertainties (resolved)
- Static build with API down bakes in null -> client fetch fallback makes the page work either way. When the API is up at build time the data is baked in, and goes stale until the next build.
- Slides can only be inserted via SQL/seed (TASK-030); no admin CRUD in this task.
- Add to Cart only updates the local store (cart server sync is the cart task).
- Product link target /products/<slug>/ does not exist yet (product pages are the next tasks).

## Skipped
- Navbar/footer (functional spec) - not in this task's file list.
