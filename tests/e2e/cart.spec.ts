import { test, expect, type Page } from '@playwright/test'

const now = '2026-02-01T12:00:00Z'
const ID = '33333333-3333-4333-8333-333333333333'
const SLUG = 'sample-product'
const summary = { id: ID, name: 'Mandala Journal', slug: SLUG, price: 200, discount_percent: 10, discounted_price: 180, attributes: [], cover_image: null, categories: [], is_active: true }
const detail = { ...summary, description: '<p>Nice.</p>', stock: 5, images: [], created_at: now, updated_at: now }

const json = (status: number, body: unknown) => ({ status, contentType: 'application/json', body: JSON.stringify(body) })

/** Stateful in-memory fake of the cart endpoints. */
function fakeCart(initialQty = 0) {
  let qty = initialQty
  const cart = () => {
    const items = qty > 0 ? [{ id: 'item-1', product: summary, quantity: qty, unit_price: 180, line_total: 180 * qty }] : []
    return { id: 'u1', items, item_count: items.length, subtotal: 180 * qty, updated_at: now }
  }
  return {
    state: () => qty,
    async install(page: Page, opts: { failUpdate?: boolean } = {}) {
      await page.route('**/api/v1/cart', (r) => {
        if (r.request().method() === 'DELETE') {
          qty = 0
          return r.fulfill(json(200, { message: 'Cart cleared.' }))
        }
        return r.fulfill(json(200, cart()))
      })
      await page.route('**/api/v1/cart/items', (r) => {
        qty += r.request().postDataJSON().quantity
        return r.fulfill(json(201, cart()))
      })
      await page.route('**/api/v1/cart/items/item-1', (r) => {
        if (r.request().method() === 'DELETE') {
          qty = 0
          return r.fulfill(json(200, cart()))
        }
        if (opts.failUpdate) return r.fulfill(json(500, { code: 500, message: 'boom' }))
        qty = r.request().postDataJSON().quantity
        return r.fulfill(json(200, cart()))
      })
    },
  }
}

async function mockCatalog(page: Page) {
  await page.route(/\/api\/v1\/products(\?.*)?$/, (r) =>
    r.fulfill(json(200, { data: [summary], meta: { total: 1, page: 1, limit: 100, total_pages: 1 }, filters: { categories: [], attributes: [], sort_options: [] } })))
  await page.route(`**/api/v1/products/${ID}`, (r) => r.fulfill(json(200, detail)))
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

test('unauthenticated visitors see a login prompt and a disabled checkout', async ({ page }) => {
  await page.goto('/cart/')
  await expect(page).toHaveTitle(/Your Cart/)
  await expect(page.getByTestId('cart-login-prompt')).toBeVisible()
  await expect(page.getByRole('button', { name: 'Proceed to Checkout' })).toHaveCount(0)
})

test('empty cart shows the empty state with a Continue Shopping link', async ({ page }) => {
  await fakeCart().install(page)
  await loginTo(page, '/cart/')
  await expect(page.getByTestId('cart-empty')).toBeVisible()
  await expect(page.getByRole('link', { name: 'Continue Shopping' })).toHaveAttribute('href', '/products/')
})

test('add from product page, view cart, change quantity, remove', async ({ page }) => {
  const cart = fakeCart()
  await cart.install(page)
  await mockCatalog(page)
  await loginTo(page, `/products/${SLUG}/`)

  await page.getByRole('button', { name: 'Add to Cart' }).click()
  await expect(page.getByTestId('toast')).toHaveText('Mandala Journal added to cart')
  expect(cart.state()).toBe(1)

  await page.getByRole('link', { name: /cart/i }).first().click()
  await expect(page.getByRole('heading', { name: 'Your Cart' })).toBeVisible()
  await expect(page.getByTestId('cart-line')).toHaveCount(1)
  await expect(page.getByTestId('summary-subtotal')).toHaveText('₹200.00')
  await expect(page.getByTestId('summary-discount')).toHaveText('-₹20.00')
  await expect(page.getByTestId('summary-total')).toHaveText('₹180.00')

  await page.getByRole('button', { name: 'Increase quantity of Mandala Journal' }).click()
  await expect(page.getByTestId('cart-qty')).toHaveText('2')
  await expect(page.getByTestId('summary-total')).toHaveText('₹360.00')
  expect(cart.state()).toBe(2)

  await expect(page.getByRole('link', { name: 'Proceed to Checkout' })).toBeVisible()

  // Quantity reaching zero removes the line.
  await page.getByRole('button', { name: 'Decrease quantity of Mandala Journal' }).click()
  await page.getByRole('button', { name: 'Decrease quantity of Mandala Journal' }).click()
  await expect(page.getByTestId('cart-empty')).toBeVisible()
  expect(cart.state()).toBe(0)
})

test('a failed quantity update rolls back and shows an error', async ({ page }) => {
  const cart = fakeCart(1)
  await cart.install(page, { failUpdate: true })
  await loginTo(page, '/cart/')
  await expect(page.getByTestId('cart-qty')).toHaveText('1')
  await page.getByRole('button', { name: 'Increase quantity of Mandala Journal' }).click()
  await expect(page.getByText('We could not update your cart')).toBeVisible()
  await expect(page.getByTestId('cart-qty')).toHaveText('1')
})
