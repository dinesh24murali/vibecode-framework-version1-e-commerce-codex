'use client'
import { Fragment, useState } from 'react'
import type { AttributeDefinition, AttributeDefinitionValue } from '@/lib/api/generated/model'
import { Button } from '@/components/ui/button'

/** Table of attribute definitions; an expanded row shows its values. */
export function AttributeTable({
  attributes,
  valuesByAttr,
  onLoadValues,
  onEdit,
  onDelete,
  onAddValue,
  onEditValue,
  onDeleteValue,
}: {
  attributes: AttributeDefinition[]
  valuesByAttr: Record<string, AttributeDefinitionValue[]>
  onLoadValues: (attributeId: string) => void
  onEdit: (a: AttributeDefinition) => void
  onDelete: (a: AttributeDefinition) => void
  onAddValue: (a: AttributeDefinition) => void
  onEditValue: (a: AttributeDefinition, v: AttributeDefinitionValue) => void
  onDeleteValue: (a: AttributeDefinition, v: AttributeDefinitionValue) => void
}) {
  const [open, setOpen] = useState<string | null>(null)

  if (attributes.length === 0) return <p className="text-sm text-gray-600">No attributes yet.</p>

  return (
    <table className="w-full text-left text-sm">
      <thead>
        <tr className="border-b">
          <th className="py-2">Name</th>
          <th>Code</th>
          <th>Filterable</th>
          <th className="text-right">Actions</th>
        </tr>
      </thead>
      <tbody>
        {attributes.map((a) => {
          const expanded = open === a.id
          const vals = valuesByAttr[a.id] ?? []
          return (
            <Fragment key={a.id}>
              <tr className="border-b">
                <td className="py-2">{a.name}</td>
                <td>{a.code}</td>
                <td>{a.is_filterable ? 'Yes' : 'No'}</td>
                <td className="space-x-2 text-right">
                  <button
                    className="underline"
                    aria-expanded={expanded}
                    onClick={() => {
                      setOpen(expanded ? null : a.id)
                      if (!expanded) onLoadValues(a.id)
                    }}
                  >
                    {expanded ? 'Hide values' : 'Values'}
                  </button>
                  <button className="underline" onClick={() => onEdit(a)}>
                    Edit
                  </button>
                  <button className="text-red-600 underline" onClick={() => onDelete(a)}>
                    Delete
                  </button>
                </td>
              </tr>
              {expanded && (
                <tr className="border-b bg-gray-50">
                  <td colSpan={4} className="space-y-2 p-3">
                    {vals.length === 0 && <p className="text-gray-600">No values yet.</p>}
                    <ul className="space-y-1">
                      {vals.map((v) => (
                        <li key={v.id} className="flex items-center justify-between">
                          <span>{v.value}</span>
                          <span className="space-x-2">
                            <button className="underline" onClick={() => onEditValue(a, v)}>
                              Rename
                            </button>
                            <button className="text-red-600 underline" onClick={() => onDeleteValue(a, v)}>
                              Remove
                            </button>
                          </span>
                        </li>
                      ))}
                    </ul>
                    <Button size="sm" variant="outline" onClick={() => onAddValue(a)}>
                      Add value
                    </Button>
                  </td>
                </tr>
              )}
            </Fragment>
          )
        })}
      </tbody>
    </table>
  )
}
