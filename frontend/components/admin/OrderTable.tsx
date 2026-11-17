'use client'
import type { Order } from '@/lib/api/generated/model'

/** Admin order list; each row opens the detail modal. */
export function OrderTable({ orders, onOpen }: { orders: Order[]; onOpen: (o: Order) => void }) {
  if (orders.length === 0) return <p className="text-sm text-gray-600">No orders found.</p>
  return (
    <table className="w-full text-left text-sm">
      <thead>
        <tr className="border-b">
          <th className="py-2">Order</th>
          <th>Customer</th>
          <th>Date</th>
          <th>Total</th>
          <th>Status</th>
          <th className="text-right">Actions</th>
        </tr>
      </thead>
      <tbody>
        {orders.map((o) => (
          <tr key={o.id} className="border-b">
            <td className="py-2">{o.order_number}</td>
            <td>{o.customer_name ?? '—'}</td>
            <td>{new Date(o.created_at).toLocaleDateString('en-IN')}</td>
            <td>₹{o.total}</td>
            <td className="capitalize">{o.status}</td>
            <td className="text-right">
              <button className="underline" onClick={() => onOpen(o)}>
                View
              </button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}
