'use client'
import Link from 'next/link'
import type { CartItem } from '@/lib/api/generated/model'

export const formatInr = (n: number) => `₹${n.toFixed(2)}`

interface Props {
  item: CartItem
  busy?: boolean
  onQuantity: (item: CartItem, quantity: number) => void
}

export function CartLineItem({ item, busy, onQuantity }: Props) {
  const { product } = item
  return (
    <li data-testid="cart-line" className="flex gap-4 border-b py-4">
      {product.cover_image ? (
        // eslint-disable-next-line @next/next/no-img-element
        <img src={product.cover_image.url} alt={product.cover_image.alt_text ?? product.name} className="h-20 w-20 rounded object-cover" />
      ) : (
        <div className="h-20 w-20 rounded bg-gray-100" aria-hidden />
      )}
      <div className="flex-1">
        <Link href={`/products/${product.slug}/`} className="font-medium hover:underline">
          {product.name}
        </Link>
        <p className="text-sm text-gray-600">{formatInr(item.unit_price)} each</p>
        <div className="mt-2 flex items-center gap-2" role="group" aria-label={`Quantity for ${product.name}`}>
          <button type="button" disabled={busy} aria-label={`Decrease quantity of ${product.name}`} className="h-8 w-8 rounded border disabled:opacity-50" onClick={() => onQuantity(item, item.quantity - 1)}>
            -
          </button>
          <span data-testid="cart-qty" className="w-8 text-center">{item.quantity}</span>
          <button type="button" disabled={busy} aria-label={`Increase quantity of ${product.name}`} className="h-8 w-8 rounded border disabled:opacity-50" onClick={() => onQuantity(item, item.quantity + 1)}>
            +
          </button>
          <button type="button" disabled={busy} className="ml-4 text-sm text-red-700 underline disabled:opacity-50" onClick={() => onQuantity(item, 0)}>
            Remove
          </button>
        </div>
      </div>
      <p className="font-medium" data-testid="cart-line-total">{formatInr(item.line_total)}</p>
    </li>
  )
}
