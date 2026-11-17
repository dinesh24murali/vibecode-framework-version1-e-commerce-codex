'use client'
import { useCallback, useEffect, useState } from 'react'
import {
  adminDeleteCategory,
  adminListCategories,
} from '@/lib/api/generated/admin-categories/admin-categories'
import type { Category } from '@/lib/api/generated/model'
import { errorMessage } from '@/lib/auth/helpers'
import { Button } from '@/components/ui/button'
import { CategoryFormModal } from '@/components/admin/CategoryFormModal'

type ModalState = { category?: Category } | null

// Create/edit/delete all happen in a modal on this page (no dynamic [id] route: static export).
export default function AdminCategoriesPage() {
  const [categories, setCategories] = useState<Category[]>([])
  const [modal, setModal] = useState<ModalState>(null)
  const [error, setError] = useState<string | null>(null)
  const [loaded, setLoaded] = useState(false)

  const load = useCallback(async () => {
    try {
      const res = await adminListCategories({ limit: 100 })
      if (res.status === 200) setCategories(res.data.data)
      else setError(errorMessage(res, 'Could not load categories'))
    } catch {
      setError('Network error. Please try again.')
    } finally {
      setLoaded(true)
    }
  }, [])

  useEffect(() => {
    void load()
  }, [load])

  async function remove(c: Category) {
    if (!window.confirm(`Delete category "${c.name}"?`)) return
    const res = await adminDeleteCategory(c.id)
    if (res.status === 200) {
      setError(null)
      await load()
    } else if (res.status === 409) {
      setError(`"${c.name}" still has products. Move or remove them first.`)
    } else setError(errorMessage(res, 'Could not delete the category'))
  }

  return (
    <main className="mx-auto max-w-3xl space-y-6 p-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Categories</h1>
        <Button onClick={() => setModal({})}>New category</Button>
      </div>
      {error && (
        <p role="alert" className="text-sm text-red-600">
          {error}
        </p>
      )}
      {loaded && categories.length === 0 && <p className="text-sm text-gray-600">No categories yet.</p>}
      {loaded && categories.length > 0 && (
        <table className="w-full text-left text-sm">
          <thead>
            <tr className="border-b">
              <th className="py-2">Name</th>
              <th>Slug</th>
              <th className="text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            {categories.map((c) => (
              <tr key={c.id} className="border-b">
                <td className="py-2">{c.name}</td>
                <td>{c.slug}</td>
                <td className="space-x-2 text-right">
                  <button className="underline" onClick={() => setModal({ category: c })}>
                    Edit
                  </button>
                  <button className="text-red-600 underline" onClick={() => void remove(c)}>
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
      {modal && (
        <CategoryFormModal
          category={modal.category}
          onClose={() => setModal(null)}
          onSaved={() => {
            setModal(null)
            setError(null)
            void load()
          }}
        />
      )}
    </main>
  )
}
