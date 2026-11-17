'use client'
import { useState } from 'react'
import { adminCreateCoupon, adminUpdateCoupon } from '@/lib/api/generated/admin-coupons/admin-coupons'
import type { Coupon } from '@/lib/api/generated/model'
import { errorMessage } from '@/lib/auth/helpers'
import { Button } from '@/components/ui/button'
import { Modal } from './Modal'

type CouponKind = 'flat' | 'percentage'

/** Create (no `coupon`) or edit a coupon. The type decides whether the flat or the percent field is used. */
export function CouponFormModal({
  coupon,
  onClose,
  onSaved,
}: {
  coupon?: Coupon
  onClose: () => void
  onSaved: () => void
}) {
  const [name, setName] = useState(coupon?.name ?? '')
  const [code, setCode] = useState(coupon?.code ?? '')
  const [type, setType] = useState<CouponKind>(coupon?.type ?? 'flat')
  const [flat, setFlat] = useState(coupon?.flat != null ? String(coupon.flat) : '')
  const [percent, setPercent] = useState(coupon?.percent != null ? String(coupon.percent) : '')
  const [expiry, setExpiry] = useState(coupon?.coupon_expiry ? coupon.coupon_expiry.slice(0, 10) : '')
  const [enabled, setEnabled] = useState(coupon?.is_enable ?? true)
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault()
    setBusy(true)
    setError(null)
    const amount = type === 'flat' ? { flat: Number(flat) } : { percent: Number(percent) }
    const coupon_expiry = expiry ? `${expiry}T23:59:59Z` : null
    try {
      const res = coupon
        ? await adminUpdateCoupon(coupon.id, { name, type, is_enable: enabled, coupon_expiry, ...amount })
        : await adminCreateCoupon({ name, code, type, is_enable: enabled, coupon_expiry, ...amount })
      if (res.status !== 200 && res.status !== 201) {
        setError(res.status === 409 ? 'That coupon code is already in use' : errorMessage(res, 'Could not save the coupon'))
        return
      }
      onSaved()
    } catch {
      setError('Network error. Please try again.')
    } finally {
      setBusy(false)
    }
  }

  const input = 'w-full rounded-md border px-3 py-2'
  return (
    <Modal title={coupon ? 'Edit coupon' : 'New coupon'} onClose={onClose}>
      <form onSubmit={onSubmit} className="space-y-4">
        <div className="space-y-1">
          <label htmlFor="cp-name" className="block text-sm font-medium">
            Name
          </label>
          <input id="cp-name" required className={input} value={name} onChange={(e) => setName(e.target.value)} />
        </div>
        <div className="space-y-1">
          <label htmlFor="cp-code" className="block text-sm font-medium">
            Code
          </label>
          <input
            id="cp-code"
            required
            disabled={!!coupon}
            className={`${input} uppercase disabled:bg-gray-100`}
            value={code}
            onChange={(e) => setCode(e.target.value)}
          />
        </div>
        <div className="space-y-1">
          <label htmlFor="cp-type" className="block text-sm font-medium">
            Type
          </label>
          <select id="cp-type" className={input} value={type} onChange={(e) => setType(e.target.value as CouponKind)}>
            <option value="flat">Flat amount (INR)</option>
            <option value="percentage">Percentage</option>
          </select>
        </div>
        {type === 'flat' ? (
          <div className="space-y-1">
            <label htmlFor="cp-flat" className="block text-sm font-medium">
              Flat amount (INR)
            </label>
            <input id="cp-flat" required type="number" min="1" className={input} value={flat} onChange={(e) => setFlat(e.target.value)} />
          </div>
        ) : (
          <div className="space-y-1">
            <label htmlFor="cp-percent" className="block text-sm font-medium">
              Percent off
            </label>
            <input id="cp-percent" required type="number" min="1" max="100" className={input} value={percent} onChange={(e) => setPercent(e.target.value)} />
          </div>
        )}
        <div className="space-y-1">
          <label htmlFor="cp-expiry" className="block text-sm font-medium">
            Expires on (optional)
          </label>
          <input id="cp-expiry" type="date" className={input} value={expiry} onChange={(e) => setExpiry(e.target.value)} />
        </div>
        <label className="flex items-center gap-2 text-sm">
          <input type="checkbox" checked={enabled} onChange={(e) => setEnabled(e.target.checked)} />
          Enabled
        </label>
        {error && (
          <p role="alert" className="text-sm text-red-600">
            {error}
          </p>
        )}
        <div className="flex justify-end gap-2">
          <Button type="button" variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" disabled={busy}>
            Save
          </Button>
        </div>
      </form>
    </Modal>
  )
}
