import { test, expect } from '@playwright/test'

async function fillValid(page: import('@playwright/test').Page) {
  await page.getByLabel('Name', { exact: true }).fill('Bob Kumar')
  await page.getByLabel('Email').fill('bob@example.com')
  await page.getByLabel('Phone (optional)').fill('+919876500000')
  await page.getByLabel('Subject').fill('Bulk order')
  await page.getByLabel('Message').fill('50 journals please')
}

test('shows address, map and social links', async ({ page }) => {
  await page.goto('/contact/')
  await expect(page).toHaveTitle('Contact Us | e-commerce-site')
  await expect(page.getByTestId('business-address')).toBeVisible()
  await expect(page.locator('iframe[src*="google.com/maps"]')).toHaveCount(1)
  await expect(page.getByRole('list', { name: 'Social media' }).first().getByRole('link')).toHaveCount(3)
})

test('happy path: submits the enquiry and shows success', async ({ page }) => {
  let payload: Record<string, unknown> | null = null
  await page.route('**/api/v1/contact', (r) => {
    payload = r.request().postDataJSON()
    return r.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ message: 'ok' }) })
  })
  await page.goto('/contact/')
  await fillValid(page)
  await page.getByRole('button', { name: 'Send message' }).click()
  await expect(page.getByRole('status')).toContainText('Your message has been sent')
  expect(payload).toMatchObject({ name: 'Bob Kumar', email: 'bob@example.com', subject: 'Bulk order', phone: '+919876500000' })
  await expect(page.getByLabel('Name', { exact: true })).toHaveValue('')
})

test('client-side validation blocks invalid input', async ({ page }) => {
  let called = false
  await page.route('**/api/v1/contact', (r) => {
    called = true
    return r.fulfill({ status: 200, body: '{}' })
  })
  await page.goto('/contact/')
  await page.getByRole('button', { name: 'Send message' }).click()
  await expect(page.getByText('Name is required')).toBeVisible()
  await expect(page.getByText('Email is required')).toBeVisible()
  await expect(page.getByText('Subject is required')).toBeVisible()
  await expect(page.getByText('Message is required')).toBeVisible()

  await page.getByLabel('Email').fill('nope')
  await page.getByLabel('Phone (optional)').fill('abc')
  await page.getByRole('button', { name: 'Send message' }).click()
  await expect(page.getByText('Enter a valid email address')).toBeVisible()
  await expect(page.getByText('Enter a valid phone number')).toBeVisible()
  expect(called).toBe(false)
})

test('shows an error message when the API fails', async ({ page }) => {
  await page.route('**/api/v1/contact', (r) =>
    r.fulfill({ status: 500, contentType: 'application/json', body: JSON.stringify({ code: 500, message: 'internal error' }) }))
  await page.goto('/contact/')
  await fillValid(page)
  await page.getByRole('button', { name: 'Send message' }).click()
  await expect(page.getByRole('alert').filter({ hasText: 'could not send' })).toBeVisible()
})
