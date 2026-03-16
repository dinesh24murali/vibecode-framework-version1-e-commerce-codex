# Prompt: Generate Functional Specification

> **How to use:** Complete Phase 1 Prompt 01 (PRD) first. Fill in every `[[PLACEHOLDER]]`, then paste this entire file into your AI tool. Save the output to `docs/02_outputs/02_functional_spec.md`.

---

## Instructions to AI

You are a senior software engineer writing a **Functional Specification** document. This document bridges the PRD and technical implementation — it defines exactly how features behave, screen by screen and flow by flow.

Use the PRD at `docs/02_outputs/01_prd.md` as input. If it is not yet available, use the project details below.

---

## Project Input

**Project name:** [[PROJECT_NAME]]
**Backend:** [[BACKEND_LANG]] / [[BACKEND_FRAMEWORK]]
**Frontend:** [[FRONTEND_FRAMEWORK]]
**Auth strategy:** [[AUTH_STRATEGY]]
**Database:** [[DATABASE]]

**v1 Features:**
[[Fill from intake 3.1]]

---

## Required Sections

Generate the functional spec with all of the following:

1. **Overview** — scope of this document, link to PRD
2. **User Flows** — for each v1 feature, a numbered step-by-step flow from the user's perspective
3. **Screen / Page Inventory** — list every screen/page/view with:
   - Route or path
   - Purpose
   - Key UI elements
   - Auth requirement (public / authenticated / role-gated)
4. **Feature Specifications** — for each feature:
   - Trigger / entry point
   - Happy path
   - Error states and messages
   - Edge cases
   - Validation rules
5. **State Diagrams** — describe state transitions for any stateful entity (e.g. order status, user account status) in Mermaid diagram format
6. **Notifications & Emails** — list any system-generated messages, their triggers, and content outline
7. **Permissions Matrix** — table of [role × action × resource] showing allow/deny
8. **Data Retention & Deletion** — what data is kept, for how long, and deletion behavior

---

## Output Format

- Use Markdown
- Use Mermaid for state diagrams (` ```mermaid `)
- Number all user flows (UF-001, UF-002)
- Number all feature specs (FS-001, FS-002)
- Target length: 2000–4000 words
