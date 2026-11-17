'use client'
import { useProductFilterUrlSync } from '@/lib/store/useProductFilters'
import { FilterSidebar } from './FilterSidebar'
import { ProductListGrid } from './ProductListGrid'

export function ProductListPage() {
  const ready = useProductFilterUrlSync()
  if (!ready) return null
  return (
    <div className="flex flex-col gap-4 md:flex-row">
      <FilterSidebar />
      <ProductListGrid />
    </div>
  )
}
