import { AdminLoginForm } from '@/components/admin/AdminLoginForm'

export default function AdminLoginPage() {
  return (
    <main className="mx-auto max-w-sm space-y-6 p-6">
      <h1 className="text-2xl font-semibold">Admin sign in</h1>
      <AdminLoginForm />
    </main>
  )
}
