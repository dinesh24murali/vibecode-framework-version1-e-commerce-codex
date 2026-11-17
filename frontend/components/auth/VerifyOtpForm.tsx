'use client'
import { useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { registerComplete, registerVerifyOTP } from '@/lib/api/generated/auth/auth'
import type { AuthTokenResponse, OTPTokenResponse } from '@/lib/api/generated/model'
import { errorMessage } from '@/lib/auth/helpers'
import { useStore } from '@/lib/store'
import { Button } from '@/components/ui/button'
import { OtpInput } from './OtpInput'

/** Signup steps 2 (verify OTP) and 3 (name + password). */
export function VerifyOtpForm() {
  const router = useRouter()
  const email = useSearchParams().get('email') ?? ''
  const login = useStore((s) => s.login)
  const [otp, setOtp] = useState('')
  const [otpToken, setOtpToken] = useState<string | null>(null)
  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  async function verify(e: React.FormEvent) {
    e.preventDefault()
    setBusy(true)
    setError(null)
    try {
      const res = await registerVerifyOTP({ email, otp })
      if (res.status !== 200) {
        setError(errorMessage(res, 'Invalid or expired code')) // otp value is kept
        return
      }
      setOtpToken((res.data as OTPTokenResponse).otp_token)
    } catch {
      setError('Network error. Please try again.')
    } finally {
      setBusy(false)
    }
  }

  async function complete(e: React.FormEvent) {
    e.preventDefault()
    if (!otpToken) return
    setBusy(true)
    setError(null)
    try {
      const res = await registerComplete({
        otp_token: otpToken,
        first_name: firstName,
        last_name: lastName,
        password,
      })
      if (res.status !== 201) {
        setError(errorMessage(res))
        return
      }
      login(res.data as AuthTokenResponse)
      router.push('/')
    } catch {
      setError('Network error. Please try again.')
    } finally {
      setBusy(false)
    }
  }

  if (!otpToken) {
    return (
      <form onSubmit={verify} className="space-y-4">
        <p className="text-sm">We sent a code to {email}.</p>
        <OtpInput value={otp} onChange={setOtp} error={error} />
        <Button type="submit" disabled={busy || otp.length !== 6} className="w-full">
          Verify
        </Button>
      </form>
    )
  }

  return (
    <form onSubmit={complete} className="space-y-4">
      <div className="space-y-1">
        <label htmlFor="first-name" className="block text-sm font-medium">
          First name
        </label>
        <input
          id="first-name"
          required
          value={firstName}
          onChange={(e) => setFirstName(e.target.value)}
          className="w-full rounded-md border px-3 py-2"
        />
      </div>
      <div className="space-y-1">
        <label htmlFor="last-name" className="block text-sm font-medium">
          Last name
        </label>
        <input
          id="last-name"
          required
          value={lastName}
          onChange={(e) => setLastName(e.target.value)}
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
          minLength={8}
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
        Create account
      </Button>
    </form>
  )
}
