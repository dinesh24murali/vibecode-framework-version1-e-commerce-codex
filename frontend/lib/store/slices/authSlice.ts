import type { StateCreator } from 'zustand'

export interface AuthUser {
  id: string
  email: string
  first_name: string
  last_name: string
  role: 'customer' | 'admin'
}

export interface AuthTokens {
  access_token: string
  refresh_token?: string
}

export interface AuthSlice {
  user: AuthUser | null
  /** In memory only: never persisted (see memory/auth.md). */
  accessToken: string | null
  refreshToken: string | null
  login: (tokens: AuthTokens) => void
  logout: () => void
  setUser: (profile: AuthUser | null) => void
}

export const createAuthSlice: StateCreator<AuthSlice, [], [], AuthSlice> = (set) => ({
  user: null,
  accessToken: null,
  refreshToken: null,
  login: (tokens) =>
    set({ accessToken: tokens.access_token, refreshToken: tokens.refresh_token ?? null }),
  logout: () => set({ user: null, accessToken: null, refreshToken: null }),
  setUser: (profile) => set({ user: profile }),
})
