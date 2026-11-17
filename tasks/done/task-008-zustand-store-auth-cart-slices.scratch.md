# Scratch: task-008 Zustand auth & cart slices

## Plan
1. authSlice: user, accessToken, refreshToken, login/logout/setUser. Tokens live in memory only (memory/auth.md, arch §5); persist `partialize` stores only `user` (non-sensitive profile).
2. cartSlice: client-side line items, addItem/removeItem/updateQuantity/clearCart/syncWithServer. Not persisted (partialize excludes it).
3. syncWithServer takes an injectable `CartApi` adapter (set via `setCartApi`), because the generated client (TASK-007) does not exist yet (008 runs before 007). TASK-007 wires the real client. Decision recorded here, no hand-written fetch.
4. Root store index.ts combines slices, exports StoreState, useStore.
5. useHydrateStore hook (per ADR-0003) calling persist.rehydrate() in useEffect; mounted via a client `StoreHydrator` component in root layout.
6. fe-build + fe-lint, then verify.

## Files
- frontend/lib/store/slices/authSlice.ts (create)
- frontend/lib/store/slices/cartSlice.ts (create)
- frontend/lib/store/index.ts (update)
- frontend/lib/store/useHydrateStore.ts (create)
- frontend/components/StoreHydrator.tsx (create), frontend/app/layout.tsx (update)

## Uncertainties (resolved)
- Task says "persist auth tokens" but architecture says tokens in memory only; spec wins -> persist user profile only.
- Cart API client absent -> adapter interface.

## Skipped
- Real API wiring (TASK-007), refresh flow (TASK-009).

## Risks
- Hydration mismatch: avoided via skipHydration + rehydrate in effect.
