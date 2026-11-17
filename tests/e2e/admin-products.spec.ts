import { test, expect, type Page } from '@playwright/test'

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
const category = { id: 'c1', name: 'Journals', slug: 'journals', image_url: null, created_at: ts, updated_at: ts }
const product = (id: string, name: string, active = true) => ({
  id, name, slug: name.toLowerCase(), description: 'd', price: 100, discount_percent: 0, discounted_price: 100,
  stock: 5, attributes: [], images: [{ id: 'i1', url: 'https://cdn.example.com/a.jpg', alt_text: null, sort_order: 1 }],
  categories: [category], is_active: active, created_at: ts, updated_at: ts,
})

async function mockApi(page: Page) {
  const items = [product('p1', 'Alpha')]
  const calls: string[] = []
  await page.route('**/api/v1/admin/auth/login', (r) => r.fulfill(json(200, admin)))
  await page.route('**/api/v1/categories', (r) => r.fulfill(json(200, { data: [category] })))
  await page.route((url) => url.pathname.startsWith('/api/v1/admin/products'), async (r) => {
    const req = r.request()
    const path = new URL(req.url()).pathname
    const m = req.method()
    calls.push(`${m} ${path}`)
    if (path.endsWith('/p1') && m === 'GET') return r.fulfill(json(200, items[0]))
    if (path.endsWith('/p1') && m === 'PUT') {
      items[0] = { ...items[0], name: req.postDataJSON().name }
      return r.fulfill(json(200, items[0]))
    }
    if (path.endsWith('/p1') && m === 'DELETE') {
      items[0] = { ...items[0], is_active: false }
      return r.fulfill(json(200, { message: 'product deleted' }))
    }
    if (path.endsWith('/products') && m === 'POST') {
      const created = product('p2', req.postDataJSON().name)
      items.push(created)
      return r.fulfill(json(201, created))
    }
    if (path.endsWith('/products') && m === 'GET')
      return r.fulfill(json(200, { data: items, meta: { total: items.length, page: 1, limit: 100, total_pages: 1 } }))
    return r.fulfill(json(404, { code: 404, message: 'nf' }))
  })
  return calls
}

test('admin creates, edits (prefetched by id) and soft-deletes a product', async ({ page }) => {
  const calls = await mockApi(page)
  await page.goto('/admin/login/')
  await page.getByLabel('Email').fill('admin@example.com')
  await page.getByLabel('Password').fill('x')
  await page.getByRole('button', { name: 'Sign in' }).click()
  await page.getByRole('link', { name: 'Products' }).click()
  await expect(page.getByRole('cell', { name: 'Alpha' })).toBeVisible()

  // Create
  await page.getByRole('button', { name: 'New product' }).click()
  await page.getByLabel('Name', { exact: true }).fill('Beta')
  await page.getByLabel('Slug').fill('beta')
  await page.getByLabel('Price (INR)').fill('250')
  await page.getByLabel('Category').selectOption('c1')
  await page.getByRole('button', { name: 'Save product' }).click()
  await expect(page.getByRole('cell', { name: 'Beta' })).toBeVisible()

  // Edit: the drawer fetches GET /admin/products/p1 and is pre-populated.
  await page.getByRole('row', { name: /Alpha/ }).getByRole('button', { name: 'Edit' }).click()
  await expect(page.getByLabel('Name', { exact: true })).toHaveValue('Alpha')
  expect(calls).toContain('GET /api/v1/admin/products/p1')
  await page.getByLabel('Name', { exact: true }).fill('Alpha 2')
  await page.getByRole('button', { name: 'Save product' }).click()
  await expect(page.getByRole('cell', { name: 'Alpha 2' })).toBeVisible()

  // Soft delete
  page.once('dialog', (d) => d.accept())
  await page.getByRole('row', { name: /Alpha 2/ }).getByRole('button', { name: 'Delete' }).click()
  await expect(page.getByRole('row', { name: /Alpha 2/ })).toContainText('Hidden')
})
