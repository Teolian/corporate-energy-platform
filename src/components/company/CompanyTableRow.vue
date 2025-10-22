<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { TrendingUp, TrendingDown, Minus } from 'lucide-vue-next'
import type { Company } from '@/types/models'

interface Props {
  company: Company
  isSelected: boolean
  isSelectionDisabled: boolean
}

interface Emits {
  (e: 'toggle-selection', companyId: string): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

function getEfficiencyBadge(score: number) {
  if (score >= 90) return { text: 'Elite', class: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300' }
  if (score >= 80) return { text: 'High', class: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300' }
  if (score >= 70) return { text: 'Medium', class: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300' }
  return { text: 'Low', class: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300' }
}

function getRenewableBadge(percentage: number) {
  if (percentage >= 80) return { text: 'Leader', class: 'bg-energy-green text-white dark:bg-energy-green/90' }
  if (percentage >= 50) return { text: 'Good', class: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300' }
  return { text: 'Fair', class: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300' }
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
  return colors[industry] || 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300'
}

function getTagColor(tag: string): string {
  const colors: Record<string, string> = {
    'ESG Leader': 'bg-green-100 dark:bg-green-900/20 text-green-700 dark:text-green-400',
    'Carbon Neutral': 'bg-emerald-100 dark:bg-emerald-900/20 text-emerald-700 dark:text-emerald-400',
    'Net Zero 2030': 'bg-teal-100 dark:bg-teal-900/20 text-teal-700 dark:text-teal-400',
    'LEED Certified': 'bg-lime-100 dark:bg-lime-900/20 text-lime-700 dark:text-lime-400',
    'ISO 50001': 'bg-blue-100 dark:bg-blue-900/20 text-blue-700 dark:text-blue-400',
    'Energy Star': 'bg-yellow-100 dark:bg-yellow-900/20 text-yellow-700 dark:text-yellow-400'
  }
  return colors[tag] || 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-400'
}
</script>

<template>
  <tr
    class="hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
    :class="{ 'bg-primary-50/50 dark:bg-primary-900/10': isSelected }"
  >
    <!-- Checkbox Cell -->
    <td class="px-4 py-4">
      <input
        type="checkbox"
        :checked="isSelected"
        @change="emit('toggle-selection', company.id)"
        :disabled="isSelectionDisabled"
        class="w-4 h-4 text-primary-600 bg-gray-100 border-gray-300 rounded focus:ring-primary-500 focus:ring-2 dark:bg-gray-700 dark:border-gray-600 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
      />
    </td>

    <!-- Company Name -->
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

    <!-- Industry with Tags -->
    <td class="px-6 py-4">
      <div class="space-y-2">
        <div>
          <span
            :class="getIndustryColor(company.industry)"
            class="inline-block px-2.5 py-1 text-xs font-semibold rounded-md"
          >
            {{ company.industry }}
          </span>
        </div>
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

    <!-- Energy Consumption -->
    <td class="px-6 py-4 text-gray-900 dark:text-white">
      {{ company.energyConsumption.toLocaleString() }}
    </td>

    <!-- Efficiency Score -->
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

    <!-- Renewable Percentage -->
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

    <!-- Trend -->
    <td class="px-6 py-4">
      <div class="flex items-center gap-1.5">
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
</template>
