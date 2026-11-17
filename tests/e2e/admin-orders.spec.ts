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
const address = {
  id: 'a1', address_line1: '42 Rose Street', address_line2: null, landmark: null, city: 'Pune', state: 'MH',
  pincode: '411001', company_name: null, gst_number: null, phone: '9999999999', is_default: false, created_at: ts, updated_at: ts,
}
const order = (id: string, num: string, status: string) => ({
  id, order_number: num, status, customer_name: 'Zed Customer', customer_email: 'zed@example.com',
  items: [{ id: 'it1', product_id: 'p1', product_name: 'Mandala Journal', product_slug: 'm', cover_image_url: null, quantity: 2, unit_price: 50, line_total: 100 }],
  shipping_address: address, subtotal: 100, coupon_code: null, coupon_discount: 0, taxable_amount: 100, gst_amount: 18,
  total: 118, created_at: ts, updated_at: ts,
})

test('admin filters orders by status tab, opens the detail and fulfils an order', async ({ page }) => {
  const orders = [order('o1', 'ORD-1', 'open'), order('o2', 'ORD-2', 'cancelled')]
  const listCalls: string[] = []
  await page.route('**/api/v1/admin/auth/login', (r) => r.fulfill(json(200, admin)))
  await page.route((url) => url.pathname.startsWith('/api/v1/admin/orders'), async (r) => {
    const req = r.request()
    const u = new URL(req.url())
    if (u.pathname.endsWith('/o1/status') && req.method() === 'PUT') {
      orders[0] = order('o1', 'ORD-1', req.postDataJSON().status)
      return r.fulfill(json(200, orders[0]))
    }
    if (u.pathname.endsWith('/o1')) return r.fulfill(json(200, orders[0]))
    listCalls.push(u.search)
    const status = u.searchParams.get('status')
    const data = status ? orders.filter((o) => o.status === status) : orders
    return r.fulfill(json(200, { data, meta: { total: data.length, page: 1, limit: 20, total_pages: 1 } }))
  })

  await page.goto('/admin/login/')
  await page.getByLabel('Email').fill('admin@example.com')
  await page.getByLabel('Password').fill('x')
  await page.getByRole('button', { name: 'Sign in' }).click()
  await page.getByRole('link', { name: 'Orders' }).click()
  await expect(page.getByRole('cell', { name: 'ORD-1' })).toBeVisible()
  await expect(page.getByRole('cell', { name: 'ORD-2' })).toBeVisible()
  expect(listCalls[0]).toContain('limit=20')

  // Status tab sends the filter and narrows the table.
  await page.getByRole('tab', { name: 'Cancelled' }).click()
  await expect(page.getByRole('cell', { name: 'ORD-1' })).toHaveCount(0)
  await expect(page.getByRole('cell', { name: 'ORD-2' })).toBeVisible()
  await page.getByRole('tab', { name: 'Open' }).click()
  await expect(page.getByRole('cell', { name: 'ORD-1' })).toBeVisible()

  // Detail modal shows customer, address, items, total.
  await page.getByRole('button', { name: 'View' }).click()
  const dialog = page.getByRole('dialog')
  await expect(dialog.getByText('Zed Customer')).toBeVisible()
  await expect(dialog.getByText('42 Rose Street')).toBeVisible()
  await expect(dialog.getByText('Mandala Journal × 2')).toBeVisible()
  await expect(dialog.getByText('₹118')).toBeVisible()

  // Change status via the dropdown.
  await dialog.getByLabel(/Status/).selectOption('fulfilled')
  await dialog.getByRole('button', { name: 'Update' }).click()
  await expect(dialog.getByText('Status (fulfilled)')).toBeVisible()
})
