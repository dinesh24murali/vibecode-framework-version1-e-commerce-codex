'use client'
import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { loginFacebook } from '@/lib/api/generated/auth/auth'
import type { AuthTokenResponse } from '@/lib/api/generated/model'
import { errorMessage } from '@/lib/auth/helpers'
import { useStore } from '@/lib/store'
import { Button } from '@/components/ui/button'

interface FBSdk {
  init: (o: { appId: string; version: string; cookie?: boolean }) => void
  login: (
    cb: (r: { authResponse?: { accessToken: string } | null }) => void,
    opts: { scope: string }
  ) => void
}
declare global {
  interface Window {
    FB?: FBSdk
  }
}

const APP_ID = process.env.NEXT_PUBLIC_FACEBOOK_APP_ID ?? ''

function loadSdk(): Promise<FBSdk> {
  return new Promise((resolve, reject) => {
    if (window.FB) return resolve(window.FB)
    const s = document.createElement('script')
    s.src = 'https://connect.facebook.net/en_US/sdk.js'
    s.async = true
    s.onload = () => {
      if (!window.FB) return reject(new Error('Facebook SDK unavailable'))
      window.FB.init({ appId: APP_ID, version: 'v19.0', cookie: true })
      resolve(window.FB)
    }
    s.onerror = () => reject(new Error('Facebook SDK failed to load'))
    document.head.appendChild(s)
  })
}

export function FacebookLoginButton({ redirect = '/' }: { redirect?: string }) {
  const router = useRouter()
  const login = useStore((s) => s.login)
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  async function onClick() {
    setBusy(true)
    setError(null)
    try {
      const fb = await loadSdk()
      const accessToken = await new Promise<string | null>((resolve) =>
        fb.login((r) => resolve(r.authResponse?.accessToken ?? null), {
          scope: 'email,public_profile',
        })
      )
      if (!accessToken) {
        setError('Facebook login was cancelled')
        return
      }
      const res = await loginFacebook({ access_token: accessToken })
      if (res.status !== 200) {
        setError(errorMessage(res, 'Facebook login failed'))
        return
      }
      login(res.data as AuthTokenResponse)
      router.push(redirect)
    } catch {
      setError('Facebook login is unavailable right now')
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="space-y-1">
      <Button type="button" variant="outline" onClick={onClick} disabled={busy} className="w-full">
        Continue with Facebook
      </Button>
      {error && (
        <p role="alert" className="text-sm text-red-600">
          {error}
        </p>
      )}
    </div>
  )
}
