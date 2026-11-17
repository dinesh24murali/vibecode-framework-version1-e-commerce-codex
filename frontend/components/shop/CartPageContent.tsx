'use client'
import { useCallback, useEffect, useState } from 'react'
import Link from 'next/link'
import { useStore } from '@/lib/store'
import { getCart, removeCartItem, updateCartItem } from '@/lib/api/generated/cart/cart'
import type { Cart, CartItem } from '@/lib/api/generated/model'
import { CartLineItem } from './CartLineItem'
import { CartSummary } from './CartSummary'

const round2 = (n: number) => Math.round(n * 100) / 100

/** Applies a quantity change locally (0 removes the line) so the UI updates before the API answers. */
function applyLocal(cart: Cart, itemId: string, quantity: number): Cart {
  const items = cart.items
    .map((i) => (i.id === itemId ? { ...i, quantity, line_total: round2(i.unit_price * quantity) } : i))
    .filter((i) => i.quantity > 0)
  return { ...cart, items, item_count: items.length, subtotal: round2(items.reduce((n, i) => n + i.line_total, 0)) }
}

export function CartPageContent() {
  const accessToken = useStore((s) => s.accessToken)
  const [cart, setCart] = useState<Cart | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  /** Mirrors the server cart into cartSlice (nav badge) and the page state. */
  const adopt = useCallback((next: Cart) => {
    setCart(next)
    useStore.setState({
      items: next.items.map((i) => ({ productId: i.product.id, quantity: i.quantity, itemId: i.id, unitPrice: i.unit_price })),
    })
  }, [])

  useEffect(() => {
    if (!accessToken) return
    let cancelled = false
    setLoading(true)
    getCart()
      .then((res) => {
        if (cancelled) return
        if (res.status !== 200) throw new Error('load failed')
        adopt(res.data as Cart)
      })
      .catch(() => !cancelled && setError('We could not load your cart. Please try again.'))
      .finally(() => !cancelled && setLoading(false))
    return () => {
      cancelled = true
    }
  }, [accessToken, adopt])

  const changeQuantity = async (item: CartItem, quantity: number) => {
    if (!cart || busy) return
    const previous = cart
    setError(null)
    setBusy(true)
    adopt(applyLocal(cart, item.id, quantity)) // optimistic
    try {
      const res = quantity <= 0 ? await removeCartItem(item.id) : await updateCartItem(item.id, { quantity })
      if (res.status !== 200) throw new Error('update failed')
      adopt(res.data as Cart) // confirmed by the API
    } catch {
      adopt(previous) // roll back
      setError('We could not update your cart. Please try again.')
    } finally {
      setBusy(false)
    }
  }

  return (
    <main className="mx-auto max-w-5xl p-4">
      <h1 className="mb-4 text-3xl font-semibold">Your Cart</h1>
      {!accessToken ? (
        <p data-testid="cart-login-prompt">
          Please <Link href="/login/?redirect=/cart/" className="underline">log in</Link> to view your cart.
        </p>
      ) : loading && !cart ? (
        <p role="status">Loading your cart…</p>
      ) : (
        <>
          {error && <p role="alert" className="mb-3 text-red-700">{error}</p>}
          {cart && cart.items.length > 0 ? (
            <div className="grid gap-8 md:grid-cols-[1fr_20rem]">
              <ul data-testid="cart-lines">
                {cart.items.map((i) => (
                  <CartLineItem key={i.id} item={i} busy={busy} onQuantity={changeQuantity} />
                ))}
              </ul>
              <CartSummary cart={cart} authenticated />
            </div>
          ) : (
            cart && (
              <div data-testid="cart-empty" className="py-12 text-center">
                <svg aria-hidden viewBox="0 0 64 64" className="mx-auto mb-4 h-24 w-24 text-gray-300" fill="none" stroke="currentColor" strokeWidth="3">
                  <path d="M6 10h8l6 32h30l6-22H16" strokeLinecap="round" strokeLinejoin="round" />
                  <circle cx="26" cy="52" r="3" />
                  <circle cx="46" cy="52" r="3" />
                </svg>
                <p className="mb-4 text-lg">Your cart is empty.</p>
                <Link href="/products/" className="rounded bg-black px-4 py-2 text-white">Continue Shopping</Link>
              </div>
            )
          )}
        </>
      )}
    </main>
  )
}
