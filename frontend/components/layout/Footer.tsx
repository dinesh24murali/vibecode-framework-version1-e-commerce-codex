import Link from 'next/link'
import { Facebook, Instagram, MessageCircle } from 'lucide-react'
import { SOCIAL_LINKS } from '@/lib/constants/social'

const INTERNAL_LINKS = [
  { name: 'Terms & Conditions', href: '/terms/' },
  { name: 'About Us', href: '/about/' },
  { name: 'Wall of Fame', href: '/wall-of-fame/' },
]

const ICONS = { Facebook, Instagram, WhatsApp: MessageCircle }

export function Footer() {
  return (
    <footer className="mt-12 border-t bg-white" data-testid="footer">
      <div className="mx-auto flex max-w-6xl flex-col gap-4 px-4 py-6 md:flex-row md:items-center md:justify-between">
        <nav aria-label="Footer">
          <ul className="flex flex-wrap gap-x-6">
            {INTERNAL_LINKS.map((l) => (
              <li key={l.href}>
                <Link href={l.href} className="inline-flex min-h-[44px] items-center text-sm hover:underline">
                  {l.name}
                </Link>
              </li>
            ))}
          </ul>
        </nav>
        <ul className="flex gap-2" aria-label="Social media">
          {SOCIAL_LINKS.map((s) => {
            const Icon = ICONS[s.name]
            return (
              <li key={s.name}>
                <a
                  href={s.href}
                  target="_blank"
                  rel="noopener noreferrer"
                  aria-label={s.name}
                  className="inline-flex h-11 w-11 items-center justify-center rounded hover:bg-muted"
                >
                  <Icon size={20} />
                </a>
              </li>
            )
          })}
        </ul>
      </div>
    </footer>
  )
}
