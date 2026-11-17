import { test, expect, type Page } from '@playwright/test'

const product = {
  id: 'p1', name: 'Mandala Journal', slug: 'mandala-journal', price: 100, discount_percent: 0,
  discounted_price: 100, attributes: [], cover_image: null, categories: [], is_active: true,
}

async function mockApi(page: Page) {
  const json = (body: unknown) => ({ status: 200, contentType: 'application/json', body: JSON.stringify(body) })
  await page.route('**/api/v1/landing', (r) => r.fulfill(json({ carousel_slides: [], featured_products: [], new_arrivals: [] })))
  await page.route('**/api/v1/categories', (r) => r.fulfill(json({ data: [] })))
  await page.route('**/api/v1/attributes', (r) => r.fulfill(json({ data: [] })))
  await page.route(/\/api\/v1\/products(\?.*)?$/, (r) =>
    r.fulfill(json({ data: [product], meta: { total: 1, page: 1, limit: 12, total_pages: 1 }, filters: { categories: [], attributes: [], sort_options: [] } })))
}

test('Book Store menu lists 12 categories and navigates on click', async ({ page }) => {
  await mockApi(page)
  await page.goto('/')
  await page.getByRole('button', { name: 'Book Store' }).hover()
  const menu = page.getByTestId('bookstore-menu')
  await expect(menu.getByRole('link')).toHaveCount(12)
  await menu.getByRole('link', { name: 'Journals' }).click()
  await expect(page).toHaveURL(/\/products\/\?category=journals$/)
})

test('Printing Services menu shows both columns with 24 items', async ({ page }) => {
  await mockApi(page)
  await page.goto('/')
  await page.getByRole('button', { name: 'Printing Services' }).hover()
  const menu = page.getByTestId('printing-menu')
  await expect(menu.getByText('Products Offered')).toBeVisible()
  await expect(menu.getByText('Inhouse Services')).toBeVisible()
  await expect(menu.getByRole('link')).toHaveCount(24)
})

test('search submits on Enter to /products/?q=', async ({ page }) => {
  await mockApi(page)
  await page.goto('/')
  await page.getByRole('searchbox', { name: 'Search store' }).fill('mandala art')
  await page.keyboard.press('Enter')
  await expect(page).toHaveURL(/\/products\/\?q=mandala%20art$/)
})

test('cart badge reflects the cart item count', async ({ page }) => {
  await mockApi(page)
  await page.goto('/products/')
  await expect(page.getByTestId('cart-badge')).toHaveCount(0)
  await page.getByRole('button', { name: 'Add to Cart' }).click()
  await page.getByRole('button', { name: 'Add to Cart' }).click()
  await expect(page.getByTestId('cart-badge')).toHaveText('2')
})

test.describe('mobile', () => {
  test.use({ viewport: { width: 375, height: 800 } })

  test('hamburger drawer has accordion sections', async ({ page }) => {
    await mockApi(page)
    await page.goto('/')
    await expect(page.getByRole('button', { name: 'Book Store' })).toBeHidden()
    await page.getByRole('button', { name: 'Open menu' }).click()
    const drawer = page.getByTestId('mobile-drawer')
    await drawer.getByRole('button', { name: /Book Store/ }).click()
    await expect(drawer.getByTestId('bookstore-menu').getByRole('link')).toHaveCount(12)
    await drawer.getByRole('button', { name: /Printing Services/ }).click()
    await expect(drawer.getByTestId('bookstore-menu')).toHaveCount(0)
    await expect(drawer.getByTestId('printing-menu').getByRole('link')).toHaveCount(24)
  })
})
