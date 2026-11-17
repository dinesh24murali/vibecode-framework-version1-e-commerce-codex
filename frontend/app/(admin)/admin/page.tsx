'use client'
import Link from 'next/link'
import { useStore } from '@/lib/store'

// Dashboard; further admin sections arrive in TASK-029.
export default function AdminHomePage() {
  const user = useStore((s) => s.adminUser)
  const logout = useStore((s) => s.adminLogout)
  return (
    <main className="mx-auto max-w-2xl space-y-4 p-6">
      <h1 className="text-2xl font-semibold">Admin dashboard</h1>
      <p className="text-sm">Signed in as {user?.email}</p>
      <nav aria-label="Admin sections">
        <ul className="list-disc pl-5 text-sm">
          <li>
            <Link href="/admin/categories/" className="underline">
              Categories
            </Link>
          </li>
          <li>
            <Link href="/admin/products/" className="underline">
              Products
            </Link>
          </li>
          <li>
            <Link href="/admin/attributes/" className="underline">
              Attributes
            </Link>
          </li>
        </ul>
      </nav>
      <button onClick={logout} className="text-sm underline">
        Sign out
      </button>
    </main>
  )
}
