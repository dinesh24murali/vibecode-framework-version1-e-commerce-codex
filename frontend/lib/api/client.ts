import { useStore } from '@/lib/store'

const BASE_URL = process.env.NEXT_PUBLIC_API_URL ?? ''

/**
 * Hand-maintained base client used as the orval mutator. The ONLY place in
 * frontend/ allowed to call fetch(). Injects the Bearer token when present.
 */
export const apiFetch = async <T>(url: string, options: RequestInit = {}): Promise<T> => {
  const state = typeof window !== 'undefined' ? useStore.getState() : null
  // Admin endpoints use the separate admin session; the login call itself sends no token.
  const isAdminCall = url.startsWith('/api/v1/admin/') && !url.startsWith('/api/v1/admin/auth/')
  const token = isAdminCall ? (state?.adminAccessToken ?? null) : (state?.accessToken ?? null)
  const headers = new Headers(options.headers)
  if (token) headers.set('Authorization', `Bearer ${token}`)
  const res = await fetch(`${BASE_URL}${url}`, { ...options, headers, credentials: 'include' })
  const text = [204, 205, 304].includes(res.status) ? '' : await res.text()
  const data = text ? JSON.parse(text) : {}
  return { data, status: res.status, headers: res.headers } as T
}
