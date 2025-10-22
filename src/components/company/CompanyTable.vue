<script setup lang="ts">
import type { Company, SortField, SortDirection } from '@/types/models'
import TableHeaderCell from './TableHeaderCell.vue'
import CompanyTableRow from './CompanyTableRow.vue'
import BaseButton from '@/components/common/BaseButton.vue'

interface Props {
  companies: Company[]
  selectedCompanies: string[]
  sortField?: SortField
  sortDirection?: SortDirection
  currentPage: number
  totalPages: number
  hasMore: boolean
  loading?: boolean
}

interface Emits {
  (e: 'sort', field: SortField): void
  (e: 'toggle-selection', companyId: string): void
  (e: 'page-change', page: number): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const isSelectionDisabled = (companyId: string) => {
  return !props.selectedCompanies.includes(companyId) && props.selectedCompanies.length >= 4
}
</script>

<template>
  <div class="card overflow-hidden">
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

            <!-- Sortable Columns -->
            <TableHeaderCell
              field="name"
              label="Company"
              :current-sort-field="sortField"
              :current-sort-direction="sortDirection"
              @sort="emit('sort', $event)"
            />
            <TableHeaderCell
              field="industry"
              label="Industry"
              :current-sort-field="sortField"
              :current-sort-direction="sortDirection"
              @sort="emit('sort', $event)"
            />
            <TableHeaderCell
              field="energyConsumption"
              label="Energy (MWh/yr)"
              :current-sort-field="sortField"
              :current-sort-direction="sortDirection"
              @sort="emit('sort', $event)"
            />
            <TableHeaderCell
              field="efficiencyScore"
              label="Efficiency"
              :current-sort-field="sortField"
              :current-sort-direction="sortDirection"
              @sort="emit('sort', $event)"
            />
            <TableHeaderCell
              field="renewablePercentage"
              label="Renewable"
              :current-sort-field="sortField"
              :current-sort-direction="sortDirection"
              @sort="emit('sort', $event)"
            />
            <TableHeaderCell
              field="trend"
              label="Trend"
              :current-sort-field="sortField"
              :current-sort-direction="sortDirection"
              @sort="emit('sort', $event)"
            />
          </tr>
        </thead>
        <tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
          <CompanyTableRow
            v-for="company in companies"
            :key="company.id"
            :company="company"
            :is-selected="selectedCompanies.includes(company.id)"
            :is-selection-disabled="isSelectionDisabled(company.id)"
            @toggle-selection="emit('toggle-selection', $event)"
          />
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div
      v-if="totalPages > 1"
      class="px-6 py-4 border-t border-gray-200 dark:border-gray-700 flex items-center justify-between"
    >
      <BaseButton
        variant="secondary"
        @click="emit('page-change', currentPage - 1)"
        :disabled="currentPage === 1"
      >
        Previous
      </BaseButton>
      <span class="text-gray-600 dark:text-gray-300">
        Page {{ currentPage }} of {{ totalPages }}
      </span>
      <BaseButton
        variant="secondary"
        @click="emit('page-change', currentPage + 1)"
        :disabled="!hasMore"
      >
        Next
      </BaseButton>
    </div>
  </div>
</template>
