# AGENTS.md — Master AI Instructions

> This file is the single source of truth for all AI tools working in this repository.
> `CLAUDE.md` and `.cursorrules` redirect here.

---

## 1. Project Overview

- **Project name:** `[[PROJECT_NAME]]`
- **Description:** `[[PROJECT_DESCRIPTION]]`
- **Backend:** `[[BACKEND_LANG]]` / `[[BACKEND_FRAMEWORK]]`
- **Frontend:** `[[FRONTEND_FRAMEWORK]]`
- **Database:** `[[DATABASE]]`
- **Deployment target:** `[[DEPLOYMENT_TARGET]]`

---

## 2. How to Use This Framework

This is a **two-phase vibe-coding framework**:

### Phase 1 — Document Generation (before writing code)
1. Fill in `docs/00_intake/intake_questionnaire.md`
2. Paste each prompt from `docs/01_prompts/` into your AI tool (Claude, Cursor, Codex, Gemini — any will work)
3. Save AI outputs to `docs/02_outputs/` (e.g. `docs/02_outputs/prd.md`)
4. Review outputs, then create your first ADRs in `adr/` and seed `memory/`

### Phase 2 — Living Project (during development)
For every feature or bugfix:
1. Create `tasks/active/<task-name>.md` from `tasks/_template.md`
2. Ask AI to write `tasks/active/<task-name>.scratch.md` **before touching code**
3. AI implements the task
4. AI updates `CHANGELOG.md` and relevant `memory/<domain>.md` files
5. Run `verify/scripts/run-e2e.sh` (frontend changes) or relevant tests
6. Move task file to `tasks/done/`

---

## 3. Universal Rules — All AI Tools Must Follow

These rules apply regardless of which AI tool is in use.

### 3.1 Before Writing Any Code
- **Always** create `tasks/active/<task-name>.scratch.md` first
- The scratch file must contain:
  - The plan (step-by-step)
  - Uncertainties or open questions
  - What will be skipped and why
  - Estimated files to be changed
- Do **not** proceed to implementation until the scratch file is written

### 3.2 After Completing Any Task
- Update `CHANGELOG.md` with a brief entry under the current date
- Move the task file from `tasks/active/` to `tasks/done/`
- Update any relevant `memory/<domain>.md` files if domain knowledge changed

### 3.3 Architectural Decisions
- When making a significant architectural decision, create an ADR in `adr/`
- Use the naming convention: `adr/NNNN-short-slug.md`
- Use `adr/template.md` as the base
- Reference the ADR in related task files

### 3.4 Memory Files
- Keep `memory/<domain>.md` files updated when domain knowledge changes
- Always check relevant memory files before starting a task
- See `memory/AGENTS.md` for the full memory maintenance protocol

### 3.5 Frontend Changes
- After any frontend change, run:
  ```
  make verify
  ```
  or manually:
  ```
  bash verify/scripts/run-e2e.sh
  npx ts-node verify/scripts/check-dom.ts
  ```
- Consult `verify/dom-checks.md` for the DOM verification checklist
- See `verify/mcp-playbook.md` for MCP-based verification steps

### 3.6 Test Coverage
- Coverage must **not drop** below the thresholds in `tests/coverage-baseline.json`
- Update `tests/coverage-baseline.json` when new packages are added
- Contract tests in `tests/contract/` must be generated from the OpenAPI spec in `docs/02_outputs/`
- See `tests/AGENTS.md` for the full testing protocol

### 3.7 API Contracts
- The frontend must be generated from (or validated against) the OpenAPI spec
- Never manually write API client code that duplicates what the spec defines
- If the spec changes, regenerate the frontend client before implementing

---

## 4. File Map

| Path | Purpose |
|------|---------|
| `docs/00_intake/` | Intake questionnaire — fill this first |
| `docs/01_prompts/` | Prompt files — paste into any AI tool |
| `docs/02_outputs/` | AI-generated docs (PRD, arch, API spec, etc.) |
| `adr/` | Architecture Decision Records |
| `memory/` | Persistent AI memory files (domain knowledge) |
| `tasks/active/` | In-progress task files + scratch files |
| `tasks/done/` | Completed task files (archive) |
| `tests/e2e/` | End-to-end tests (Playwright) |
| `tests/contract/` | Contract tests (OpenAPI-derived) |
| `verify/` | Verification scripts and playbooks |
| `backend/` | Backend source code |
| `frontend/` | Frontend source code |

Each folder has its own `AGENTS.md` with folder-specific instructions.

---

## 5. Folder-Level AGENTS.md Files

- `docs/AGENTS.md` — docs folder purpose and ordering
- `adr/AGENTS.md` — ADR creation rules and naming
- `memory/AGENTS.md` — memory maintenance protocol
- `tasks/AGENTS.md` — task workflow and scratch file convention
- `tests/AGENTS.md` — coverage gates and test conventions
- `verify/AGENTS.md` — verification protocol
- `backend/AGENTS.md` — backend conventions (`[[BACKEND_LANG]]`)
- `frontend/AGENTS.md` — frontend conventions (`[[FRONTEND_FRAMEWORK]]`)

---

## 6. Placeholders Reference

When you clone this framework, find-and-replace these tokens:

| Token | Replace with |
|-------|-------------|
| `[[PROJECT_NAME]]` | Your project name |
| `[[PROJECT_DESCRIPTION]]` | One-sentence description |
| `[[BACKEND_LANG]]` | e.g. `Go`, `Python`, `Node.js` |
| `[[BACKEND_FRAMEWORK]]` | e.g. `Gin`, `FastAPI`, `Express` |
| `[[FRONTEND_FRAMEWORK]]` | e.g. `React`, `Vue`, `SvelteKit` |
| `[[DATABASE]]` | e.g. `PostgreSQL`, `SQLite`, `MongoDB` |
| `[[DEPLOYMENT_TARGET]]` | e.g. `AWS`, `GCP`, `Fly.io`, `Vercel` |
| `[[TEAM_NAME]]` | Your team or company name |
| `[[AUTH_STRATEGY]]` | e.g. `JWT`, `OAuth2`, `Session` |
| `[[UI_LIBRARY]]` | e.g. `shadcn/ui`, `Material UI`, `Tailwind`, `none` |
| `[[CODEGEN_COMMAND]]` | Command to regenerate frontend API client from OpenAPI spec |
| `[[CI_CD_TOOL]]` | e.g. `GitHub Actions`, `GitLab CI`, `CircleCI` |
| `[[CI_TEST_COMMAND]]` | e.g. `npm test`, `pytest`, `go test ./...` |
| `[[CI_LINT_COMMAND]]` | e.g. `npm run lint`, `ruff check .`, `golangci-lint run` |
| `[[REPO_URL]]` | Full URL of the git repository |
| `[[APP_PORT]]` | Local dev server port, e.g. `3000`, `8080` |
| `[[STAGING_URL]]` | Staging environment URL |
| `[[PROD_URL]]` | Production URL |
| `[[PROJECT_START_DATE]]` | Date project was started (YYYY-MM-DD) |
