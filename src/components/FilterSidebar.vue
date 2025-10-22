<script setup lang="ts">
import { ref, computed } from 'vue'
import type { CompanyFilters } from '@/types/models'
import { X, Search } from 'lucide-vue-next'
import { FILTER_RANGES } from '@/constants/filters'
import BaseButton from '@/components/common/BaseButton.vue'
import BaseInput from '@/components/common/BaseInput.vue'

interface Props {
  filters: CompanyFilters
  isOpen: boolean
}

interface Emits {
  (e: 'update:filters', filters: CompanyFilters): void
  (e: 'update:isOpen', value: boolean): void
  (e: 'reset'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Industry options (can be loaded dynamically from API in the future)
const industryOptions = [
  'Technology',
  'Automotive',
  'Energy',
  'Manufacturing',
  'Finance',
  'Healthcare',
  'Retail',
  'Transportation'
]

// Local filter state
const localFilters = computed({
  get: () => props.filters,
  set: (value) => emit('update:filters', value)
})

const searchQuery = computed({
  get: () => props.filters.search,
  set: (value: string | number) => {
    emit('update:filters', {
      ...localFilters.value,
      search: String(value)
    })
  }
})

const toggleIndustry = (industry: string) => {
  const industries = [...localFilters.value.industries]
  const index = industries.indexOf(industry)

  if (index > -1) {
    industries.splice(index, 1)
  } else {
    industries.push(industry)
  }

  emit('update:filters', {
    ...localFilters.value,
    industries
  })
}

const updateEfficiencyRange = (min: number, max: number) => {
  emit('update:filters', {
    ...localFilters.value,
    efficiencyRange: [min, max]
  })
}

const updateRenewableRange = (min: number, max: number) => {
  emit('update:filters', {
    ...localFilters.value,
    renewableRange: [min, max]
  })
}

const closeSidebar = () => {
  emit('update:isOpen', false)
}
</script>

<template>
  <!-- Backdrop -->
  <Transition name="backdrop">
    <div
      v-if="isOpen"
      class="fixed inset-0 bg-black/50 z-40"
      @click="closeSidebar"
    ></div>
  </Transition>

  <!-- Sidebar -->
  <Transition name="sidebar">
    <div
      v-if="isOpen"
      class="fixed top-0 right-0 h-full w-96 bg-white dark:bg-gray-800 shadow-2xl z-50 overflow-y-auto"
    >
      <!-- Header -->
      <div class="sticky top-0 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-6 py-4 flex items-center justify-between">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white">Filters</h2>
        <BaseButton variant="ghost" size="sm" @click="closeSidebar">
          <X :size="20" />
        </BaseButton>
      </div>

      <!-- Filters Content -->
      <div class="p-6 space-y-6">
        <!-- Search -->
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Search Company
          </label>
          <BaseInput
            v-model="searchQuery"
            type="search"
            placeholder="Name, location, industry..."
            :icon="Search"
          />
        </div>

        <!-- Industries -->
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
            Industry ({{ localFilters.industries.length }})
          </label>
          <div class="flex flex-wrap gap-2">
            <BaseButton
              v-for="industry in industryOptions"
              :key="industry"
              variant="chip"
              size="sm"
              :active="localFilters.industries.includes(industry)"
              @click="toggleIndustry(industry)"
            >
              {{ industry }}
            </BaseButton>
          </div>
        </div>

        <!-- Efficiency Score Range -->
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
            Efficiency Score
          </label>

          <div class="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-4 space-y-3">
            <!-- Min Range -->
            <div>
              <div class="flex items-center justify-between mb-2">
                <span class="text-xs text-gray-600 dark:text-gray-400">Minimum</span>
                <span class="text-sm font-semibold text-primary-600 dark:text-primary-400">
                  {{ localFilters.efficiencyRange[0] }}%
                </span>
              </div>
              <input
                type="range"
                :value="localFilters.efficiencyRange[0]"
                @input="updateEfficiencyRange(parseInt(($event.target as HTMLInputElement).value), localFilters.efficiencyRange[1])"
                :min="FILTER_RANGES.EFFICIENCY.MIN"
                :max="FILTER_RANGES.EFFICIENCY.MAX"
                class="w-full accent-primary-500"
              />
            </div>

            <!-- Max Range -->
            <div>
              <div class="flex items-center justify-between mb-2">
                <span class="text-xs text-gray-600 dark:text-gray-400">Maximum</span>
                <span class="text-sm font-semibold text-primary-600 dark:text-primary-400">
                  {{ localFilters.efficiencyRange[1] }}%
                </span>
              </div>
              <input
                type="range"
                :value="localFilters.efficiencyRange[1]"
                @input="updateEfficiencyRange(localFilters.efficiencyRange[0], parseInt(($event.target as HTMLInputElement).value))"
                :min="FILTER_RANGES.EFFICIENCY.MIN"
                :max="FILTER_RANGES.EFFICIENCY.MAX"
                class="w-full accent-primary-500"
              />
            </div>

            <!-- Range Display -->
            <div class="pt-2 border-t border-gray-200 dark:border-gray-600">
              <div class="text-center text-sm text-gray-700 dark:text-gray-300">
                Range: <span class="font-semibold">{{ localFilters.efficiencyRange[0] }}% - {{ localFilters.efficiencyRange[1] }}%</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Renewable Percentage Range -->
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
            Renewable Energy
          </label>

          <div class="bg-green-50 dark:bg-green-900/20 rounded-lg p-4 space-y-3">
            <!-- Min Range -->
            <div>
              <div class="flex items-center justify-between mb-2">
                <span class="text-xs text-gray-600 dark:text-gray-400">Minimum</span>
                <span class="text-sm font-semibold text-energy-green">
                  {{ localFilters.renewableRange[0] }}%
                </span>
              </div>
              <input
                type="range"
                :value="localFilters.renewableRange[0]"
                @input="updateRenewableRange(parseInt(($event.target as HTMLInputElement).value), localFilters.renewableRange[1])"
                :min="FILTER_RANGES.RENEWABLE.MIN"
                :max="FILTER_RANGES.RENEWABLE.MAX"
                class="w-full accent-energy-green"
              />
            </div>

            <!-- Max Range -->
            <div>
              <div class="flex items-center justify-between mb-2">
                <span class="text-xs text-gray-600 dark:text-gray-400">Maximum</span>
                <span class="text-sm font-semibold text-energy-green">
                  {{ localFilters.renewableRange[1] }}%
                </span>
              </div>
              <input
                type="range"
                :value="localFilters.renewableRange[1]"
                @input="updateRenewableRange(localFilters.renewableRange[0], parseInt(($event.target as HTMLInputElement).value))"
                :min="FILTER_RANGES.RENEWABLE.MIN"
                :max="FILTER_RANGES.RENEWABLE.MAX"
                class="w-full accent-energy-green"
              />
            </div>

            <!-- Range Display -->
            <div class="pt-2 border-t border-green-200 dark:border-green-800">
              <div class="text-center text-sm text-gray-700 dark:text-gray-300">
                Range: <span class="font-semibold text-energy-green">{{ localFilters.renewableRange[0] }}% - {{ localFilters.renewableRange[1] }}%</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Reset Button -->
        <BaseButton
          variant="secondary"
          class="w-full"
          @click="emit('reset')"
        >
          Reset All Filters
        </BaseButton>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
/* Backdrop transitions */
.backdrop-enter-active,
.backdrop-leave-active {
  transition: opacity 0.3s ease;
}

.backdrop-enter-from,
.backdrop-leave-to {
  opacity: 0;
}

/* Sidebar transitions */
.sidebar-enter-active,
.sidebar-leave-active {
  transition: transform 0.3s ease;
}

.sidebar-enter-from,
.sidebar-leave-to {
  transform: translateX(100%);
}
</style>
