import Link from 'next/link'
import { BOOK_CATEGORIES } from '@/lib/constants/navigation'

export function BookStoreMegaMenu({ onNavigate }: { onNavigate?: () => void }) {
  return (
    <ul className="grid grid-cols-2 gap-x-6 gap-y-1 md:grid-cols-3" data-testid="bookstore-menu">
      {BOOK_CATEGORIES.map((c) => (
        <li key={c.href}>
          <Link href={c.href} onClick={onNavigate} className="block py-1 text-sm hover:underline">
            {c.name}
          </Link>
        </li>
      ))}
    </ul>
  )
}
