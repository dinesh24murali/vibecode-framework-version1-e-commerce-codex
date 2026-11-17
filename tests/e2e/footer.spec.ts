import { test, expect, type Page } from '@playwright/test'

async function mockApi(page: Page) {
  const json = (body: unknown) => ({ status: 200, contentType: 'application/json', body: JSON.stringify(body) })
  await page.route('**/api/v1/landing', (r) => r.fulfill(json({ carousel_slides: [], featured_products: [], new_arrivals: [] })))
  await page.route('**/api/v1/categories', (r) => r.fulfill(json({ data: [] })))
  await page.route('**/api/v1/attributes', (r) => r.fulfill(json({ data: [] })))
  await page.route(/\/api\/v1\/products(\?.*)?$/, (r) =>
    r.fulfill(json({ data: [], meta: { total: 0, page: 1, limit: 12, total_pages: 0 }, filters: { categories: [], attributes: [], sort_options: [] } })))
}

for (const path of ['/', '/products/']) {
  test(`footer is present on ${path} with internal and social links`, async ({ page }) => {
    await mockApi(page)
    await page.goto(path)
    const footer = page.getByTestId('footer')
    await expect(footer).toBeVisible()
    await expect(footer.getByRole('link', { name: 'Terms & Conditions' })).toHaveAttribute('href', '/terms/')
    await expect(footer.getByRole('link', { name: 'About Us' })).toHaveAttribute('href', '/about/')
    await expect(footer.getByRole('link', { name: 'Wall of Fame' })).toHaveAttribute('href', '/wall-of-fame/')
    for (const name of ['Facebook', 'Instagram', 'WhatsApp']) {
      const link = footer.getByRole('link', { name })
      await expect(link).toHaveAttribute('target', '_blank')
      await expect(link).toHaveAttribute('rel', /noopener/)
    }
  })
}

test('social icons open in a new tab', async ({ page, context }) => {
  await mockApi(page)
  await context.route(/facebook\.com/, (r) => r.fulfill({ status: 200, body: 'ok' }))
  await page.goto('/')
  const [popup] = await Promise.all([
    context.waitForEvent('page'),
    page.getByTestId('footer').getByRole('link', { name: 'Facebook' }).click(),
  ])
  expect(popup.url()).toContain('facebook.com')
})
