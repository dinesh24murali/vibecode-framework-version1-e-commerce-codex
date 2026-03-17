# Prompt: Generate Phased Implementation Plan

> **How to use:** Complete Prompts 01–04 first. Fill in every `[[PLACEHOLDER]]`, then paste this entire file into your AI tool. Save the output to `docs/02_outputs/05_implementation_plan.md`.

---

## Instructions to AI

You are a senior engineering lead. Generate a **phased implementation plan** that takes the project from zero to a production-ready v1. The plan must be granular enough that an AI coding agent can execute each task independently.

---

## Project Input

**Project name:** e-commerce-site
**Backend:** Go 1.26 / Gin Web Framework
**Frontend:** NextJS 16
**Database:** Postgres
**Auth strategy:** JWT
**Deployment target:** AWS
**CI/CD:** [[Fill from intake 6.3]]

**v1 Features:**
[[Fill from intake 3.1]]

---

## Required Output

Generate a phased plan with the following structure:

### Phase 0 — Project Setup
Tasks to scaffold the project, configure tooling, set up CI/CD, and establish the development environment.

### Phase 1 — Core Infrastructure
Tasks to build the foundational layers: database schema and migrations, auth system, base API structure, frontend scaffolding and routing.

### Phase 2 — Feature Implementation
One sub-phase per v1 feature. Each sub-phase contains:
- Backend tasks (models, handlers, tests)
- Frontend tasks (components, pages, API client integration)
- Integration tasks

### Phase 3 — Quality & Hardening
Tasks for: error handling, input validation, rate limiting, logging, test coverage gaps, security review.

### Phase 4 — Deployment & Launch
Tasks for: production infrastructure, CI/CD pipeline, monitoring setup, documentation, launch checklist.

---

## Task Format

For each task, use this format:

```
#### TASK-NNN: <Task name>
- **Phase:** <phase name>
- **Depends on:** <TASK-NNN list or "none">
- **Files:** <list of files to create or modify>
- **Description:** <2-3 sentence description of what to do>
- **Acceptance criteria:**
  - [ ] ...
```

---

## Additional Requirements

- Every task must have clear acceptance criteria
- Flag tasks that require an ADR with `> **ADR needed:** <topic>`
- Flag tasks that must be done sequentially vs. in parallel
- Total task count should be between 20 and 50 for a typical v1
- Include a dependency graph summary at the end (which phases can overlap)
