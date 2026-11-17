'use client'
import { useEffect, useState } from 'react'
import { usePathname, useRouter, useSearchParams } from 'next/navigation'
import { useStore } from './index'
import { filtersToParams, type FilterState } from './slices/filterSlice'

export { filtersToParams }

/** Current filters from the store. */
export function useStoreFilters(): FilterState {
  return useStore((s) => s.filters)
}

/**
 * Two-way sync between the filter slice and the URL:
 * - on mount and on back/forward the slice is hydrated from useSearchParams()
 * - every filter change pushes the new URL (no page reload).
 * Returns true once the slice has been hydrated from the URL; render consumers only then.
 */
export function useProductFilterUrlSync(): boolean {
  const router = useRouter()
  const pathname = usePathname()
  const searchParams = useSearchParams()
  const filters = useStore((s) => s.filters)
  const hydrate = useStore((s) => s.hydrateFilters)
  const [ready, setReady] = useState(false)

  const urlKey = searchParams.toString()
  const stateKey = filtersToParams(filters).toString()

  // URL -> store (initial load, back/forward).
  useEffect(() => {
    if (!ready || urlKey !== stateKey) hydrate(new URLSearchParams(urlKey))
    setReady(true)
    // only react to URL changes
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [urlKey])

  // store -> URL.
  useEffect(() => {
    if (!ready || urlKey === stateKey) return
    router.push(stateKey ? `${pathname}?${stateKey}` : pathname)
    // only react to store changes
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [stateKey, ready])

  return ready
}
