'use client'
import { useEffect, useState } from 'react'
import { adminGetOrder, adminUpdateOrderStatus } from '@/lib/api/generated/admin-orders/admin-orders'
import type { Order } from '@/lib/api/generated/model'
import { errorMessage } from '@/lib/auth/helpers'
import { Button } from '@/components/ui/button'
import { Modal } from './Modal'

/** Order detail: customer, delivery address, line items, total, and a status dropdown. */
export function OrderDetailModal({
  orderId,
  onClose,
  onChanged,
}: {
  orderId: string
  onClose: () => void
  onChanged: () => void
}) {
  const [order, setOrder] = useState<Order | null>(null)
  const [status, setStatus] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  useEffect(() => {
    let cancelled = false
    adminGetOrder(orderId)
      .then((res) => {
        if (cancelled) return
        if (res.status === 200) setOrder(res.data)
        else setError(errorMessage(res, 'Could not load the order'))
      })
      .catch(() => !cancelled && setError('Network error. Please try again.'))
    return () => {
      cancelled = true
    }
  }, [orderId])

  async function save() {
    if (!status) return
    setBusy(true)
    setError(null)
    try {
      const res = await adminUpdateOrderStatus(orderId, { status: status as 'fulfilled' | 'cancelled' })
      if (res.status !== 200) {
        setError(errorMessage(res, 'Could not update the status'))
        return
      }
      setOrder(res.data)
      setStatus('')
      onChanged()
    } catch {
      setError('Network error. Please try again.')
    } finally {
      setBusy(false)
    }
  }

  const a = order?.shipping_address
  return (
    <Modal title={order ? `Order ${order.order_number}` : 'Order'} onClose={onClose}>
      {!order && !error && <p className="text-sm">Loading…</p>}
      {order && (
        <div className="space-y-4 text-sm">
          <div>
            <p className="font-medium">{order.customer_name ?? '—'}</p>
            {order.customer_email && <p className="text-gray-600">{order.customer_email}</p>}
          </div>
          {a && (
            <address className="not-italic text-gray-700">
              {a.address_line1}
              {a.address_line2 ? `, ${a.address_line2}` : ''}
              <br />
              {a.city}, {a.state} {a.pincode}
            </address>
          )}
          <ul className="divide-y border-y">
            {order.items.map((it) => (
              <li key={it.id} className="flex justify-between py-1">
                <span>
                  {it.product_name} × {it.quantity}
                </span>
                <span>₹{it.line_total}</span>
              </li>
            ))}
          </ul>
          <p className="flex justify-between font-medium">
            <span>Total</span>
            <span>₹{order.total}</span>
          </p>
          <div className="flex items-end gap-2">
            <div className="flex-1 space-y-1">
              <label htmlFor="order-status" className="block font-medium">
                Status ({order.status})
              </label>
              <select
                id="order-status"
                value={status}
                onChange={(e) => setStatus(e.target.value)}
                disabled={order.status !== 'open'}
                className="w-full rounded-md border px-3 py-2"
              >
                <option value="">Change status…</option>
                <option value="fulfilled">Fulfilled</option>
                <option value="cancelled">Cancelled</option>
              </select>
            </div>
            <Button onClick={() => void save()} disabled={busy || !status}>
              Update
            </Button>
          </div>
        </div>
      )}
      {error && (
        <p role="alert" className="text-sm text-red-600">
          {error}
        </p>
      )}
    </Modal>
  )
}
