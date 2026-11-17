'use client'
import { useCallback, useEffect, useState } from 'react'
import {
  adminDeleteAttributeDefinition,
  adminDeleteAttributeValue,
  adminListAttributeDefinitions,
  adminListAttributeValues,
} from '@/lib/api/generated/admin-attributes/admin-attributes'
import type { AttributeDefinition, AttributeDefinitionValue } from '@/lib/api/generated/model'
import { errorMessage } from '@/lib/auth/helpers'
import { Button } from '@/components/ui/button'
import { AttributeFormModal } from '@/components/admin/AttributeFormModal'
import { AttributeTable } from '@/components/admin/AttributeTable'
import { AttributeValueFormModal } from '@/components/admin/AttributeValueFormModal'

type ModalState =
  | { kind: 'attr'; attribute?: AttributeDefinition }
  | { kind: 'value'; attribute: AttributeDefinition; value?: AttributeDefinitionValue }
  | null

// All CRUD happens in modals on this one page (no dynamic [id] route: static export).
export default function AdminAttributesPage() {
  const [attributes, setAttributes] = useState<AttributeDefinition[]>([])
  const [values, setValues] = useState<Record<string, AttributeDefinitionValue[]>>({})
  const [modal, setModal] = useState<ModalState>(null)
  const [error, setError] = useState<string | null>(null)
  const [loaded, setLoaded] = useState(false)

  const loadAttributes = useCallback(async () => {
    try {
      const res = await adminListAttributeDefinitions({ limit: 100 })
      if (res.status === 200) setAttributes(res.data.data)
      else setError(errorMessage(res, 'Could not load attributes'))
    } catch {
      setError('Network error. Please try again.')
    } finally {
      setLoaded(true)
    }
  }, [])

  const loadValues = useCallback(async (id: string) => {
    const res = await adminListAttributeValues(id)
    if (res.status === 200) setValues((v) => ({ ...v, [id]: res.data.data }))
  }, [])

  useEffect(() => {
    void loadAttributes()
  }, [loadAttributes])

  async function removeAttribute(a: AttributeDefinition) {
    if (!window.confirm(`Delete "${a.name}" and all its values?`)) return
    const res = await adminDeleteAttributeDefinition(a.id)
    if (res.status === 204) {
      setError(null)
      await loadAttributes()
    } else setError(errorMessage(res, 'Could not delete the attribute'))
  }

  async function removeValue(a: AttributeDefinition, v: AttributeDefinitionValue) {
    if (!window.confirm(`Remove value "${v.value}"?`)) return
    const res = await adminDeleteAttributeValue(a.id, v.id)
    if (res.status === 204) {
      setError(null)
      await loadValues(a.id)
    } else if (res.status === 409) {
      setError(`"${v.value}" is assigned to products and cannot be removed.`)
    } else setError(errorMessage(res, 'Could not remove the value'))
  }

  return (
    <main className="mx-auto max-w-3xl space-y-6 p-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Attributes</h1>
        <Button onClick={() => setModal({ kind: 'attr' })}>New attribute</Button>
      </div>
      {error && (
        <p role="alert" className="text-sm text-red-600">
          {error}
        </p>
      )}
      {loaded && (
        <AttributeTable
          attributes={attributes}
          valuesByAttr={values}
          onLoadValues={(id) => void loadValues(id)}
          onEdit={(a) => setModal({ kind: 'attr', attribute: a })}
          onDelete={(a) => void removeAttribute(a)}
          onAddValue={(a) => setModal({ kind: 'value', attribute: a })}
          onEditValue={(a, v) => setModal({ kind: 'value', attribute: a, value: v })}
          onDeleteValue={(a, v) => void removeValue(a, v)}
        />
      )}
      {modal?.kind === 'attr' && (
        <AttributeFormModal
          attribute={modal.attribute}
          onClose={() => setModal(null)}
          onSaved={() => {
            setModal(null)
            setError(null)
            void loadAttributes()
          }}
        />
      )}
      {modal?.kind === 'value' && (
        <AttributeValueFormModal
          attributeId={modal.attribute.id}
          value={modal.value}
          onClose={() => setModal(null)}
          onSaved={() => {
            const id = modal.attribute.id
            setModal(null)
            setError(null)
            void loadValues(id)
          }}
        />
      )}
    </main>
  )
}
