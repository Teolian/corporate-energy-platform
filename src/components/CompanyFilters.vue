<script setup lang="ts">
import { ref, computed } from 'vue'
import { FILTER_RANGES } from '@/constants/filters'

interface Filters {
  search: string
  industries: string[]
  efficiencyRange: [number, number]
  renewableRange: [number, number]
  sizeRange: string[]
}

interface Props {
  filters: Filters
  availableIndustries: string[]
}

const props = defineProps<Props>()
const emit = defineEmits<{
  update: [filters: Partial<Filters>]
  reset: []
}>()

const showAdvanced = ref(false)

const industryOptions = computed(() => props.availableIndustries)

function updateSearch(value: string) {
  emit('update', { search: value })
}

function toggleIndustry(industry: string) {
  const industries = props.filters.industries.includes(industry)
    ? props.filters.industries.filter((i) => i !== industry)
    : [...props.filters.industries, industry]
  emit('update', { industries })
}

function updateEfficiencyRange(min: number, max: number) {
  emit('update', { efficiencyRange: [min, max] })
}

function updateRenewableRange(min: number, max: number) {
  emit('update', { renewableRange: [min, max] })
}
</script>

<template>
  <div class="space-y-4">
    <!-- Search Bar -->
    <div class="flex gap-4">
      <div class="relative flex-1">
        <input
          :value="filters.search"
          @input="updateSearch(($event.target as HTMLInputElement).value)"
          type="text"
          placeholder="Search companies, industries, locations..."
          class="input pl-10"
        />
        <svg
          class="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 20 20"
          fill="currentColor"
        >
          <path
            fill-rule="evenodd"
            d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z"
            clip-rule="evenodd"
          />
        </svg>
      </div>
      <button @click="showAdvanced = !showAdvanced" class="btn btn-secondary whitespace-nowrap">
        {{ showAdvanced ? 'Hide Filters' : 'Show Filters' }}
      </button>
      <button @click="emit('reset')" class="btn btn-secondary">Reset</button>
    </div>

    <!-- Advanced Filters -->
    <div v-if="showAdvanced" class="card p-4 space-y-4">
      <!-- Industry Filter -->
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Industries
        </label>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="industry in industryOptions"
            :key="industry"
            @click="toggleIndustry(industry)"
            :class="[
              'px-3 py-1.5 rounded-lg text-sm font-medium transition-colors',
              filters.industries.includes(industry)
                ? 'bg-primary-600 text-white'
                : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
            ]"
          >
            {{ industry }}
          </button>
        </div>
      </div>

      <!-- Efficiency Range -->
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Efficiency Score: {{ filters.efficiencyRange[0] }}% - {{ filters.efficiencyRange[1] }}%
        </label>
        <div class="flex items-center gap-3">
          <input
            type="range"
            :value="filters.efficiencyRange[0]"
            @input="
              updateEfficiencyRange(
                Number(($event.target as HTMLInputElement).value),
                filters.efficiencyRange[1]
              )
            "
            :min="FILTER_RANGES.EFFICIENCY.MIN"
            :max="FILTER_RANGES.EFFICIENCY.MAX"
            class="flex-1"
          />
          <input
            type="range"
            :value="filters.efficiencyRange[1]"
            @input="
              updateEfficiencyRange(
                filters.efficiencyRange[0],
                Number(($event.target as HTMLInputElement).value)
              )
            "
            :min="FILTER_RANGES.EFFICIENCY.MIN"
            :max="FILTER_RANGES.EFFICIENCY.MAX"
            class="flex-1"
          />
        </div>
      </div>

      <!-- Renewable Energy Range -->
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Renewable Energy: {{ filters.renewableRange[0] }}% - {{ filters.renewableRange[1] }}%
        </label>
        <div class="flex items-center gap-3">
          <input
            type="range"
            :value="filters.renewableRange[0]"
            @input="
              updateRenewableRange(
                Number(($event.target as HTMLInputElement).value),
                filters.renewableRange[1]
              )
            "
            :min="FILTER_RANGES.RENEWABLE.MIN"
            :max="FILTER_RANGES.RENEWABLE.MAX"
            class="flex-1"
          />
          <input
            type="range"
            :value="filters.renewableRange[1]"
            @input="
              updateRenewableRange(
                filters.renewableRange[0],
                Number(($event.target as HTMLInputElement).value)
              )
            "
            :min="FILTER_RANGES.RENEWABLE.MIN"
            :max="FILTER_RANGES.RENEWABLE.MAX"
            class="flex-1"
          />
        </div>
      </div>
    </div>
  </div>
</template>
