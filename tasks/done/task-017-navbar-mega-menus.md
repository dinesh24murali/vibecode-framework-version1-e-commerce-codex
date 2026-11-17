# Task: TASK-017 — Navbar with mega-menus

**Status:** done
**Created:** 2026-11-17
**ADR refs:** ADR-0003

## Goal
Responsive navbar on all shop pages with Book Store / Printing Services mega-menus, search and cart badge.

## Acceptance Criteria
- [x] Hovering "Book Store" reveals all 12 categories; clicking one navigates to `/products/?category=<slug>`
- [x] Hovering "Printing Services" reveals both columns with all 24 sub-items
- [x] Search bar submits on Enter and navigates to `/products/?q=<query>`
- [x] Cart icon badge reflects cartSlice item count
- [x] Mobile: mega-menus collapse into a hamburger drawer with accordion sections
- [x] E2e test covers category click and search submission
- [x] Static build, e2e and DOM checks pass (docker equivalents of make verify)

## Dependencies
Tasks: 003. Memory: memory/product.md.
