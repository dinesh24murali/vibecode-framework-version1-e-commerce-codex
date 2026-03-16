# verify/ — AGENTS.md

## Purpose

Verification is the step between "code written" and "task done". These scripts and playbooks ensure that what was built actually works from the user's perspective.

## When to Run Verification

**Always run `make verify` after:**
- Any frontend component change
- Any API endpoint change that affects the UI
- Any change to routing, auth flows, or form submissions

You may skip verification for:
- Pure backend changes with no UI impact (but run unit tests)
- Documentation-only changes
- Test file changes

## Verification Methods

### 1. Automated Scripts (always run first)
```bash
make verify
# Equivalent to:
bash verify/scripts/run-e2e.sh
npx ts-node verify/scripts/check-dom.ts
```

### 2. DOM Checklist (for complex UI changes)
Consult `verify/dom-checks.md` for a checklist of what to verify in the DOM after frontend changes.

### 3. MCP Playbook (for AI tools with browser access)
If you have access to Playwright MCP or Chrome DevTools MCP, follow `verify/mcp-playbook.md` for interactive browser verification.

## Rules for AI Tools

- Do not mark a task as done if `make verify` fails
- If e2e tests fail, fix the code — do not comment out the tests
- If the DOM check script finds issues, resolve them before completing the task
- Log any persistent verification failures in the task's scratch file
