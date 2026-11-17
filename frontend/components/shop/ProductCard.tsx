'use client'
import Link from 'next/link'
import type { ProductSummary } from '@/lib/api/generated/model'
import { useStore } from '@/lib/store'
import { Button } from '@/components/ui/button'

const inr = (n: number) => `₹${n.toFixed(2)}`

export function ProductCard({ product }: { product: ProductSummary }) {
  const addItem = useStore((s) => s.addItem)
  const hasDiscount = product.discount_percent > 0

  return (
    <article className="relative flex flex-col rounded-md border p-3" data-testid="product-card">
      {hasDiscount && (
        <span className="absolute right-2 top-2 rounded bg-red-600 px-2 py-0.5 text-xs text-white">
          {Math.round(product.discount_percent)}% off
        </span>
      )}
      <Link href={`/products/${product.slug}/`} className="block">
        {product.cover_image ? (
          // eslint-disable-next-line @next/next/no-img-element
          <img
            src={product.cover_image.url}
            alt={product.cover_image.alt_text ?? product.name}
            className="mb-2 h-40 w-full rounded object-cover"
          />
        ) : (
          <div className="mb-2 h-40 w-full rounded bg-muted" />
        )}
        <h3 className="font-medium">{product.name}</h3>
      </Link>
      <p className="mt-1 text-sm">
        <span className="font-semibold">{inr(product.discounted_price)}</span>
        {hasDiscount && <span className="ml-2 text-muted-foreground line-through">{inr(product.price)}</span>}
      </p>
      <Button className="mt-3" onClick={() => addItem(product.id, 1)}>
        Add to Cart
      </Button>
    </article>
  )
}
