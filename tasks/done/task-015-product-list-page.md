# Task: TASK-015 — Product list page with filters

**Status:** done
**Created:** 2026-11-17
**ADR refs:** ADR-0003

## Goal
Product list page with sidebar filters synced to the URL.

## Acceptance Criteria
- [x] Active filters are reflected in the URL (`/products/?category=journals&attr=<value uuid>&q=..&sort=..&page=..`; plan's `attributeValueIds=3,7` form replaced by repeated params with UUIDs, since the API takes attribute value UUIDs)
- [x] Navigating directly to a filtered URL pre-checks sidebar filters (first request already carries the filters)
- [x] Filter changes update the grid without a full page reload (e2e asserts the same document)
- [x] Mobile: sidebar collapses behind a Filters toggle
- [x] E2e product-list.spec.ts verifies filtering, URL sync, back button, direct URL, mobile toggle, error state
- [x] Static build and DOM checks pass (make not installed on host; equivalent docker commands run)

## Dependencies
Tasks: 012, 011, 007.
