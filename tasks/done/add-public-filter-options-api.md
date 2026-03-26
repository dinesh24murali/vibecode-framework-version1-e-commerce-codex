# Task: add-public-filter-options-api

**Status:** active
**Created:** 2026-03-26
**ADR refs:** none

---

## Goal

Add public API endpoints to retrieve product filter options (level, size, theme) so the frontend can render filter UI without having to call the product list endpoint first.

## Background

`GET /api/v1/products` returns a `ProductListResponse` which embeds a `filters` field containing `AttributeFilterGroup` entries for the *current result set*. However, there are no standalone public endpoints to fetch all available filter options (all attribute definitions and their complete value lists). Without these, a frontend cannot pre-populate filter dropdowns on page load before any products are queried.

The admin-only endpoints (`/api/v1/admin/attributes/...`) cover this for admin UI, but are auth-protected and not appropriate for the public storefront.

## Acceptance Criteria

- [ ] `GET /api/v1/attributes` added to spec — returns all filterable attribute definitions with their values
- [ ] `GET /api/v1/attributes/{attrDefId}/values` added to spec — returns values for a single attribute definition
- [ ] New `attributes` tag added to the tags list
- [ ] New `AttributeWithValues` schema added to components/schemas
- [ ] New `PathAttrDefId` parameter reused (already exists in components/parameters)
- [ ] Both endpoints are public (no auth required)

## Dependencies

- **Tasks:** none
- **Memory files:** none
- **Docs:** `docs/02_outputs/04_api_spec.yaml`

## Files Expected to Change

- `docs/02_outputs/04_api_spec.yaml`

## Notes

- Reuse existing `AttributeDefinitionValue` schema for the values array items
- The `AttributeFilterGroup` schema is semantically tied to "values present in the current result set"; prefer a new `AttributeWithValues` schema for the standalone endpoint
- Only return attributes where `is_filterable: true` in the public endpoint

---

## Scratch File

See `tasks/active/add-public-filter-options-api.scratch.md`
