import axios from 'axios'
import type { CompanyListParams, PaginatedResponse, CompanyDetailResponse } from '@/types/api'
import type { Company } from '@/types/models'
import { realCompaniesData } from '@/data/companies'

/**
 * API client configuration for future backend integration.
 * Currently using mock data, but ready to switch to real API.
 *
 * @example
 * ```ts
 * const response = await apiClient.get('/companies', { params })
 * ```
 */
export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api',
  headers: {
    'Content-Type': 'application/json'
  },
  timeout: 10000
})

// Real company data from Fortune 500
const mockCompanies: Company[] = realCompaniesData

export const companyService = {
  async getCompanies(params: CompanyListParams): Promise<PaginatedResponse<Company>> {
    // TODO: Replace with real API call
    // const response = await apiClient.get('/companies', { params })
    // return response.data

    // Mock implementation
    await new Promise((resolve) => setTimeout(resolve, 500))

    let filtered = [...mockCompanies]

    // Apply search filter
    if (params.search) {
      const search = params.search.toLowerCase()
      filtered = filtered.filter(
        (c) =>
          c.name.toLowerCase().includes(search) ||
          c.industry.toLowerCase().includes(search) ||
          c.location.toLowerCase().includes(search)
      )
    }

    // Apply industry filter
    if (params.industries?.length) {
      const industries = params.industries
      filtered = filtered.filter((c) => industries.includes(c.industry))
    }

    // Apply efficiency range filter
    const { minEfficiency, maxEfficiency, minRenewable, maxRenewable } = params

    if (minEfficiency !== undefined) {
      filtered = filtered.filter((c) => c.efficiencyScore >= minEfficiency)
    }
    if (maxEfficiency !== undefined) {
      filtered = filtered.filter((c) => c.efficiencyScore <= maxEfficiency)
    }

    // Apply renewable range filter
    if (minRenewable !== undefined) {
      filtered = filtered.filter((c) => c.renewablePercentage >= minRenewable)
    }
    if (maxRenewable !== undefined) {
      filtered = filtered.filter((c) => c.renewablePercentage <= maxRenewable)
    }

    // Apply sorting
    if (params.sortBy) {
      filtered.sort((a, b) => {
        const aVal = a[params.sortBy as keyof Company] as number | string
        const bVal = b[params.sortBy as keyof Company] as number | string
        if (aVal === undefined || bVal === undefined) return 0
        const comparison = aVal > bVal ? 1 : aVal < bVal ? -1 : 0
        return params.sortOrder === 'desc' ? -comparison : comparison
      })
    }

    // Apply pagination
    const page = params.page || 1
    const pageSize = params.pageSize || 30
    const start = (page - 1) * pageSize
    const paginated = filtered.slice(start, start + pageSize)

    return {
      data: paginated,
      total: filtered.length,
      page,
      pageSize,
      hasMore: start + pageSize < filtered.length
    }
  },

  async getCompanyById(id: string): Promise<CompanyDetailResponse> {
    // TODO: Replace with real API call
    // const response = await apiClient.get(`/companies/${id}`)
    // return response.data

    // Mock implementation
    await new Promise((resolve) => setTimeout(resolve, 300))

    const company = mockCompanies.find((c) => c.id === id)
    if (!company) {
      throw new Error('Company not found')
    }

    // Generate mock historical data for charts
    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
    const forecasts = months.map((month, index) => ({
      companyId: id,
      month,
      predicted: 95 + Math.random() * 20 + index * 1.5,
      actual: index < 10 ? 95 + Math.random() * 20 + index * 1.5 : undefined,
      accuracy: index < 10 ? 96 + Math.random() * 3 : undefined,
      confidence: 92 + Math.random() * 5
    }))

    return {
      company,
      forecasts,
      benchmark: {
        industry: company.industry,
        avgConsumption: 1500,
        avgEfficiency: 76.5,
        topPerformers: mockCompanies.slice(0, 3)
      }
    }
  }
}
