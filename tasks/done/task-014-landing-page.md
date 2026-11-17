# Task: TASK-014 — Landing page — backend + frontend

**Status:** done
**Created:** 2026-11-17
**ADR refs:** none

## Goal
GET /api/v1/landing and a statically generated landing page with hero carousel and featured product cards.

## Acceptance Criteria
- [x] GET /api/v1/landing returns LandingPageResponse per spec (carousel_slides, featured_products); contract test backend/tests/contract/landing_test.go
- [x] Carousel (embla, autoplay 4s, touch/swipe via embla) renders up to 5 slides (query LIMIT 5); e2e verifies 3 slides and auto-advance. Swipe is embla's built-in pointer handling, not separately tested.
- [x] Product cards show name, price, discount badge, Add to Cart button (adds to the client cart store)
- [x] out/index.html exists after next build
- [x] E2e landing.spec.ts: carousel visible, cards, auto-advance, API-failure fallbacks

## Dependencies
Tasks: 012, 007.
