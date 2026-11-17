import type { Metadata } from 'next'
import { getProduct, listProducts } from '@/lib/api/generated/products/products'
import type { Product } from '@/lib/api/generated/model'
import { ProductDetail } from '@/components/shop/ProductDetail'

// Used only when the API is unreachable/empty at build time so the static export still builds.
const FALLBACK_SLUG = 'sample-product'

async function listAllSlugs(): Promise<{ slug: string; id: string }[]> {
  const out: { slug: string; id: string }[] = []
  try {
    for (let page = 1; page < 500; page++) {
      const res = await listProducts({ page, limit: 100 })
      if (res.status !== 200) break
      out.push(...res.data.data.map((p) => ({ slug: p.slug, id: p.id })))
      if (page >= res.data.meta.total_pages) break
    }
  } catch {
    // API unreachable during build
  }
  return out
}

async function loadProduct(slug: string): Promise<Product | null> {
  const hit = (await listAllSlugs()).find((p) => p.slug === slug)
  if (!hit) return null
  try {
    const res = await getProduct(hit.id)
    return res.status === 200 ? (res.data as Product) : null
  } catch {
    return null
  }
}

export async function generateStaticParams() {
  const slugs = await listAllSlugs()
  return slugs.length > 0 ? slugs.map((p) => ({ slug: p.slug })) : [{ slug: FALLBACK_SLUG }]
}

export async function generateMetadata({ params }: { params: Promise<{ slug: string }> }): Promise<Metadata> {
  const { slug } = await params
  const product = await loadProduct(slug)
  if (!product) return {}
  const description = product.description.replace(/<[^>]*>/g, '').slice(0, 160)
  return { title: product.name, description }
}

export default async function ProductPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params
  const initial = await loadProduct(slug)
  return (
    <main className="mx-auto max-w-6xl p-4">
      <ProductDetail slug={slug} initial={initial} />
    </main>
  )
}
