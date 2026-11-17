import { test, expect, type Page } from '@playwright/test'

const now = '2026-02-01T12:00:00Z'
const ID = '33333333-3333-4333-8333-333333333333'
const SLUG = 'sample-product'
const summary = { id: ID, name: 'Mandala Journal', slug: SLUG, price: 200, discount_percent: 10, discounted_price: 180, attributes: [], cover_image: null, categories: [], is_active: true }
const detail = { ...summary, description: '<p>Nice.</p>', stock: 5, images: [], created_at: now, updated_at: now }
const address = {
  id: 'a1', address_line1: '42 Rose St', address_line2: null, landmark: null, city: 'Chennai', state: 'TN',
  pincode: '600001', company_name: null, gst_number: null, phone: '+919876543210', is_default: true, created_at: now, updated_at: now,
}

const json = (status: number, body: unknown) => ({ status, contentType: 'application/json', body: JSON.stringify(body) })

interface Opts { qty: number; addresses?: unknown[]; placeStatus?: number }

async function mockApi(page: Page, opts: Opts) {
  let qty = opts.qty
  let placed: Record<string, unknown> | null = null
  const cart = () => {
    const items = qty > 0 ? [{ id: 'item-1', product: summary, quantity: qty, unit_price: 180, line_total: 180 * qty }] : []
    return { id: 'u1', items, item_count: items.length, subtotal: 180 * qty, updated_at: now }
  }
  await page.route('**/api/v1/cart', (r) => r.fulfill(json(200, cart())))
  await page.route('**/api/v1/cart/items', (r) => {
    qty += r.request().postDataJSON().quantity
    return r.fulfill(json(201, cart()))
  })
  await page.route('**/api/v1/account/addresses', (r) => r.fulfill(json(200, { data: opts.addresses ?? [address] })))
  await page.route(/\/api\/v1\/products(\?.*)?$/, (r) =>
    r.fulfill(json(200, { data: [summary], meta: { total: 1, page: 1, limit: 100, total_pages: 1 }, filters: { categories: [], attributes: [], sort_options: [] } })))
  await page.route(`**/api/v1/products/${ID}`, (r) => r.fulfill(json(200, detail)))

  await page.route('**/api/v1/checkout/coupons/validate', (r) => {
    const code = String(r.request().postDataJSON().coupon_code).toUpperCase()
    return r.fulfill(json(200, code === 'SUPERHIT'
      ? { is_valid: true, coupon_code: 'SUPERHIT', discount_type: 'flat', discount_value: 150, message: 'Coupon applied — INR 150 off your order' }
      : { is_valid: false, coupon_code: code, discount_type: null, discount_value: null, message: 'Invalid coupon code' }))
  })
  await page.route('**/api/v1/checkout/calculate', (r) => {
    const coupon = r.request().postDataJSON().coupon_code as string | null
    const subtotal = 180 * qty
    const discount = coupon === 'SUPERHIT' ? 150 : 0
    const taxable = subtotal - discount
    const gst = Math.round(taxable * 18) / 100
    return r.fulfill(json(200, {
      items: cart().items, subtotal, coupon_code: coupon, coupon_discount: discount,
      taxable_amount: taxable, gst_amount: gst, total: taxable + gst, address,
    }))
  })
  await page.route('**/api/v1/checkout/orders', (r) => {
    placed = r.request().postDataJSON()
    if (opts.placeStatus && opts.placeStatus !== 201) return r.fulfill(json(opts.placeStatus, { code: opts.placeStatus, message: 'Insufficient stock' }))
    qty = 0
    return r.fulfill(json(201, { id: 'o1', order_number: 'ORD-1', status: 'open', items: [], shipping_address: address, subtotal: 0, coupon_discount: 0, taxable_amount: 0, gst_amount: 0, total: 0, created_at: now, updated_at: now }))
  })
  await page.route('**/api/v1/orders*', (r) => r.fulfill(json(200, { data: [], meta: { total: 0, page: 1, limit: 20, total_pages: 0 } })))
  await page.route('**/api/v1/account/profile', (r) => r.fulfill(json(200, { id: 'u1', email: 'alice@example.com', first_name: 'Alice', last_name: 'Doe', role: 'customer', created_at: now, updated_at: now })))
  return { placed: () => placed }
}

