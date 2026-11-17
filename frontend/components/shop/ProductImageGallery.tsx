'use client'
import { useState } from 'react'
import type { ProductImage } from '@/lib/api/generated/model'
import { ProductZoom } from './ProductZoom'

export function ProductImageGallery({ images, name }: { images: ProductImage[]; name: string }) {
  const sorted = [...images].sort((a, b) => a.sort_order - b.sort_order)
  const [active, setActive] = useState(0)
  if (sorted.length === 0) return <div className="aspect-square w-full rounded bg-muted" />
  const current = sorted[active] ?? sorted[0]

  return (
    <div>
      <ProductZoom src={current.url} alt={current.alt_text ?? name} />
      {sorted.length > 1 && (
        <ul className="mt-3 flex gap-2 overflow-x-auto" aria-label="Product images">
          {sorted.map((img, i) => (
            <li key={img.id}>
              <button
                type="button"
                data-testid="thumbnail"
                aria-label={`Show image ${i + 1}`}
                aria-current={i === active}
                onClick={() => setActive(i)}
                className={`h-16 w-16 overflow-hidden rounded border-2 ${i === active ? 'border-black' : 'border-transparent'}`}
              >
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img src={img.url} alt="" className="h-full w-full object-cover" />
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
