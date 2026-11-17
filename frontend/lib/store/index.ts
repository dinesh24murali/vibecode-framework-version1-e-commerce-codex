import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'
import { createAuthSlice, type AuthSlice } from './slices/authSlice'
import { createCartSlice, type CartSlice } from './slices/cartSlice'

export type { AuthSlice, AuthUser, AuthTokens } from './slices/authSlice'
export type { CartSlice, CartLine, CartApi } from './slices/cartSlice'

export type StoreState = AuthSlice & CartSlice

/**
 * SSG-safe root store (ADR-0003). `skipHydration: true` keeps localStorage out of
 * the static build; `useHydrateStore` rehydrates after mount.
 * Only the non-sensitive `user` profile is persisted: tokens stay in memory
 * (memory/auth.md) and the cart is server-authoritative.
 */
export const useStore = create<StoreState>()(
  persist(
    (...a) => ({
      ...createAuthSlice(...a),
      ...createCartSlice(...a),
    }),
    {
      name: 'e-commerce-store',
      storage: createJSONStorage(() =>
        typeof window !== 'undefined'
          ? localStorage
          : { getItem: () => null, setItem: () => {}, removeItem: () => {} }
      ),
      partialize: (s) => ({ user: s.user }),
      skipHydration: true,
    }
  )
)
