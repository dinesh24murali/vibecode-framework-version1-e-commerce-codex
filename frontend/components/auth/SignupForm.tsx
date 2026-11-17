'use client'
import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { registerInit } from '@/lib/api/generated/auth/auth'
import { errorMessage } from '@/lib/auth/helpers'
import { Button } from '@/components/ui/button'

export function SignupForm() {
  const router = useRouter()
  const [email, setEmail] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault()
    setBusy(true)
    setError(null)
    try {
      const res = await registerInit({ email })
      if (res.status !== 200) {
        setError(errorMessage(res))
        return
      }
      router.push(`/signup/verify-otp/?email=${encodeURIComponent(email)}`)
    } catch {
      setError('Network error. Please try again.')
    } finally {
      setBusy(false)
    }
  }

  return (
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
      {error && (
        <p role="alert" className="text-sm text-red-600">
          {error}
        </p>
      )}
      <Button type="submit" disabled={busy} className="w-full">
        Send verification code
      </Button>
    </form>
  )
}
