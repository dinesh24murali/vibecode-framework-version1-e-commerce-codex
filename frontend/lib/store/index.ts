import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'

/**
 * Root store type — extend with slice interfaces as feature tasks are implemented.
 * See tasks/active/task-008-zustand-store-auth-cart-slices.md for auth and cart slices.
 */
export interface RootState {
  _hydrated: boolean
  setHydrated: () => void
}

/**
 * SSG-safe store.
 *
 * `skipHydration: true` prevents Zustand from reading localStorage during the
 * Next.js static build, which runs in a Node.js environment where localStorage
 * does not exist. The store is manually hydrated on the client after mount via
 * `useHydrateStore` (see lib/store/useHydrateStore.ts, added in TASK-008).
 *
 * See ADR-0003 for the full rationale.
 */
export const useStore = create<RootState>()(
  persist(
    (set) => ({
      _hydrated: false,
      setHydrated: () => set({ _hydrated: true }),
    }),
    {
      name: 'e-commerce-store',
      storage: createJSONStorage(() =>
        typeof window !== 'undefined' ? localStorage : ({} as Storage)
      ),
      skipHydration: true,
    }
  )
)
