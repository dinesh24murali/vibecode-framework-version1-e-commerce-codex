# ADR-0003 — Zustand SSR/SSG Hydration Strategy

| Field | Value |
|-------|-------|
| **Status** | Accepted |
| **Date** | 2026-03-31 |
| **Deciders** | Feather Tech |

---

## Context

The storefront is a Next.js 16 App Router application configured with `output: 'export'`. Static generation runs in a Node.js build environment where `localStorage` does not exist. Zustand's `persist` middleware reads from `localStorage` by default on store creation, which causes a build-time crash if called during SSG.

Additionally, App Router renders Server Components on the server; any store that reads `localStorage` at module load time will throw in that context.

## Decision

Use Zustand `persist` middleware with `skipHydration: true`.

- The store is created with `skipHydration: true`, which tells the middleware not to call `storage.getItem` during store initialization.
- A `createJSONStorage` factory guards the `localStorage` access behind a `typeof window !== 'undefined'` check, providing a no-op fallback for the build environment.
- On the client, after the component mounts, `useStore.persist.rehydrate()` is called to restore persisted state. A `useHydrateStore` hook (added in TASK-008) will encapsulate this pattern and be placed in the root layout.

## Alternatives Considered

### Check `typeof window !== 'undefined'` inline in each slice
Rejected — requires every slice author to remember the guard. `skipHydration` centralizes the concern in the store definition.

### Use cookies instead of localStorage for persistence
Rejected for the initial phase — the client-only state (cart preview, UI preferences) does not need to be server-readable. LocalStorage is simpler and sufficient.

### Use `zustand/react` with `createStore` + React Context
Rejected — adds boilerplate for a static site that does not need server-per-request store instances. `create()` with `skipHydration` is sufficient for SSG.

## Consequences

- Store state is not available during the first server render (expected — the site is static).
- Developers adding new persisted slices do not need to add any SSG guards themselves.
- `useHydrateStore` must be mounted in the root client layout before any slice reads persisted state (enforced in TASK-008).
