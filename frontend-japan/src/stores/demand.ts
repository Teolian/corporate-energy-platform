import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Area, DemandResponse, DemandMetrics, ChartDataPoint } from '@/types/demand'
import { getDemandData, getCurrentDataMode } from '@/services/dataClient'

export const useDemandStore = defineStore('demand', () => {
  // State
  const tokyoData = ref<DemandResponse | null>(null)
  const kansaiData = ref<DemandResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const currentDate = ref(formatDate(new Date()))

  // Data mode
  const dataMode = ref(getCurrentDataMode())

  // Getters
  const tokyoMetrics = computed((): DemandMetrics | null => {
    if (!tokyoData.value) return null
    return calculateMetrics(tokyoData.value)
  })

  const kansaiMetrics = computed((): DemandMetrics | null => {
    if (!kansaiData.value) return null
    return calculateMetrics(kansaiData.value)
  })

  const tokyoChartData = computed((): ChartDataPoint[] => {
    if (!tokyoData.value) return []
    return transformToChartData(tokyoData.value)
  })

  const kansaiChartData = computed((): ChartDataPoint[] => {
    if (!kansaiData.value) return []
    return transformToChartData(kansaiData.value)
  })

  // Actions
  async function fetchDemandData(area: Area, date?: string) {
    loading.value = true
    error.value = null

    try {
      const targetDate = date || currentDate.value
      const data = await getDemandData(area, targetDate)

      if (area === 'tokyo') {
        tokyoData.value = data
      } else {
        kansaiData.value = data
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : `Failed to fetch ${area} demand`
      console.error(`Error fetching ${area} demand:`, err)
    } finally {
      loading.value = false
    }
  }

  async function fetchAllDemandData(date?: string) {
    const targetDate = date || currentDate.value
    await Promise.all([
      fetchDemandData('tokyo', targetDate),
      fetchDemandData('kansai', targetDate)
    ])
  }

  function setDate(date: string) {
    currentDate.value = date
  }

  function nextDay() {
    const d = new Date(currentDate.value)
    d.setDate(d.getDate() + 1)
    currentDate.value = formatDate(d)
    fetchAllDemandData()
  }

  function prevDay() {
    const d = new Date(currentDate.value)
    d.setDate(d.getDate() - 1)
    currentDate.value = formatDate(d)
    fetchAllDemandData()
  }

  return {
    // State
    tokyoData,
    kansaiData,
    loading,
    error,
    currentDate,
    dataMode,
    // Getters
    tokyoMetrics,
    kansaiMetrics,
    tokyoChartData,
    kansaiChartData,
    // Actions
    fetchDemandData,
    fetchAllDemandData,
    setDate,
    nextDay,
    prevDay
  }
})

// Helper functions
function formatDate(date: Date): string {
  return date.toISOString().split('T')[0]
}

function calculateMetrics(data: DemandResponse): DemandMetrics {
  const actualValues = data.series.map(p => p.demand_mw)
  const peak = Math.max(...actualValues)
  const average = actualValues.reduce((sum, val) => sum + val, 0) / actualValues.length

  // Calculate forecast accuracy if forecast data exists
  let forecastAccuracy: number | undefined
  const forecastValues = data.series.filter(p => p.forecast_mw !== null && p.forecast_mw !== undefined)

  if (forecastValues.length > 0) {
    const errors = forecastValues.map(p => {
      const error = Math.abs((p.forecast_mw! - p.demand_mw) / p.demand_mw)
      return error
    })
    const avgError = errors.reduce((sum, err) => sum + err, 0) / errors.length
    forecastAccuracy = Math.round((1 - avgError) * 100)
  }

  return { peak, average, forecastAccuracy }
}

function transformToChartData(data: DemandResponse): ChartDataPoint[] {
  return data.series.map(point => {
    const hour = new Date(point.ts).getHours()
    return {
      time: `${hour.toString().padStart(2, '0')}:00`,
      actual: point.demand_mw,
      forecast: point.forecast_mw ?? undefined
    }
  })
}
