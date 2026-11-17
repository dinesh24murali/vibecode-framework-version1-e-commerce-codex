'use client'
import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { adminLogin } from '@/lib/api/generated/admin-auth/admin-auth'
import type { AuthTokenResponse } from '@/lib/api/generated/model'
import { errorMessage } from '@/lib/auth/helpers'
import { useStore } from '@/lib/store'
import { Button } from '@/components/ui/button'

export function AdminLoginForm() {
  const router = useRouter()
  const adminLoginStore = useStore((s) => s.adminLogin)
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault()
    setBusy(true)
    setError(null)
    try {
      const res = await adminLogin({ email, password })
      if (res.status !== 200) {
        setError(errorMessage(res, 'Invalid email or password'))
        return
      }
      const data = res.data as AuthTokenResponse
      if (data.user?.role !== 'admin') {
        setError('This account does not have admin access')
        return
      }
      adminLoginStore(data, data.user as never)
      router.push('/admin/')
    } catch {
      setError('Network error. Please try again.')
    } finally {
      setBusy(false)
    }
  }

  return (
    <form onSubmit={onSubmit} className="space-y-4">
      <div className="space-y-1">
        <label htmlFor="admin-email" className="block text-sm font-medium">
          Email
        </label>
        <input
          id="admin-email"
          type="email"
          required
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="w-full rounded-md border px-3 py-2"
        />
      </div>
      <div className="space-y-1">
        <label htmlFor="admin-password" className="block text-sm font-medium">
          Password
        </label>
        <input
          id="admin-password"
          type="password"
          required
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="w-full rounded-md border px-3 py-2"
        />
      </div>
      {error && (
        <p role="alert" className="text-sm text-red-600">
          {error}
        </p>
      )}
      <Button type="submit" disabled={busy} className="w-full">
        Sign in
      </Button>
    </form>
  )
}
