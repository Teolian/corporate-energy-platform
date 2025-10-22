<script setup lang="ts">
import { onMounted, watch, ref } from 'vue'
import { watchDebounced } from '@vueuse/core'
import { useCompanyStore } from '@/stores/companies'
import { RouterLink, useRouter } from 'vue-router'
import ThemeToggle from '@/components/ThemeToggle.vue'
import FilterSidebar from '@/components/FilterSidebar.vue'
import ActiveFiltersChips from '@/components/ActiveFiltersChips.vue'
import CompanyTable from '@/components/company/CompanyTable.vue'
import type { CompanyFilters as FiltersType, SortField, SortDirection } from '@/types/models'
import { Filter, GitCompare, X, Search, SearchX } from 'lucide-vue-next'

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

// Debounced filters watch for smooth UX
watchDebounced(
  () => [store.filters, store.sortField, store.sortDirection, store.currentPage],
  () => {
    store.fetchCompanies()
  },
  { debounce: 400, deep: true }
)

// Debounced search with automatic cleanup on unmount
watchDebounced(
  searchQuery,
  (newValue) => {
    store.setFilters({ search: newValue })
  },
  { debounce: 300 }
)
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
      <CompanyTable
        :companies="store.companies"
        :selected-companies="selectedCompanies"
        :sort-field="store.sortField"
        :sort-direction="store.sortDirection"
        :current-page="store.currentPage"
        :total-pages="Math.ceil(store.total / store.pageSize)"
        :has-more="store.hasMore"
        :loading="store.loading"
        @sort="handleSort"
        @toggle-selection="toggleCompanySelection"
        @page-change="store.setPage"
      />
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
