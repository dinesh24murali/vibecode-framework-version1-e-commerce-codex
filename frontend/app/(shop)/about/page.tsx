import type { Metadata } from 'next'
import { PlaceholderImage, Section, StaticPage } from '@/components/shop/StaticPage'

export const metadata: Metadata = {
  title: 'About Us | e-commerce-site',
  description: 'Learn about our story, what we make and the values behind our books and printing services.',
}

export default function AboutPage() {
  return (
    <StaticPage title="About Us">
      <div className="grid items-center gap-6 md:grid-cols-2">
        <PlaceholderImage label="Our workshop" />
        <Section heading="Our Story">
          We are a printing house that grew into a bookstore. Colouring books, journals, planners and more are designed
          and produced in-house. This text is placeholder content.
        </Section>
      </div>
      <Section heading="What We Do">
        From activity books to packaging and corporate printing, we handle design, printing, finishing and binding under one roof.
      </Section>
      <Section heading="Our Values">
        Quality materials, honest pricing and dependable delivery.
      </Section>
    </StaticPage>
  )
}
