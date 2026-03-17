# Task: 000-project-setup

**Status:** done
**Created:** 16-3-2026
**ADR refs:** ADR-0001

---

## Goal

Scaffold the project from the vibecode framework template: replace all placeholders, set up the repository, configure the development environment, and verify the toolchain works end-to-end.

## Background

This is always the first task for any project using the vibecode framework. It must be completed before any feature work begins. See ADR-0001 for the decision to use this framework.

## Acceptance Criteria

- [x] All `[[PLACEHOLDER]]` tokens replaced (run `make init` to verify none remain)
- [x] `.env` created from `.env.example` and all required vars populated
- [x] `make check-env` passes
- [x] Repository pushed to remote with branch protection on `main`
- [x] CI pipeline (`.github/workflows/ci.yml`) runs and passes on first push
- [x] `docs/00_intake/intake_questionnaire.md` fully filled in
- [x] At least one team member has cloned the repo and run the dev server successfully
- [x] `adr/0001-use-vibecode-framework.md` date and team fields filled in
- [x] `memory/project.md` seeded with initial stack and key decisions

## Dependencies

- **Tasks:** none (this is the first task)
- **Memory files:** none yet
- **Docs:** `docs/00_intake/intake_questionnaire.md`

## Files Changed

| File | Action |
|------|--------|
| `AGENTS.md` | Replaced all `[[PLACEHOLDER]]` tokens |
| `README.md` | Replaced `e-commerce-site` |
| `Makefile` | Replaced `e-commerce-site` |
| `.env.example` | Replaced `JWT` |
| `backend/AGENTS.md` | Replaced stack placeholders |
| `frontend/AGENTS.md` | Replaced stack + codegen command |
| `adr/0001-use-vibecode-framework.md` | Filled in date and team |
| `memory/project.md` | Created with initial project context |
| `.env` | Created from `.env.example` (not committed) |

## Notes

- Placeholder replacement: use `make init` to get a full list of remaining tokens before closing this task
- The `.env` file is gitignored — store secrets in your team's secret manager and document the location in `memory/infra.md`
