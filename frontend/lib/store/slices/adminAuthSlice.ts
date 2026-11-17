import type { StateCreator } from 'zustand'
import type { AuthUser } from './authSlice'

/** Admin session: separate from the customer session, in memory only (memory/auth.md). */
export interface AdminAuthSlice {
  adminUser: AuthUser | null
  adminAccessToken: string | null
  adminRefreshToken: string | null
  adminLogin: (tokens: { access_token: string; refresh_token?: string }, user: AuthUser) => void
  adminLogout: () => void
}

export const createAdminAuthSlice: StateCreator<AdminAuthSlice, [], [], AdminAuthSlice> = (set) => ({
  adminUser: null,
  adminAccessToken: null,
  adminRefreshToken: null,
  adminLogin: (tokens, user) =>
    set({
      adminUser: user,
      adminAccessToken: tokens.access_token,
      adminRefreshToken: tokens.refresh_token ?? null,
    }),
  adminLogout: () => set({ adminUser: null, adminAccessToken: null, adminRefreshToken: null }),
})
