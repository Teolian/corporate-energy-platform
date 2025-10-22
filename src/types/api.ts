import type { Company, EnergyForecast, Benchmark } from './models'

export interface ApiResponse<T> {
  data: T
  message?: string
  success: boolean
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  pageSize: number
  hasMore: boolean
}

export interface CompanyListParams {
  page?: number
  pageSize?: number
  search?: string
  industries?: string[]
  minEfficiency?: number
  maxEfficiency?: number
  minRenewable?: number
  maxRenewable?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}

export interface CompanyDetailResponse {
  company: Company
  forecasts: EnergyForecast[]
  benchmark: Benchmark
}
