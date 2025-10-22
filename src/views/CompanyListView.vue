<script setup lang="ts">
import { onMounted, watch, ref } from 'vue'
import { watchDebounced } from '@vueuse/core'
import { useCompanyStore } from '@/stores/companies'
import { RouterLink, useRouter } from 'vue-router'
import ThemeToggle from '@/components/ThemeToggle.vue'
import FilterSidebar from '@/components/FilterSidebar.vue'
import ActiveFiltersChips from '@/components/ActiveFiltersChips.vue'
import type { CompanyFilters as FiltersType, SortField, SortDirection } from '@/types/models'
import { TrendingUp, TrendingDown, Minus, Filter, ArrowUpDown, ArrowUp, ArrowDown, GitCompare, X, Search, SearchX } from 'lucide-vue-next'

const store = useCompanyStore()
const router = useRouter()
const isFilterSidebarOpen = ref(false)
const selectedCompanies = ref<string[]>([])
const searchQuery = ref(store.filters.search)

function handleSort(field: SortField) {
  if (store.sortField === field) {
    // Toggle direction if same field
    const newDirection: SortDirection = store.sortDirection === 'asc' ? 'desc' : 'asc'
    store.setSorting(field, newDirection)
  } else {
    // Default to desc for new field
    store.setSorting(field, 'desc')
  }
}

function toggleCompanySelection(companyId: string) {
  const index = selectedCompanies.value.indexOf(companyId)
  if (index > -1) {
    selectedCompanies.value.splice(index, 1)
  } else {
    if (selectedCompanies.value.length < 4) {
      selectedCompanies.value.push(companyId)
    }
  }
}

function clearSelection() {
  selectedCompanies.value = []
}

function goToCompare() {
  if (selectedCompanies.value.length < 2) {
    return // Already handled by disabled state
  }

  router.push({
    name: 'compare',
    query: { ids: selectedCompanies.value.join(',') }
  })
}

function handleFilterUpdate(filters: Partial<FiltersType>) {
  store.setFilters(filters)
}

function handleFilterReset() {
  store.resetFilters()
}

function handleRemoveIndustry(industry: string) {
  const newIndustries = store.filters.industries.filter((i) => i !== industry)
  store.setFilters({
    industries: newIndustries
  })
}

function handleResetSearch() {
  store.setFilters({
    search: ''
  })
}

function handleResetEfficiency() {
  store.setFilters({
    efficiencyRange: [0, 100]
  })
}

function handleResetRenewable() {
  store.setFilters({
    renewableRange: [0, 100]
  })
}

onMounted(() => {
  store.fetchCompanies()
})

watch(
  () => [store.filters, store.sortField, store.sortDirection, store.currentPage],
  () => {
    store.fetchCompanies()
  },
  { deep: true }
)

// Debounced search with automatic cleanup on unmount
watchDebounced(
  searchQuery,
  (newValue) => {
    store.setFilters({ search: newValue })
  },
  { debounce: 300 }
)

function getEfficiencyBadge(score: number) {
  if (score >= 90) return { text: 'Elite', class: 'bg-green-100 text-green-800' }
  if (score >= 80) return { text: 'High', class: 'bg-blue-100 text-blue-800' }
  if (score >= 70) return { text: 'Medium', class: 'bg-yellow-100 text-yellow-800' }
  return { text: 'Low', class: 'bg-red-100 text-red-800' }
}

function getRenewableBadge(percentage: number) {
  if (percentage >= 80) return { text: 'Leader', class: 'bg-energy-green text-white' }
  if (percentage >= 50) return { text: 'Good', class: 'bg-green-100 text-green-800' }
  return { text: 'Fair', class: 'bg-yellow-100 text-yellow-800' }
}

