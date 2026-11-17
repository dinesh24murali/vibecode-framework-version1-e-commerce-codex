'use client'
import { useState } from 'react'
import Link from 'next/link'
import { Menu, ShoppingCart, X } from 'lucide-react'
import { useStore } from '@/lib/store'
import { SIMPLE_LINKS } from '@/lib/constants/navigation'
import { BookStoreMegaMenu } from './BookStoreMegaMenu'
import { PrintingServicesMegaMenu } from './PrintingServicesMegaMenu'
import { SearchBar } from './SearchBar'

type MenuKey = 'books' | 'printing'

const linkCls = 'inline-flex min-h-[44px] items-center px-3 text-sm font-medium hover:underline'

/** Desktop trigger: opens on hover and on focus/click (keyboard + touch laptops). */
function DesktopMenu({ id, label, open, setOpen, children }: {
  id: MenuKey
  label: string
  open: MenuKey | null
  setOpen: (k: MenuKey | null) => void
  children: React.ReactNode
}) {
  const isOpen = open === id
  return (
    <div
      className="relative"
      onMouseEnter={() => setOpen(id)}
      onMouseLeave={() => setOpen(null)}
      onBlur={(e) => {
        if (!e.currentTarget.contains(e.relatedTarget)) setOpen(null)
      }}
    >
      <button
        type="button"
        aria-expanded={isOpen}
        aria-haspopup="true"
        className={linkCls}
        onClick={() => setOpen(isOpen ? null : id)}
        onKeyDown={(e) => e.key === 'Escape' && setOpen(null)}
      >
        {label}
      </button>
      {isOpen && (
        <div className="absolute left-0 top-full z-50 w-[28rem] rounded border bg-white p-4 shadow-lg">
          {children}
        </div>
      )}
    </div>
  )
}

export function Navbar() {
  const count = useStore((s) => s.items.reduce((n, i) => n + i.quantity, 0))
  const [open, setOpen] = useState<MenuKey | null>(null)
  const [drawer, setDrawer] = useState(false)
  const [section, setSection] = useState<MenuKey | null>(null)
  const closeDrawer = () => {
    setDrawer(false)
    setSection(null)
  }

  return (
    <header className="border-b bg-white">
      <nav aria-label="Main" className="mx-auto flex max-w-6xl items-center gap-2 px-4">
        <Link href="/" className="mr-4 inline-flex min-h-[44px] items-center font-semibold">
          e-commerce-site
        </Link>

        <div className="hidden items-center md:flex">
          <DesktopMenu id="books" label="Book Store" open={open} setOpen={setOpen}>
            <BookStoreMegaMenu onNavigate={() => setOpen(null)} />
          </DesktopMenu>
          <Link href={SIMPLE_LINKS[0].href} className={linkCls}>{SIMPLE_LINKS[0].name}</Link>
          <Link href={SIMPLE_LINKS[1].href} className={linkCls}>{SIMPLE_LINKS[1].name}</Link>
          <DesktopMenu id="printing" label="Printing Services" open={open} setOpen={setOpen}>
            <PrintingServicesMegaMenu onNavigate={() => setOpen(null)} />
          </DesktopMenu>
          <Link href={SIMPLE_LINKS[2].href} className={linkCls}>{SIMPLE_LINKS[2].name}</Link>
        </div>

        <div className="ml-auto hidden w-64 md:block">
          <SearchBar />
        </div>

        <Link
          href="/cart/"
          aria-label={`Cart, ${count} items`}
          className="relative ml-auto inline-flex h-11 w-11 items-center justify-center md:ml-2"
        >
          <ShoppingCart size={22} />
          {count > 0 && (
            <span data-testid="cart-badge" className="absolute right-0 top-0 rounded-full bg-red-600 px-1.5 text-xs text-white">
              {count}
            </span>
          )}
        </Link>

        <button
          type="button"
          aria-label={drawer ? 'Close menu' : 'Open menu'}
          aria-expanded={drawer}
          className="inline-flex h-11 w-11 items-center justify-center md:hidden"
          onClick={() => (drawer ? closeDrawer() : setDrawer(true))}
        >
          {drawer ? <X size={22} /> : <Menu size={22} />}
        </button>
      </nav>

      {drawer && (
        <div className="border-t px-4 py-2 md:hidden" data-testid="mobile-drawer">
          <SearchBar onSubmit={closeDrawer} />
          {([['books', 'Book Store'], ['printing', 'Printing Services']] as const).map(([key, label]) => (
            <div key={key} className="border-b">
              <button
                type="button"
                aria-expanded={section === key}
                className="flex min-h-[44px] w-full items-center justify-between text-left text-sm font-medium"
                onClick={() => setSection(section === key ? null : key)}
              >
                {label}
                <span aria-hidden>{section === key ? '-' : '+'}</span>
              </button>
              {section === key && (
                <div className="pb-2">
                  {key === 'books' ? (
                    <BookStoreMegaMenu onNavigate={closeDrawer} />
                  ) : (
                    <PrintingServicesMegaMenu onNavigate={closeDrawer} />
                  )}
                </div>
              )}
            </div>
          ))}
          {SIMPLE_LINKS.map((l) => (
            <Link key={l.href} href={l.href} onClick={closeDrawer} className="flex min-h-[44px] items-center border-b text-sm font-medium">
              {l.name}
            </Link>
          ))}
        </div>
      )}
    </header>
  )
}
