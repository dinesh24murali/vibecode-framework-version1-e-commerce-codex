import Link from 'next/link'
import type { Cart } from '@/lib/api/generated/model'
import { formatInr } from './CartLineItem'

export function CartSummary({ cart, authenticated }: { cart: Cart; authenticated: boolean }) {
  const mrp = cart.items.reduce((n, i) => n + i.product.price * i.quantity, 0)
  const discount = Math.max(0, mrp - cart.subtotal)
  return (
    <aside className="h-fit rounded border p-4" aria-label="Order summary">
      <h2 className="mb-3 text-lg font-semibold">Summary</h2>
      <dl className="space-y-1 text-sm">
        <div className="flex justify-between"><dt>Subtotal</dt><dd data-testid="summary-subtotal">{formatInr(mrp)}</dd></div>
        <div className="flex justify-between"><dt>Discount</dt><dd data-testid="summary-discount">-{formatInr(discount)}</dd></div>
        <div className="flex justify-between border-t pt-2 font-semibold"><dt>Total</dt><dd data-testid="summary-total">{formatInr(cart.subtotal)}</dd></div>
      </dl>
      <p className="mt-2 text-xs text-gray-500">Taxes are calculated at checkout.</p>
      {authenticated ? (
        <Link href="/checkout/" className="mt-4 block rounded bg-black px-4 py-2 text-center text-white">
          Proceed to Checkout
        </Link>
      ) : (
        <div className="mt-4 space-y-2">
          <button type="button" disabled className="w-full rounded bg-black px-4 py-2 text-white opacity-50">Proceed to Checkout</button>
          <p className="text-sm">
            <Link href="/login/?redirect=/cart/" className="underline">Log in</Link> to check out.
          </p>
        </div>
      )}
    </aside>
  )
}
