import type { CheckoutSummary as Summary } from '@/lib/api/generated/model'
import { formatInr } from './CartLineItem'

interface Props {
  summary: Summary | null
  placing: boolean
  canPlace: boolean
  onPlace: () => void
}

export function CheckoutSummary({ summary, placing, canPlace, onPlace }: Props) {
  return (
    <aside className="h-fit rounded border p-4" aria-label="Order summary">
      <h2 className="mb-3 text-lg font-semibold">Order summary</h2>
      {summary ? (
        <>
          <ul className="mb-3 space-y-1 text-sm">
            {summary.items.map((i) => (
              <li key={i.id} className="flex justify-between">
                <span>{i.product.name} × {i.quantity}</span>
                <span>{formatInr(i.line_total)}</span>
              </li>
            ))}
          </ul>
          <dl className="space-y-1 border-t pt-2 text-sm">
            <div className="flex justify-between"><dt>Subtotal</dt><dd data-testid="co-subtotal">{formatInr(summary.subtotal)}</dd></div>
            {summary.coupon_discount > 0 && (
              <div className="flex justify-between"><dt>Coupon ({summary.coupon_code})</dt><dd data-testid="co-discount">-{formatInr(summary.coupon_discount)}</dd></div>
            )}
            <div className="flex justify-between"><dt>Taxable amount</dt><dd>{formatInr(summary.taxable_amount)}</dd></div>
            <div className="flex justify-between"><dt>GST (18%)</dt><dd data-testid="co-gst">{formatInr(summary.gst_amount)}</dd></div>
            <div className="flex justify-between border-t pt-2 font-semibold"><dt>Total</dt><dd data-testid="co-total">{formatInr(summary.total)}</dd></div>
          </dl>
        </>
      ) : (
        <p className="text-sm text-gray-500">Calculating…</p>
      )}
      <button
        type="button"
        onClick={onPlace}
        disabled={!canPlace || placing}
        className="mt-4 w-full rounded bg-black px-4 py-2 text-white disabled:opacity-50"
      >
        {placing ? 'Placing order…' : 'Place Order'}
      </button>
    </aside>
  )
}
