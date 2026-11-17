import type { ProductSummary } from '@/lib/api/generated/model'
import { ProductCard } from './ProductCard'

export function ProductGrid({ products }: { products: ProductSummary[] }) {
  if (products.length === 0) {
    return <p className="p-6 text-center text-sm">Products are temporarily unavailable.</p>
  }
  return (
    <div className="grid grid-cols-2 gap-4 p-4 md:grid-cols-4" data-testid="product-grid">
      {products.map((p) => (
        <ProductCard key={p.id} product={p} />
      ))}
    </div>
  )
}
