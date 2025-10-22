<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { useCompanyStore } from '@/stores/companies'
import EnergyTrendChart from '@/components/charts/EnergyTrendChart.vue'
import BenchmarkChart from '@/components/charts/BenchmarkChart.vue'
import ForecastAccuracyChart from '@/components/charts/ForecastAccuracyChart.vue'
import ThemeToggle from '@/components/ThemeToggle.vue'
import { TrendingUp, Zap, Star, Leaf, Users, BarChart3, Target, Activity, Bookmark, Lightbulb, CheckCircle, Clock, DollarSign } from 'lucide-vue-next'

const route = useRoute()
const store = useCompanyStore()
const companyDetail = ref<any>(null)

onMounted(async () => {
  const id = route.params.id as string
  companyDetail.value = await store.fetchCompanyById(id)
})

const trendChartData = computed(() => {
  if (!companyDetail.value?.forecasts) return null

  const forecasts = companyDetail.value.forecasts
  return {
    actual: forecasts.map((f: any) => f.actual || null),
    forecast: forecasts.map((f: any) => f.predicted),
    industryAvg: forecasts.map(() => companyDetail.value.benchmark.avgConsumption / 12)
  }
})

const benchmarkData = computed(() => {
  if (!companyDetail.value?.benchmark) return []

  const topPerformers = companyDetail.value.benchmark.topPerformers.map((c: any) => ({
    name: c.name,
    score: c.efficiencyScore,
    isYou: c.id === store.currentCompany?.id
  }))

  return topPerformers
})

