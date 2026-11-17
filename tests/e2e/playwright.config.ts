import { defineConfig } from '@playwright/test'

export default defineConfig({
  testDir: '.',
  use: { baseURL: process.env.APP_URL ?? 'http://static:3000' },
})
