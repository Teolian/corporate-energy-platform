// TypeScript types matching backend JSON structure (AGENT_TECH_SPEC §3.1)

export type Area = 'tokyo' | 'kansai'

export type Timescale = 'hourly'

export interface Source {
  name: string
  url: string
}

export interface SeriesPoint {
  ts: string // ISO8601 timestamp
  demand_mw: number
  forecast_mw?: number | null
}

export interface Meta {
  warning?: string
}

// Backend JSON response structure
export interface DemandResponse {
  area: Area
  date: string // YYYY-MM-DD
  timezone: string // Always "Asia/Tokyo"
  timescale: Timescale
  series: SeriesPoint[]
  source: Source
  meta?: Meta
}

// UI-friendly chart data
export interface ChartDataPoint {
  time: string // HH:00 format
  actual: number
  forecast?: number
}

export interface DemandMetrics {
  peak: number
  average: number
  forecastAccuracy?: number
}

export type DataMode = 'mock' | 'live'
