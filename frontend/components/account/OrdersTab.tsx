'use client'
import { useCallback, useEffect, useState } from 'react'
import { cancelOrder, listOrders } from '@/lib/api/generated/orders/orders'
import type { OrderListResponse } from '@/lib/api/generated/model'
import { errorMessage } from '@/lib/auth/helpers'
import { Button } from '@/components/ui/button'

const PAGE_SIZE = 10

const badge: Record<string, string> = {
  open: 'bg-blue-100 text-blue-800',
  fulfilled: 'bg-green-100 text-green-800',
  cancelled: 'bg-gray-200 text-gray-700',
}

export function OrdersTab() {
  const [page, setPage] = useState(1)
  const [data, setData] = useState<OrderListResponse | null>(null)
  const [error, setError] = useState<string | null>(null)

  const load = useCallback(async (p: number) => {
    const res = await listOrders({ page: p, limit: PAGE_SIZE })
    if (res.status === 200) {
      setData(res.data as OrderListResponse)
      setError(null)
    } else setError(errorMessage(res))
  }, [])

  useEffect(() => {
    void load(page)
  }, [load, page])

  async function cancel(id: string) {
    const res = await cancelOrder(id)
    if (res.status === 200) await load(page)
    else setError(errorMessage(res))
  }

  const totalPages = data?.meta.total_pages ?? 1

  return (
    <div className="space-y-4">
      {error && (
        <p role="alert" className="text-sm text-red-600">
          {error}
        </p>
      )}
      <ul className="space-y-3">
        {data?.data.map((o) => (
          <li key={o.id} className="rounded-md border p-3 text-sm" data-testid="order-row">
            <div className="flex items-center justify-between">
              <span className="font-medium">{o.order_number}</span>
              <span className={`rounded px-2 py-0.5 text-xs ${badge[o.status] ?? ''}`}>{o.status}</span>
            </div>
            <p>
              {o.items.length} item(s) · ₹{o.total.toFixed(2)} · {new Date(o.created_at).toLocaleDateString()}
            </p>
            {o.status === 'open' && (
              <Button size="sm" variant="outline" className="mt-2" onClick={() => void cancel(o.id)}>
                Cancel
              </Button>
            )}
          </li>
        ))}
        {data && data.data.length === 0 && <li className="text-sm">No orders yet.</li>}
      </ul>
      <div className="flex items-center gap-3 text-sm">
        <Button size="sm" variant="outline" disabled={page <= 1} onClick={() => setPage(page - 1)}>
          Previous
        </Button>
        <span>
          Page {page} of {totalPages}
        </span>
        <Button size="sm" variant="outline" disabled={page >= totalPages} onClick={() => setPage(page + 1)}>
          Next
        </Button>
      </div>
    </div>
  )
}
