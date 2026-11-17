import { test, expect, type Page } from '@playwright/test'

const json = (status: number, body: unknown) => ({
  status,
  contentType: 'application/json',
  body: JSON.stringify(body),
})

const admin = {
  access_token: 'admin-access',
  refresh_token: 'admin-refresh',
  token_type: 'bearer',
  expires_in: 900,
  user: { id: 'u1', email: 'admin@example.com', first_name: 'A', last_name: 'B', role: 'admin' },
}
const ts = '2026-01-01T00:00:00Z'

// In-memory fake of the admin attributes API.
async function mockApi(page: Page) {
  const defs: Record<string, unknown>[] = [
    { id: 'a1', code: 'level', name: 'Level', is_filterable: true, sort_order: 1, created_at: ts, updated_at: ts },
  ]
  const vals: Record<string, unknown>[] = [
    { id: 'v1', attribute_id: 'a1', value: 'Beginner', value_slug: 'beginner', sort_order: 1, created_at: ts, updated_at: ts },
  ]
  await page.route('**/api/v1/admin/auth/login', (r) => r.fulfill(json(200, admin)))
  await page.route((url) => url.pathname.startsWith('/api/v1/admin/attributes'), async (r) => {
    const req = r.request()
    const path = new URL(req.url()).pathname
    const m = req.method()
    if (path.endsWith('/values/v1') && m === 'DELETE') return r.fulfill(json(409, { code: 409, message: 'conflict' }))
    if (path.endsWith('/values') && m === 'GET') return r.fulfill(json(200, { data: vals }))
    if (path.endsWith('/values') && m === 'POST') {
      const body = req.postDataJSON()
      vals.push({ id: 'v2', attribute_id: 'a1', value: body.value, value_slug: 'x', sort_order: 0, created_at: ts, updated_at: ts })
      return r.fulfill(json(201, vals[vals.length - 1]))
    }
    if (path.endsWith('/attributes') && m === 'GET')
      return r.fulfill(json(200, { data: defs, meta: { total: defs.length, page: 1, limit: 100, total_pages: 1 } }))
    if (path.endsWith('/attributes') && m === 'POST') {
      const body = req.postDataJSON()
      defs.push({ id: 'a2', code: body.code, name: body.name, is_filterable: true, sort_order: 0, created_at: ts, updated_at: ts })
      return r.fulfill(json(201, defs[defs.length - 1]))
    }
    return r.fulfill(json(404, { code: 404, message: 'nf' }))
  })
}

test('admin manages attributes and values via modals', async ({ page }) => {
  await mockApi(page)
  await page.goto('/admin/login/')
  await page.getByLabel('Email').fill('admin@example.com')
  await page.getByLabel('Password').fill('x')
  await page.getByRole('button', { name: 'Sign in' }).click()
  // The admin session is memory-only, so navigate client-side via the dashboard link.
  await page.getByRole('link', { name: 'Attributes' }).click()
  await expect(page.getByRole('heading', { name: 'Attributes' })).toBeVisible()
  await expect(page.getByRole('cell', { name: 'Level', exact: true })).toBeVisible()

  // Create an attribute in a modal.
  await page.getByRole('button', { name: 'New attribute' }).click()
  await page.getByLabel('Code', { exact: true }).fill('theme')
  await page.getByLabel('Name', { exact: true }).fill('Theme')
  await page.getByRole('button', { name: 'Save' }).click()
  await expect(page.getByRole('cell', { name: 'Theme', exact: true })).toBeVisible()

  // Values: add one, and a 409 on delete shows a clear message.
  await page.getByRole('button', { name: 'Values' }).first().click()
  await expect(page.getByText('Beginner')).toBeVisible()
  await page.getByRole('button', { name: 'Add value' }).click()
  await page.getByLabel('Value', { exact: true }).fill('Advanced')
  await page.getByRole('button', { name: 'Save' }).click()
  page.once('dialog', (d) => d.accept())
  await page.getByRole('button', { name: 'Remove' }).first().click()
  await expect(page.getByText('is assigned to products')).toBeVisible()
})
