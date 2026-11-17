'use client'
import { useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { login as loginRequest } from '@/lib/api/generated/auth/auth'
import type { AuthTokenResponse } from '@/lib/api/generated/model'
import { errorMessage, safeRedirect } from '@/lib/auth/helpers'
import { useStore } from '@/lib/store'
import { Button } from '@/components/ui/button'
import { FacebookLoginButton } from './FacebookLoginButton'

export function LoginForm() {
  const router = useRouter()
  const redirect = safeRedirect(useSearchParams().get('redirect'))
  const login = useStore((s) => s.login)
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault()
    setBusy(true)
    setError(null)
    try {
      const res = await loginRequest({ email, password })
      if (res.status !== 200) {
        setError(errorMessage(res, 'Invalid email or password'))
        return
      }
      login(res.data as AuthTokenResponse)
      router.push(redirect)
    } catch {
      setError('Network error. Please try again.')
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="space-y-4">
      <form onSubmit={onSubmit} className="space-y-4">
        <div className="space-y-1">
          <label htmlFor="email" className="block text-sm font-medium">
            Email
          </label>
          <input
            id="email"
            type="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full rounded-md border px-3 py-2"
          />
        </div>
        <div className="space-y-1">
          <label htmlFor="password" className="block text-sm font-medium">
            Password
          </label>
          <input
            id="password"
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
          Log in
        </Button>
      </form>
      <FacebookLoginButton redirect={redirect} />
    </div>
  )
}
