# e-commerce-site — Vibecode Framework

> A provider-agnostic vibe-coding framework. Works with Claude Code, Cursor, Codex, Gemini, or any AI tool.

---

## Quick Start

```bash
# 1. Clone this template
git clone <this-repo> my-project
cd my-project

# 2. Find-and-replace all placeholders
#    e-commerce-site, Go, NextJS 16 Static Site Generation, etc.
#    See AGENTS.md §6 for the full placeholder list

# 3. Scaffold directories
make scaffold

# 4. Fill in the intake questionnaire
open docs/00_intake/intake_questionnaire.md
```

---

## Two-Phase Workflow

### Phase 1 — Generate Your Documents (do this before writing code)

Work through `docs/01_prompts/` in order. For each file:
1. Open the `.prompt.md` file
2. Fill in the `[[PLACEHOLDER]]` tokens from your intake questionnaire
3. Paste the entire prompt into your AI tool of choice
4. Save the output to `docs/02_outputs/` (e.g. `docs/02_outputs/01_prd.md`)

| Prompt | Output |
|--------|--------|
| `01_prd.prompt.md` | Product Requirements Document |
| `02_functional_spec.prompt.md` | Functional Specifications |
| `03_tech_architecture.prompt.md` | Technical Architecture |
| `04_api_spec.prompt.md` | OpenAPI 3.1 Spec |
| `05_implementation_plan.prompt.md` | Phased Implementation Plan |
| `06_dev_setup.prompt.md` | Developer Setup Guide |

After generating docs:
- Review and adjust outputs
- Create your first ADRs from architectural decisions: `make adr SLUG=initial-stack`
- Seed memory files from the architecture doc: `memory/<domain>.md`

### Phase 2 — Build Features (the living project loop)

For every feature or bugfix:

```
1. make task NAME=my-feature
2. Ask AI: "Write the scratch file for tasks/active/my-feature.md"
3. Review the scratch file — adjust if needed
4. Tell AI to implement
5. AI updates CHANGELOG.md and memory files
6. make verify   (for frontend changes)
7. AI moves task to tasks/done/
```

---

## Project Structure

```
.
├── docs/
│   ├── 00_intake/          # Fill this first
│   ├── 01_prompts/         # Paste into any AI tool
│   └── 02_outputs/         # AI-generated docs land here
├── adr/                    # Architecture Decision Records
├── memory/                 # AI memory / domain knowledge
├── tasks/
│   ├── active/             # In-progress tasks + scratch files
│   └── done/               # Completed tasks
├── tests/
│   ├── e2e/                # Playwright end-to-end tests
│   └── contract/           # OpenAPI contract tests
├── verify/                 # Verification scripts and playbooks
├── backend/                # Backend source code
├── frontend/               # Frontend source code
├── AGENTS.md               # Master AI instructions (all tools read this)
├── CLAUDE.md               # Redirects to AGENTS.md
└── .cursorrules            # Redirects to AGENTS.md
```

---

## AI Tool Setup

This framework is **provider-agnostic**. The master instructions live in `AGENTS.md`.

| Tool | Config file | Action |
|------|-------------|--------|
| Claude Code | `CLAUDE.md` | Redirects to `AGENTS.md` |
| Cursor | `.cursorrules` | Redirects to `AGENTS.md` |
| Codex / GPT | — | Paste `AGENTS.md` content as system prompt |
| Gemini | — | Paste `AGENTS.md` content as system prompt |

---

## Makefile Commands

```bash
make scaffold   # Create directories (run once)
make verify     # Run e2e + DOM checks
make adr SLUG=my-decision   # Create a new ADR
make task NAME=my-feature   # Create a new task file
make help       # Show all commands
```

---

## Philosophy

- **Docs first, code second** — AI-generated documents become the spec; code follows
- **Scratch files before code** — AI must plan before implementing
- **CHANGELOG discipline** — every task leaves a trail
- **Provider-agnostic** — prompts are plain markdown; paste anywhere
- **Living memory** — domain knowledge is captured in `memory/`, not lost between sessions
