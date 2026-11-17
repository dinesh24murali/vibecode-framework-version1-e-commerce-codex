import type { StateCreator } from 'zustand'

export type SortBy = 'price_asc' | 'price_desc' | 'newest' | 'popular'
const SORTS: SortBy[] = ['price_asc', 'price_desc', 'newest', 'popular']

export interface FilterState {
  categories: string[] // category slugs
  valueIds: string[] // attribute value UUIDs
  search: string
  sort: SortBy
  page: number
}

export interface FilterSlice {
  filters: FilterState
  toggleCategory: (slug: string) => void
  toggleValue: (id: string) => void
  setSearch: (q: string) => void
  setSort: (s: SortBy) => void
  setPage: (p: number) => void
  clearFilters: () => void
  /** Replaces the filters from URL search params (called on mount and on back/forward). */
  hydrateFilters: (params: URLSearchParams) => void
}

export const emptyFilters: FilterState = { categories: [], valueIds: [], search: '', sort: 'newest', page: 1 }

/** URL form: ?category=a&category=b&attr=<uuid>&q=text&sort=price_asc&page=2 (defaults omitted). */
export function filtersToParams(f: FilterState): URLSearchParams {
  const p = new URLSearchParams()
  f.categories.forEach((c) => p.append('category', c))
  f.valueIds.forEach((v) => p.append('attr', v))
  if (f.search) p.set('q', f.search)
  if (f.sort !== 'newest') p.set('sort', f.sort)
  if (f.page > 1) p.set('page', String(f.page))
  return p
}

export function paramsToFilters(p: URLSearchParams): FilterState {
  const sort = p.get('sort') as SortBy | null
  const page = Number(p.get('page'))
  return {
    categories: p.getAll('category').filter(Boolean),
    valueIds: p.getAll('attr').filter(Boolean),
    search: p.get('q') ?? '',
    sort: sort && SORTS.includes(sort) ? sort : 'newest',
    page: Number.isInteger(page) && page > 1 ? page : 1,
  }
}

const toggle = (list: string[], v: string) => (list.includes(v) ? list.filter((x) => x !== v) : [...list, v])

// Any filter change resets pagination to page 1.
export const createFilterSlice: StateCreator<FilterSlice, [], [], FilterSlice> = (set) => ({
  filters: emptyFilters,
  toggleCategory: (slug) =>
    set((s) => ({ filters: { ...s.filters, categories: toggle(s.filters.categories, slug), page: 1 } })),
  toggleValue: (id) =>
    set((s) => ({ filters: { ...s.filters, valueIds: toggle(s.filters.valueIds, id), page: 1 } })),
  setSearch: (q) => set((s) => ({ filters: { ...s.filters, search: q, page: 1 } })),
  setSort: (sort) => set((s) => ({ filters: { ...s.filters, sort, page: 1 } })),
  setPage: (page) => set((s) => ({ filters: { ...s.filters, page } })),
  clearFilters: () => set({ filters: emptyFilters }),
  hydrateFilters: (params) => set({ filters: paramsToFilters(params) }),
})
