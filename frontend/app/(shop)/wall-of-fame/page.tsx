import type { Metadata } from 'next'
import { PlaceholderImage, StaticPage } from '@/components/shop/StaticPage'

export const metadata: Metadata = {
  title: 'Wall of Fame | e-commerce-site',
  description: 'Celebrating the artists, customers and partners who make our community special.',
}

const ENTRIES = [
  { name: 'Featured Artist of the Month', text: 'Beautiful mandala colouring work shared by our community.' },
  { name: 'Customer Spotlight', text: 'A journal collection that turned into a daily creative habit.' },
  { name: 'Partner Showcase', text: 'A custom printing project delivered from our in-house press.' },
]

export default function WallOfFamePage() {
  return (
    <StaticPage title="Wall of Fame">
      <p>A showcase of the creativity and achievements of our community. Content below is placeholder and will be replaced.</p>
      <div className="grid gap-6 md:grid-cols-3">
        {ENTRIES.map((e) => (
          <article key={e.name}>
            <PlaceholderImage label={e.name} />
            <h2 className="mt-2 font-semibold">{e.name}</h2>
            <p>{e.text}</p>
          </article>
        ))}
      </div>
    </StaticPage>
  )
}
