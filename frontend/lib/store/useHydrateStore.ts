'use client'
import { useEffect } from 'react'
import { useStore } from './index'

/** Rehydrates the persisted store on the client after mount (ADR-0003). */
export function useHydrateStore() {
  useEffect(() => {
    void useStore.persist.rehydrate()
  }, [])
}
