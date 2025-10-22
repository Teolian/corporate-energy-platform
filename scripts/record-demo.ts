/**
 * Automated demo video recording script using Playwright
 *
 * Records a video walkthrough of the key features of the application.
 * The video can be converted to GIF for README.
 *
 * Run: npx tsx scripts/record-demo.ts
 */

import { chromium } from 'playwright'
import { mkdirSync } from 'fs'
import { join } from 'path'

const VIDEOS_DIR = join(process.cwd(), 'demo')
const BASE_URL = 'http://localhost:5173'

// Ensure videos directory exists
mkdirSync(VIDEOS_DIR, { recursive: true })

async function recordDemo() {
  console.log('🎬 Starting demo recording...\n')

  const browser = await chromium.launch({
    headless: true
  })

  const context = await browser.newContext({
    viewport: { width: 1920, height: 1080 },
    recordVideo: {
      dir: VIDEOS_DIR,
      size: { width: 1920, height: 1080 }
    }
  })

  const page = await context.newPage()

  try {
    console.log('📍 Step 1: Home page')
    await page.goto(BASE_URL, { waitUntil: 'networkidle' })
    await page.waitForTimeout(2000)

    // Scroll to features section
    await page.evaluate(() => window.scrollTo({ top: 800, behavior: 'smooth' }))
    await page.waitForTimeout(2000)

    console.log('📍 Step 2: Navigate to Companies list')
    await page.click('text=Explore Database')
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(2000)

    console.log('📍 Step 3: Search for companies')
    await page.fill('input[placeholder*="Search"]', 'Technology')
    await page.waitForTimeout(1500)

    console.log('📍 Step 4: Open filters')
    await page.click('button:has-text("Filters")')
    await page.waitForTimeout(1000)

    console.log('📍 Step 5: Apply filter')
    await page.click('button:has-text("Technology")')
    await page.waitForTimeout(1000)

    // Close sidebar
    await page.keyboard.press('Escape')
    await page.waitForTimeout(1000)

    console.log('📍 Step 6: View company details')
    // Click on first company row
    const firstRow = page.locator('tbody tr').first()
    await firstRow.click()
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(2000)

    // Scroll to see charts
    await page.evaluate(() => window.scrollTo({ top: 600, behavior: 'smooth' }))
    await page.waitForTimeout(2000)

    console.log('📍 Step 7: Back to list and toggle dark mode')
    await page.goBack()
    await page.waitForLoadState('networkidle')
    await page.waitForTimeout(1000)

    // Toggle dark mode
    await page.click('button[aria-label*="dark"]')
    await page.waitForTimeout(2000)

    console.log('📍 Step 8: Select companies for comparison')
    // Select first 3 companies
    const checkboxes = page.locator('input[type="checkbox"]')
    await checkboxes.nth(0).click()
    await page.waitForTimeout(500)
    await checkboxes.nth(1).click()
    await page.waitForTimeout(500)
    await checkboxes.nth(2).click()
    await page.waitForTimeout(1500)

    console.log('✅ Demo recording completed!')
  } catch (error) {
    console.error('❌ Error during recording:', error)
  } finally {
    await page.close()
    await context.close()
    await browser.close()
  }

  console.log('\n✨ Video saved!')
  console.log(`📁 Location: ${VIDEOS_DIR}`)
  console.log('\n💡 To convert to GIF, use:')
  console.log('   ffmpeg -i demo/video.webm -vf "fps=10,scale=1280:-1:flags=lanczos" demo/demo.gif')
}

// Run the script
recordDemo().catch((error) => {
  console.error('❌ Error recording demo:', error)
  process.exit(1)
})
