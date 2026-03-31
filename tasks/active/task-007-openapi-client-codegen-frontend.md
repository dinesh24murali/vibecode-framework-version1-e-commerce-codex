# Task: TASK-007 — OpenAPI client codegen — frontend

**Status:** active
**Created:** 2026-03-31
**ADR refs:** none

---

## Goal

Generate a fully typed TypeScript API client from `docs/02_outputs/04_api_spec.yaml` using `orval` or `openapi-typescript-codegen`, wire the base URL from env, inject the auth header from Zustand, and add a `make codegen` target.

## Background

The frontend must never contain hand-written `fetch` calls that duplicate what the OpenAPI spec defines. All API interaction goes through the generated client in `frontend/lib/api/`. This is enforced by `make verify`. The generated client must be regenerated every time the spec changes, before any frontend implementation work begins.

See `frontend/AGENTS.md` (API Client convention), `docs/02_outputs/05_implementation_plan.md` TASK-007, and `docs/02_outputs/04_api_spec.yaml`.

## Acceptance Criteria

- [ ] `make codegen` regenerates `frontend/lib/api/` from `docs/02_outputs/04_api_spec.yaml` without errors
- [ ] Generated types match all schemas defined in the API spec
- [ ] Auth header (`Authorization: Bearer <token>`) is injected automatically for protected endpoints using the Zustand auth store
- [ ] No hand-written `fetch` call exists anywhere in `frontend/` (enforced by `make verify`)
- [ ] `frontend/AGENTS.md` `[[CODEGEN_COMMAND]]` placeholder is replaced with the actual codegen command
- [ ] `make verify` passes

## Dependencies

- **Tasks:** TASK-003 (frontend scaffold must exist), TASK-008 (Zustand auth store must exist for token injection)
- **Memory files:** `memory/auth.md` (understand token storage to wire the header injector correctly)
- **Docs:** `docs/02_outputs/04_api_spec.yaml` (the spec being consumed), `frontend/AGENTS.md`

## Files Expected to Change

- `frontend/lib/api/` *(generated — do not hand-edit after generation)*
- `frontend/lib/api/client.ts` *(create — base client with auth header injection; this file is hand-maintained, not generated)*
- `Makefile` *(update `codegen` target with actual command)*
- `frontend/AGENTS.md` *(update `[[CODEGEN_COMMAND]]` placeholder)*
- `frontend/package.json` *(update — add codegen tool as devDependency)*

## Notes

- Tooling choice: prefer `orval` (generates React Query hooks + types) or `openapi-typescript-codegen`; document the choice with a brief rationale in the scratch file
- Base URL: read from `NEXT_PUBLIC_API_URL` environment variable
- Auth header injection: the base client reads the access token from the Zustand auth store; in SSG context (no token) the header is omitted
- The `frontend/lib/api/` directory is auto-generated — add it to `.gitignore` OR commit it with a clear comment; decide in the scratch file
- `make verify` must confirm no raw `fetch(` or `axios(` calls exist outside of `frontend/lib/api/client.ts`

---

## Scratch File

AI: before implementing, create `tasks/active/task-007-openapi-client-codegen-frontend.scratch.md` with your plan.
See `tasks/AGENTS.md` for the required scratch file format.
