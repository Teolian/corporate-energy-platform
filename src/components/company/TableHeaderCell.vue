<script setup lang="ts">
import { ArrowUp, ArrowDown, ArrowUpDown } from 'lucide-vue-next'
import type { SortField, SortDirection } from '@/types/models'

interface Props {
  field: SortField
  label: string
  currentSortField?: SortField
  currentSortDirection?: SortDirection
}

interface Emits {
  (e: 'sort', field: SortField): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const isActive = (field: SortField) => props.currentSortField === field
</script>

<template>
  <th class="px-6 py-3 text-left">
    <button
      @click="emit('sort', field)"
      class="flex items-center gap-2 text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider hover:text-gray-700 dark:hover:text-gray-100 transition-colors group"
    >
      <span>{{ label }}</span>
      <ArrowUp
        v-if="isActive(field) && currentSortDirection === 'asc'"
        :size="14"
        class="text-primary-600 dark:text-primary-400"
      />
      <ArrowDown
        v-else-if="isActive(field) && currentSortDirection === 'desc'"
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
</template>
