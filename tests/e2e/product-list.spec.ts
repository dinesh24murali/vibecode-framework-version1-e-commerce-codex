import { test, expect, type Page } from '@playwright/test'

const now = '2026-02-01T12:00:00Z'
const cats = [
  { id: 'c1', name: 'Journals', slug: 'journals', image_url: null, created_at: now, updated_at: now },
  { id: 'c2', name: 'Planners', slug: 'planners', image_url: null, created_at: now, updated_at: now },
]
const BEGINNER = '11111111-1111-4111-8111-111111111111'
const MANDALA = '22222222-2222-4222-8222-222222222222'
const attrs = [
  {
    id: 'a1', code: 'level', name: 'Skill Level', is_filterable: true, sort_order: 1,
    values: [{ id: BEGINNER, attribute_id: 'a1', value: 'Beginner', value_slug: 'beginner', sort_order: 1, created_at: now, updated_at: now }],
  },
  {
    id: 'a2', code: 'theme', name: 'Theme', is_filterable: true, sort_order: 2,
    values: [{ id: MANDALA, attribute_id: 'a2', value: 'Mandala', value_slug: 'mandala', sort_order: 1, created_at: now, updated_at: now }],
  },
]
const product = (id: string, name: string) => ({
  id, name, slug: id, price: 100, discount_percent: 0, discounted_price: 100,
  attributes: [], cover_image: null, categories: [], is_active: true,
})

const ALL = [product('p1', 'Mandala Journal'), product('p2', 'Weekly Planner'), product('p3', 'Plain Journal')]
const FILTERED = [product('p1', 'Mandala Journal')]
const filters = { categories: ['journals'], attributes: [], sort_options: ['newest'] }

/** Mocks the catalogue API and records each /products request's query string. */
async function mockCatalogue(page: Page): Promise<URLSearchParams[]> {
  const requests: URLSearchParams[] = []
  const json = (body: unknown) => ({ status: 200, contentType: 'application/json', body: JSON.stringify(body) })
  await page.route('**/api/v1/categories', (r) => r.fulfill(json({ data: cats })))
  await page.route('**/api/v1/attributes', (r) => r.fulfill(json({ data: attrs })))
  await page.route(/\/api\/v1\/products(\?.*)?$/, (r) => {
    const q = new URL(r.request().url()).searchParams
    requests.push(q)
    const narrowed = q.getAll('category[]').length > 0 || q.getAll('attr_value_id[]').length > 0 || q.get('q')
    const data = narrowed ? FILTERED : ALL
    return r.fulfill(json({ data, meta: { total: data.length, page: 1, limit: 12, total_pages: 1 }, filters }))
  })
  return requests
}

test('filter changes update the grid and the URL without a page reload', async ({ page }) => {
  const requests = await mockCatalogue(page)
  await page.goto('/products/')
  await expect(page.getByTestId('product-card')).toHaveCount(3)

  await page.evaluate(() => ((window as unknown as { __marker: string }).__marker = 'still-here'))

  await page.getByLabel('Journals').check()
  await expect(page).toHaveURL(/\/products\/\?category=journals$/)
  await expect(page.getByTestId('product-card')).toHaveCount(1)
  expect(requests.at(-1)?.getAll('category[]')).toEqual(['journals'])

  await page.getByLabel('Beginner').check()
  await expect(page).toHaveURL(new RegExp(`category=journals&attr=${BEGINNER}$`))
  expect(requests.at(-1)?.getAll('attr_value_id[]')).toEqual([BEGINNER])

  // Same document: no full reload happened.
  expect(await page.evaluate(() => (window as unknown as { __marker: string }).__marker)).toBe('still-here')

  await page.getByRole('button', { name: 'Clear filters' }).click()
  await expect(page).toHaveURL(/\/products\/$/)
  await expect(page.getByTestId('product-card')).toHaveCount(3)

  // Search
  await page.getByLabel('Search products').fill('mandala')
  await page.getByRole('button', { name: 'Search' }).click()
  await expect(page).toHaveURL(/q=mandala/)
  await expect(page.getByTestId('product-card')).toHaveCount(1)

  // Back button restores the previous filter state
  await page.goBack()
  await expect(page).toHaveURL(/\/products\/$/)
  await expect(page.getByTestId('product-card')).toHaveCount(3)
})

test('direct navigation to a filtered URL pre-checks the sidebar', async ({ page }) => {
  const requests = await mockCatalogue(page)
  await page.goto(`/products/?category=journals&category=planners&attr=${MANDALA}&sort=price_asc`)

  await expect(page.getByLabel('Journals')).toBeChecked()
  await expect(page.getByLabel('Planners')).toBeChecked()
  await expect(page.getByLabel('Mandala')).toBeChecked()
  await expect(page.getByLabel('Beginner')).not.toBeChecked()
  await expect(page.getByLabel('Sort by')).toHaveValue('price_asc')

  await expect(page.getByTestId('product-card')).toHaveCount(1)
  const last = requests.at(-1)
  expect(last?.getAll('category[]')).toEqual(['journals', 'planners'])
  expect(last?.getAll('attr_value_id[]')).toEqual([MANDALA])
  expect(last?.get('sort_by')).toBe('price_asc')
  // The first request already carries the URL filters (no flash of unfiltered data).
  expect(requests[0].getAll('category[]')).toEqual(['journals', 'planners'])
})

test('mobile: filters collapse behind a toggle', async ({ page }) => {
  await page.setViewportSize({ width: 375, height: 800 })
  await mockCatalogue(page)
  await page.goto('/products/')
  await expect(page.getByTestId('product-card')).toHaveCount(3)

  await expect(page.getByTestId('filter-panel')).toBeHidden()
  await page.getByRole('button', { name: /^Filters/ }).click()
  await expect(page.getByTestId('filter-panel')).toBeVisible()
  await page.getByLabel('Journals').check()
  await expect(page.getByRole('button', { name: 'Filters (1)' })).toBeVisible()
})

test('shows an error state when the API fails', async ({ page }) => {
  await page.route('**/api/v1/categories', (r) => r.fulfill({ status: 200, contentType: 'application/json', body: '{"data":[]}' }))
  await page.route('**/api/v1/attributes', (r) => r.fulfill({ status: 200, contentType: 'application/json', body: '{"data":[]}' }))
  await page.route(/\/api\/v1\/products/, (r) => r.fulfill({ status: 500, contentType: 'application/json', body: '{"code":500,"message":"x"}' }))
  await page.goto('/products/')
  await expect(page.getByText('Products are temporarily unavailable.')).toBeVisible()
})
