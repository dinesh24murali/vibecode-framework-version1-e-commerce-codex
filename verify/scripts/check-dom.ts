/**
 * check-dom.ts — DOM assertion script
 *
 * Runs Playwright checks against the running dev server.
 * Extend this file with project-specific assertions.
 *
 * Usage:
 *   npx ts-node verify/scripts/check-dom.ts
 *
 * Prerequisites:
 *   npm install -D playwright @playwright/test ts-node typescript
 *   npx playwright install chromium
 */

import { chromium, Browser, Page } from "playwright";

const BASE_URL = process.env.APP_URL ?? "http://localhost:3000";

interface CheckResult {
  name: string;
  passed: boolean;
  message?: string;
}

const results: CheckResult[] = [];

function pass(name: string) {
  results.push({ name, passed: true });
  console.log(`  ✓ ${name}`);
}

function fail(name: string, message: string) {
  results.push({ name, passed: false, message });
  console.error(`  ✗ ${name}: ${message}`);
}

async function checkPageLoad(page: Page) {
  console.log("\n[Page Load]");
  try {
    const response = await page.goto(BASE_URL, { waitUntil: "networkidle" });
    if (response && response.status() < 400) {
      pass("Home page returns 2xx/3xx");
    } else {
      fail("Home page returns 2xx/3xx", `Status: ${response?.status()}`);
    }
  } catch (e) {
    fail("Home page loads", String(e));
    return;
  }

  // Check for console errors
  const errors: string[] = [];
  page.on("console", (msg) => {
    if (msg.type() === "error") errors.push(msg.text());
  });

  await page.waitForTimeout(1000);

  if (errors.length === 0) {
    pass("No console errors on load");
  } else {
    fail("No console errors on load", errors.join("; "));
  }

  // Check page title is not empty
  const title = await page.title();
  if (title && title.length > 0) {
    pass(`Page title present: "${title}"`);
  } else {
    fail("Page title present", "Title is empty");
  }
}

async function checkNoHorizontalOverflow(page: Page) {
  console.log("\n[Layout]");
  try {
    await page.setViewportSize({ width: 375, height: 812 });
    await page.goto(BASE_URL, { waitUntil: "networkidle" });

    const overflow = await page.evaluate(() => {
      return document.documentElement.scrollWidth > window.innerWidth;
    });

    if (!overflow) {
      pass("No horizontal overflow on mobile (375px)");
    } else {
      fail(
        "No horizontal overflow on mobile (375px)",
        `scrollWidth=${await page.evaluate(() => document.documentElement.scrollWidth)} > innerWidth=375`
      );
    }
  } catch (e) {
    fail("Layout check", String(e));
  }
}

// TODO: Add project-specific checks below
// async function checkLoginPage(page: Page) { ... }
// async function checkDashboard(page: Page) { ... }

async function main() {
  console.log(`Running DOM checks against ${BASE_URL}...\n`);

  let browser: Browser | undefined;
  try {
    browser = await chromium.launch({ headless: true });
    const page = await browser.newPage();

    await checkPageLoad(page);
    await checkNoHorizontalOverflow(page);

    // TODO: Add more check functions here
  } finally {
    await browser?.close();
  }

  const passed = results.filter((r) => r.passed).length;
  const failed = results.filter((r) => !r.passed).length;

  console.log(`\n─────────────────────────────────`);
  console.log(`DOM checks: ${passed} passed, ${failed} failed`);

  if (failed > 0) {
    console.error("\nFailed checks:");
    results
      .filter((r) => !r.passed)
      .forEach((r) => console.error(`  - ${r.name}: ${r.message}`));
    process.exit(1);
  }
}

main().catch((e) => {
  console.error("Unexpected error:", e);
  process.exit(1);
});
