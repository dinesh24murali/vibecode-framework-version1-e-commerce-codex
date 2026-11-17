'use client'
import { useEffect, useState } from 'react'
import {
  adminCreateProduct,
  adminGetProduct,
  adminUpdateProduct,
} from '@/lib/api/generated/admin-products/admin-products'
import { listCategories } from '@/lib/api/generated/categories/categories'
import type { Category } from '@/lib/api/generated/model'
import { errorMessage } from '@/lib/auth/helpers'
import { Button } from '@/components/ui/button'
import { ProductImageUploader, type UploadedImage } from './ProductImageUploader'

interface FormState {
  name: string
  slug: string
  description: string
  price: string
  discount: string
  stock: string
  categoryId: string
  isActive: boolean
  images: UploadedImage[]
}

const empty: FormState = {
  name: '',
  slug: '',
  description: '',
  price: '',
  discount: '0',
  stock: '0',
  categoryId: '',
  isActive: true,
  images: [],
}

/**
 * Right-hand drawer for creating (no `productId`) or editing a product. When editing it
 * fetches GET /admin/products/{id} on open, so the form never trusts stale list data.
 */
export function ProductFormModal({
  productId,
  onClose,
  onSaved,
}: {
  productId?: string
  onClose: () => void
  onSaved: () => void
}) {
  const [form, setForm] = useState<FormState>(empty)
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(!!productId)
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  useEffect(() => {
    let cancelled = false
    async function load() {
      try {
        const cats = await listCategories()
        if (!cancelled && cats.status === 200) setCategories(cats.data.data)
        if (productId) {
          const res = await adminGetProduct(productId)
          if (cancelled) return
          if (res.status !== 200) {
            setError(errorMessage(res, 'Could not load the product'))
            return
          }
          const p = res.data
          setForm({
            name: p.name,
            slug: p.slug,
            description: p.description,
            price: String(p.price),
            discount: String(p.discount_percent),
            stock: String(p.stock),
            categoryId: p.categories[0]?.id ?? '',
            isActive: p.is_active,
            images: p.images.map((i) => ({ url: i.url, alt_text: i.alt_text ?? null })),
          })
        }
      } catch {
        if (!cancelled) setError('Network error. Please try again.')
      } finally {
        if (!cancelled) setLoading(false)
      }
    }
    void load()
    return () => {
      cancelled = true
    }
  }, [productId])

  const set = <K extends keyof FormState>(k: K, v: FormState[K]) => setForm((f) => ({ ...f, [k]: v }))

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault()
    setBusy(true)
    setError(null)
    const payload = {
      name: form.name,
      slug: form.slug,
      description: form.description,
      price: Number(form.price),
      discount_percent: Number(form.discount),
      stock: Number(form.stock),
      is_active: form.isActive,
      category_ids: form.categoryId ? [form.categoryId] : [],
      image_urls: form.images.map((img, i) => ({ url: img.url, alt_text: img.alt_text, sort_order: i + 1 })),
    }
    try {
      const res = productId ? await adminUpdateProduct(productId, payload) : await adminCreateProduct(payload)
      if (res.status !== 200 && res.status !== 201) {
        setError(res.status === 409 ? 'That slug is already in use' : errorMessage(res, 'Could not save the product'))
        return
      }
      onSaved()
    } catch {
      setError('Network error. Please try again.')
    } finally {
      setBusy(false)
    }
  }

  const input = 'w-full rounded-md border px-3 py-2'
  return (
    <div className="fixed inset-0 z-50 flex justify-end bg-black/40">
      <aside
        role="dialog"
        aria-modal="true"
        aria-label={productId ? 'Edit product' : 'New product'}
        className="h-full w-full max-w-lg overflow-y-auto bg-white p-6 shadow-lg"
      >
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-lg font-semibold">{productId ? 'Edit product' : 'New product'}</h2>
          <button type="button" onClick={onClose} aria-label="Close" className="text-xl leading-none">
            ×
          </button>
        </div>
        {loading ? (
          <p className="text-sm">Loading…</p>
        ) : (
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="space-y-1">
              <label htmlFor="p-name" className="block text-sm font-medium">
                Name
              </label>
              <input id="p-name" required className={input} value={form.name} onChange={(e) => set('name', e.target.value)} />
            </div>
            <div className="space-y-1">
              <label htmlFor="p-slug" className="block text-sm font-medium">
                Slug
              </label>
              <input id="p-slug" required className={input} value={form.slug} onChange={(e) => set('slug', e.target.value)} />
            </div>
            <div className="space-y-1">
              <label htmlFor="p-desc" className="block text-sm font-medium">
                Description
              </label>
              <textarea id="p-desc" className={input} rows={3} value={form.description} onChange={(e) => set('description', e.target.value)} />
            </div>
            <div className="grid grid-cols-3 gap-2">
              <div className="space-y-1">
                <label htmlFor="p-price" className="block text-sm font-medium">
                  Price (INR)
                </label>
                <input id="p-price" required type="number" min="0" step="0.01" className={input} value={form.price} onChange={(e) => set('price', e.target.value)} />
              </div>
              <div className="space-y-1">
                <label htmlFor="p-discount" className="block text-sm font-medium">
                  Discount %
                </label>
                <input id="p-discount" type="number" min="0" max="100" className={input} value={form.discount} onChange={(e) => set('discount', e.target.value)} />
              </div>
              <div className="space-y-1">
                <label htmlFor="p-stock" className="block text-sm font-medium">
                  Stock
                </label>
                <input id="p-stock" type="number" min="0" className={input} value={form.stock} onChange={(e) => set('stock', e.target.value)} />
              </div>
            </div>
            <div className="space-y-1">
              <label htmlFor="p-cat" className="block text-sm font-medium">
                Category
              </label>
              <select id="p-cat" required className={input} value={form.categoryId} onChange={(e) => set('categoryId', e.target.value)}>
                <option value="">Select a category</option>
                {categories.map((c) => (
                  <option key={c.id} value={c.id}>
                    {c.name}
                  </option>
                ))}
              </select>
            </div>
            <label className="flex items-center gap-2 text-sm">
              <input type="checkbox" checked={form.isActive} onChange={(e) => set('isActive', e.target.checked)} />
              Visible in the store
            </label>
            <ProductImageUploader images={form.images} onChange={(images) => set('images', images)} />
            {error && (
              <p role="alert" className="text-sm text-red-600">
                {error}
              </p>
            )}
            <div className="flex justify-end gap-2">
              <Button type="button" variant="outline" onClick={onClose}>
                Cancel
              </Button>
              <Button type="submit" disabled={busy}>
                Save product
              </Button>
            </div>
          </form>
        )}
        {loading === false && error && productId && form.name === '' && (
          <p className="mt-2 text-sm">Close and try again.</p>
        )}
      </aside>
    </div>
  )
}
