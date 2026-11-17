# Scratch: TASK-016

## Plan
1. Spec only exposes GET /products/{id} (UUID); pages are routed by slug. Build time: page through listProducts to build slug->id, then getProduct(id) for initial data.
2. `[slug]/page.tsx`: generateStaticParams (all slugs; fallback `sample-product` if API unreachable/empty so the export still builds), generateMetadata (title/description from product).
3. `ProductDetail` client component: uses `initial` product, else resolves slug from usePathname via listProducts({q}) + getProduct. Sets document.title when resolved client-side.
4. ProductImageGallery (thumbnail strip), ProductZoom (hover zoom via transform-origin; two-finger pinch via touch handlers; no new dependency), AddToCartButton (cartSlice.addItem + role=status toast).
5. e2e spec with mocked API; run docker build + e2e.

## Files
frontend/app/(shop)/products/[slug]/page.tsx, 4 components, tests/e2e/product-detail.spec.ts, CHANGELOG, memory/product.md, PROGRESS.

## Uncertainties / decisions
- react-image-magnify skipped (no React 19 support); custom zoom, no dependency.
- Slug->id lookup uses the list API since the spec has no by-slug endpoint (spec unchanged).

## Skipped
Server cart sync beyond the existing slice; reviews; related products.

## Risks
Products created after build 404 until rebuild (SSG by design).
