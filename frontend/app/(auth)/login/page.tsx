import { Suspense } from 'react'
import Link from 'next/link'
import { LoginForm } from '@/components/auth/LoginForm'

export default function LoginPage() {
  return (
    <main className="mx-auto max-w-sm space-y-6 p-6">
      <h1 className="text-2xl font-semibold">Log in</h1>
      <Suspense>
        <LoginForm />
      </Suspense>
      <p className="text-sm">
        New here?{' '}
        <Link href="/signup/" className="underline">
          Create an account
        </Link>
      </p>
    </main>
  )
}
