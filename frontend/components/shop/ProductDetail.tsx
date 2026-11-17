'use client'
import { useEffect, useState } from 'react'
import { usePathname } from 'next/navigation'
import { getProduct, listProducts } from '@/lib/api/generated/products/products'
import type { Product } from '@/lib/api/generated/model'
import { ProductImageGallery } from './ProductImageGallery'
import { AddToCartButton } from './AddToCartButton'

const inr = (n: number) => `₹${n.toFixed(2)}`

/** Resolves a slug to a product: the spec only offers lookup by id, so search then fetch. */
async function fetchBySlug(slug: string): Promise<Product | null> {
  const list = await listProducts({ q: slug, limit: 100 })
  if (list.status !== 200) return null
  const hit = list.data.data.find((p) => p.slug === slug)
  if (!hit) return null
  const res = await getProduct(hit.id)
  return res.status === 200 ? (res.data as Product) : null
}

export function ProductDetail({ slug, initial }: { slug: string; initial: Product | null }) {
  const pathname = usePathname()
  const [product, setProduct] = useState<Product | null>(initial)
  const [state, setState] = useState<'ready' | 'loading' | 'missing'>(initial ? 'ready' : 'loading')

  useEffect(() => {
    if (initial) return
    // Not pre-rendered with data (API was down at build): resolve from the URL on the client.
    const real = pathname.split('/').filter(Boolean).pop() ?? slug
    let cancelled = false
    fetchBySlug(decodeURIComponent(real))
      .then((p) => {
        if (cancelled) return
        setProduct(p)
        setState(p ? 'ready' : 'missing')
        if (p) document.title = p.name
      })
      .catch(() => {
        if (!cancelled) setState('missing')
      })
    return () => {
      cancelled = true
    }
  }, [initial, pathname, slug])

  if (state === 'loading') return <p>Loading...</p>
  if (!product) return <p role="alert">Product not found.</p>

  const hasDiscount = product.discount_percent > 0
  return (
    <div className="grid gap-8 md:grid-cols-2">
      <ProductImageGallery images={product.images} name={product.name} />
      <div>
        <h1 className="text-2xl font-semibold">{product.name}</h1>
        <p className="mt-2 text-xl">
          <span className="font-semibold">{inr(product.discounted_price)}</span>
          {hasDiscount && (
            <>
              <span className="ml-2 text-muted-foreground line-through">{inr(product.price)}</span>
              <span className="ml-2 text-sm text-red-600">{Math.round(product.discount_percent)}% off</span>
            </>
          )}
        </p>
        <dl className="mt-4 text-sm">
          {product.attributes.map((a) => (
            <div key={a.attribute_id + a.value_slug} className="flex gap-2">
              <dt className="font-medium">{a.attribute_name}:</dt>
              <dd>{a.value}</dd>
            </div>
          ))}
        </dl>
        <p className="mt-4 whitespace-pre-line text-sm">{product.description.replace(/<[^>]*>/g, '')}</p>
        <div className="mt-6">
          {product.stock > 0 ? (
            <AddToCartButton productId={product.id} name={product.name} />
          ) : (
            <p className="text-red-600">Out of stock</p>
          )}
        </div>
      </div>
    </div>
  )
}
