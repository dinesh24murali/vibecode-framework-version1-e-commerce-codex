# Task: TASK-008 — Zustand v5 store — auth & cart slices

**Status:** active
**Created:** 2026-03-31
**ADR refs:** ADR-NNNN-zustand-ssr-hydration-strategy (created in TASK-003)

---

## Goal

Implement the `authSlice` and `cartSlice` using Zustand v5's slice pattern; persist auth tokens to `localStorage` safely in SSG context; expose a `syncWithServer()` action on the cart slice.

## Background

The auth slice is a prerequisite for TASK-007 (API client needs access to the token for header injection) and TASK-009 (auth pages read/write auth state). The cart slice is the prerequisite for all cart and checkout UI in Phase 2. Both slices must follow the SSG hydration strategy decided in TASK-003's ADR — no `localStorage` access during static generation.

See `docs/02_outputs/05_implementation_plan.md` TASK-008 and the ADR created in TASK-003.

## Acceptance Criteria

- [ ] `authSlice` hydrates from `localStorage` on the first client render without React hydration mismatch warnings
- [ ] `authSlice` exposes: `user`, `accessToken`, `refreshToken`, `login(tokens)`, `logout()`, `setUser(profile)` actions
- [ ] `cartSlice` accumulates line items client-side; exposes `addItem()`, `removeItem()`, `updateQuantity()`, `clearCart()`, `syncWithServer()` actions
- [ ] Both slices are combined in the root store (`frontend/lib/store/index.ts`) with correct TypeScript types exported
- [ ] `make verify` passes with no SSR hydration errors in the browser console
- [ ] Zustand `persist` middleware is used for `authSlice`; `cartSlice` does not persist (server state is authoritative on load)

## Dependencies

- **Tasks:** TASK-003 (frontend scaffold and Zustand installed; hydration ADR decided)
- **Memory files:** `memory/auth.md` (token storage and lifetime conventions)
- **Docs:** `docs/02_outputs/03_tech_architecture.md` §5 (Auth — access token in memory, refresh in HttpOnly cookie), `docs/02_outputs/05_implementation_plan.md` TASK-008, ADR from TASK-003

## Files Expected to Change

- `frontend/lib/store/slices/authSlice.ts` *(create)*
- `frontend/lib/store/slices/cartSlice.ts` *(create)*
- `frontend/lib/store/index.ts` *(update — combine slices, export types)*

## Notes

- Architecture mandates access tokens in memory (not `localStorage`); refresh tokens in `HttpOnly` cookie (managed by the browser, not JS). The `authSlice` should store the access token in-memory (Zustand state) — the Zustand `persist` middleware should only persist non-sensitive user profile data (e.g. `user.name`, `user.email`) so the UI can render without a re-fetch on page load
- The `syncWithServer()` cart action should call the cart API endpoints from TASK-007's generated client to reconcile local cart state with the server
- For SSG safety: use Zustand's `skipHydration: true` option or wrap `localStorage` access in a `useEffect` / `typeof window !== 'undefined'` guard — follow the pattern documented in TASK-003's ADR
- TypeScript: export a `StoreState` type and `useStore` hook from `frontend/lib/store/index.ts`

---

## Scratch File

AI: before implementing, create `tasks/active/task-008-zustand-store-auth-cart-slices.scratch.md` with your plan.
See `tasks/AGENTS.md` for the required scratch file format.
