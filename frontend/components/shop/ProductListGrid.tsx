'use client'
import { useEffect, useState } from 'react'
import { listProducts } from '@/lib/api/generated/products/products'
import type { ProductListResponse } from '@/lib/api/generated/model'
import { filtersToParams, useStoreFilters } from '@/lib/store/useProductFilters'
import { useStore } from '@/lib/store'
import { Button } from '@/components/ui/button'
import { ProductGrid } from './ProductGrid'

const PAGE_SIZE = 12

/** Fetches and renders the filtered product grid; refetches on every filter change without a page reload. */
export function ProductListGrid() {
  const filters = useStoreFilters()
  const setPage = useStore((s) => s.setPage)
  const setSearch = useStore((s) => s.setSearch)
  const setSort = useStore((s) => s.setSort)
  const [data, setData] = useState<ProductListResponse | null>(null)
  const [error, setError] = useState(false)
  const [searchText, setSearchText] = useState(filters.search)

  const key = filtersToParams(filters).toString()

  useEffect(() => {
    let cancelled = false
    setError(false)
    listProducts({
      q: filters.search || undefined,
      'category[]': filters.categories.length ? filters.categories : undefined,
      'attr_value_id[]': filters.valueIds.length ? filters.valueIds : undefined,
      sort_by: filters.sort,
      page: filters.page,
      limit: PAGE_SIZE,
    } as Parameters<typeof listProducts>[0])
      .then((res) => {
        if (cancelled) return
        if (res.status === 200) setData(res.data as ProductListResponse)
        else setError(true)
      })
      .catch(() => !cancelled && setError(true))
    return () => {
      cancelled = true
    }
    // key captures every filter field
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [key])

  useEffect(() => setSearchText(filters.search), [filters.search])

  return (
    <div className="flex-1 space-y-4">
      <div className="flex flex-wrap items-center gap-2 px-4">
        <form
          className="flex gap-2"
          onSubmit={(e) => {
            e.preventDefault()
            setSearch(searchText.trim())
          }}
        >
          <input
            aria-label="Search products"
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
            className="rounded-md border px-3 py-1.5 text-sm"
          />
          <Button type="submit" size="sm">
            Search
          </Button>
        </form>
        <select
          aria-label="Sort by"
          value={filters.sort}
          onChange={(e) => setSort(e.target.value as typeof filters.sort)}
          className="rounded-md border px-2 py-1.5 text-sm"
        >
          <option value="newest">Newest</option>
          <option value="price_asc">Price: low to high</option>
          <option value="price_desc">Price: high to low</option>
          <option value="popular">Popular</option>
        </select>
      </div>

      {error ? (
        <p role="alert" className="p-6 text-center text-sm">
          Products are temporarily unavailable.
        </p>
      ) : data && data.data.length === 0 ? (
        <p className="p-6 text-center text-sm">No products match your filters.</p>
      ) : (
        <ProductGrid products={data?.data ?? []} />
      )}

      {data && data.meta.total_pages > 1 && (
        <div className="flex items-center justify-center gap-3 text-sm">
          <Button size="sm" variant="outline" disabled={filters.page <= 1} onClick={() => setPage(filters.page - 1)}>
            Previous
          </Button>
          <span>
            Page {data.meta.page} of {data.meta.total_pages}
          </span>
          <Button
            size="sm"
            variant="outline"
            disabled={filters.page >= data.meta.total_pages}
            onClick={() => setPage(filters.page + 1)}
          >
            Next
          </Button>
        </div>
      )}
    </div>
  )
}
