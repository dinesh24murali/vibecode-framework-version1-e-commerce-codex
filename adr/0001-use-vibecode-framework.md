# ADR-0001: Use Vibecode Framework Structure

**Status:** Accepted
**Date:** [[PROJECT_START_DATE]]
**Deciders:** [[TEAM_NAME]]

---

## Context

Starting a new web application project requires agreeing on how the team (and AI coding agents) will coordinate: where documentation lives, how tasks are tracked, when decisions get recorded, and how code quality is enforced.

Without a shared structure, AI-assisted development degrades quickly — agents lose context between sessions, documentation drifts from the code, and there is no record of why architectural choices were made.

We evaluated three approaches:

1. **No structure** — each developer (and AI agent) improvises their own organization
2. **Wiki-based** — documentation in Confluence/Notion, tasks in Jira, decisions in Slack threads
3. **Vibecode Framework** — a git-native, file-based structure with AI-readable AGENTS.md files at every level

## Decision

We will use the **Vibecode Framework** as the project scaffold. All documentation, task tracking, architectural decisions, and AI instructions live in the repository alongside the code.

The framework provides:
- `docs/` for Phase 1 (document generation before coding)
- `tasks/` for Phase 2 (task-by-task development with scratch files)
- `adr/` for architectural decisions
- `memory/` for AI persistent context
- `AGENTS.md` files at each level to guide AI tools (Claude, Cursor, Codex, Gemini)

## Consequences

### Positive
- AI agents have full context without needing external system access
- All decisions are version-controlled and traceable via git history
- Framework is provider-agnostic — not locked to any single AI tool
- Onboarding a new contributor (human or AI) is a matter of reading `AGENTS.md`
- The scratch-file-before-coding rule prevents AI agents from diving in without a plan

### Negative
- Team must maintain discipline around updating `CHANGELOG.md`, `memory/`, and moving tasks to `done/`
- Adds file overhead to the repository (mitigated: all files are small markdown)
- Requires developers to learn the folder conventions upfront

### Neutral
- The framework is tech-stack agnostic — all `[[PLACEHOLDER]]` tokens must be replaced at project start

## Alternatives Considered

| Option | Reason rejected |
|--------|----------------|
| No structure | AI agents lose context between sessions; no audit trail of decisions |
| Wiki + Jira | External systems AI agents cannot easily read; documentation drifts from code |
| Custom structure | Requires design time; Vibecode Framework solves the known problems already |

## References

- `AGENTS.md` — master AI instructions
- `docs/00_intake/intake_questionnaire.md` — project intake form
- `README.md` — framework usage guide
