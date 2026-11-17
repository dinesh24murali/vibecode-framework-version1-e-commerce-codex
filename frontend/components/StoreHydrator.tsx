'use client'
import { useEffect } from 'react'
import { useStore } from '@/lib/store'
import { useHydrateStore } from '@/lib/store/useHydrateStore'
import { cartApi } from '@/lib/api/cartApi'

export function StoreHydrator() {
  useHydrateStore()
  useEffect(() => {
    useStore.getState().setCartApi(cartApi)
  }, [])
  return null
}
