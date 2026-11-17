# Scratch: task-007 OpenAPI client codegen

## Plan
1. Add orval as devDependency (node container). Rationale: spec-driven, supports a custom mutator so one hand-written file (client.ts) owns fetch + auth header. Chose plain fetch client (no react-query) to keep deps minimal; hooks can be added later.
2. frontend/orval.config.ts: input ../docs/02_outputs/04_api_spec.yaml, output frontend/lib/api/generated (mode tags-split, client fetch), override mutator -> lib/api/client.ts `apiFetch`.
3. client.ts: base URL from NEXT_PUBLIC_API_URL, injects Bearer from useStore.getState().accessToken (omitted if none), credentials include for refresh cookie.
4. Makefile codegen runs in frontend dir (already NODE-based); update AGENTS.md placeholder.
5. Wire CartApi adapter from TASK-008 to generated cart functions in lib/api/cartApi.ts (setCartApi).
6. Commit generated code (decision: committed, header comment from orval) so builds do not require codegen.
7. Add no-raw-fetch check to verify (script grep) .

## Files
frontend/package.json(+lock), frontend/orval.config.ts, frontend/lib/api/client.ts, frontend/lib/api/generated/*, frontend/lib/api/cartApi.ts, Makefile, frontend/AGENTS.md, verify/scripts/check-no-fetch.sh

## Uncertainties
- orval version/fetch client output shape; resolve by inspection.

## Skipped
- React Query hooks.

## Risks
- Spec quirks failing orval validation.
