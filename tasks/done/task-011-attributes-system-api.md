# Task: TASK-011 — Attributes system API

**Status:** done
**Created:** 2026-11-17
**ADR refs:** none

## Goal
Public endpoints listing filterable attributes and their values for the product filter sidebar.

## Acceptance Criteria
- [x] GET /api/v1/attributes returns filterable attribute definitions (id, code, name, is_filterable, sort_order, values) per the spec (spec uses `code`, not `slug`)
- [x] GET /api/v1/attributes/{attrDefId}/values returns values for the definition (404 for unknown id)
- [x] Both endpoints are public (registered outside the auth middleware)
- [x] Contract test backend/tests/contract/attributes_test.go validates response shapes
- [x] Real-DB query path verified by repository integration test (TASK-030 seed data not yet present; test inserts its own rows)

## Dependencies
Tasks: 005. Docs: 04_api_spec.yaml (attributes).
