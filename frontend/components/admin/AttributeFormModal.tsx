'use client'
import { useState } from 'react'
import {
  adminCreateAttributeDefinition,
  adminUpdateAttributeDefinition,
} from '@/lib/api/generated/admin-attributes/admin-attributes'
import type { AttributeDefinition } from '@/lib/api/generated/model'
import { errorMessage } from '@/lib/auth/helpers'
import { Button } from '@/components/ui/button'
import { Modal } from './Modal'

/** Create (no `attribute`) or rename/edit an attribute definition. The code is immutable once created. */
export function AttributeFormModal({
  attribute,
  onClose,
  onSaved,
}: {
  attribute?: AttributeDefinition
  onClose: () => void
  onSaved: () => void
}) {
  const [code, setCode] = useState(attribute?.code ?? '')
  const [name, setName] = useState(attribute?.name ?? '')
  const [filterable, setFilterable] = useState(attribute?.is_filterable ?? true)
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault()
    setBusy(true)
    setError(null)
    try {
      const res = attribute
        ? await adminUpdateAttributeDefinition(attribute.id, { name, is_filterable: filterable })
        : await adminCreateAttributeDefinition({ code, name, is_filterable: filterable })
      if (res.status !== 200 && res.status !== 201) {
        setError(errorMessage(res, 'Could not save the attribute'))
        return
      }
      onSaved()
    } catch {
      setError('Network error. Please try again.')
    } finally {
      setBusy(false)
    }
  }

  return (
    <Modal title={attribute ? 'Edit attribute' : 'New attribute'} onClose={onClose}>
      <form onSubmit={onSubmit} className="space-y-4">
        <div className="space-y-1">
          <label htmlFor="attr-code" className="block text-sm font-medium">
            Code
          </label>
          <input
            id="attr-code"
            required
            disabled={!!attribute}
            value={code}
            onChange={(e) => setCode(e.target.value)}
            placeholder="e.g. level"
            className="w-full rounded-md border px-3 py-2 disabled:bg-gray-100"
          />
        </div>
        <div className="space-y-1">
          <label htmlFor="attr-name" className="block text-sm font-medium">
            Name
          </label>
          <input
            id="attr-name"
            required
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full rounded-md border px-3 py-2"
          />
        </div>
        <label className="flex items-center gap-2 text-sm">
          <input type="checkbox" checked={filterable} onChange={(e) => setFilterable(e.target.checked)} />
          Show as a filter
        </label>
        {error && (
          <p role="alert" className="text-sm text-red-600">
            {error}
          </p>
        )}
        <div className="flex justify-end gap-2">
          <Button type="button" variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" disabled={busy}>
            Save
          </Button>
        </div>
      </form>
    </Modal>
  )
}
