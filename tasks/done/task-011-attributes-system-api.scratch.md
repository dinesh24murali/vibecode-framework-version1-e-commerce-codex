# Scratch: task-011 attributes API

## Plan
1. Query GetAttributeByID; regenerate sqlc (v1.30.0 header restored).
2. models/attribute.go DTOs, AttributeRepository, AttributeService, AttributeHandler; public routes.
3. Contract test with fakes + repository integration test on compose Postgres.

## Uncertainties (resolved)
- Plan says `slug`; spec says `code` and `value_slug` -> spec wins.
- GET /attributes returns only is_filterable=true, each with values (spec AttributeWithValues).
- Abstract Factory for filter predicates belongs with product listing filters (TASK-012), skipped here.

## Skipped
- Frontend work (generated client already contains listFilterableAttributes/listAttributeValues).

## Risks
- N+1 value queries (one per attribute); attribute count is small. Revisit if it grows.
