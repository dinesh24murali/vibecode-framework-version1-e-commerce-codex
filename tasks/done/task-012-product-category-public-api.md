# Task: TASK-012 — Product & category public API

**Status:** done
**Created:** 2026-11-17
**ADR refs:** none

## Goal
Public catalogue endpoints: product list/detail with filters and pagination, category list/detail.

## Acceptance Criteria
- [x] Filter by category, attribute value ids and search return correct results (verified by real-DB integration test with its own seed rows: OR within a dimension, AND across dimensions, inactive excluded)
- [x] Pagination metadata present in list responses (spec fields: total, page, limit, total_pages)
- [~] GET /categories returns all 12 seeded categories (returns all categories in the table; the 12 are inserted by TASK-030, not yet done)
- [x] Inactive products (no deleted_at column exists; is_active=false) do not appear in any public endpoint
- [x] Contract tests cover all 4 product/category endpoints (spec defines 4, not 5): backend/tests/contract/products_test.go

## Dependencies
Tasks: 011. Docs: 04_api_spec.yaml.