const forecastData = computed(() => {
  if (!companyDetail.value?.forecasts) return []
  return companyDetail.value.forecasts
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <!-- Loading State -->
    <div v-if="store.loading" class="flex items-center justify-center min-h-screen">
      <div class="text-gray-600 dark:text-gray-300">Loading company details...</div>
    </div>

    <!-- Error State -->
    <div v-else-if="store.error" class="flex items-center justify-center min-h-screen">
      <div class="card p-6 text-center">
        <div class="text-red-600 dark:text-red-400 mb-4">{{ store.error }}</div>
        <RouterLink to="/companies" class="btn btn-primary">Back to Companies</RouterLink>
      </div>
    </div>

    <!-- Company Detail -->
    <div v-else-if="store.currentCompany">
      <!-- Header -->
      <header class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div class="flex items-center gap-2 mb-4 text-sm">
            <RouterLink to="/" class="text-primary-600 hover:underline">
              Home
            </RouterLink>
            <span class="text-gray-400">/</span>
            <RouterLink to="/companies" class="text-primary-600 hover:underline">
              Companies
            </RouterLink>
            <span class="text-gray-400">/</span>
            <span class="text-gray-600 dark:text-gray-400">{{ store.currentCompany.name }}</span>
          </div>
          <div class="flex items-start justify-between">
            <div>
              <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
                {{ store.currentCompany.name }}
              </h1>
              <p class="text-gray-600 dark:text-gray-300 mt-1">
                {{ store.currentCompany.location }} | {{ store.currentCompany.industry }}
              </p>
            </div>
            <div class="flex items-center gap-3">
              <ThemeToggle />
              <button class="btn btn-primary flex items-center gap-2">
                <Bookmark :size="18" />
                <span>Add to Watch</span>
              </button>
            </div>
          </div>
        </div>
      </header>

      <!-- Main Content -->
      <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <!-- Key Metrics -->
        <div class="grid md:grid-cols-4 gap-6 mb-8">
          <div class="card p-6">
            <div class="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 mb-2">
              <Zap :size="16" class="text-primary-500" />
              <span>Energy/Year</span>
            </div>
            <div class="text-2xl font-bold text-gray-900 dark:text-white">
              {{ store.currentCompany.energyConsumption.toLocaleString() }} MWh
            </div>
            <div class="flex items-center gap-1 text-sm text-green-600 mt-1">
              <TrendingUp :size="14" />
              <span>Industry</span>
            </div>
          </div>

          <div class="card p-6">
            <div class="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 mb-2">
              <Activity :size="16" class="text-primary-500" />
              <span>Efficiency Score</span>
            </div>
            <div class="text-2xl font-bold text-gray-900 dark:text-white">
              {{ store.currentCompany.efficiencyScore.toFixed(1) }}%
            </div>
            <div class="flex items-center gap-1 text-sm text-yellow-600 mt-1">
              <Star :size="14" />
              <span>Elite</span>
            </div>
          </div>

          <div class="card p-6">
            <div class="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 mb-2">
              <Leaf :size="16" class="text-primary-500" />
              <span>Renewable Energy</span>
            </div>
            <div class="text-2xl font-bold text-gray-900 dark:text-white">
              {{ store.currentCompany.renewablePercentage }}%
            </div>
            <div class="flex items-center gap-1 text-sm text-energy-green mt-1">
              <Leaf :size="14" />
              <span>Leader</span>
            </div>
          </div>

          <div class="card p-6">
            <div class="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 mb-2">
              <Users :size="16" class="text-primary-500" />
              <span>Employees</span>
            </div>
            <div class="text-2xl font-bold text-gray-900 dark:text-white">
              {{ store.currentCompany.employees.toLocaleString() }}
            </div>
            <div class="text-sm text-gray-600 dark:text-gray-400 mt-1">Large Enterprise</div>
          </div>
        </div>

        <!-- Energy Consumption Trend Chart -->
        <div class="card p-6 mb-8">
          <div class="flex items-center gap-2 mb-4">
            <BarChart3 :size="24" class="text-primary-500" />
            <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
              Energy Consumption Trend (Last 12 Months)
            </h2>
          </div>
          <div v-if="trendChartData" class="h-80">
            <EnergyTrendChart
              :company-name="store.currentCompany.name"
              :data="trendChartData"
              :labels="companyDetail.forecasts.map((f: any) => f.month)"
            />
          </div>
          <div v-else class="h-64 flex items-center justify-center bg-gray-50 dark:bg-gray-700 rounded">
            <p class="text-gray-500 dark:text-gray-400">Loading chart data...</p>
          </div>
        </div>

        <!-- Insights & Benchmark -->
        <div class="grid md:grid-cols-2 gap-6">
          <!-- Benchmark -->
          <div class="card p-6">
            <div class="flex items-center gap-2 mb-4">
              <Target :size="24" class="text-primary-500" />
              <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
                Benchmark Comparison
              </h2>
            </div>
            <div v-if="benchmarkData.length > 0" class="h-64 mb-4">
              <BenchmarkChart :companies="benchmarkData" />
            </div>
            <div class="space-y-3 mt-4">
              <div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Your Rank</div>
                <div class="text-lg font-semibold text-gray-900 dark:text-white">
                  #12 / 1,234 companies
                </div>
              </div>
              <div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Industry</div>
                <div class="text-lg text-gray-900 dark:text-white">
                  {{ store.currentCompany.industry }}
                </div>
              </div>
              <div class="text-sm text-green-600">Better than 99% of companies in sector</div>
            </div>
          </div>

          <!-- Forecast -->
          <div class="card p-6">
            <div class="flex items-center gap-2 mb-4">
              <Activity :size="24" class="text-primary-500" />
              <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
                Forecast vs Actual
              </h2>
            </div>
            <div v-if="forecastData.length > 0" class="h-64 mb-4">
              <ForecastAccuracyChart :data="forecastData" />
            </div>
            <div class="space-y-3 mt-4">
              <div v-if="forecastData.length > 0">
                <div class="text-sm text-gray-600 dark:text-gray-400">Last Month</div>
                <div class="flex items-center gap-2">
                  <span class="text-gray-900 dark:text-white">Forecast: {{ forecastData[forecastData.length - 3]?.predicted.toFixed(1) }} MWh</span>
                </div>
                <div class="flex items-center gap-2">
                  <span class="text-gray-900 dark:text-white">Actual: {{ forecastData[forecastData.length - 3]?.actual?.toFixed(1) }} MWh</span>
                </div>
                <div class="text-sm text-green-600 mt-1">Accuracy: {{ forecastData[forecastData.length - 3]?.accuracy?.toFixed(1) }}%</div>
              </div>
              <div v-if="forecastData.length > 0">
                <div class="text-sm text-gray-600 dark:text-gray-400">This Month (Forecast)</div>
                <div class="text-lg font-semibold text-gray-900 dark:text-white">
                  Predicted: {{ forecastData[forecastData.length - 1]?.predicted.toFixed(1) }} MWh
                </div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Confidence: {{ forecastData[forecastData.length - 1]?.confidence.toFixed(1) }}%</div>
              </div>
            </div>
          </div>
        </div>

        <!-- AI Insights -->
        <div class="card p-6 mt-6">
          <div class="flex items-center gap-2 mb-4">
            <Lightbulb :size="24" class="text-yellow-500" />
            <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
              AI Insights & Recommendations
            </h2>
          </div>
          <div class="space-y-3">
            <div class="flex items-start gap-3">
              <CheckCircle :size="20" class="text-green-600 mt-0.5 flex-shrink-0" />
              <p class="text-gray-700 dark:text-gray-300">
                Your efficiency is 15% above industry average
              </p>
            </div>
            <div class="flex items-start gap-3">
              <Clock :size="20" class="text-primary-600 mt-0.5 flex-shrink-0" />
              <p class="text-gray-700 dark:text-gray-300">
                Peak usage: Monday-Friday 9am-5pm (optimize scheduling)
              </p>
            </div>
            <div class="flex items-start gap-3">
              <Leaf :size="20" class="text-energy-green mt-0.5 flex-shrink-0" />
              <p class="text-gray-700 dark:text-gray-300">
                100% renewable - maintaining leadership position
              </p>
            </div>
            <div class="flex items-start gap-3">
              <DollarSign :size="20" class="text-yellow-600 mt-0.5 flex-shrink-0" />
              <p class="text-gray-700 dark:text-gray-300">
                Potential savings: $124k/year if reduce peak by 10%
              </p>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>
