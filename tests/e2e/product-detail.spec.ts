import { test, expect, type Page } from '@playwright/test'

const now = '2026-02-01T12:00:00Z'
const ID = '33333333-3333-4333-8333-333333333333'
const SLUG = 'sample-product'
const summary = { id: ID, name: 'Mandala Journal', slug: SLUG, price: 200, discount_percent: 10, discounted_price: 180, attributes: [], cover_image: null, categories: [], is_active: true }
const detail = {
  ...summary,
  description: '<p>A lovely journal.</p>',
  stock: 5,
  images: [
    { id: 'i1', url: 'data:image/svg+xml,%3Csvg xmlns=%22http://www.w3.org/2000/svg%22 width=%2240%22 height=%2240%22/%3E', alt_text: 'Front', sort_order: 1 },
    { id: 'i2', url: 'data:image/svg+xml,%3Csvg xmlns=%22http://www.w3.org/2000/svg%22 width=%2250%22 height=%2250%22/%3E', alt_text: 'Back', sort_order: 2 },
  ],
  attributes: [{ attribute_id: 'a1', code: 'level', attribute_name: 'Skill Level', value: 'Beginner', value_slug: 'beginner' }],
  created_at: now,
  updated_at: now,
}

async function mockApi(page: Page) {
  const json = (body: unknown) => ({ status: 200, contentType: 'application/json', body: JSON.stringify(body) })
  await page.route(/\/api\/v1\/products(\?.*)?$/, (r) =>
    r.fulfill(json({ data: [summary], meta: { total: 1, page: 1, limit: 100, total_pages: 1 }, filters: { categories: [], attributes: [], sort_options: [] } })))
  await page.route(`**/api/v1/products/${ID}`, (r) => r.fulfill(json(detail)))
}

test('loads product, switches gallery image, zooms on hover and adds to cart', async ({ page }) => {
  await mockApi(page)
  await page.goto(`/products/${SLUG}/`)

  await expect(page.getByRole('heading', { name: 'Mandala Journal' })).toBeVisible()
  await expect(page.getByText('₹180.00')).toBeVisible()
  await expect(page.getByText('Skill Level:')).toBeVisible()
  await expect(page).toHaveTitle('Mandala Journal')

  const zoomImg = page.getByTestId('zoom-image')
  await expect(zoomImg).toHaveAttribute('alt', 'Front')
  await page.getByTestId('thumbnail').nth(1).click()
  await expect(zoomImg).toHaveAttribute('alt', 'Back')

  await expect(zoomImg).toHaveAttribute('data-scale', '1')
  await page.getByTestId('product-zoom').hover()
  await expect(zoomImg).toHaveAttribute('data-scale', '2')

  await page.getByRole('button', { name: 'Add to Cart' }).click()
  await expect(page.getByTestId('toast')).toHaveText('Mandala Journal added to cart')
})

test('shows not found for an unknown slug', async ({ page }) => {
  await mockApi(page)
  await page.goto(`/products/${SLUG}/`)
  // Pre-rendered fallback shell resolves from the URL; an unmatched slug shows not found.
  await page.route(/\/api\/v1\/products(\?.*)?$/, (r) =>
    r.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ data: [], meta: { total: 0, page: 1, limit: 100, total_pages: 0 }, filters: { categories: [], attributes: [], sort_options: [] } }) }))
  await page.reload()
  await expect(page.getByText('Product not found.')).toBeVisible()
})
