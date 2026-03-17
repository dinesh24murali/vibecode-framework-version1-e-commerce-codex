# Prompt: Generate Product Requirements Document (PRD)

> **How to use:** Fill in every `[[PLACEHOLDER]]` from your intake questionnaire, then paste this entire file into your AI tool. Save the output to `docs/02_outputs/01_prd.md`.

---

## Instructions to AI

You are a senior product manager. Generate a comprehensive, production-ready **Product Requirements Document (PRD)** based on the project details below.

The PRD must be structured for an engineering team to begin implementation immediately. Be specific and opinionated — do not leave ambiguous sections. Where the input is vague, make reasonable assumptions and note them explicitly.

---

## Project Input

**Project name:** e-commerce-site
**Description:** This is a simple e-commerce-site for selling books
**Team:** Feather Tech

**Primary users:** [[Fill from intake 2.1]]
**Problem being solved:** [[Fill from intake 2.2]]
**Success metrics (6 months):** [[Fill from intake 2.3]]

**Must-have features (v1):**
[[Fill from intake 3.1 — paste the list]]

**Out of scope (v1):**
[[Fill from intake 3.2]]

**External integrations:**
[[Fill from intake 3.3]]

**Target launch date:** [[Fill from intake 1.4]]

---

## Required PRD Sections

Generate the PRD with all of the following sections:

1. **Executive Summary** — 2-3 sentences: what, who, why
2. **Problem Statement** — current pain, root cause, impact
3. **Goals & Success Metrics** — SMART goals with measurable KPIs
4. **Non-Goals** — explicit scope exclusions
5. **User Personas** — 2-3 personas with needs and pain points
6. **User Stories** — written as "As a [persona], I want [action] so that [benefit]" — cover all v1 features
7. **Functional Requirements** — numbered, grouped by feature area
8. **Non-Functional Requirements** — performance, security, scalability, availability
9. **Assumptions** — list every assumption made
10. **Open Questions** — what must be resolved before or during implementation
11. **Out of Scope** — confirm what is explicitly excluded
12. **Timeline** — high-level milestone breakdown toward launch date

---

## Output Format

- Use Markdown with clear `##` and `###` headings
- Number all functional requirements (e.g. FR-001, FR-002)
- Flag assumptions with `> **Assumption:**`
- Flag open questions with `> **Open question:**`
- Target length: 1500–3000 words
