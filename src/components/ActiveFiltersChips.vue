<script setup lang="ts">
import { computed } from 'vue'
import type { CompanyFilters } from '@/types/models'
import { Search, X } from 'lucide-vue-next'
import { FILTER_RANGES } from '@/constants/filters'
import FilterChip from '@/components/common/FilterChip.vue'
import BaseButton from '@/components/common/BaseButton.vue'

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
    <FilterChip v-if="filters.search" @remove="emit('resetSearch')">
      <Search :size="16" />
      <span>{{ filters.search }}</span>
    </FilterChip>

    <!-- Industry chips -->
    <FilterChip
      v-for="industry in filters.industries"
      :key="industry"
      @remove="emit('removeIndustry', industry)"
    >
      {{ industry }}
    </FilterChip>

    <!-- Efficiency range chip -->
    <FilterChip v-if="hasEfficiencyFilter" @remove="emit('resetEfficiency')">
      Efficiency: {{ filters.efficiencyRange[0] }}-{{ filters.efficiencyRange[1] }}%
    </FilterChip>

    <!-- Renewable range chip -->
    <FilterChip
      v-if="hasRenewableFilter"
      variant="green"
      @remove="emit('resetRenewable')"
    >
      Renewable: {{ filters.renewableRange[0] }}-{{ filters.renewableRange[1] }}%
    </FilterChip>

    <!-- Clear all button -->
    <BaseButton variant="ghost" size="sm" @click="emit('resetAll')">
      <X :size="16" />
      <span>Clear All</span>
    </BaseButton>
  </div>

  <!-- No filters active, just show total count -->
  <div v-else class="flex items-center mb-4">
    <div class="text-sm font-medium text-gray-700 dark:text-gray-300">
      Total Companies: <span class="text-primary-600 dark:text-primary-400">{{ totalCount }}</span>
    </div>
  </div>
</template>
