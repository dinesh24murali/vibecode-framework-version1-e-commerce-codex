'use client'
import { useEffect, useState } from 'react'
import { listFilterableAttributes } from '@/lib/api/generated/attributes/attributes'
import { listCategories } from '@/lib/api/generated/categories/categories'
import type { AttributeWithValues, Category } from '@/lib/api/generated/model'
import { useStore } from '@/lib/store'
import { Button } from '@/components/ui/button'

/** Filter checkboxes. Collapses behind a "Filters" toggle below the md breakpoint. */
export function FilterSidebar() {
  const filters = useStore((s) => s.filters)
  const toggleCategory = useStore((s) => s.toggleCategory)
  const toggleValue = useStore((s) => s.toggleValue)
  const clearFilters = useStore((s) => s.clearFilters)
  const [attrs, setAttrs] = useState<AttributeWithValues[]>([])
  const [cats, setCats] = useState<Category[]>([])
  const [open, setOpen] = useState(false)

  useEffect(() => {
    void listFilterableAttributes().then((r) => {
      if (r.status === 200) setAttrs((r.data as { data: AttributeWithValues[] }).data)
    })
    void listCategories().then((r) => {
      if (r.status === 200) setCats((r.data as { data: Category[] }).data)
    })
  }, [])

  const active = filters.categories.length + filters.valueIds.length

  return (
    <aside className="md:w-56 md:shrink-0">
      <Button
        variant="outline"
        className="mb-3 w-full md:hidden"
        aria-expanded={open}
        aria-controls="filter-panel"
        onClick={() => setOpen(!open)}
      >
        Filters{active > 0 ? ` (${active})` : ''}
      </Button>
      <div id="filter-panel" data-testid="filter-panel" className={`${open ? 'block' : 'hidden'} space-y-5 md:block`}>
        <fieldset>
          <legend className="mb-1 text-sm font-semibold">Category</legend>
          {cats.map((c) => (
            <label key={c.id} className="flex items-center gap-2 py-0.5 text-sm">
              <input
                type="checkbox"
                checked={filters.categories.includes(c.slug)}
                onChange={() => toggleCategory(c.slug)}
              />
              {c.name}
            </label>
          ))}
        </fieldset>
        {attrs.map((a) => (
          <fieldset key={a.id}>
            <legend className="mb-1 text-sm font-semibold">{a.name}</legend>
            {a.values.map((v) => (
              <label key={v.id} className="flex items-center gap-2 py-0.5 text-sm">
                <input
                  type="checkbox"
                  checked={filters.valueIds.includes(v.id)}
                  onChange={() => toggleValue(v.id)}
                />
                {v.value}
              </label>
            ))}
          </fieldset>
        ))}
        {active > 0 && (
          <Button variant="outline" size="sm" onClick={clearFilters}>
            Clear filters
          </Button>
        )}
      </div>
    </aside>
  )
}
