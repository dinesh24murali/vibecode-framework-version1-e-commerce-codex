import Link from 'next/link'
import { PRINTING_PRODUCTS, PRINTING_SERVICES, type NavItem } from '@/lib/constants/navigation'

function Column({ title, items, onNavigate }: { title: string; items: NavItem[]; onNavigate?: () => void }) {
  return (
    <div>
      <h3 className="mb-1 text-sm font-semibold">{title}</h3>
      <ul>
        {items.map((i) => (
          <li key={i.href}>
            <Link href={i.href} onClick={onNavigate} className="block py-1 text-sm hover:underline">
              {i.name}
            </Link>
          </li>
        ))}
      </ul>
    </div>
  )
}

export function PrintingServicesMegaMenu({ onNavigate }: { onNavigate?: () => void }) {
  return (
    <div className="grid gap-6 md:grid-cols-2" data-testid="printing-menu">
      <Column title="Products Offered" items={PRINTING_PRODUCTS} onNavigate={onNavigate} />
      <Column title="Inhouse Services" items={PRINTING_SERVICES} onNavigate={onNavigate} />
    </div>
  )
}
