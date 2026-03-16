# tests/ — AGENTS.md

## Purpose

This folder holds all automated tests. Tests are the safety net that lets AI agents make changes confidently.

## Directory Layout

```
tests/
├── coverage-baseline.json   ← Coverage thresholds; never let coverage fall below these
├── e2e/                     ← End-to-end tests (Playwright)
└── contract/                ← Contract tests (generated from OpenAPI spec)
```

## Coverage Gate

**Coverage must not drop below the thresholds in `tests/coverage-baseline.json`.**

- Before merging any task, run the test suite and verify coverage
- If you add a new package, add its threshold to `coverage-baseline.json`
- If coverage genuinely cannot be maintained (e.g. generated code), document the exception in the JSON file

## Contract Testing

Contract tests ensure the frontend and backend agree on the API shape.

- The contract test source of truth is `docs/02_outputs/04_api_spec.yaml`
- When the API spec changes, regenerate the contract tests before implementing
- Tools to consider: `openapi-typescript` (frontend types), `oapi-codegen` (Go server stubs), `schemathesis` (automated API fuzzing)
- Place generated client code in `frontend/src/api/` (or equivalent)
- Place contract test files in `tests/contract/`

## E2E Tests

- Use **Playwright** for e2e tests (see `verify/mcp-playbook.md` for MCP-based verification)
- E2e tests live in `tests/e2e/`
- Run e2e tests with: `make verify` or `bash verify/scripts/run-e2e.sh`
- Name test files: `<feature>.spec.ts`
- E2e tests should cover all happy paths and critical error states

## Rules for AI Tools

- Do not delete tests to make coverage pass — fix the code
- Do not mock away the database in integration tests — use a real test database
- When adding a new API endpoint, add a contract test for it
- When adding a frontend page, add an e2e test for its happy path
- Run `make verify` after any frontend change before considering a task done
