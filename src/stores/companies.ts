import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Company, CompanyFilters, SortField, SortDirection } from '@/types/models'
import type { CompanyListParams } from '@/types/api'
import { companyService } from '@/services/api'
import { logger } from '@/utils/logger'
import { FILTER_RANGES } from '@/constants/filters'

export const useCompanyStore = defineStore('companies', () => {
  // State
  const companies = ref<Company[]>([])
  const currentCompany = ref<Company | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const total = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(30)

  // Filters
  const filters = ref<CompanyFilters>({
    search: '',
    industries: [],
    efficiencyRange: [FILTER_RANGES.EFFICIENCY.MIN, FILTER_RANGES.EFFICIENCY.MAX],
    renewableRange: [FILTER_RANGES.RENEWABLE.MIN, FILTER_RANGES.RENEWABLE.MAX],
    sizeRange: []
  })

  // Sorting
  const sortField = ref<SortField>('efficiencyScore')
  const sortDirection = ref<SortDirection>('desc')

  // Getters
  const hasMore = computed(() => currentPage.value * pageSize.value < total.value)
  const industries = computed(() => {
    const uniqueIndustries = new Set(companies.value.map((c) => c.industry))
    return Array.from(uniqueIndustries).sort()
  })

  // Actions
  async function fetchCompanies() {
    loading.value = true
    error.value = null

    try {
      const params: CompanyListParams = {
        page: currentPage.value,
        pageSize: pageSize.value,
        search: filters.value.search || undefined,
        industries: filters.value.industries.length > 0 ? filters.value.industries : undefined,
        minEfficiency: filters.value.efficiencyRange[0],
        maxEfficiency: filters.value.efficiencyRange[1],
        minRenewable: filters.value.renewableRange[0],
        maxRenewable: filters.value.renewableRange[1],
        sortBy: sortField.value,
        sortOrder: sortDirection.value
      }

      const response = await companyService.getCompanies(params)
      companies.value = response.data
      total.value = response.total
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch companies'
      logger.error('Error fetching companies', err, {
        page: currentPage.value,
        filters: filters.value
      })
    } finally {
      loading.value = false
    }
  }

  async function fetchCompanyById(id: string) {
    loading.value = true
    error.value = null

    try {
      const response = await companyService.getCompanyById(id)
      currentCompany.value = response.company
      return response
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch company details'
      logger.error('Error fetching company details', err, { companyId: id })
      throw err
    } finally {
      loading.value = false
    }
  }

  function setFilters(newFilters: Partial<CompanyFilters>) {
    filters.value = { ...filters.value, ...newFilters }
    currentPage.value = 1 // Reset to first page on filter change
  }

  function setSorting(field: SortField, direction: SortDirection) {
    sortField.value = field
    sortDirection.value = direction
    currentPage.value = 1
  }

  function resetFilters() {
    filters.value = {
      search: '',
      industries: [],
      efficiencyRange: [FILTER_RANGES.EFFICIENCY.MIN, FILTER_RANGES.EFFICIENCY.MAX],
      renewableRange: [FILTER_RANGES.RENEWABLE.MIN, FILTER_RANGES.RENEWABLE.MAX],
      sizeRange: []
    }
    currentPage.value = 1
  }

  function setPage(page: number) {
    currentPage.value = page
  }

  return {
    // State
    companies,
    currentCompany,
    loading,
    error,
    total,
    currentPage,
    pageSize,
    filters,
    sortField,
    sortDirection,
    // Getters
    hasMore,
    industries,
    // Actions
    fetchCompanies,
    fetchCompanyById,
    setFilters,
    setSorting,
    resetFilters,
    setPage
  }
})
