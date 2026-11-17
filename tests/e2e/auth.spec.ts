import { test, expect, type Page } from '@playwright/test'

const tokens = {
  access_token: 'access-abc',
  refresh_token: 'refresh-abc',
  token_type: 'bearer',
  expires_in: 900,
}
const json = (status: number, body: unknown) => ({
  status,
  contentType: 'application/json',
  body: JSON.stringify(body),
})

// The static verify stack has no API; the network is mocked at the browser level.
async function mockAuth(page: Page) {
  await page.route('**/api/v1/auth/register/init', (r) => r.fulfill(json(200, { message: 'sent' })))
  await page.route('**/api/v1/auth/register/verify-otp', (r) => {
    const { otp } = r.request().postDataJSON()
    return r.fulfill(
      otp === '123456'
        ? json(200, { otp_token: 'otp-tok', expires_in: 600 })
        : json(400, { code: 400, message: 'Invalid OTP' })
    )
  })
  await page.route('**/api/v1/auth/register/complete', (r) => r.fulfill(json(201, tokens)))
  await page.route('**/api/v1/auth/login', (r) => r.fulfill(json(200, tokens)))
}

test('signup happy path; invalid OTP keeps the input', async ({ page }) => {
  await mockAuth(page)
  await page.goto('/signup/')
  await page.getByLabel('Email').fill('alice@example.com')
  await page.getByRole('button', { name: 'Send verification code' }).click()
  await expect(page).toHaveURL(/verify-otp/)

  await page.getByLabel('6-digit code').fill('000000')
  await page.getByRole('button', { name: 'Verify' }).click()
  await expect(page.getByText('Invalid OTP')).toBeVisible()
  await expect(page.getByLabel('6-digit code')).toHaveValue('000000')

  await page.getByLabel('6-digit code').fill('123456')
  await page.getByRole('button', { name: 'Verify' }).click()
  await page.getByLabel('First name').fill('Alice')
  await page.getByLabel('Last name').fill('Doe')
  await page.getByLabel('Password').fill('Secure@1234')
  await page.getByRole('button', { name: 'Create account' }).click()
  await expect(page).toHaveURL(/\/$/)
})

test('login honours ?redirect=', async ({ page }) => {
  await mockAuth(page)
  await page.goto('/login/?redirect=/signup/')
  await page.getByLabel('Email').fill('alice@example.com')
  await page.getByLabel('Password').fill('Secure@1234')
  await page.getByRole('button', { name: 'Log in' }).click()
  await expect(page).toHaveURL(/\/signup\/$/)
})
