import { test, expect } from '@playwright/test'

const json = (status: number, body: unknown) => ({
  status,
  contentType: 'application/json',
  body: JSON.stringify(body),
})

const adminResponse = {
  access_token: 'admin-access',
  refresh_token: 'admin-refresh',
  token_type: 'bearer',
  expires_in: 900,
  user: {
    id: '3fa85f64-5717-4562-b3fc-2c963f66afa6',
    email: 'admin@example.com',
    first_name: 'Ada',
    last_name: 'Admin',
    role: 'admin',
  },
}

test('unauthenticated /admin/ redirects to the admin login', async ({ page }) => {
  await page.goto('/admin/')
  await expect(page).toHaveURL(/\/admin\/login\/?$/)
  await expect(page.getByRole('heading', { name: 'Admin sign in' })).toBeVisible()
})

test('admin login page has no public site chrome', async ({ page }) => {
  await page.goto('/admin/login/')
  await expect(page.getByRole('navigation')).toHaveCount(0)
  await expect(page.getByRole('contentinfo')).toHaveCount(0)
})

test('admin login happy path lands on the dashboard', async ({ page }) => {
  await page.route('**/api/v1/admin/auth/login', (r) => r.fulfill(json(200, adminResponse)))
  await page.goto('/admin/login/')
  await page.getByLabel('Email').fill('admin@example.com')
  await page.getByLabel('Password').fill('Secure@1234')
  await page.getByRole('button', { name: 'Sign in' }).click()
  await expect(page).toHaveURL(/\/admin\/?$/)
  await expect(page.getByRole('heading', { name: 'Admin dashboard' })).toBeVisible()
})

test('invalid admin credentials show an error and stay on the login page', async ({ page }) => {
  await page.route('**/api/v1/admin/auth/login', (r) =>
    r.fulfill(json(401, { code: 401, message: 'invalid email or password' }))
  )
  await page.goto('/admin/login/')
  await page.getByLabel('Email').fill('x@example.com')
  await page.getByLabel('Password').fill('wrong')
  await page.getByRole('button', { name: 'Sign in' }).click()
  await expect(page.getByText('invalid email or password')).toBeVisible()
  await expect(page).toHaveURL(/\/admin\/login\/?$/)
})
