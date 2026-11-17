import { test, expect, type Page } from '@playwright/test'

const json = (status: number, body: unknown) => ({
  status,
  contentType: 'application/json',
  body: JSON.stringify(body),
})

const now = '2026-02-01T12:00:00Z'
const profile = {
  id: 'u1', email: 'alice@example.com', first_name: 'Alice', last_name: 'Doe',
  role: 'customer', created_at: now, updated_at: now,
}
const addr = {
  id: 'a1', address_line1: '42 Rose St', address_line2: null, landmark: 'Mall', city: 'Chennai',
  state: 'TN', pincode: '600001', company_name: 'Acme', gst_number: '33AAAAA0000A1Z5',
  phone: '+919876543210', is_default: true, created_at: now, updated_at: now,
}
const order = (id: string, n: string, status: string) => ({
  id, order_number: n, status, items: [], shipping_address: addr, subtotal: 100, coupon_discount: 0,
  taxable_amount: 100, gst_amount: 5, total: 105, created_at: now, updated_at: now,
})

async function login(page: Page) {
  await page.route('**/api/v1/auth/login', (r) =>
    r.fulfill(json(200, { access_token: 'tok', refresh_token: 'ref', token_type: 'bearer', expires_in: 900 }))
  )
  await page.goto('/login/?redirect=/account/')
  await page.getByLabel('Email').fill('alice@example.com')
  await page.getByLabel('Password').fill('Secure@1234')
  await page.getByRole('button', { name: 'Log in' }).click()
  await expect(page.getByRole('heading', { name: 'My account' })).toBeVisible()
}

test('unauthenticated /account redirects to login with redirect param', async ({ page }) => {
  await page.goto('/account/')
  await expect(page).toHaveURL(/\/login\/\?redirect=\/account\/$/)
})

test('profile, address and orders tabs', async ({ page }) => {
  let profileBody: Record<string, unknown> | null = null
  let addressBody: Record<string, unknown> | null = null
  let cancelled = false

  await page.route('**/api/v1/account/profile', (r) => {
    if (r.request().method() === 'PUT') {
      profileBody = r.request().postDataJSON()
      return r.fulfill(json(200, { ...profile, first_name: profileBody?.first_name }))
    }
    return r.fulfill(json(200, profile))
  })
  await page.route('**/api/v1/account/profile/email-change/request', (r) =>
    r.fulfill(json(200, { message: 'sent' }))
  )
  await page.route('**/api/v1/account/profile/email-change/confirm', (r) =>
    r.fulfill(json(200, { ...profile, email: 'new@example.com' }))
  )
  await page.route('**/api/v1/account/addresses', (r) => {
    if (r.request().method() === 'POST') {
      addressBody = r.request().postDataJSON()
      return r.fulfill(json(201, addr))
    }
    return r.fulfill(json(200, { data: addressBody ? [addr] : [] }))
  })
  await page.route(/\/api\/v1\/orders\/o1\/cancel$/, (r) => {
    cancelled = true
    return r.fulfill(json(200, order('o1', 'ORD-1', 'cancelled')))
  })
  await page.route(/\/api\/v1\/orders\?/, (r) => {
    const p = new URL(r.request().url()).searchParams
    const meta = { total: 12, page: Number(p.get('page')), limit: 10, total_pages: 2 }
    if (p.get('limit') !== '10') return r.fulfill(json(400, { code: 400, message: 'bad limit' }))
    const data =
      meta.page === 1
        ? [order('o1', 'ORD-1', cancelled ? 'cancelled' : 'open'), order('o2', 'ORD-2', 'fulfilled')]
        : [order('o3', 'ORD-3', 'cancelled')]
    return r.fulfill(json(200, { data, meta }))
  })

  await login(page)

  // Profile tab
  await expect(page.getByLabel('First name')).toHaveValue('Alice')
  await page.getByLabel('First name').fill('Alicia')
  await page.getByRole('button', { name: 'Save profile' }).click()
  await expect(page.getByText('Profile saved')).toBeVisible()
  expect(profileBody).toMatchObject({ first_name: 'Alicia', last_name: 'Doe' })

  await page.getByLabel('New email').fill('new@example.com')
  await page.getByRole('button', { name: 'Send verification code' }).click()
  await page.getByLabel('Verification code').fill('123456')
  await page.getByRole('button', { name: 'Confirm email change' }).click()
  await expect(page.getByText('Email updated')).toBeVisible()

  // Address tab: all fields incl. GST are sent
  await page.getByRole('tab', { name: 'Addresses' }).click()
  await page.getByRole('button', { name: 'Add address' }).click()
  await page.getByLabel('Address line 1').fill('42 Rose St')
  await page.getByLabel('Landmark').fill('Mall')
  await page.getByLabel('City').fill('Chennai')
  await page.getByLabel('State').fill('TN')
  await page.getByLabel('Pincode').fill('600001')
  await page.getByLabel('Company name').fill('Acme')
  await page.getByLabel('GST number').fill('33AAAAA0000A1Z5')
  await page.getByLabel('Phone').fill('+919876543210')
  await page.getByRole('button', { name: 'Save address' }).click()
  await expect(page.getByTestId('address-row')).toContainText('GST: 33AAAAA0000A1Z5')
  expect(addressBody).toMatchObject({ gst_number: '33AAAAA0000A1Z5', company_name: 'Acme', landmark: 'Mall' })

  // Orders tab
  await page.getByRole('tab', { name: 'Orders' }).click()
  const rows = page.getByTestId('order-row')
  await expect(rows).toHaveCount(2)
  await expect(rows.nth(0).getByRole('button', { name: 'Cancel' })).toBeVisible()
  await expect(rows.nth(1).getByRole('button', { name: 'Cancel' })).toHaveCount(0)
  await expect(page.getByText('Page 1 of 2')).toBeVisible()

  await rows.nth(0).getByRole('button', { name: 'Cancel' }).click()
  await expect(rows.nth(0)).toContainText('cancelled')
  expect(cancelled).toBe(true)

  await page.getByRole('button', { name: 'Next' }).click()
  await expect(page.getByText('Page 2 of 2')).toBeVisible()
  await expect(rows).toHaveCount(1)
})
