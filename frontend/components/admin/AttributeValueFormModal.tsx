'use client'
import { useState } from 'react'
import {
  adminCreateAttributeValue,
  adminUpdateAttributeValue,
} from '@/lib/api/generated/admin-attributes/admin-attributes'
import type { AttributeDefinitionValue } from '@/lib/api/generated/model'
import { errorMessage } from '@/lib/auth/helpers'
import { Button } from '@/components/ui/button'
import { Modal } from './Modal'

/** Add (no `value`) or rename a value of the attribute `attributeId`. */
export function AttributeValueFormModal({
  attributeId,
  value,
  onClose,
  onSaved,
}: {
  attributeId: string
  value?: AttributeDefinitionValue
  onClose: () => void
  onSaved: () => void
}) {
  const [label, setLabel] = useState(value?.value ?? '')
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault()
    setBusy(true)
    setError(null)
    try {
      const res = value
        ? await adminUpdateAttributeValue(attributeId, value.id, { value: label })
        : await adminCreateAttributeValue(attributeId, { value: label })
      if (res.status !== 200 && res.status !== 201) {
        setError(errorMessage(res, 'Could not save the value'))
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
    <Modal title={value ? 'Edit value' : 'New value'} onClose={onClose}>
      <form onSubmit={onSubmit} className="space-y-4">
        <div className="space-y-1">
          <label htmlFor="attr-value" className="block text-sm font-medium">
            Value
          </label>
          <input
            id="attr-value"
            required
            value={label}
            onChange={(e) => setLabel(e.target.value)}
            className="w-full rounded-md border px-3 py-2"
          />
        </div>
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
