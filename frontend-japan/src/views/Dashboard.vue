<script setup lang="ts">
import { onMounted } from 'vue'
import { useDemandStore } from '@/stores/demand'
import { useDarkMode } from '@/composables/useDarkMode'
import BaseCard from '@/components/common/BaseCard.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import DemandChart from '@/components/demand/DemandChart.vue'
import { Moon, Sun, ChevronLeft, ChevronRight } from 'lucide-vue-next'

const demandStore = useDemandStore()
const { isDark, toggleDarkMode } = useDarkMode()

onMounted(() => {
  demandStore.fetchAllDemandData()
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 p-8">
    <!-- Header -->
    <header class="mb-8">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-4xl font-bold text-gray-900 dark:text-white">
            Japan Energy Dashboard
          </h1>
          <p class="text-gray-600 dark:text-gray-400 mt-2">
            Real-time demand data visualization • {{ demandStore.dataMode.toUpperCase() }} mode
          </p>
        </div>
        <button @click="toggleDarkMode" class="p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700">
          <Moon v-if="!isDark" :size="24" />
          <Sun v-else :size="24" />
        </button>
      </div>

      <!-- Date Navigation -->
      <div class="flex items-center gap-4 mt-6">
        <BaseButton size="sm" @click="demandStore.prevDay">
          <ChevronLeft :size="16" />
        </BaseButton>
        <div class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ demandStore.currentDate }}
        </div>
        <BaseButton size="sm" @click="demandStore.nextDay">
          <ChevronRight :size="16" />
        </BaseButton>
      </div>
    </header>

    <!-- Loading State -->
    <div v-if="demandStore.loading" class="text-center py-12">
      <div class="text-gray-600 dark:text-gray-400">Loading...</div>
    </div>

    <!-- Error State -->
    <div v-else-if="demandStore.error" class="text-center py-12">
      <div class="text-red-600">{{ demandStore.error }}</div>
    </div>

    <!-- Content -->
    <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- Tokyo -->
      <BaseCard>
        <DemandChart
          title="Tokyo Demand (TEPCO)"
          :data="demandStore.tokyoChartData"
        />
        <div v-if="demandStore.tokyoMetrics" class="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700 grid grid-cols-3 gap-4">
          <div>
            <div class="text-sm text-gray-600 dark:text-gray-400">Peak</div>
            <div class="text-xl font-bold">{{ demandStore.tokyoMetrics.peak.toFixed(0) }} MW</div>
          </div>
          <div>
            <div class="text-sm text-gray-600 dark:text-gray-400">Average</div>
            <div class="text-xl font-bold">{{ demandStore.tokyoMetrics.average.toFixed(0) }} MW</div>
          </div>
          <div v-if="demandStore.tokyoMetrics.forecastAccuracy">
            <div class="text-sm text-gray-600 dark:text-gray-400">Accuracy</div>
            <div class="text-xl font-bold">{{ demandStore.tokyoMetrics.forecastAccuracy }}%</div>
          </div>
        </div>
      </BaseCard>

      <!-- Kansai -->
      <BaseCard>
        <DemandChart
          title="Kansai Demand"
          :data="demandStore.kansaiChartData"
        />
        <div v-if="demandStore.kansaiMetrics" class="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700 grid grid-cols-3 gap-4">
          <div>
            <div class="text-sm text-gray-600 dark:text-gray-400">Peak</div>
            <div class="text-xl font-bold">{{ demandStore.kansaiMetrics.peak.toFixed(0) }} MW</div>
          </div>
          <div>
            <div class="text-sm text-gray-600 dark:text-gray-400">Average</div>
            <div class="text-xl font-bold">{{ demandStore.kansaiMetrics.average.toFixed(0) }} MW</div>
          </div>
          <div v-if="demandStore.kansaiMetrics.forecastAccuracy">
            <div class="text-sm text-gray-600 dark:text-gray-400">Accuracy</div>
            <div class="text-xl font-bold">{{ demandStore.kansaiMetrics.forecastAccuracy }}%</div>
          </div>
        </div>
      </BaseCard>
    </div>
  </div>
</template>
