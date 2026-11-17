import { test, expect } from '@playwright/test'

const json = (status: number, body: unknown) => ({
  status,
  contentType: 'application/json',
  body: JSON.stringify(body),
})
const ts = '2026-01-01T00:00:00Z'
const admin = {
  access_token: 'admin-access',
  refresh_token: 'r',
  token_type: 'bearer',
  expires_in: 900,
  user: { id: 'u1', email: 'admin@example.com', first_name: 'A', last_name: 'B', role: 'admin' },
}
const cat = (id: string, name: string, slug: string) => ({ id, name, slug, image_url: null, created_at: ts, updated_at: ts })

test('admin creates, renames and deletes categories; 409 shows a clear message', async ({ page }) => {
  const items = [cat('c1', 'Journals', 'journals'), cat('c2', 'Planners', 'planners')]
  await page.route('**/api/v1/admin/auth/login', (r) => r.fulfill(json(200, admin)))
  await page.route((url) => url.pathname.startsWith('/api/v1/admin/categories'), async (r) => {
    const req = r.request()
    const path = new URL(req.url()).pathname
    const m = req.method()
    if (path.endsWith('/c1') && m === 'DELETE') return r.fulfill(json(409, { code: 409, message: 'conflict' }))
    if (path.endsWith('/c2') && m === 'DELETE') {
      items.splice(1, 1)
      return r.fulfill(json(200, { message: 'category deleted' }))
    }
    if (path.endsWith('/c1') && m === 'PUT') {
      items[0] = cat('c1', req.postDataJSON().name, 'journals')
      return r.fulfill(json(200, items[0]))
    }
    if (path.endsWith('/categories') && m === 'POST') {
      const b = req.postDataJSON()
      items.push(cat('c3', b.name, b.slug))
      return r.fulfill(json(201, items[items.length - 1]))
    }
    if (path.endsWith('/categories') && m === 'GET')
      return r.fulfill(json(200, { data: items, meta: { total: items.length, page: 1, limit: 100, total_pages: 1 } }))
    return r.fulfill(json(404, { code: 404, message: 'nf' }))
  })

  await page.goto('/admin/login/')
  await page.getByLabel('Email').fill('admin@example.com')
  await page.getByLabel('Password').fill('x')
  await page.getByRole('button', { name: 'Sign in' }).click()
  await page.getByRole('link', { name: 'Categories' }).click()
  await expect(page.getByRole('cell', { name: 'Journals', exact: true })).toBeVisible()

  // Create: slug auto-fills from the name; the list refreshes without a reload.
  await page.getByRole('button', { name: 'New category' }).click()
  await page.getByLabel('Name', { exact: true }).fill('Sticker Books')
  await expect(page.getByLabel('Slug')).toHaveValue('sticker-books')
  await page.getByRole('button', { name: 'Save' }).click()
  await expect(page.getByRole('cell', { name: 'Sticker Books', exact: true })).toBeVisible()

  // Rename
  await page.getByRole('row', { name: /Journals/ }).getByRole('button', { name: 'Edit' }).click()
  await page.getByLabel('Name', { exact: true }).fill('Journals 2')
  await page.getByRole('button', { name: 'Save' }).click()
  await expect(page.getByRole('cell', { name: 'Journals 2', exact: true })).toBeVisible()

  // Delete with products -> 409 message; delete without -> removed.
  page.once('dialog', (d) => d.accept())
  await page.getByRole('row', { name: /Journals 2/ }).getByRole('button', { name: 'Delete' }).click()
  await expect(page.getByText('still has products')).toBeVisible()
  page.once('dialog', (d) => d.accept())
  await page.getByRole('row', { name: /Planners/ }).getByRole('button', { name: 'Delete' }).click()
  await expect(page.getByRole('cell', { name: 'Planners', exact: true })).toHaveCount(0)
})
