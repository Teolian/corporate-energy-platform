<script setup lang="ts">
import { computed } from 'vue'
import type { CompanyFilters } from '@/types/models'
import { Search, X } from 'lucide-vue-next'
import { FILTER_RANGES } from '@/constants/filters'

interface Props {
  filters: CompanyFilters
  totalResults: number
  totalCount: number
}

interface Emits {
  (e: 'removeIndustry', industry: string): void
  (e: 'resetSearch'): void
  (e: 'resetEfficiency'): void
  (e: 'resetRenewable'): void
  (e: 'resetAll'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Compute active filters count
const activeFiltersCount = computed(() => {
  let count = 0
  if (props.filters.search) count++
  if (props.filters.industries.length > 0) count++
  if (
    props.filters.efficiencyRange[0] > FILTER_RANGES.EFFICIENCY.MIN ||
    props.filters.efficiencyRange[1] < FILTER_RANGES.EFFICIENCY.MAX
  )
    count++
  if (
    props.filters.renewableRange[0] > FILTER_RANGES.RENEWABLE.MIN ||
    props.filters.renewableRange[1] < FILTER_RANGES.RENEWABLE.MAX
  )
    count++
  return count
})

const hasEfficiencyFilter = computed(
  () =>
    props.filters.efficiencyRange[0] > FILTER_RANGES.EFFICIENCY.MIN ||
    props.filters.efficiencyRange[1] < FILTER_RANGES.EFFICIENCY.MAX
)

const hasRenewableFilter = computed(
  () =>
    props.filters.renewableRange[0] > FILTER_RANGES.RENEWABLE.MIN ||
    props.filters.renewableRange[1] < FILTER_RANGES.RENEWABLE.MAX
)
</script>

<template>
  <div v-if="activeFiltersCount > 0" class="flex flex-wrap items-center gap-3 mb-4">
    <!-- Results Counter -->
    <div class="text-sm font-medium text-gray-700 dark:text-gray-300">
      Found: <span class="text-primary-600 dark:text-primary-400">{{ totalResults }}</span> of {{ totalCount }}
    </div>

    <!-- Separator -->
    <div class="h-4 w-px bg-gray-300 dark:bg-gray-600"></div>

    <!-- Search chip -->
    <div
      v-if="filters.search"
      class="inline-flex items-center gap-2 px-3 py-1.5 bg-primary-50 dark:bg-primary-900/20 text-primary-700 dark:text-primary-300 rounded-full text-sm font-medium"
    >
      <Search :size="16" />
      <span>{{ filters.search }}</span>
      <button
        @click="emit('resetSearch')"
        class="hover:bg-primary-100 dark:hover:bg-primary-800/50 rounded-full p-0.5 transition-colors"
      >
        <X :size="16" />
      </button>
    </div>

    <!-- Industry chips -->
    <div
      v-for="industry in filters.industries"
      :key="industry"
      class="inline-flex items-center gap-2 px-3 py-1.5 bg-primary-50 dark:bg-primary-900/20 text-primary-700 dark:text-primary-300 rounded-full text-sm font-medium"
    >
      <span>{{ industry }}</span>
      <button
        @click="emit('removeIndustry', industry)"
        class="hover:bg-primary-100 dark:hover:bg-primary-800/50 rounded-full p-0.5 transition-colors"
      >
        <X :size="16" />
      </button>
    </div>

    <!-- Efficiency range chip -->
    <div
      v-if="hasEfficiencyFilter"
      class="inline-flex items-center gap-2 px-3 py-1.5 bg-primary-50 dark:bg-primary-900/20 text-primary-700 dark:text-primary-300 rounded-full text-sm font-medium"
    >
      <span>Efficiency: {{ filters.efficiencyRange[0] }}-{{ filters.efficiencyRange[1] }}%</span>
      <button
        @click="emit('resetEfficiency')"
        class="hover:bg-primary-100 dark:hover:bg-primary-800/50 rounded-full p-0.5 transition-colors"
      >
        <X :size="16" />
      </button>
    </div>

    <!-- Renewable range chip -->
    <div
      v-if="hasRenewableFilter"
      class="inline-flex items-center gap-2 px-3 py-1.5 bg-energy-green/10 dark:bg-energy-green/20 text-energy-green dark:text-energy-green rounded-full text-sm font-medium"
    >
      <span>Renewable: {{ filters.renewableRange[0] }}-{{ filters.renewableRange[1] }}%</span>
      <button
        @click="emit('resetRenewable')"
        class="hover:bg-energy-green/20 dark:hover:bg-energy-green/30 rounded-full p-0.5 transition-colors"
      >
        <X :size="16" />
      </button>
    </div>

    <!-- Clear all button -->
    <button
      @click="emit('resetAll')"
      class="inline-flex items-center gap-1.5 px-3 py-1.5 text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200 text-sm font-medium transition-colors"
    >
      <X :size="16" />
      <span>Clear All</span>
    </button>
  </div>

  <!-- No filters active, just show total count -->
  <div v-else class="flex items-center mb-4">
    <div class="text-sm font-medium text-gray-700 dark:text-gray-300">
      Total Companies: <span class="text-primary-600 dark:text-primary-400">{{ totalCount }}</span>
    </div>
  </div>
</template>
