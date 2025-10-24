// Data client with dual-mode support (MOCK/LIVE)
// Per AGENT_INTEGRATION_GUIDE.md §2-3

import type { Area, DemandResponse, DataMode } from '@/types/demand'

const STORAGE_KEY = 'jp-energy-data-mode'

// Get data mode from URL query, localStorage, or env
function getDataMode(): DataMode {
  // 1. Check URL query ?mode=live or ?mode=mock
  const urlParams = new URLSearchParams(window.location.search)
  const urlMode = urlParams.get('mode')
  if (urlMode === 'live' || urlMode === 'mock') {
    // Persist to localStorage
    localStorage.setItem(STORAGE_KEY, urlMode)
    return urlMode
  }

  // 2. Check localStorage
  const storedMode = localStorage.getItem(STORAGE_KEY)
  if (storedMode === 'live' || storedMode === 'mock') {
    return storedMode
  }

  // 3. Check env at build time
  const envMode = import.meta.env.VITE_DATA_MODE
  if (envMode === 'live' || envMode === 'mock') {
    return envMode
  }

  // 4. Default to mock
  return 'mock'
}

// Generate mock data matching backend JSON structure
function generateMockData(area: Area, date: string): DemandResponse {
  const baseValues = area === 'tokyo' ?
    { min: 25000, max: 40000 } :
    { min: 12000, max: 19000 }

  const series = Array.from({ length: 24 }, (_, hour) => {
    const hourFactor = Math.sin((hour - 6) * Math.PI / 12) * 0.3 + 0.7
    const randomness = Math.random() * 0.1 + 0.95
    const actual = Math.round((baseValues.min + (baseValues.max - baseValues.min) * hourFactor) * randomness)
    const forecast = Math.round(actual * (Math.random() * 0.1 + 0.95))

    return {
      ts: `${date}T${hour.toString().padStart(2, '0')}:00:00+09:00`,
      demand_mw: actual,
      forecast_mw: forecast
    }
  })

  return {
    area,
    date,
    timezone: 'Asia/Tokyo',
    timescale: 'hourly',
    series,
    source: {
      name: area === 'tokyo' ? 'TEPCO (Mock)' : 'Kansai (Mock)',
      url: 'https://example.com/mock'
    }
  }
}

// Fetch live data from static JSON files
async function fetchLiveData(area: Area, date: string): Promise<DemandResponse> {
  const response = await fetch(`/data/jp/${area}/demand-${date}.json`)

  if (!response.ok) {
    throw new Error(`Failed to fetch ${area} demand for ${date}: ${response.statusText}`)
  }

  const data: DemandResponse = await response.json()
  return data
}

// Main data fetcher with fallback
export async function getDemandData(area: Area, date: string): Promise<DemandResponse> {
  const mode = getDataMode()

  try {
    if (mode === 'live') {
      try {
        return await fetchLiveData(area, date)
      } catch (error) {
        console.warn(`[live-data] Failed to fetch live data for ${area}/${date}, falling back to mock`, error)
        // Graceful fallback to mock
        const mockData = generateMockData(area, date)
        mockData.meta = {
          warning: 'Live data unavailable, showing mock data'
        }
        return mockData
      }
    }
  } catch {
    // Mode check failed, use mock
  }

  return generateMockData(area, date)
}

// Get current data mode (for UI display)
export function getCurrentDataMode(): DataMode {
  return getDataMode()
}

// Set data mode programmatically
export function setDataMode(mode: DataMode): void {
  localStorage.setItem(STORAGE_KEY, mode)
  window.location.reload() // Reload to apply new mode
}
