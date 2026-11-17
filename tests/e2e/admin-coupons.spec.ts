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
const coupon = (id: string, code: string, extra: Record<string, unknown> = {}) => ({
  id, name: code, description: null, code, is_enable: true, type: 'flat', flat: 50, percent: null,
  coupon_expiry: null, created_at: ts, updated_at: ts, ...extra,
})

test('admin creates, toggles and deletes coupons; type gates the amount field', async ({ page }) => {
  const items = [coupon('c1', 'SAVE50'), coupon('c2', 'USED10', { type: 'percentage', flat: null, percent: 10 })]
  let created: Record<string, unknown> | null = null
  await page.route('**/api/v1/admin/auth/login', (r) => r.fulfill(json(200, admin)))
  await page.route((url) => url.pathname.startsWith('/api/v1/admin/coupons'), async (r) => {
    const req = r.request()
    const path = new URL(req.url()).pathname
    const m = req.method()
    if (path.endsWith('/c2') && m === 'DELETE') return r.fulfill(json(409, { code: 409, message: 'conflict' }))
    if (path.endsWith('/c1') && m === 'DELETE') {
      items.splice(0, 1)
      return r.fulfill(json(200, { message: 'coupon deleted' }))
    }
    if (path.endsWith('/c1') && m === 'PUT') {
      items[0] = { ...items[0], ...req.postDataJSON() }
      return r.fulfill(json(200, items[0]))
    }
    if (path.endsWith('/coupons') && m === 'POST') {
      created = req.postDataJSON()
      items.push(coupon('c3', String(created!.code).toUpperCase(), created!))
      return r.fulfill(json(201, items[items.length - 1]))
    }
    if (path.endsWith('/coupons') && m === 'GET')
      return r.fulfill(json(200, { data: items, meta: { total: items.length, page: 1, limit: 100, total_pages: 1 } }))
    return r.fulfill(json(404, { code: 404, message: 'nf' }))
  })

  await page.goto('/admin/login/')
  await page.getByLabel('Email').fill('admin@example.com')
  await page.getByLabel('Password').fill('x')
  await page.getByRole('button', { name: 'Sign in' }).click()
  await page.getByRole('link', { name: 'Coupons' }).click()
  await expect(page.getByRole('cell', { name: 'SAVE50' })).toBeVisible()

  // Create a percentage coupon: only the percent field is offered and sent.
  await page.getByRole('button', { name: 'New coupon' }).click()
  await page.getByLabel('Name', { exact: true }).fill('Welcome')
  await page.getByLabel('Code', { exact: true }).fill('welcome15')
  await expect(page.getByLabel('Flat amount (INR)')).toBeVisible()
  await page.getByLabel('Type').selectOption('percentage')
  await expect(page.getByLabel('Flat amount (INR)')).toHaveCount(0)
  await page.getByLabel('Percent off').fill('15')
  await page.getByRole('button', { name: 'Save' }).click()
  await expect(page.getByRole('cell', { name: 'WELCOME15' })).toBeVisible()
  expect(created).toMatchObject({ type: 'percentage', percent: 15 })
  expect(created).not.toHaveProperty('flat')

  // Disable via the toggle.
  await page.getByRole('row', { name: /SAVE50/ }).getByRole('button', { name: 'Disable' }).click()
  await expect(page.getByRole('row', { name: /SAVE50/ })).toContainText('Disabled')

  // Deleting a redeemed coupon explains why; deleting an unused one removes it.
  page.once('dialog', (d) => d.accept())
  await page.getByRole('row', { name: /USED10/ }).getByRole('button', { name: 'Delete' }).click()
  await expect(page.getByText('already been redeemed')).toBeVisible()
  page.once('dialog', (d) => d.accept())
  await page.getByRole('row', { name: /SAVE50/ }).getByRole('button', { name: 'Delete' }).click()
  await expect(page.getByRole('cell', { name: 'SAVE50' })).toHaveCount(0)
})
