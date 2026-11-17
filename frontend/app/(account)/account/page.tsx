'use client'
import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useStore } from '@/lib/store'
import { ProfileTab } from '@/components/account/ProfileTab'
import { AddressTab } from '@/components/account/AddressTab'
import { OrdersTab } from '@/components/account/OrdersTab'

const tabs = [
  { id: 'profile', label: 'Profile' },
  { id: 'addresses', label: 'Addresses' },
  { id: 'orders', label: 'Orders' },
] as const

export default function AccountPage() {
  const router = useRouter()
  const accessToken = useStore((s) => s.accessToken)
  const [tab, setTab] = useState<(typeof tabs)[number]['id']>('profile')

  // ?tab=orders deep-links to a tab (read after mount: the page is statically exported).
  useEffect(() => {
    const wanted = new URLSearchParams(window.location.search).get('tab')
    const match = tabs.find((t) => t.id === wanted)
    if (match) setTab(match.id)
  }, [])

  // Access tokens live in memory only (memory/auth.md): no token means not signed in.
  useEffect(() => {
    if (!accessToken) router.replace('/login/?redirect=/account/')
  }, [accessToken, router])

  if (!accessToken) return null

  return (
    <main className="mx-auto max-w-2xl space-y-6 p-6">
      <h1 className="text-2xl font-semibold">My account</h1>
      <div role="tablist" className="flex gap-2 border-b">
        {tabs.map((t) => (
          <button
            key={t.id}
            role="tab"
            aria-selected={tab === t.id}
            onClick={() => setTab(t.id)}
            className={`px-3 py-2 text-sm ${tab === t.id ? 'border-b-2 border-black font-medium' : ''}`}
          >
            {t.label}
          </button>
        ))}
      </div>
      <div role="tabpanel">
        {tab === 'profile' && <ProfileTab />}
        {tab === 'addresses' && <AddressTab />}
        {tab === 'orders' && <OrdersTab />}
      </div>
    </main>
  )
}
