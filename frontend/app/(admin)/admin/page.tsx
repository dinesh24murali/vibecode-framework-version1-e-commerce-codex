'use client'
import { useStore } from '@/lib/store'

// Dashboard placeholder; admin sections arrive in TASK-026..029.
export default function AdminHomePage() {
  const user = useStore((s) => s.adminUser)
  const logout = useStore((s) => s.adminLogout)
  return (
    <main className="mx-auto max-w-2xl space-y-4 p-6">
      <h1 className="text-2xl font-semibold">Admin dashboard</h1>
      <p className="text-sm">Signed in as {user?.email}</p>
      <button onClick={logout} className="text-sm underline">
        Sign out
      </button>
    </main>
  )
}
