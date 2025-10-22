export interface Company {
  id: string
  name: string
  industry: string
  employees: number
  energyConsumption: number // MWh/year
  efficiencyScore: number // 0-100
  renewablePercentage: number // 0-100
  location: string
  tags: string[]
  trend: number // percentage change
  costPerEmployee?: number
}

export interface EnergyForecast {
  companyId: string
  month: string
  predicted: number
  actual?: number
  accuracy?: number
  confidence: number
}

export interface Benchmark {
  industry: string
  avgConsumption: number
  avgEfficiency: number
  topPerformers: Company[]
}

export interface CompanyFilters {
  search: string
  industries: string[]
  efficiencyRange: [number, number]
  renewableRange: [number, number]
  sizeRange: string[]
}

export type SortField = 'name' | 'industry' | 'energyConsumption' | 'efficiencyScore' | 'renewablePercentage' | 'trend'
export type SortDirection = 'asc' | 'desc'
