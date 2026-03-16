# memory/ — AGENTS.md

## Purpose

Memory files capture domain knowledge that would otherwise be lost between AI sessions. They are the project's long-term memory — facts, decisions, and context that an AI agent needs to do good work, but that aren't obvious from reading the code.

## When to Create a Memory File

Create a memory file when you learn:
- How a domain concept works (e.g. how the billing cycle is calculated)
- A non-obvious constraint (e.g. "the payments API has a 100ms SLA")
- A convention that isn't enforced by tooling (e.g. "all dates are stored as UTC ISO 8601")
- Context behind a decision that isn't in an ADR
- User preferences or feedback that should persist across sessions

## When to Update a Memory File

Update an existing memory file when:
- A fact changes (update in place, note the change)
- A new piece of knowledge belongs to an existing domain
- A memory is found to be wrong

## File Naming

Use `<domain>.md` for the filename:
- `memory/auth.md` — authentication and session behavior
- `memory/billing.md` — billing and subscription logic
- `memory/users.md` — user model and roles
- `memory/api.md` — API conventions and quirks
- `memory/infra.md` — infrastructure and deployment notes

## Index

Always update `memory/_index.md` when creating or significantly changing a memory file.

## Rules for AI Tools

- **Read relevant memory files before starting any task** — check `_index.md` first
- **Never store ephemeral state** in memory files (current task progress, etc.)
- Do not duplicate information that is already in ADRs or code comments
- Keep each memory file focused on one domain — split if it grows beyond ~200 lines
- Write memory entries as permanent facts, not personal notes
