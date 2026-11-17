# Task: TASK-016 — Product detail page

**Status:** done
**Created:** 2026-11-17
**ADR refs:** ADR-0003

## Goal
Product detail page pre-rendered per slug, with image gallery, zoom and Add to Cart.

## Acceptance Criteria
- [x] All product slugs pre-rendered at build time (`generateStaticParams` fetches from API)
- [x] Image zoom on desktop hover; pinch-to-zoom on mobile
- [x] "Add to Cart" updates `cartSlice` and shows a toast
- [x] `<title>` and `<meta name="description">` from product data
- [x] E2e test covers load, gallery interaction, add-to-cart
- [x] Static build and DOM checks pass (docker equivalents of `make verify`)

## Dependencies
Tasks: 012, 007. Memory: memory/product.md.
Files: frontend/app/(shop)/products/[slug]/page.tsx, components/shop/{ProductDetail,ProductImageGallery,ProductZoom,AddToCartButton}.tsx, tests/e2e/product-detail.spec.ts
