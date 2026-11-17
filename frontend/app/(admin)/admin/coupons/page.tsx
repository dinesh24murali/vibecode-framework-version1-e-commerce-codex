'use client'
import { useCallback, useEffect, useState } from 'react'
import {
  adminDeleteCoupon,
  adminListCoupons,
  adminUpdateCoupon,
} from '@/lib/api/generated/admin-coupons/admin-coupons'
import type { Coupon } from '@/lib/api/generated/model'
import { errorMessage } from '@/lib/auth/helpers'
import { Button } from '@/components/ui/button'
import { CouponFormModal } from '@/components/admin/CouponFormModal'
import { CouponTable } from '@/components/admin/CouponTable'

type ModalState = { coupon?: Coupon } | null

// Create/edit happen in a modal on this page (no dynamic [id] route: static export).
export default function AdminCouponsPage() {
  const [coupons, setCoupons] = useState<Coupon[]>([])
  const [modal, setModal] = useState<ModalState>(null)
  const [error, setError] = useState<string | null>(null)
  const [loaded, setLoaded] = useState(false)

  const load = useCallback(async () => {
    try {
      const res = await adminListCoupons({ limit: 100 })
      if (res.status === 200) setCoupons(res.data.data)
      else setError(errorMessage(res, 'Could not load coupons'))
    } catch {
      setError('Network error. Please try again.')
    } finally {
      setLoaded(true)
    }
  }, [])

  useEffect(() => {
    void load()
  }, [load])

  async function toggle(c: Coupon) {
    const res = await adminUpdateCoupon(c.id, { is_enable: !c.is_enable })
    if (res.status === 200) {
      setError(null)
      await load()
    } else setError(errorMessage(res, 'Could not update the coupon'))
  }

  async function remove(c: Coupon) {
    if (!window.confirm(`Delete coupon "${c.code}"?`)) return
    const res = await adminDeleteCoupon(c.id)
    if (res.status === 200) {
      setError(null)
      await load()
    } else if (res.status === 409) {
      setError(`"${c.code}" has already been redeemed and cannot be deleted. Disable it instead.`)
    } else setError(errorMessage(res, 'Could not delete the coupon'))
  }

  return (
    <main className="mx-auto max-w-3xl space-y-6 p-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Coupons</h1>
        <Button onClick={() => setModal({})}>New coupon</Button>
      </div>
      {error && (
        <p role="alert" className="text-sm text-red-600">
          {error}
        </p>
      )}
      {loaded && (
        <CouponTable
          coupons={coupons}
          onEdit={(c) => setModal({ coupon: c })}
          onToggle={(c) => void toggle(c)}
          onDelete={(c) => void remove(c)}
        />
      )}
      {modal && (
        <CouponFormModal
          coupon={modal.coupon}
          onClose={() => setModal(null)}
          onSaved={() => {
            setModal(null)
            setError(null)
            void load()
          }}
        />
      )}
    </main>
  )
}
