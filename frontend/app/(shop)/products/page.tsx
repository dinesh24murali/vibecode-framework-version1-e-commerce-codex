import { Suspense } from 'react'
import { ProductListPage } from '@/components/shop/ProductListPage'

export default function ProductsPage() {
  return (
    <main className="mx-auto max-w-6xl p-4">
      <h1 className="mb-4 text-2xl font-semibold">Products</h1>
      {/* useSearchParams requires a Suspense boundary in a static export */}
      <Suspense>
        <ProductListPage />
      </Suspense>
    </main>
  )
}
