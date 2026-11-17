'use client'
import type { Product } from '@/lib/api/generated/model'

/** Admin product list. Inactive (soft-deleted) products are shown greyed with a badge. */
export function ProductTable({
  products,
  onEdit,
  onDelete,
}: {
  products: Product[]
  onEdit: (p: Product) => void
  onDelete: (p: Product) => void
}) {
  if (products.length === 0) return <p className="text-sm text-gray-600">No products found.</p>
  return (
    <table className="w-full text-left text-sm">
      <thead>
        <tr className="border-b">
          <th className="py-2">Name</th>
          <th>Price</th>
          <th>Stock</th>
          <th>Status</th>
          <th className="text-right">Actions</th>
        </tr>
      </thead>
      <tbody>
        {products.map((p) => (
          <tr key={p.id} className={`border-b ${p.is_active ? '' : 'text-gray-400'}`}>
            <td className="py-2">{p.name}</td>
            <td>₹{p.discounted_price}</td>
            <td>{p.stock}</td>
            <td>{p.is_active ? 'Active' : 'Hidden'}</td>
            <td className="space-x-2 text-right">
              <button className="underline" onClick={() => onEdit(p)}>
                Edit
              </button>
              {p.is_active && (
                <button className="text-red-600 underline" onClick={() => onDelete(p)}>
                  Delete
                </button>
              )}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}
