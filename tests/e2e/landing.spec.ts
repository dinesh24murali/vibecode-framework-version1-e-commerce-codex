import { test, expect } from '@playwright/test'

const now = '2026-02-01T12:00:00Z'
const cat = { id: 'c1', name: 'Journals', slug: 'journals', image_url: null, created_at: now, updated_at: now }
// 1x1 transparent PNG so slide images load without any network.
const PNG = Buffer.from(
  'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==',
  'base64'
)

const slides = [1, 2, 3].map((n) => ({
  id: `s${n}`, title: `Slide ${n}`, subtitle: n === 1 ? 'New arrivals' : null,
  image_url: `http://img.test/s${n}.png`, cta_label: n === 1 ? 'Shop Now' : null,
  cta_url: n === 1 ? '/products/' : null, sort_order: n,
}))
const product = (id: string, name: string, price: number, disc: number) => ({
  id, name, slug: id, price, discount_percent: disc,
  discounted_price: Math.round(price * (100 - disc)) / 100,
  attributes: [], cover_image: null, categories: [cat], is_active: true,
})

test.beforeEach(async ({ page }) => {
  await page.route('http://img.test/**', (r) => r.fulfill({ status: 200, contentType: 'image/png', body: PNG }))
})

test('landing renders carousel, product cards and add to cart', async ({ page }) => {
  // The static build had no API, so the page fetches /landing on load.
  await page.route('**/api/v1/landing', (r) =>
    r.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        carousel_slides: slides,
        featured_products: [product('p1', 'Mandala Journal', 500, 10), product('p2', 'Sketch Pad', 200, 0)],
      }),
    })
  )
  await page.goto('/')

  await expect(page.getByTestId('hero-carousel')).toBeVisible()
  await expect(page.getByTestId('carousel-slide')).toHaveCount(3)
  await expect(page.getByRole('link', { name: 'Shop Now' })).toBeVisible()

  const cards = page.getByTestId('product-card')
  await expect(cards).toHaveCount(2)
  await expect(cards.nth(0)).toContainText('Mandala Journal')
  await expect(cards.nth(0)).toContainText('450.00')
  await expect(cards.nth(0)).toContainText('10% off')
  await expect(cards.nth(1)).not.toContainText('% off')
  await expect(cards.nth(1).getByRole('button', { name: 'Add to Cart' })).toBeVisible()
})

test('carousel auto-advances', async ({ page }) => {
  await page.route('**/api/v1/landing', (r) =>
    r.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ carousel_slides: slides, featured_products: [] }) })
  )
  await page.goto('/')
  const dots = page.getByRole('button', { name: /Go to slide/ })
  await expect(dots.nth(0)).toHaveAttribute('aria-current', 'true')
  await expect(dots.nth(1)).toHaveAttribute('aria-current', 'true', { timeout: 8000 })
})

test('shows fallbacks when the API fails', async ({ page }) => {
  await page.route('**/api/v1/landing', (r) => r.fulfill({ status: 500, contentType: 'application/json', body: '{"code":500,"message":"x"}' }))
  await page.goto('/')
  await expect(page.getByTestId('carousel-fallback')).toBeVisible()
  await expect(page.getByText('Products are temporarily unavailable.')).toBeVisible()
})
