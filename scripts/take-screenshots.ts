/**
 * Automated screenshot capture script using Playwright
 *
 * This script navigates through all key pages of the application
 * and captures high-quality screenshots for documentation.
 *
 * Run: npx tsx scripts/take-screenshots.ts
 */

import { chromium } from 'playwright'
import { mkdirSync } from 'fs'
import { join } from 'path'

const SCREENSHOTS_DIR = join(process.cwd(), 'screenshots')
const BASE_URL = 'http://localhost:5173'

// Ensure screenshots directory exists
mkdirSync(SCREENSHOTS_DIR, { recursive: true })

interface Screenshot {
  name: string
  path: string
  action?: (page: any) => Promise<void>
  viewport?: { width: number; height: number }
}

const screenshots: Screenshot[] = [
  // Home Page (Landing) - Full page with all sections
  {
    name: '01-home-hero',
    path: '/',
    viewport: { width: 1920, height: 1080 }
  },
  {
    name: '01-home-features',
    path: '/',
    viewport: { width: 1920, height: 1080 },
    action: async (page) => {
      await page.evaluate(() => window.scrollTo({ top: 900, behavior: 'smooth' }))
      await page.waitForTimeout(500)
    }
  },
  {
    name: '01-home-use-cases',
    path: '/',
    viewport: { width: 1920, height: 1080 },
    action: async (page) => {
      await page.evaluate(() => window.scrollTo({ top: 2400, behavior: 'smooth' }))
      await page.waitForTimeout(500)
    }
  },
  {
    name: '01-home-full',
    path: '/',
    viewport: { width: 1920, height: 1080 },
    action: async (page) => {
      // Take full page screenshot
      await page.evaluate(() => window.scrollTo({ top: 0, behavior: 'instant' }))
      await page.waitForTimeout(500)
    }
  },
  {
    name: '01-home-mobile',
    path: '/',
    viewport: { width: 375, height: 812 }
  },
  // Companies Database Page - List with table
  {
    name: '02-companies-list',
    path: '/companies',
    viewport: { width: 1920, height: 1080 }
  },
  {
    name: '02-companies-filters-open',
    path: '/companies',
    viewport: { width: 1920, height: 1080 },
    action: async (page) => {
      // Open filter sidebar
      await page.click('button:has-text("Filters")')
      await page.waitForTimeout(500)

      // Select Technology industry
      await page.click('button:has-text("Technology")')
      await page.waitForTimeout(300)
    }
  },
  {
    name: '02-companies-search',
    path: '/companies',
    viewport: { width: 1920, height: 1080 },
    action: async (page) => {
      // Type in search
      await page.fill('input[placeholder*="Search"]', 'Apple')
      await page.waitForTimeout(800) // Wait for debounce
    }
  },
  // Company Detail Page - Individual company view
  {
    name: '03-company-detail-top',
    path: '/companies/1',
    viewport: { width: 1920, height: 1080 }
  },
  {
    name: '03-company-detail-charts',
    path: '/companies/1',
    viewport: { width: 1920, height: 1080 },
    action: async (page) => {
      await page.evaluate(() => window.scrollTo(0, 800))
      await page.waitForTimeout(300)
    }
  },
  // Dark Mode Examples
  {
    name: '04-dark-mode-home',
    path: '/',
    viewport: { width: 1920, height: 1080 },
    action: async (page) => {
      // Toggle dark mode
      await page.click('button[aria-label*="dark"]')
      await page.waitForTimeout(300)
    }
  },
  {
    name: '04-dark-mode-companies',
    path: '/companies',
    viewport: { width: 1920, height: 1080 },
    action: async (page) => {
      // Toggle dark mode
      await page.click('button[aria-label*="dark"]')
      await page.waitForTimeout(300)
    }
  }
]

async function takeScreenshots() {
  console.log('🚀 Starting screenshot capture...\n')

  const browser = await chromium.launch()
  const context = await browser.newContext()

  for (const screenshot of screenshots) {
    console.log(`📸 Capturing: ${screenshot.name}`)

    const page = await context.newPage()

    // Set viewport
    if (screenshot.viewport) {
      await page.setViewportSize(screenshot.viewport)
    }

    // Navigate to page
    await page.goto(`${BASE_URL}${screenshot.path}`, {
      waitUntil: 'networkidle'
    })

    // Wait for page to be fully loaded
    await page.waitForTimeout(1000)

    // Execute custom action if provided
    if (screenshot.action) {
      await screenshot.action(page)
    }

    // Take screenshot (full page for home-full)
    const filename = `${screenshot.name}.png`
    const isFullPage = screenshot.name.includes('full')
    await page.screenshot({
      path: join(SCREENSHOTS_DIR, filename),
      fullPage: isFullPage
    })

    console.log(`   ✅ Saved: ${filename}`)

    await page.close()
  }

  await browser.close()

  console.log('\n✨ All screenshots captured successfully!')
  console.log(`📁 Location: ${SCREENSHOTS_DIR}`)
}

// Run the script
takeScreenshots().catch((error) => {
  console.error('❌ Error taking screenshots:', error)
  process.exit(1)
})
