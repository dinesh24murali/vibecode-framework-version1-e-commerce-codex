# Scratch: add-public-filter-options-api

## Plan

1. Add `attributes` tag to the `tags` list in `04_api_spec.yaml`
2. Add `GET /api/v1/attributes` path entry (public, tag: attributes)
   - No auth required
   - Response: `{ data: AttributeWithValues[] }`
   - Only returns attributes where `is_filterable: true`
3. Add `GET /api/v1/attributes/{attrDefId}/values` path entry (public, tag: attributes)
   - Path param: `attrDefId` (reuse `$ref: "#/components/parameters/PathAttrDefId"`)
   - Response: `{ data: AttributeDefinitionValue[] }`
4. Add `AttributeWithValues` schema to `components/schemas`
   - Fields: id, code, name, is_filterable, sort_order, created_at, updated_at (from AttributeDefinition)
   - Plus: values (array of AttributeDefinitionValue)

## Files to Change

- `docs/02_outputs/04_api_spec.yaml` — only file changed

## Uncertainties

- None; the existing `AttributeDefinition`, `AttributeDefinitionValue`, and `PathAttrDefId` schemas/params all exist and can be referenced directly.

## What Will Be Skipped

- No backend code changes (spec only)
- No new query parameters on `GET /api/v1/attributes` (could add `?filterable_only=true` but it's always filterable for public — keep it simple)

## Risks

- `PathAttrDefId` parameter name — need to verify it matches the name used in admin endpoints before reusing it