function getIndustryColor(industry: string): string {
  const colors: Record<string, string> = {
    'Technology': 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 border border-blue-200 dark:border-blue-800',
    'Automotive': 'bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-300 border border-purple-200 dark:border-purple-800',
    'Energy': 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-700 dark:text-yellow-300 border border-yellow-200 dark:border-yellow-800',
    'Manufacturing': 'bg-gray-100 dark:bg-gray-700/50 text-gray-700 dark:text-gray-300 border border-gray-200 dark:border-gray-600',
    'Finance': 'bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300 border border-green-200 dark:border-green-800',
    'Healthcare': 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300 border border-red-200 dark:border-red-800',
    'Retail': 'bg-pink-100 dark:bg-pink-900/30 text-pink-700 dark:text-pink-300 border border-pink-200 dark:border-pink-800',
    'Transportation': 'bg-indigo-100 dark:bg-indigo-900/30 text-indigo-700 dark:text-indigo-300 border border-indigo-200 dark:border-indigo-800'
  }
  return colors[industry] || 'bg-gray-100 dark:bg-gray-700/50 text-gray-700 dark:text-gray-300 border border-gray-200 dark:border-gray-600'
}

function getTagColor(tag: string): string {
  const colors: Record<string, string> = {
    'Tech': 'bg-primary-100 dark:bg-primary-900/30 text-primary-700 dark:text-primary-300',
    'Cloud': 'bg-sky-100 dark:bg-sky-900/30 text-sky-700 dark:text-sky-300',
    'BtoB': 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-300',
    'BtoC': 'bg-orange-100 dark:bg-orange-900/30 text-orange-700 dark:text-orange-300',
    'Auto': 'bg-violet-100 dark:bg-violet-900/30 text-violet-700 dark:text-violet-300',
    'Hardware': 'bg-slate-100 dark:bg-slate-700/50 text-slate-700 dark:text-slate-300',
    'Manufacturing': 'bg-stone-100 dark:bg-stone-700/50 text-stone-700 dark:text-stone-300'
  }
  return colors[tag] || 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300'
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <!-- Header -->
    <header class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <!-- Breadcrumbs -->
        <div class="flex items-center gap-2 mb-4 text-sm">
          <RouterLink to="/" class="text-primary-600 hover:underline">
            Home
          </RouterLink>
          <span class="text-gray-400">/</span>
          <span class="text-gray-600 dark:text-gray-400">Companies</span>
        </div>

        <!-- Main Header Content -->
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
              Corporate Energy Database
            </h1>
            <p class="text-gray-600 dark:text-gray-300 mt-1">
              {{ store.total }} companies | Showing {{ store.companies.length }}
            </p>
          </div>
          <div class="flex items-center gap-3">
            <!-- Quick Search -->
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <Search :size="18" class="text-gray-400" />
              </div>
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Search companies..."
                class="pl-10 pr-10 py-2 w-64 text-sm border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
              />
              <button
                v-if="searchQuery"
                @click="searchQuery = ''"
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
              >
                <X :size="16" />
              </button>
            </div>

            <ThemeToggle />
            <button
              @click="isFilterSidebarOpen = true"
              class="btn btn-primary flex items-center gap-2"
            >
              <Filter :size="20" />
              Filters
            </button>
          </div>
        </div>
      </div>
    </header>

    <!-- Active Filters Chips -->
    <div class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <ActiveFiltersChips
          :filters="store.filters"
          :total-results="store.companies.length"
          :total-count="store.total"
          @remove-industry="handleRemoveIndustry"
          @reset-search="handleResetSearch"
          @reset-efficiency="handleResetEfficiency"
          @reset-renewable="handleResetRenewable"
          @reset-all="handleFilterReset"
        />
      </div>
    </div>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Loading State -->
      <div v-if="store.loading" class="text-center py-12">
        <div class="text-gray-600 dark:text-gray-300">Loading companies...</div>
      </div>

      <!-- Error State -->
      <div v-else-if="store.error" class="card p-6 text-center">
        <div class="text-red-600 dark:text-red-400">{{ store.error }}</div>
      </div>

      <!-- Empty State -->
      <div v-else-if="store.companies.length === 0" class="card p-12 text-center">
        <div class="flex flex-col items-center max-w-md mx-auto">
          <div class="p-4 bg-gray-100 dark:bg-gray-700 rounded-full mb-4">
            <SearchX :size="48" class="text-gray-400 dark:text-gray-500" />
          </div>
          <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">
            No companies found
          </h3>
          <p class="text-gray-600 dark:text-gray-400 mb-6">
            We couldn't find any companies matching your current filters. Try adjusting your search criteria or clearing filters.
          </p>
          <button
            @click="handleFilterReset"
            class="btn btn-primary flex items-center gap-2"
          >
            <X :size="18" />
            <span>Clear All Filters</span>
          </button>
        </div>
      </div>

      <!-- Company Table -->
      <div v-else class="card overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-gray-50 dark:bg-gray-700">
              <tr>
                <!-- Checkbox Column -->
                <th class="px-4 py-3 w-12">
                  <div class="text-xs font-medium text-gray-500 dark:text-gray-400">
                    {{ selectedCompanies.length > 0 ? selectedCompanies.length : '' }}
                  </div>
                </th>

                <!-- Sortable: Company Name -->
                <th class="px-6 py-3 text-left">
                  <button
                    @click="handleSort('name')"
                    class="flex items-center gap-2 text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider hover:text-gray-700 dark:hover:text-gray-100 transition-colors group"
                  >
                    <span>Company</span>
                    <ArrowUp
                      v-if="store.sortField === 'name' && store.sortDirection === 'asc'"
                      :size="14"
                      class="text-primary-600 dark:text-primary-400"
                    />
                    <ArrowDown
                      v-else-if="store.sortField === 'name' && store.sortDirection === 'desc'"
                      :size="14"
                      class="text-primary-600 dark:text-primary-400"
                    />
                    <ArrowUpDown
                      v-else
                      :size="14"
                      class="text-gray-400 opacity-40 group-hover:opacity-100 transition-opacity"
                    />
                  </button>
                </th>

                <!-- Sortable: Industry -->
                <th class="px-6 py-3 text-left">
                  <button
                    @click="handleSort('industry')"
                    class="flex items-center gap-2 text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider hover:text-gray-700 dark:hover:text-gray-100 transition-colors group"
                  >
                    <span>Industry</span>
                    <ArrowUp
                      v-if="store.sortField === 'industry' && store.sortDirection === 'asc'"
                      :size="14"
                      class="text-primary-600 dark:text-primary-400"
                    />
                    <ArrowDown
                      v-else-if="store.sortField === 'industry' && store.sortDirection === 'desc'"
                      :size="14"
                      class="text-primary-600 dark:text-primary-400"
                    />
                    <ArrowUpDown
                      v-else
                      :size="14"
                      class="text-gray-400 opacity-40 group-hover:opacity-100 transition-opacity"
                    />
                  </button>
                </th>

                <!-- Sortable: Energy Consumption -->
                <th class="px-6 py-3 text-left">
                  <button
                    @click="handleSort('energyConsumption')"
                    class="flex items-center gap-2 text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider hover:text-gray-700 dark:hover:text-gray-100 transition-colors group"
                  >
                    <span>Energy (MWh/yr)</span>
                    <ArrowUp
                      v-if="store.sortField === 'energyConsumption' && store.sortDirection === 'asc'"
                      :size="14"
                      class="text-primary-600 dark:text-primary-400"
                    />
                    <ArrowDown
                      v-else-if="store.sortField === 'energyConsumption' && store.sortDirection === 'desc'"
                      :size="14"
                      class="text-primary-600 dark:text-primary-400"
                    />
                    <ArrowUpDown
                      v-else
                      :size="14"
                      class="text-gray-400 opacity-40 group-hover:opacity-100 transition-opacity"
                    />
                  </button>
                </th>

                <!-- Sortable: Efficiency -->
                <th class="px-6 py-3 text-left">
                  <button
                    @click="handleSort('efficiencyScore')"
                    class="flex items-center gap-2 text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider hover:text-gray-700 dark:hover:text-gray-100 transition-colors group"
                  >
                    <span>Efficiency</span>
                    <ArrowUp
                      v-if="store.sortField === 'efficiencyScore' && store.sortDirection === 'asc'"
                      :size="14"
                      class="text-primary-600 dark:text-primary-400"
                    />
                    <ArrowDown
                      v-else-if="store.sortField === 'efficiencyScore' && store.sortDirection === 'desc'"
                      :size="14"
                      class="text-primary-600 dark:text-primary-400"
                    />
                    <ArrowUpDown
                      v-else
                      :size="14"
                      class="text-gray-400 opacity-40 group-hover:opacity-100 transition-opacity"
                    />
                  </button>
                </th>

                <!-- Sortable: Renewable -->
                <th class="px-6 py-3 text-left">
                  <button
                    @click="handleSort('renewablePercentage')"
                    class="flex items-center gap-2 text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider hover:text-gray-700 dark:hover:text-gray-100 transition-colors group"
                  >
                    <span>Renewable</span>
                    <ArrowUp
                      v-if="store.sortField === 'renewablePercentage' && store.sortDirection === 'asc'"
                      :size="14"
                      class="text-primary-600 dark:text-primary-400"
                    />
                    <ArrowDown
                      v-else-if="store.sortField === 'renewablePercentage' && store.sortDirection === 'desc'"
                      :size="14"
                      class="text-primary-600 dark:text-primary-400"
                    />
                    <ArrowUpDown
                      v-else
                      :size="14"
                      class="text-gray-400 opacity-40 group-hover:opacity-100 transition-opacity"
                    />
                  </button>
                </th>

                <!-- Sortable: Trend -->
                <th class="px-6 py-3 text-left">
                  <button
                    @click="handleSort('trend')"
                    class="flex items-center gap-2 text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider hover:text-gray-700 dark:hover:text-gray-100 transition-colors group"
                  >
                    <span>Trend</span>
                    <ArrowUp
                      v-if="store.sortField === 'trend' && store.sortDirection === 'asc'"
                      :size="14"
                      class="text-primary-600 dark:text-primary-400"
                    />
                    <ArrowDown
                      v-else-if="store.sortField === 'trend' && store.sortDirection === 'desc'"
                      :size="14"
                      class="text-primary-600 dark:text-primary-400"
                    />
                    <ArrowUpDown
                      v-else
                      :size="14"
                      class="text-gray-400 opacity-40 group-hover:opacity-100 transition-opacity"
                    />
                  </button>
                </th>
              </tr>
            </thead>
            <tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
              <tr
                v-for="company in store.companies"
                :key="company.id"
                class="hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                :class="{ 'bg-primary-50/50 dark:bg-primary-900/10': selectedCompanies.includes(company.id) }"
              >
                <!-- Checkbox Cell -->
                <td class="px-4 py-4">
                  <input
                    type="checkbox"
                    :checked="selectedCompanies.includes(company.id)"
                    @change="toggleCompanySelection(company.id)"
                    :disabled="!selectedCompanies.includes(company.id) && selectedCompanies.length >= 4"
                    class="w-4 h-4 text-primary-600 bg-gray-100 border-gray-300 rounded focus:ring-primary-500 focus:ring-2 dark:bg-gray-700 dark:border-gray-600 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                  />
                </td>

                <td class="px-6 py-4">
                  <div>
                    <RouterLink
                      :to="`/companies/${company.id}`"
                      class="text-primary-600 dark:text-primary-400 font-medium hover:underline"
                    >
                      {{ company.name }}
                    </RouterLink>
                    <div class="text-sm text-gray-500 dark:text-gray-400">
                      {{ company.location }}
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4">
                  <div class="space-y-2">
                    <!-- Industry Badge -->
                    <div>
                      <span
                        :class="getIndustryColor(company.industry)"
                        class="inline-block px-2.5 py-1 text-xs font-semibold rounded-md"
                      >
                        {{ company.industry }}
                      </span>
                    </div>
                    <!-- Tags -->
                    <div class="flex flex-wrap gap-1">
                      <span
                        v-for="tag in company.tags"
                        :key="tag"
                        :class="getTagColor(tag)"
                        class="inline-block px-2 py-0.5 text-xs rounded"
                      >
                        {{ tag }}
                      </span>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4 text-gray-900 dark:text-white">
                  {{ company.energyConsumption.toLocaleString() }}
                </td>
                <td class="px-6 py-4">
                  <div class="flex items-center gap-2">
                    <span class="text-gray-900 dark:text-white font-medium">
                      {{ company.efficiencyScore.toFixed(1) }}%
                    </span>
                    <span
                      :class="getEfficiencyBadge(company.efficiencyScore).class"
                      class="px-2 py-1 text-xs rounded"
                    >
                      {{ getEfficiencyBadge(company.efficiencyScore).text }}
                    </span>
                  </div>
                </td>
                <td class="px-6 py-4">
                  <div class="flex items-center gap-2">
                    <span class="text-gray-900 dark:text-white">
                      {{ company.renewablePercentage }}%
                    </span>
                    <span
                      :class="getRenewableBadge(company.renewablePercentage).class"
                      class="px-2 py-1 text-xs rounded"
                    >
                      {{ getRenewableBadge(company.renewablePercentage).text }}
                    </span>
                  </div>
                </td>
                <td class="px-6 py-4">
                  <div class="flex items-center gap-1.5">
                    <!-- Lucide Icon -->
                    <TrendingUp
                      v-if="company.trend > 0"
                      :size="20"
                      class="text-green-600 dark:text-green-400"
                    />
                    <TrendingDown
                      v-else-if="company.trend < 0"
                      :size="20"
                      class="text-red-600 dark:text-red-400"
                    />
                    <Minus
                      v-else
                      :size="20"
                      class="text-gray-400 dark:text-gray-500"
                    />

                    <!-- Percentage Value -->
                    <span
                      :class="[
                        'font-semibold',
                        company.trend > 0 ? 'text-green-600 dark:text-green-400' :
                        company.trend < 0 ? 'text-red-600 dark:text-red-400' :
                        'text-gray-600 dark:text-gray-400'
                      ]"
                    >
                      {{ company.trend > 0 ? '+' : '' }}{{ company.trend }}%
                    </span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div
          v-if="store.total > store.pageSize"
          class="px-6 py-4 border-t border-gray-200 dark:border-gray-700 flex items-center justify-between"
        >
          <button
            @click="store.setPage(store.currentPage - 1)"
            :disabled="store.currentPage === 1"
            class="btn btn-secondary disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Previous
          </button>
          <span class="text-gray-600 dark:text-gray-300">
            Page {{ store.currentPage }} of {{ Math.ceil(store.total / store.pageSize) }}
          </span>
          <button
            @click="store.setPage(store.currentPage + 1)"
            :disabled="!store.hasMore"
            class="btn btn-secondary disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Next
          </button>
        </div>
      </div>
    </main>

    <!-- Filter Sidebar -->
    <FilterSidebar
      :filters="store.filters"
      :is-open="isFilterSidebarOpen"
      @update:filters="handleFilterUpdate"
      @update:is-open="(value) => isFilterSidebarOpen = value"
      @reset="handleFilterReset"
    />

    <!-- Floating Compare Action Bar -->
    <Transition
      enter-active-class="transition-all duration-300 ease-out"
      enter-from-class="translate-y-full opacity-0"
      enter-to-class="translate-y-0 opacity-100"
      leave-active-class="transition-all duration-200 ease-in"
      leave-from-class="translate-y-0 opacity-100"
      leave-to-class="translate-y-full opacity-0"
    >
      <div
        v-if="selectedCompanies.length > 0"
        class="fixed bottom-0 left-0 right-0 z-50"
      >
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-6">
          <div class="card bg-white/95 dark:bg-gray-800/95 backdrop-blur-lg border-2 border-primary-500 dark:border-primary-400 shadow-2xl">
            <div class="px-6 py-4 flex items-center justify-between">
              <div class="flex items-center gap-4">
                <div class="flex items-center gap-2">
                  <div class="w-8 h-8 bg-primary-600 text-white rounded-full flex items-center justify-center font-semibold text-sm">
                    {{ selectedCompanies.length }}
                  </div>
                  <div>
                    <div class="text-sm font-medium text-gray-900 dark:text-white">
                      {{ selectedCompanies.length }} {{ selectedCompanies.length === 1 ? 'company' : 'companies' }} selected
                    </div>
                    <div class="text-xs text-gray-600 dark:text-gray-400">
                      {{ selectedCompanies.length < 4 ? `Select up to ${4 - selectedCompanies.length} more` : 'Maximum reached' }}
                    </div>
                  </div>
                </div>
              </div>

              <div class="flex items-center gap-3">
                <button
                  @click="clearSelection"
                  class="btn btn-secondary flex items-center gap-2"
                >
                  <X :size="18" />
                  <span>Clear</span>
                </button>
                <button
                  @click="goToCompare"
                  :disabled="selectedCompanies.length < 2"
                  class="btn btn-primary flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <GitCompare :size="18" />
                  <span>Compare Selected</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>
