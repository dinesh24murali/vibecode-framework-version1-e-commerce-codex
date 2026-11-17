import { getLandingPage } from '@/lib/api/generated/landing/landing'
import type { LandingPageResponse } from '@/lib/api/generated/model'
import { LandingContent } from '@/components/shop/LandingContent'

// Statically generated: the landing payload is fetched at build time. If the API is
// unreachable during the build, the client fetches it on first load instead.
async function loadLanding(): Promise<LandingPageResponse | null> {
  try {
    const res = await getLandingPage()
    return res.status === 200 ? (res.data as LandingPageResponse) : null
  } catch {
    return null
  }
}

export default async function HomePage() {
  const initial = await loadLanding()
  return (
    <main>
      <LandingContent initial={initial} />
    </main>
  )
}
