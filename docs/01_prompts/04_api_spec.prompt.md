# Prompt: Generate OpenAPI 3.1 Specification

> **How to use:** Complete Prompts 01–03 first. Fill in every `[[PLACEHOLDER]]`, then paste this entire file into your AI tool. Save the output to `docs/02_outputs/04_api_spec.yaml`. This file becomes the contract for `tests/contract/`.

---

## Instructions to AI

You are a senior API designer. Generate a complete **OpenAPI 3.1 specification** in YAML format for the project below.

This spec is the **single source of truth** for the API. The frontend client will be generated from it. Contract tests in `tests/contract/` will validate against it. Do not leave placeholder paths — generate real, complete endpoint definitions.

---

## Project Input

**Project name:** e-commerce-site
**Backend:** Go 1.26 / Gin Web Framework
**Auth strategy:** JWT
**Database:** Postgres

**v1 Features:**
[[Fill from intake 3.1]]

**User roles:**
[[Fill from intake 5.3]]

---

## Requirements

Generate a complete OpenAPI 3.1 YAML spec with:

1. **Info block** — title, version `1.0.0`, description
2. **Servers** — `http://localhost:8080` for dev; add staging/prod stubs
3. **Security schemes** — matching `JWT` (e.g. BearerAuth for JWT)
4. **Tags** — one per feature area
5. **Paths** — for every v1 feature, define:
   - All CRUD endpoints that make sense
   - Request body schema (inline `$ref` components)
   - Response schemas for 200, 201, 400, 401, 403, 404, 422, 500
   - `security` field on protected endpoints
6. **Components / Schemas** — reusable schemas for all entities
7. **Components / Parameters** — reusable path/query params (e.g. `id`, `page`, `limit`)
8. **Components / Responses** — reusable error response shapes
9. **Pagination** — use cursor-based or offset pagination consistently; define once, reference everywhere

---

## Output Format

- Output **only valid YAML** — no markdown code fences, no explanation text
- Start with `openapi: "3.1.0"`
- Use `$ref` for all reusable schemas — do not inline complex objects
- All schema properties must have `type` and `description`
- Include at least one example per schema
