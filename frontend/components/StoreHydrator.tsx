'use client'
import { useHydrateStore } from '@/lib/store/useHydrateStore'

export function StoreHydrator() {
  useHydrateStore()
  return null
}
