'use client'
import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useStore } from '@/lib/store'
import { getCart } from '@/lib/api/generated/cart/cart'
import { listAddresses } from '@/lib/api/generated/account/account'
import { calculateCheckout, placeOrder } from '@/lib/api/generated/checkout/checkout'
import type { Address, Cart, CheckoutSummary as Summary } from '@/lib/api/generated/model'
import { errorMessage } from '@/lib/auth/helpers'
import { AddressSelector } from './AddressSelector'
import { CouponInput } from './CouponInput'
import { CheckoutSummary } from './CheckoutSummary'

export function CheckoutPageContent() {
  const router = useRouter()
  const accessToken = useStore((s) => s.accessToken)
  const [ready, setReady] = useState(false)
  const [addresses, setAddresses] = useState<Address[]>([])
  const [addressId, setAddressId] = useState<string | null>(null)
  const [coupon, setCoupon] = useState<string | null>(null)
  const [summary, setSummary] = useState<Summary | null>(null)
  const [placing, setPlacing] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Guards. Access tokens are memory-only, so no token means not signed in.
  useEffect(() => {
    if (!accessToken) router.replace('/login/?redirect=/checkout/')
  }, [accessToken, router])

  // The slice may be empty right after login, so the server cart decides the empty-cart redirect.
  useEffect(() => {
    if (!accessToken) return
    let cancelled = false
    void (async () => {
      try {
        const [cartRes, addrRes] = await Promise.all([getCart(), listAddresses()])
        if (cancelled) return
        if (cartRes.status === 200 && (cartRes.data as Cart).items.length === 0) {
          router.replace('/cart/')
          return
        }
        if (addrRes.status === 200) {
          const rows = (addrRes.data as { data: Address[] }).data
          setAddresses(rows)
          setAddressId((rows.find((a) => a.is_default) ?? rows[0])?.id ?? null)
        }
        setReady(true)
      } catch {
        if (!cancelled) setError('We could not load checkout. Please try again.')
      }
    })()
    return () => {
      cancelled = true
    }
  }, [accessToken, router])

  // Recalculate whenever the address or the applied coupon changes.
  useEffect(() => {
    if (!ready || !addressId) return
    let cancelled = false
    void (async () => {
      try {
        const res = await calculateCheckout({ address_id: addressId, coupon_code: coupon })
        if (cancelled) return
        if (res.status === 200) {
          setSummary(res.data as Summary)
          setError(null)
        } else if (res.status === 400) {
          router.replace('/cart/')
        } else {
          setError(errorMessage(res, 'We could not calculate your total.'))
        }
      } catch {
        if (!cancelled) setError('We could not calculate your total.')
      }
    })()
    return () => {
      cancelled = true
    }
  }, [ready, addressId, coupon, router])

  async function place() {
    if (!addressId) return
    setPlacing(true)
    setError(null)
    try {
      const res = await placeOrder({ address_id: addressId, coupon_code: coupon })
      if (res.status === 201) {
        useStore.setState({ items: [] }) // the server cleared the cart
        router.push('/account/?tab=orders')
        return
      }
      setError(errorMessage(res, 'We could not place your order.'))
    } catch {
      setError('We could not place your order.')
    } finally {
      setPlacing(false)
    }
  }

  if (!accessToken) return null

  return (
    <main className="mx-auto max-w-5xl p-4">
      <h1 className="mb-4 text-3xl font-semibold">Checkout</h1>
      {error && <p role="alert" className="mb-3 text-red-700">{error}</p>}
      {!ready ? (
        <p role="status">Loading checkout…</p>
      ) : (
        <div className="grid gap-8 md:grid-cols-[1fr_20rem]">
          <div className="space-y-6">
            <AddressSelector addresses={addresses} selectedId={addressId} onChange={setAddressId} />
            <CouponInput onApply={setCoupon} />
          </div>
          <CheckoutSummary summary={summary} placing={placing} canPlace={!!summary && !!addressId} onPlace={place} />
        </div>
      )}
    </main>
  )
}
