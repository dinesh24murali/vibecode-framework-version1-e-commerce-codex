import Link from 'next/link'
import { SignupForm } from '@/components/auth/SignupForm'

export default function SignupPage() {
  return (
    <main className="mx-auto max-w-sm space-y-6 p-6">
      <h1 className="text-2xl font-semibold">Create your account</h1>
      <SignupForm />
      <p className="text-sm">
        Already registered?{' '}
        <Link href="/login/" className="underline">
          Log in
        </Link>
      </p>
    </main>
  )
}
