import type { Metadata } from 'next'
import { ContactForm } from '@/components/shop/ContactForm'
import { BUSINESS_ADDRESS, MAP_EMBED_URL } from '@/lib/constants/business'
import { SOCIAL_LINKS } from '@/lib/constants/social'

export const metadata: Metadata = {
  title: 'Contact Us | e-commerce-site',
  description: 'Get in touch with us: find our address on the map, message us on social media or send an enquiry.',
}

export default function ContactPage() {
  return (
    <main className="mx-auto max-w-6xl p-4">
      <h1 className="mb-4 text-3xl font-semibold">Contact Us</h1>
      <div className="grid gap-8 md:grid-cols-2">
        <div className="space-y-4">
          <address className="text-sm not-italic" data-testid="business-address">
            <h2 className="mb-1 text-lg font-semibold">Visit us</h2>
            {BUSINESS_ADDRESS}
          </address>
          <iframe
            title="Map showing our location"
            src={MAP_EMBED_URL}
            className="h-64 w-full rounded border"
            loading="lazy"
            referrerPolicy="no-referrer-when-downgrade"
          />
          <ul className="flex flex-wrap gap-4 text-sm" aria-label="Social media">
            {SOCIAL_LINKS.map((s) => (
              <li key={s.name}>
                <a href={s.href} target="_blank" rel="noopener noreferrer" className="inline-flex min-h-[44px] items-center underline">
                  {s.name}
                </a>
              </li>
            ))}
          </ul>
        </div>
        <div>
          <h2 className="mb-2 text-lg font-semibold">Send us a message</h2>
          <ContactForm />
        </div>
      </div>
    </main>
  )
}
