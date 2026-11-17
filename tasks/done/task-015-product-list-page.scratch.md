# Scratch: task-015 product list page

## Plan
1. filterSlice (categories, valueIds, search, sort, page) in the root store, not persisted; URL (de)serialisation helpers.
2. useProductFilterUrlSync: URL -> store on mount/back-forward; store -> URL via router.push; returns ready flag so nothing fetches before hydration (avoids an unfiltered first request and a wipe-the-URL race).
3. FilterSidebar (attributes via GET /attributes, categories via GET /categories, mobile toggle), ProductListGrid (generated listProducts, search, sort, pagination 12/page, error/empty states), page wrapped in Suspense.
4. e2e with mocked catalogue API.

## Uncertainties (resolved)
- Static export: useSearchParams needs Suspense; the page is a client component fetching on the client (filters are user-driven, not build-time).
- Attribute filter param name in the URL is `attr` (short) and maps to the API's `attr_value_id[]`.

## Skipped
- Product detail page (TASK-016). Product cards link to /products/<slug>/ which does not exist yet.
