'use client'
import { useEffect, useState } from 'react'
import { useStore } from '@/lib/store'
import { Button } from '@/components/ui/button'

export function AddToCartButton({ productId, name }: { productId: string; name: string }) {
  const addItem = useStore((s) => s.addItem)
  const [toast, setToast] = useState<string | null>(null)

  useEffect(() => {
    if (!toast) return
    const t = setTimeout(() => setToast(null), 3000)
    return () => clearTimeout(t)
  }, [toast])

  return (
    <>
      <Button
        onClick={() => {
          addItem(productId, 1)
          setToast(`${name} added to cart`)
        }}
      >
        Add to Cart
      </Button>
      {toast && (
        <div role="status" data-testid="toast" className="fixed bottom-4 right-4 rounded bg-black px-4 py-2 text-white shadow">
          {toast}
        </div>
      )}
    </>
  )
}
