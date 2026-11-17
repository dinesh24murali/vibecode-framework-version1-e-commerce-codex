'use client'
import Link from 'next/link'
import type { Address } from '@/lib/api/generated/model'

interface Props {
  addresses: Address[]
  selectedId: string | null
  onChange: (id: string) => void
}

export function AddressSelector({ addresses, selectedId, onChange }: Props) {
  if (addresses.length === 0) {
    return (
      <p data-testid="no-addresses">
        You have no saved addresses. <Link href="/account/" className="underline">Add one in your account</Link> to continue.
      </p>
    )
  }
  return (
    <fieldset className="space-y-2">
      <legend className="mb-2 text-lg font-semibold">Delivery address</legend>
      {addresses.map((a) => (
        <label key={a.id} className="flex cursor-pointer gap-3 rounded border p-3 text-sm">
          <input type="radio" name="address" checked={selectedId === a.id} onChange={() => onChange(a.id)} />
          <span>
            {a.address_line1}
            {a.address_line2 ? `, ${a.address_line2}` : ''}, {a.city}, {a.state} {a.pincode}
            {a.is_default && <em className="ml-2 text-xs text-gray-500">(default)</em>}
          </span>
        </label>
      ))}
    </fieldset>
  )
}
