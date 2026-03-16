# MCP Verification Playbook

Use this playbook when verifying frontend changes with an AI tool that has browser access via MCP.

## Prerequisites

### Option A: Playwright MCP
```bash
# Install the Playwright MCP server
npx @playwright/mcp@latest install

# Start it (Claude Code will connect automatically if configured)
npx @playwright/mcp@latest start
```

### Option B: Chrome DevTools MCP
```bash
# Launch Chrome with remote debugging enabled
google-chrome --remote-debugging-port=9222 --no-first-run --no-default-browser-check

# Or on macOS:
open -a "Google Chrome" --args --remote-debugging-port=9222
```

Configure the MCP server in your AI tool's settings to point to `localhost:9222`.

---

## Standard Verification Flow

Run these checks after any significant frontend change.

### 1. Page Load Check
```
Navigate to: http://localhost:3000 (or your dev server URL)
Verify:
- Page renders without blank screen
- No error banners visible
- Browser console has 0 errors
```

### 2. Authentication Flow
```
Test: login with valid credentials
Expected: redirect to dashboard/home
Test: visit protected route while logged out
Expected: redirect to login page
Test: logout
Expected: session cleared, redirect to login
```

### 3. Core Feature Flows
For each v1 feature, run through the happy path:
```
Navigate to the feature's entry point
Complete the primary user action
Verify the success state
Verify the data persists (refresh the page and check)
```

### 4. Form Submission
```
Submit form with empty required fields
Expected: validation errors shown inline
Submit form with invalid data (e.g. bad email format)
Expected: specific validation message
Submit form with valid data
Expected: success state, no console errors
```

### 5. Error State Testing
```
Disconnect network (DevTools → Network → Offline)
Trigger an API call
Expected: user-facing error message (not a blank screen or raw error)
Reconnect and retry
Expected: operation succeeds
```

### 6. API Call Verification
```
Open Network tab
Trigger the feature action
Verify:
- Correct endpoint called (matches OpenAPI spec)
- Request payload is correct
- Response status is 2xx
- Response data matches expected schema
```

---

## Running from Claude Code

With Playwright MCP configured, you can ask Claude:

> "Use the browser to navigate to localhost:3000, log in, and verify the dashboard loads without errors. Check the browser console for errors."

> "Run the form submission flow for the [feature] page and verify all validation states work correctly."

## Running from Cursor

In Cursor's composer, with browser MCP enabled:

> "@browser Navigate to localhost:3000 and run through the standard verification checklist in verify/mcp-playbook.md"

---

## Failure Handling

If verification fails:
1. Capture a screenshot (ask the AI tool to take one)
2. Copy the console errors
3. Add failure details to the task's scratch file
4. Fix the issue before marking the task done
