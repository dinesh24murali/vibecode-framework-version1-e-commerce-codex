'use client'
import { useEffect } from 'react'
import { usePathname, useRouter } from 'next/navigation'
import { useStore } from '@/lib/store'

/** Redirects every /admin/* page except the login page to /admin/login when there is no admin session. */
export function AdminGuard({ children }: { children: React.ReactNode }) {
  const router = useRouter()
  const pathname = usePathname() ?? ''
  const token = useStore((s) => s.adminAccessToken)
  const isLogin = pathname.replace(/\/$/, '') === '/admin/login'

  useEffect(() => {
    if (!isLogin && !token) router.replace('/admin/login/')
  }, [isLogin, token, router])

  if (!isLogin && !token) return null
  return <>{children}</>
}
