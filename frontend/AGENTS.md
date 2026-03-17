# frontend/ — AGENTS.md

## Stack

- **Framework:** `NextJS 16`
- **UI library:** `shadcn/ui`
- **API client:** generated from `docs/02_outputs/04_api_spec.yaml`

> Replace these placeholders with your actual stack after cloning.

---

## Project Structure

<!-- Document your frontend folder structure here once it is established.
     Example for a Next.js project:
     frontend/
     ├── src/
     │   ├── app/           ← Next.js app router pages
     │   ├── components/    ← Reusable UI components
     │   ├── api/           ← Generated API client (from OpenAPI spec)
     │   ├── hooks/         ← Custom React hooks
     │   └── lib/           ← Utilities
     ├── public/
     └── package.json
-->

---

## Conventions

### API Client
- The API client in `src/api/` (or equivalent) is **generated from** `docs/02_outputs/04_api_spec.yaml`
- Do not manually write API fetch calls that duplicate the spec
- To regenerate: `[[CODEGEN_COMMAND]]` (fill in your actual codegen command)

### Component Structure
<!-- Document component conventions here. Example:
- Each component in its own file
- Co-locate tests: `MyComponent.test.tsx` next to `MyComponent.tsx`
- Use named exports, not default exports
-->

### State Management
<!-- Document state management approach. Example:
- Server state: React Query / SWR
- Client state: Zustand / Context API
-->

### Routing
<!-- Document routing conventions. -->

---

## Running Locally

```bash
# Fill in with actual commands for NextJS 16
# Example for Next.js:
# cd frontend && npm run dev
```

## Running Tests

```bash
# Fill in with actual test commands
# Example:
# cd frontend && npm test
```

---

## Rules for AI Tools

- **Never write API client code manually** — regenerate from the OpenAPI spec
- After any change to a page or component, run `make verify`
- Read `memory/auth.md` before touching auth-related components
- All new pages need a route defined and an e2e test in `tests/e2e/`
- Do not import from `backend/` — all data flows through the API layer