/** Logs in through the UI (tokens are memory-only, so later navigation must be client-side). */
async function loginTo(page: Page, redirect: string) {
  await page.route('**/api/v1/auth/login', (r) =>
    r.fulfill(json(200, { access_token: 'tok', refresh_token: 'ref', token_type: 'bearer', expires_in: 900 })))
  await page.goto(`/login/?redirect=${redirect}`)
  await page.getByLabel('Email').fill('alice@example.com')
  await page.getByLabel('Password').fill('Secure@1234')
  await page.getByRole('button', { name: 'Log in' }).click()
}

test('unauthenticated visitors are sent to login with a redirect back', async ({ page }) => {
  await page.goto('/checkout/')
  await expect(page).toHaveURL(/\/login\/\?redirect=\/checkout\/$/)
})

test('an empty cart redirects to /cart', async ({ page }) => {
  await mockApi(page, { qty: 0 })
  await loginTo(page, '/checkout/')
  await expect(page).toHaveURL(/\/cart\/$/)
  await expect(page.getByTestId('cart-empty')).toBeVisible()
})

test('happy path: add to cart, checkout, coupon, place order', async ({ page }) => {
  const api = await mockApi(page, { qty: 0 })
  await loginTo(page, `/products/${SLUG}/`)
  await page.getByRole('button', { name: 'Add to Cart' }).click()
  await expect(page.getByTestId('toast')).toBeVisible()
  await page.getByRole('link', { name: /cart/i }).first().click()
  await page.getByRole('link', { name: 'Proceed to Checkout' }).click()

  await expect(page.getByRole('heading', { name: 'Checkout' })).toBeVisible()
  await expect(page.getByLabel(/42 Rose St/)).toBeChecked()
  await expect(page.getByTestId('co-subtotal')).toHaveText('₹180.00')
  await expect(page.getByTestId('co-gst')).toHaveText('₹32.40')
  await expect(page.getByTestId('co-total')).toHaveText('₹212.40')

  // Invalid coupon: red feedback, total unchanged.
  await page.getByLabel('Coupon code').fill('nope')
  await page.getByLabel('Coupon code').blur()
  await expect(page.getByTestId('coupon-feedback')).toHaveText('Invalid coupon code')
  await expect(page.getByTestId('co-total')).toHaveText('₹212.40')

  // Valid coupon: discount applied before tax, total updates live.
  await page.getByLabel('Coupon code').fill('superhit')
  await page.getByLabel('Coupon code').blur()
  await expect(page.getByTestId('coupon-feedback')).toContainText('INR 150 off')
  await expect(page.getByTestId('co-discount')).toHaveText('-₹150.00')
  await expect(page.getByTestId('co-gst')).toHaveText('₹5.40')
  await expect(page.getByTestId('co-total')).toHaveText('₹35.40')

  await page.getByRole('button', { name: 'Place Order' }).click()
  await expect(page).toHaveURL(/\/account\/\?tab=orders$/)
  await expect(page.getByRole('tab', { name: 'Orders' })).toHaveAttribute('aria-selected', 'true')
  expect(api.placed()).toEqual({ address_id: 'a1', coupon_code: 'SUPERHIT' })
})

test('without a saved address the user is asked to add one', async ({ page }) => {
  await mockApi(page, { qty: 1, addresses: [] })
  await loginTo(page, '/checkout/')
  await expect(page.getByTestId('no-addresses')).toBeVisible()
  await expect(page.getByRole('button', { name: 'Place Order' })).toBeDisabled()
})

test('a failed order shows the API message and stays on checkout', async ({ page }) => {
  await mockApi(page, { qty: 1, placeStatus: 422 })
  await loginTo(page, '/checkout/')
  await page.getByRole('button', { name: 'Place Order' }).click()
  await expect(page.getByText('Insufficient stock')).toBeVisible()
  await expect(page).toHaveURL(/\/checkout\/$/)
})
