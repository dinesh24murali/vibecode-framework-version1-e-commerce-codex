'use client'
import { useCallback, useEffect, useState } from 'react'
import {
  createAddress,
  deleteAddress,
  listAddresses,
  updateAddress,
} from '@/lib/api/generated/account/account'
import type { Address, CreateAddressRequest } from '@/lib/api/generated/model'
import { errorMessage } from '@/lib/auth/helpers'
import { Button } from '@/components/ui/button'

const empty: CreateAddressRequest = {
  address_line1: '',
  address_line2: '',
  landmark: '',
  city: '',
  state: '',
  pincode: '',
  company_name: '',
  gst_number: '',
  phone: '',
  is_default: false,
}

const fields: { key: keyof CreateAddressRequest; label: string; required?: boolean }[] = [
  { key: 'address_line1', label: 'Address line 1', required: true },
  { key: 'address_line2', label: 'Address line 2' },
  { key: 'landmark', label: 'Landmark' },
  { key: 'city', label: 'City', required: true },
  { key: 'state', label: 'State', required: true },
  { key: 'pincode', label: 'Pincode', required: true },
  { key: 'company_name', label: 'Company name' },
  { key: 'gst_number', label: 'GST number' },
  { key: 'phone', label: 'Phone', required: true },
]

// Empty optional strings are sent as null.
const clean = (f: CreateAddressRequest): CreateAddressRequest => ({
  ...f,
  address_line2: f.address_line2 || null,
  landmark: f.landmark || null,
  company_name: f.company_name || null,
  gst_number: f.gst_number || null,
})

export function AddressTab() {
  const [rows, setRows] = useState<Address[]>([])
  const [form, setForm] = useState<CreateAddressRequest | null>(null)
  const [editId, setEditId] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)

  const load = useCallback(async () => {
    const res = await listAddresses()
    if (res.status === 200) setRows((res.data as { data: Address[] }).data)
    else setError(errorMessage(res))
  }, [])

  useEffect(() => {
    void load()
  }, [load])

  async function save(e: React.FormEvent) {
    e.preventDefault()
    if (!form) return
    const body = clean(form)
    const res = editId ? await updateAddress(editId, body) : await createAddress(body)
    if (res.status === 200 || res.status === 201) {
      setForm(null)
      setEditId(null)
      setError(null)
      await load()
    } else setError(errorMessage(res))
  }

  async function remove(id: string) {
    const res = await deleteAddress(id)
    if (res.status === 200) await load()
    else setError(errorMessage(res))
  }

  function edit(a: Address) {
    setEditId(a.id)
    setForm({
      address_line1: a.address_line1,
      address_line2: a.address_line2 ?? '',
      landmark: a.landmark ?? '',
      city: a.city,
      state: a.state,
      pincode: a.pincode,
      company_name: a.company_name ?? '',
      gst_number: a.gst_number ?? '',
      phone: a.phone,
      is_default: a.is_default,
    })
  }

  return (
    <div className="space-y-6">
      {error && (
        <p role="alert" className="text-sm text-red-600">
          {error}
        </p>
      )}
      {!form && (
        <Button
          onClick={() => {
            setEditId(null)
            setForm({ ...empty })
          }}
        >
          Add address
        </Button>
      )}
      {form && (
        <form onSubmit={save} className="space-y-3">
          {fields.map((f) => (
            <div key={f.key} className="space-y-1">
              <label htmlFor={`addr-${f.key}`} className="block text-sm font-medium">
                {f.label}
              </label>
              <input
                id={`addr-${f.key}`}
                required={f.required}
                value={(form[f.key] as string | null | undefined) ?? ''}
                onChange={(e) => setForm({ ...form, [f.key]: e.target.value })}
                className="w-full rounded-md border px-3 py-2"
              />
            </div>
          ))}
          <label className="flex items-center gap-2 text-sm">
            <input
              type="checkbox"
              checked={form.is_default}
              onChange={(e) => setForm({ ...form, is_default: e.target.checked })}
            />
            Default address
          </label>
          <div className="flex gap-2">
            <Button type="submit">Save address</Button>
            <Button type="button" variant="outline" onClick={() => setForm(null)}>
              Cancel
            </Button>
          </div>
        </form>
      )}
      <ul className="space-y-3">
        {rows.map((a) => (
          <li key={a.id} className="rounded-md border p-3 text-sm" data-testid="address-row">
            <p>
              {a.address_line1}
              {a.address_line2 ? `, ${a.address_line2}` : ''}
              {a.landmark ? `, ${a.landmark}` : ''}
            </p>
            <p>
              {a.city}, {a.state} {a.pincode} · {a.phone}
            </p>
            {a.company_name && <p>{a.company_name}</p>}
            {a.gst_number && <p>GST: {a.gst_number}</p>}
            {a.is_default && <p className="font-medium">Default</p>}
            <div className="mt-2 flex gap-2">
              <Button size="sm" variant="outline" onClick={() => edit(a)}>
                Edit
              </Button>
              <Button size="sm" variant="outline" onClick={() => void remove(a.id)}>
                Delete
              </Button>
            </div>
          </li>
        ))}
        {rows.length === 0 && <li className="text-sm">No saved addresses.</li>}
      </ul>
    </div>
  )
}
