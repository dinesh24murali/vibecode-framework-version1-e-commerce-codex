# adr/ — AGENTS.md

## Purpose

Architecture Decision Records (ADRs) document significant technical decisions made during the project. They create a permanent, searchable history of *why* things were built the way they were.

## When to Create an ADR

Create an ADR when a decision:
- Is hard to reverse (framework choice, database schema, auth strategy)
- Affects multiple teams or components
- Has meaningful trade-offs worth recording
- Was debated and a specific option was chosen over alternatives
- Will confuse future contributors if unexplained

Do **not** create ADRs for routine implementation details.

## Naming Convention

```
adr/NNNN-short-slug.md
```

- `NNNN` is a zero-padded 4-digit number (0001, 0002, ...)
- `short-slug` is a kebab-case description (3-5 words)
- Examples: `0001-use-postgresql.md`, `0002-jwt-auth-strategy.md`

## How to Create an ADR

```bash
make adr SLUG=my-decision
```

This copies `adr/template.md` to `adr/NNNN-my-decision.md`.

## Status Values

- `Accepted` — decided and in effect (**default for every new ADR**)
- `Proposed` — reserved for humans who explicitly want a decision left open; AI tools never use it
- `Deprecated` — superseded by a newer ADR
- `Superseded by ADR-NNNN` — link to replacement

## Auto-Acceptance Policy

ADRs are **auto-approved**. No human sign-off is required for work to continue.

- AI-created ADRs are written with `**Status:** Accepted`, `**Deciders:** AI agent (auto-accepted)`, and today's date.
- **Never stop or wait for ADR approval.** Create the ADR, reference it in the task and scratch files, and continue implementing in the same session.
- Humans can override an auto-accepted decision at any time by creating a new ADR that supersedes it.

## Rules for AI Tools

- Reference ADRs in task files when implementing a decision
- If an implementation needs to contradict an existing ADR, do not block: create a new auto-accepted ADR that supersedes it, mark the old one `Superseded by ADR-NNNN`, note the change in the scratch file, and continue
- Do not modify `Accepted` ADRs — create a new one that supersedes the old one (the status field of the old ADR is the only permitted edit)
- Update the status field if a decision changes
