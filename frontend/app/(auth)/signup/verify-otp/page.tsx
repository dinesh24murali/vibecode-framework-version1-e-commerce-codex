import { Suspense } from 'react'
import { VerifyOtpForm } from '@/components/auth/VerifyOtpForm'

export default function VerifyOtpPage() {
  return (
    <main className="mx-auto max-w-sm space-y-6 p-6">
      <h1 className="text-2xl font-semibold">Verify your email</h1>
      <Suspense>
        <VerifyOtpForm />
      </Suspense>
    </main>
  )
}
