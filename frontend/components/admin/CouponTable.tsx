'use client'
import type { Coupon } from '@/lib/api/generated/model'

/** Admin coupon list with edit, enable/disable toggle and delete actions. */
export function CouponTable({
  coupons,
  onEdit,
  onToggle,
  onDelete,
}: {
  coupons: Coupon[]
  onEdit: (c: Coupon) => void
  onToggle: (c: Coupon) => void
  onDelete: (c: Coupon) => void
}) {
  if (coupons.length === 0) return <p className="text-sm text-gray-600">No coupons yet.</p>
  return (
    <table className="w-full text-left text-sm">
      <thead>
        <tr className="border-b">
          <th className="py-2">Code</th>
          <th>Discount</th>
          <th>Expires</th>
          <th>Status</th>
          <th className="text-right">Actions</th>
        </tr>
      </thead>
      <tbody>
        {coupons.map((c) => (
          <tr key={c.id} className={`border-b ${c.is_enable ? '' : 'text-gray-400'}`}>
            <td className="py-2">{c.code}</td>
            <td>{c.type === 'flat' ? `₹${c.flat} off` : `${c.percent}% off`}</td>
            <td>{c.coupon_expiry ? new Date(c.coupon_expiry).toLocaleDateString('en-IN') : '—'}</td>
            <td>{c.is_enable ? 'Enabled' : 'Disabled'}</td>
            <td className="space-x-2 text-right">
              <button className="underline" onClick={() => onToggle(c)}>
                {c.is_enable ? 'Disable' : 'Enable'}
              </button>
              <button className="underline" onClick={() => onEdit(c)}>
                Edit
              </button>
              <button className="text-red-600 underline" onClick={() => onDelete(c)}>
                Delete
              </button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}
