'use client'
import { useState } from 'react'
import { adminCreateCategory, adminUpdateCategory } from '@/lib/api/generated/admin-categories/admin-categories'
import type { Category } from '@/lib/api/generated/model'
import { errorMessage } from '@/lib/auth/helpers'
import { Button } from '@/components/ui/button'
import { Modal } from './Modal'

const slugify = (s: string) =>
  s
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')

/** Create (no `category`) or edit a category. */
export function CategoryFormModal({
  category,
  onClose,
  onSaved,
}: {
  category?: Category
  onClose: () => void
  onSaved: () => void
}) {
  const [name, setName] = useState(category?.name ?? '')
  const [slug, setSlug] = useState(category?.slug ?? '')
  const [slugTouched, setSlugTouched] = useState(!!category)
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault()
    setBusy(true)
    setError(null)
    try {
      const res = category
        ? await adminUpdateCategory(category.id, { name, slug })
        : await adminCreateCategory({ name, slug })
      if (res.status !== 200 && res.status !== 201) {
        setError(res.status === 409 ? 'That slug is already in use' : errorMessage(res, 'Could not save the category'))
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
    <Modal title={category ? 'Edit category' : 'New category'} onClose={onClose}>
      <form onSubmit={onSubmit} className="space-y-4">
        <div className="space-y-1">
          <label htmlFor="cat-name" className="block text-sm font-medium">
            Name
          </label>
          <input
            id="cat-name"
            required
            value={name}
            onChange={(e) => {
              setName(e.target.value)
              if (!slugTouched) setSlug(slugify(e.target.value))
            }}
            className="w-full rounded-md border px-3 py-2"
          />
        </div>
        <div className="space-y-1">
          <label htmlFor="cat-slug" className="block text-sm font-medium">
            Slug
          </label>
          <input
            id="cat-slug"
            required
            value={slug}
            onChange={(e) => {
              setSlug(e.target.value)
              setSlugTouched(true)
            }}
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
