import { chromium } from 'playwright'
import * as path from 'path'
import * as fs from 'fs'

const SCREENSHOTS_DIR = path.resolve(process.cwd(), 'screenshots')
const BASE_URL = 'http://localhost:5173'

async function main() {
  // Ensure screenshots directory exists
  if (!fs.existsSync(SCREENSHOTS_DIR)) {
    fs.mkdirSync(SCREENSHOTS_DIR, { recursive: true })
  }

  const browser = await chromium.launch()

  console.log('📸 Taking screenshots...\n')

  // Light mode screenshots
  console.log('=== LIGHT MODE ===')
  const lightPage = await browser.newPage({
    viewport: { width: 1920, height: 1080 }
  })

  // 1. Home page (landing)
  console.log('1/6 Home page (light)...')
  await lightPage.goto(`${BASE_URL}/`, { waitUntil: 'networkidle' })
  await lightPage.screenshot({
    path: path.join(SCREENSHOTS_DIR, '01-home.png'),
    fullPage: true
  })
  console.log('✓ Saved: 01-home.png')

  // 2. Companies list (table)
  console.log('2/6 Companies table (light)...')
  await lightPage.goto(`${BASE_URL}/companies`, { waitUntil: 'networkidle' })
  await lightPage.waitForTimeout(1000)
  await lightPage.screenshot({
    path: path.join(SCREENSHOTS_DIR, '02-companies.png'),
    fullPage: true
  })
  console.log('✓ Saved: 02-companies.png')

  // 3. Company detail
  console.log('3/6 Company detail (light)...')
  await lightPage.goto(`${BASE_URL}/companies/1`, { waitUntil: 'networkidle' })
  await lightPage.waitForTimeout(800)
  await lightPage.screenshot({
    path: path.join(SCREENSHOTS_DIR, '03-company-detail.png'),
    fullPage: true
  })
  console.log('✓ Saved: 03-company-detail.png')

  await lightPage.close()

  // Dark mode screenshots (viewport only)
  console.log('\n=== DARK MODE ===')
  const darkPage = await browser.newPage({
    viewport: { width: 1920, height: 1080 }
  })

  // Set dark mode in localStorage before loading pages
  await darkPage.goto(BASE_URL)
  await darkPage.evaluate(() => {
    localStorage.setItem('darkMode', 'true')
    document.documentElement.classList.add('dark')
  })

  // 4. Home page dark
  console.log('4/6 Home page (dark)...')
  await darkPage.goto(`${BASE_URL}/`, { waitUntil: 'networkidle' })
  await darkPage.waitForTimeout(500)
  // Ensure dark mode is applied
  await darkPage.evaluate(() => {
    document.documentElement.classList.add('dark')
  })
  await darkPage.waitForTimeout(200)
  await darkPage.screenshot({
    path: path.join(SCREENSHOTS_DIR, '04-home-dark.png'),
    fullPage: false // Viewport only
  })
  console.log('✓ Saved: 04-home-dark.png')

  // 5. Companies dark
  console.log('5/6 Companies table (dark)...')
  await darkPage.goto(`${BASE_URL}/companies`, { waitUntil: 'networkidle' })
  await darkPage.waitForTimeout(1000)
  // Ensure dark mode is applied
  await darkPage.evaluate(() => {
    document.documentElement.classList.add('dark')
  })
  await darkPage.waitForTimeout(200)
  await darkPage.screenshot({
    path: path.join(SCREENSHOTS_DIR, '05-companies-dark.png'),
    fullPage: false // Viewport only
  })
  console.log('✓ Saved: 05-companies-dark.png')

  // 6. Company detail dark
  console.log('6/6 Company detail (dark)...')
  await darkPage.goto(`${BASE_URL}/companies/1`, { waitUntil: 'networkidle' })
  await darkPage.waitForTimeout(800)
  // Ensure dark mode is applied
  await darkPage.evaluate(() => {
    document.documentElement.classList.add('dark')
  })
  await darkPage.waitForTimeout(200)
  await darkPage.screenshot({
    path: path.join(SCREENSHOTS_DIR, '06-company-detail-dark.png'),
    fullPage: false // Viewport only
  })
  console.log('✓ Saved: 06-company-detail-dark.png')

  await darkPage.close()
  await browser.close()

  console.log('\n✅ All 6 screenshots saved to screenshots/')
}

main().catch((error) => {
  console.error('Error taking screenshots:', error)
  process.exit(1)
})
