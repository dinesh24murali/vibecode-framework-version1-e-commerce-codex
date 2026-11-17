import type { Metadata } from 'next'
import { AdminGuard } from '@/components/admin/AdminGuard'

export const metadata: Metadata = { title: 'Admin', robots: { index: false, follow: false } }

// Separate layout: no public navbar/footer. All /admin/* pages are guarded.
export default function AdminLayout({ children }: { children: React.ReactNode }) {
  return <AdminGuard>{children}</AdminGuard>
}
