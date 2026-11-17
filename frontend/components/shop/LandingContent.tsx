'use client'
import { useEffect, useState } from 'react'
import { getLandingPage } from '@/lib/api/generated/landing/landing'
import type { LandingPageResponse } from '@/lib/api/generated/model'
import { HeroCarousel } from './HeroCarousel'
import { ProductGrid } from './ProductGrid'

/**
 * Renders the landing page. `initial` is the payload fetched at build time; when the
 * API was unreachable during the build it is null and the data is fetched on mount.
 */
export function LandingContent({ initial }: { initial: LandingPageResponse | null }) {
  const [data, setData] = useState<LandingPageResponse | null>(initial)
  const [failed, setFailed] = useState(false)

  useEffect(() => {
    if (initial) return
    getLandingPage()
      .then((res) => {
        if (res.status === 200) setData(res.data as LandingPageResponse)
        else setFailed(true)
      })
      .catch(() => setFailed(true))
  }, [initial])

  if (!data && !failed) return <p className="p-6 text-center text-sm">Loading...</p>

  return (
    <>
      <HeroCarousel slides={data?.carousel_slides ?? []} />
      <ProductGrid products={data?.featured_products ?? []} />
    </>
  )
}
