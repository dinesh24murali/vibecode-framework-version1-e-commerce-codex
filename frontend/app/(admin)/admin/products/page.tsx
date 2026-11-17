'use client'
import { useCallback, useEffect, useState } from 'react'
import { adminDeleteProduct, adminListProducts } from '@/lib/api/generated/admin-products/admin-products'
import type { Product } from '@/lib/api/generated/model'
import { errorMessage } from '@/lib/auth/helpers'
import { Button } from '@/components/ui/button'
import { ProductFormModal } from '@/components/admin/ProductFormModal'
import { ProductTable } from '@/components/admin/ProductTable'

type Drawer = { productId?: string } | null

// No dynamic [id] route (static export): editing opens ProductFormModal, which fetches the product by id.
export default function AdminProductsPage() {
  const [products, setProducts] = useState<Product[]>([])
  const [search, setSearch] = useState('')
  const [query, setQuery] = useState('')
  const [drawer, setDrawer] = useState<Drawer>(null)
  const [error, setError] = useState<string | null>(null)
  const [loaded, setLoaded] = useState(false)

  const load = useCallback(async () => {
    try {
      const res = await adminListProducts({ q: query || undefined, limit: 100 })
      if (res.status === 200) setProducts(res.data.data)
      else setError(errorMessage(res, 'Could not load products'))
    } catch {
      setError('Network error. Please try again.')
    } finally {
      setLoaded(true)
    }
  }, [query])

  useEffect(() => {
    void load()
  }, [load])

  async function remove(p: Product) {
    if (!window.confirm(`Delete "${p.name}"? It will be hidden from the store.`)) return
    const res = await adminDeleteProduct(p.id)
    if (res.status === 200) {
      setError(null)
      await load()
    } else setError(errorMessage(res, 'Could not delete the product'))
  }

  return (
    <main className="mx-auto max-w-4xl space-y-6 p-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Products</h1>
        <Button onClick={() => setDrawer({})}>New product</Button>
      </div>
      <form
        role="search"
        className="flex gap-2"
        onSubmit={(e) => {
          e.preventDefault()
          setQuery(search.trim())
        }}
      >
        <input
          aria-label="Search products"
          placeholder="Search by name, slug or SKU"
          className="w-full rounded-md border px-3 py-2"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        <Button type="submit" variant="outline">
          Search
        </Button>
      </form>
      {error && (
        <p role="alert" className="text-sm text-red-600">
          {error}
        </p>
      )}
      {loaded && (
        <ProductTable products={products} onEdit={(p) => setDrawer({ productId: p.id })} onDelete={(p) => void remove(p)} />
      )}
      {drawer && (
        <ProductFormModal
          productId={drawer.productId}
          onClose={() => setDrawer(null)}
          onSaved={() => {
            setDrawer(null)
            setError(null)
            void load()
          }}
        />
      )}
    </main>
  )
}
