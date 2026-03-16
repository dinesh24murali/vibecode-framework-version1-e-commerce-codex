# tasks/ — AGENTS.md

## Purpose

Task files track units of work. They are the AI agent's contract with the developer: here is what I will do, here is how you'll know it's done.

## Directory Layout

```
tasks/
├── _template.md       ← Copy this to create a new task
├── active/            ← In-progress tasks and scratch files
└── done/              ← Completed tasks (archive)
```

## Workflow

### Creating a Task
```bash
make task NAME=my-feature
# Creates: tasks/active/my-feature.md
```

### The Scratch File Rule (MANDATORY)

**Before writing any code**, the AI must create a scratch file:

```
tasks/active/<task-name>.scratch.md
```

The scratch file must contain:
1. **Plan** — numbered steps for implementation
2. **Files to change** — list every file that will be created or modified
3. **Uncertainties** — open questions or things to verify before coding
4. **What will be skipped** — explicitly note any out-of-scope items
5. **Risks** — anything that could go wrong

The developer reviews the scratch file before implementation begins. Do not skip this step.

### Completing a Task
1. Ensure all acceptance criteria in the task file are checked off
2. Update `CHANGELOG.md` with a brief entry
3. Update any relevant `memory/<domain>.md` files
4. Move task file and scratch file to `tasks/done/`:
   ```bash
   mv tasks/active/my-feature.md tasks/done/
   mv tasks/active/my-feature.scratch.md tasks/done/
   ```

## Rules for AI Tools

- **Never start implementing before writing the scratch file**
- One task file per unit of work — don't combine unrelated changes
- Keep task files updated as work progresses (check off criteria as they're met)
- If a task grows significantly larger than expected, split it and create sub-tasks
- Reference ADR numbers in task files when implementing architectural decisions
