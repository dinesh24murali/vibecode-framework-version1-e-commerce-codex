'use client'
import { useCallback, useEffect, useState } from 'react'
import { adminListOrders } from '@/lib/api/generated/admin-orders/admin-orders'
import type { Order } from '@/lib/api/generated/model'
import { errorMessage } from '@/lib/auth/helpers'
import { Button } from '@/components/ui/button'
import { OrderDetailModal } from '@/components/admin/OrderDetailModal'
import { OrderTable } from '@/components/admin/OrderTable'

const PAGE_SIZE = 20
const tabs = [
  { id: 'all', label: 'All' },
  { id: 'open', label: 'Open' },
  { id: 'fulfilled', label: 'Fulfilled' },
  { id: 'cancelled', label: 'Cancelled' },
] as const
type Tab = (typeof tabs)[number]['id']

// The detail view is a modal that fetches the order by id (no dynamic [id] route: static export).
export default function AdminOrdersPage() {
  const [tab, setTab] = useState<Tab>('all')
  const [page, setPage] = useState(1)
  const [orders, setOrders] = useState<Order[]>([])
  const [totalPages, setTotalPages] = useState(1)
  const [openId, setOpenId] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loaded, setLoaded] = useState(false)

  const load = useCallback(async () => {
    try {
      const res = await adminListOrders({ page, limit: PAGE_SIZE, status: tab === 'all' ? undefined : tab })
      if (res.status === 200) {
        setOrders(res.data.data)
        setTotalPages(Math.max(res.data.meta.total_pages, 1))
        setError(null)
      } else setError(errorMessage(res, 'Could not load orders'))
    } catch {
      setError('Network error. Please try again.')
    } finally {
      setLoaded(true)
    }
  }, [tab, page])

  useEffect(() => {
    void load()
  }, [load])

  return (
    <main className="mx-auto max-w-4xl space-y-6 p-6">
      <h1 className="text-2xl font-semibold">Orders</h1>
      <div role="tablist" className="flex gap-2 border-b">
        {tabs.map((t) => (
          <button
            key={t.id}
            role="tab"
            aria-selected={tab === t.id}
            onClick={() => {
              setTab(t.id)
              setPage(1)
            }}
            className={`px-3 py-2 text-sm ${tab === t.id ? 'border-b-2 border-black font-medium' : ''}`}
          >
            {t.label}
          </button>
        ))}
      </div>
      {error && (
        <p role="alert" className="text-sm text-red-600">
          {error}
        </p>
      )}
      {loaded && <OrderTable orders={orders} onOpen={(o) => setOpenId(o.id)} />}
      <div className="flex items-center justify-between text-sm">
        <Button variant="outline" disabled={page <= 1} onClick={() => setPage((p) => p - 1)}>
          Previous
        </Button>
        <span>
          Page {page} of {totalPages}
        </span>
        <Button variant="outline" disabled={page >= totalPages} onClick={() => setPage((p) => p + 1)}>
          Next
        </Button>
      </div>
      {openId && (
        <OrderDetailModal orderId={openId} onClose={() => setOpenId(null)} onChanged={() => void load()} />
      )}
    </main>
  )
}
