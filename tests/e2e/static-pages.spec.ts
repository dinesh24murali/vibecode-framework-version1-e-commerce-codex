import { test, expect } from '@playwright/test'

const PAGES = [
  { path: '/wall-of-fame/', title: 'Wall of Fame | e-commerce-site', h1: 'Wall of Fame' },
  { path: '/about/', title: 'About Us | e-commerce-site', h1: 'About Us' },
  { path: '/terms/', title: 'Terms & Conditions | e-commerce-site', h1: 'Terms & Conditions' },
  { path: '/privacy/', title: 'Privacy Policy | e-commerce-site', h1: 'Privacy Policy' },
]

for (const p of PAGES) {
  test(`${p.path} renders with unique title and no JS errors`, async ({ page }) => {
    const errors: string[] = []
    page.on('pageerror', (e) => errors.push(e.message))
    // Ignore failed resource loads: navbar links prefetch pages that are built in later tasks.
    page.on('console', (m) => m.type() === 'error' && !m.text().startsWith('Failed to load resource') && errors.push(m.text()))
    await page.goto(p.path)
    await expect(page).toHaveTitle(p.title)
    await expect(page.getByRole('heading', { level: 1, name: p.h1 })).toBeVisible()
    expect(errors).toEqual([])
  })
}

test('footer links reach the static pages', async ({ page }) => {
  await page.goto('/wall-of-fame/')
  await page.getByTestId('footer').getByRole('link', { name: 'About Us' }).click()
  await expect(page).toHaveURL(/\/about\/$/)
  await page.getByTestId('footer').getByRole('link', { name: 'Terms & Conditions' }).click()
  await expect(page).toHaveURL(/\/terms\/$/)
})
