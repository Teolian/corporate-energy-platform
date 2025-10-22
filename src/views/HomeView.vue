<script setup lang="ts">
import { RouterLink } from 'vue-router'
import ThemeToggle from '@/components/ThemeToggle.vue'
import { realCompaniesData } from '@/data/companies'
import {
  Database,
  TrendingUp,
  Target,
  Zap,
  Leaf,
  Building2,
  Users,
  DollarSign,
  LineChart,
  Shield,
  Clock,
  Award,
  ArrowRight,
  Search,
  BarChart3,
  TrendingDown
} from 'lucide-vue-next'
import {
  PREVIEW_LIMITS,
  EFFICIENCY_THRESHOLDS,
  RENEWABLE_THRESHOLD
} from '@/constants/filters'

// Get top companies for preview based on efficiency score
const previewCompanies = realCompaniesData
  .sort((a, b) => b.efficiencyScore - a.efficiencyScore)
  .slice(0, PREVIEW_LIMITS.HOME_PAGE)
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <!-- Navigation -->
    <nav class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center gap-2">
            <Zap :size="28" class="text-primary-600" />
            <span class="text-xl font-bold text-gray-900 dark:text-white">EnergyBench</span>
          </div>
          <div class="flex items-center gap-4">
            <ThemeToggle />
            <RouterLink to="/companies" class="btn btn-primary">
              Get Started
            </RouterLink>
          </div>
        </div>
      </div>
    </nav>

    <!-- Hero Section -->
    <header class="relative overflow-hidden bg-gradient-to-br from-primary-50 via-white to-primary-50 dark:from-gray-800 dark:via-gray-900 dark:to-gray-800">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 relative z-10">
        <div class="grid lg:grid-cols-2 gap-12 items-center">
          <!-- Left: Text Content -->
          <div>
            <div class="inline-flex items-center gap-2 px-4 py-2 bg-primary-100 dark:bg-primary-900/30 text-primary-700 dark:text-primary-300 rounded-full text-sm font-medium mb-6">
              <Award :size="16" />
              <span>Trusted by Fortune 500 Companies</span>
            </div>

            <h1 class="text-5xl lg:text-6xl font-bold text-gray-900 dark:text-white mb-6 leading-tight">
              Benchmark Your
              <span class="text-transparent bg-clip-text bg-gradient-to-r from-primary-600 to-primary-800">
                Energy Performance
              </span>
            </h1>

            <p class="text-xl text-gray-600 dark:text-gray-300 mb-8 leading-relaxed">
              Compare your company's energy efficiency with 50+ industry leaders. Get AI-powered insights, forecasts, and recommendations to reduce costs and carbon footprint.
            </p>

            <div class="flex flex-col sm:flex-row gap-4">
              <RouterLink to="/companies" class="btn btn-primary btn-lg inline-flex items-center gap-2 justify-center">
                <Search :size="20" />
                <span>Explore Database</span>
                <ArrowRight :size="20" />
              </RouterLink>
              <RouterLink to="/companies/1" class="btn btn-secondary btn-lg inline-flex items-center gap-2 justify-center">
                <BarChart3 :size="20" />
                <span>View Demo</span>
              </RouterLink>
            </div>

            <!-- Quick Stats -->
            <div class="mt-12 grid grid-cols-3 gap-6">
              <div>
                <div class="text-3xl font-bold text-gray-900 dark:text-white">50+</div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Companies</div>
              </div>
              <div>
                <div class="text-3xl font-bold text-gray-900 dark:text-white">8</div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Industries</div>
              </div>
              <div>
                <div class="text-3xl font-bold text-gray-900 dark:text-white">98%</div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Accuracy</div>
              </div>
            </div>
          </div>

          <!-- Right: Live Database Preview -->
          <div class="hidden lg:block">
            <div class="relative">
              <div class="card overflow-hidden bg-white dark:bg-gray-800 shadow-xl">
                <!-- Preview Header -->
                <div class="px-4 py-3 bg-gradient-to-r from-primary-50 to-primary-100 dark:from-primary-900/20 dark:to-primary-800/20 border-b border-gray-200 dark:border-gray-700">
                  <div class="flex items-center justify-between">
                    <div class="flex items-center gap-2">
                      <Database :size="16" class="text-primary-600 dark:text-primary-400" />
                      <span class="text-sm font-semibold text-gray-900 dark:text-white">Live Preview</span>
                    </div>
                    <div class="text-xs text-gray-600 dark:text-gray-400">Top 5 Companies</div>
                  </div>
                </div>

                <!-- Mini Table -->
                <div class="divide-y divide-gray-100 dark:divide-gray-700">
                  <div
                    v-for="(company, index) in previewCompanies"
                    :key="company.id"
                    class="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors cursor-pointer group"
                  >
                    <div class="flex items-center justify-between mb-2">
                      <div class="flex items-center gap-2 flex-1 min-w-0">
                        <div class="w-6 h-6 rounded-full bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center flex-shrink-0">
                          <span class="text-xs font-semibold text-primary-600 dark:text-primary-400">{{ index + 1 }}</span>
                        </div>
                        <div class="min-w-0 flex-1">
                          <div class="text-sm font-semibold text-gray-900 dark:text-white truncate">
                            {{ company.name }}
                          </div>
                          <div class="text-xs text-gray-500 dark:text-gray-400">{{ company.industry }}</div>
                        </div>
                      </div>
                      <div class="flex items-center gap-3 flex-shrink-0">
                        <div class="text-right">
                          <div class="text-xs text-gray-500 dark:text-gray-400">Efficiency</div>
                          <div class="text-sm font-bold text-gray-900 dark:text-white">{{ company.efficiencyScore.toFixed(1) }}%</div>
                        </div>
                        <div class="w-8 h-8 rounded-lg bg-green-100 dark:bg-green-900/30 flex items-center justify-center">
                          <Leaf
                            :size="16"
                            class="text-energy-green"
                            v-if="company.renewablePercentage >= RENEWABLE_THRESHOLD.HIGH"
                          />
                          <Zap :size="16" class="text-yellow-600" v-else />
                        </div>
                      </div>
                    </div>

                    <!-- Progress bar -->
                    <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-1.5">
                      <div
                        class="h-1.5 rounded-full transition-all"
                        :class="
                          company.efficiencyScore >= EFFICIENCY_THRESHOLDS.EXCELLENT
                            ? 'bg-green-500'
                            : company.efficiencyScore >= EFFICIENCY_THRESHOLDS.GOOD
                              ? 'bg-blue-500'
                              : 'bg-yellow-500'
                        "
                        :style="{ width: company.efficiencyScore + '%' }"
                      ></div>
                    </div>
                  </div>
                </div>

                <!-- View All Link -->
                <RouterLink
                  to="/companies"
                  class="block px-4 py-3 text-center text-sm font-medium text-primary-600 dark:text-primary-400 hover:bg-gray-50 dark:hover:bg-gray-700/50 border-t border-gray-200 dark:border-gray-700 transition-colors group"
                >
                  <span class="flex items-center justify-center gap-2">
                    <span>View All 50 Companies</span>
                    <ArrowRight :size="16" class="group-hover:translate-x-1 transition-transform" />
                  </span>
                </RouterLink>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Background Pattern -->
      <div class="absolute inset-0 opacity-5 dark:opacity-10">
        <div class="absolute top-0 right-0 w-96 h-96 bg-primary-500 rounded-full blur-3xl"></div>
        <div class="absolute bottom-0 left-0 w-96 h-96 bg-primary-300 rounded-full blur-3xl"></div>
      </div>
    </header>

    <!-- Features Section -->
    <section class="py-20 bg-white dark:bg-gray-800">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-16">
          <h2 class="text-4xl font-bold text-gray-900 dark:text-white mb-4">
            Everything You Need to Optimize Energy
          </h2>
          <p class="text-xl text-gray-600 dark:text-gray-300 max-w-3xl mx-auto">
            Comprehensive tools and insights to benchmark, analyze, and improve your company's energy performance
          </p>
        </div>

        <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
          <!-- Feature 1 -->
          <div class="card p-8 hover:shadow-xl transition-shadow">
            <div class="p-3 bg-primary-100 dark:bg-primary-900/30 rounded-xl w-fit mb-4">
              <Database :size="32" class="text-primary-600" />
            </div>
            <h3 class="text-xl font-semibold mb-3 text-gray-900 dark:text-white">
              Comprehensive Database
            </h3>
            <p class="text-gray-600 dark:text-gray-300 mb-4">
              Access real energy data from 50+ Fortune 500 companies across 8 industries. Search, filter, and compare with advanced tools.
            </p>
            <ul class="space-y-2 text-sm text-gray-600 dark:text-gray-400">
              <li class="flex items-center gap-2">
                <div class="w-1.5 h-1.5 bg-primary-500 rounded-full"></div>
                <span>Real-time data updates</span>
              </li>
              <li class="flex items-center gap-2">
                <div class="w-1.5 h-1.5 bg-primary-500 rounded-full"></div>
                <span>Industry filtering</span>
              </li>
              <li class="flex items-center gap-2">
                <div class="w-1.5 h-1.5 bg-primary-500 rounded-full"></div>
                <span>Efficiency scoring</span>
              </li>
            </ul>
          </div>

          <!-- Feature 2 -->
          <div class="card p-8 hover:shadow-xl transition-shadow">
            <div class="p-3 bg-green-100 dark:bg-green-900/30 rounded-xl w-fit mb-4">
              <LineChart :size="32" class="text-green-600" />
            </div>
            <h3 class="text-xl font-semibold mb-3 text-gray-900 dark:text-white">
              AI-Powered Forecasting
            </h3>
            <p class="text-gray-600 dark:text-gray-300 mb-4">
              Predict future energy consumption with 98% accuracy. Get monthly forecasts and trend analysis powered by machine learning.
            </p>
            <ul class="space-y-2 text-sm text-gray-600 dark:text-gray-400">
              <li class="flex items-center gap-2">
                <div class="w-1.5 h-1.5 bg-green-500 rounded-full"></div>
                <span>12-month predictions</span>
              </li>
              <li class="flex items-center gap-2">
                <div class="w-1.5 h-1.5 bg-green-500 rounded-full"></div>
                <span>Accuracy tracking</span>
              </li>
              <li class="flex items-center gap-2">
                <div class="w-1.5 h-1.5 bg-green-500 rounded-full"></div>
                <span>Trend visualization</span>
              </li>
            </ul>
          </div>

          <!-- Feature 3 -->
          <div class="card p-8 hover:shadow-xl transition-shadow">
            <div class="p-3 bg-yellow-100 dark:bg-yellow-900/30 rounded-xl w-fit mb-4">
              <Target :size="32" class="text-yellow-600" />
            </div>
            <h3 class="text-xl font-semibold mb-3 text-gray-900 dark:text-white">
              Industry Benchmarks
            </h3>
            <p class="text-gray-600 dark:text-gray-300 mb-4">
              Compare your performance against industry averages and top performers. Identify gaps and opportunities for improvement.
            </p>
            <ul class="space-y-2 text-sm text-gray-600 dark:text-gray-400">
              <li class="flex items-center gap-2">
                <div class="w-1.5 h-1.5 bg-yellow-500 rounded-full"></div>
                <span>Peer comparisons</span>
              </li>
              <li class="flex items-center gap-2">
                <div class="w-1.5 h-1.5 bg-yellow-500 rounded-full"></div>
                <span>Ranking system</span>
              </li>
              <li class="flex items-center gap-2">
                <div class="w-1.5 h-1.5 bg-yellow-500 rounded-full"></div>
                <span>Best practices</span>
              </li>
            </ul>
          </div>

          <!-- Feature 4 -->
          <div class="card p-8 hover:shadow-xl transition-shadow">
            <div class="p-3 bg-blue-100 dark:bg-blue-900/30 rounded-xl w-fit mb-4">
              <Shield :size="32" class="text-blue-600" />
            </div>
            <h3 class="text-xl font-semibold mb-3 text-gray-900 dark:text-white">
              Renewable Tracking
            </h3>
            <p class="text-gray-600 dark:text-gray-300 mb-4">
              Monitor renewable energy adoption across industries. Track progress towards sustainability goals and ESG compliance.
            </p>
            <ul class="space-y-2 text-sm text-gray-600 dark:text-gray-400">
              <li class="flex items-center gap-2">
                <div class="w-1.5 h-1.5 bg-blue-500 rounded-full"></div>
                <span>Renewable percentage</span>
              </li>
              <li class="flex items-center gap-2">
                <div class="w-1.5 h-1.5 bg-blue-500 rounded-full"></div>
                <span>Carbon footprint</span>
              </li>
              <li class="flex items-center gap-2">
                <div class="w-1.5 h-1.5 bg-blue-500 rounded-full"></div>
                <span>ESG reporting</span>
              </li>
            </ul>
          </div>

          <!-- Feature 5 -->
          <div class="card p-8 hover:shadow-xl transition-shadow">
            <div class="p-3 bg-purple-100 dark:bg-purple-900/30 rounded-xl w-fit mb-4">
              <Clock :size="32" class="text-purple-600" />
            </div>
            <h3 class="text-xl font-semibold mb-3 text-gray-900 dark:text-white">
              Real-Time Insights
            </h3>
            <p class="text-gray-600 dark:text-gray-300 mb-4">
              Get instant recommendations based on usage patterns. AI analyzes peak times and suggests optimization strategies.
            </p>
            <ul class="space-y-2 text-sm text-gray-600 dark:text-gray-400">
              <li class="flex items-center gap-2">
                <div class="w-1.5 h-1.5 bg-purple-500 rounded-full"></div>
                <span>Peak usage alerts</span>
              </li>
              <li class="flex items-center gap-2">
                <div class="w-1.5 h-1.5 bg-purple-500 rounded-full"></div>
                <span>Cost optimization</span>
              </li>
              <li class="flex items-center gap-2">
                <div class="w-1.5 h-1.5 bg-purple-500 rounded-full"></div>
                <span>Scheduling tips</span>
              </li>
            </ul>
          </div>

          <!-- Feature 6 -->
          <div class="card p-8 hover:shadow-xl transition-shadow">
            <div class="p-3 bg-red-100 dark:bg-red-900/30 rounded-xl w-fit mb-4">
              <DollarSign :size="32" class="text-red-600" />
            </div>
            <h3 class="text-xl font-semibold mb-3 text-gray-900 dark:text-white">
              Cost Analysis
            </h3>
            <p class="text-gray-600 dark:text-gray-300 mb-4">
              Calculate potential savings with detailed cost breakdowns. Identify ROI for energy efficiency investments.
            </p>
            <ul class="space-y-2 text-sm text-gray-600 dark:text-gray-400">
              <li class="flex items-center gap-2">
                <div class="w-1.5 h-1.5 bg-red-500 rounded-full"></div>
                <span>Savings calculator</span>
              </li>
              <li class="flex items-center gap-2">
                <div class="w-1.5 h-1.5 bg-red-500 rounded-full"></div>
                <span>ROI projections</span>
              </li>
              <li class="flex items-center gap-2">
                <div class="w-1.5 h-1.5 bg-red-500 rounded-full"></div>
                <span>Budget planning</span>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <!-- Use Cases Section -->
    <section class="py-20 bg-gray-50 dark:bg-gray-900">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-16">
          <h2 class="text-4xl font-bold text-gray-900 dark:text-white mb-4">
            Who Benefits from EnergyBench?
          </h2>
          <p class="text-xl text-gray-600 dark:text-gray-300">
            Trusted by sustainability leaders across industries
          </p>
        </div>

        <div class="grid md:grid-cols-2 gap-8">
          <!-- Use Case 1 -->
          <div class="card p-8">
            <div class="flex items-start gap-4">
              <div class="p-3 bg-primary-100 dark:bg-primary-900/30 rounded-xl flex-shrink-0">
                <Building2 :size="28" class="text-primary-600" />
              </div>
              <div>
                <h3 class="text-xl font-semibold mb-2 text-gray-900 dark:text-white">
                  Corporate ESG Managers
                </h3>
                <p class="text-gray-600 dark:text-gray-300 mb-3">
                  Track your company's energy performance vs competitors. Generate ESG reports and demonstrate progress towards sustainability goals.
                </p>
                <div class="text-sm text-primary-600 dark:text-primary-400 font-medium">
                  Use Case: Quarterly ESG reporting & benchmarking
                </div>
              </div>
            </div>
          </div>

          <!-- Use Case 2 -->
          <div class="card p-8">
            <div class="flex items-start gap-4">
              <div class="p-3 bg-green-100 dark:bg-green-900/30 rounded-xl flex-shrink-0">
                <TrendingUp :size="28" class="text-green-600" />
              </div>
              <div>
                <h3 class="text-xl font-semibold mb-2 text-gray-900 dark:text-white">
                  Investors & Analysts
                </h3>
                <p class="text-gray-600 dark:text-gray-300 mb-3">
                  Identify companies with strong ESG performance. Build sustainable investment portfolios based on real energy data.
                </p>
                <div class="text-sm text-green-600 font-medium">
                  Use Case: ESG portfolio construction & due diligence
                </div>
              </div>
            </div>
          </div>

          <!-- Use Case 3 -->
          <div class="card p-8">
            <div class="flex items-start gap-4">
              <div class="p-3 bg-yellow-100 dark:bg-yellow-900/30 rounded-xl flex-shrink-0">
                <Users :size="28" class="text-yellow-600" />
              </div>
              <div>
                <h3 class="text-xl font-semibold mb-2 text-gray-900 dark:text-white">
                  Energy Consultants
                </h3>
                <p class="text-gray-600 dark:text-gray-300 mb-3">
                  Provide data-driven recommendations to clients. Show benchmarks and identify quick wins for efficiency improvements.
                </p>
                <div class="text-sm text-yellow-600 font-medium">
                  Use Case: Client audits & optimization proposals
                </div>
              </div>
            </div>
          </div>

          <!-- Use Case 4 -->
          <div class="card p-8">
            <div class="flex items-start gap-4">
              <div class="p-3 bg-blue-100 dark:bg-blue-900/30 rounded-xl flex-shrink-0">
                <Shield :size="28" class="text-blue-600" />
              </div>
              <div>
                <h3 class="text-xl font-semibold mb-2 text-gray-900 dark:text-white">
                  Regulators & Policy Makers
                </h3>
                <p class="text-gray-600 dark:text-gray-300 mb-3">
                  Monitor industry compliance with climate targets. Track renewable energy adoption trends across sectors.
                </p>
                <div class="text-sm text-blue-600 font-medium">
                  Use Case: Policy impact assessment & compliance monitoring
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="py-20 bg-gradient-to-r from-primary-600 to-primary-800 text-white">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <h2 class="text-4xl font-bold mb-6">
          Ready to Benchmark Your Energy Performance?
        </h2>
        <p class="text-xl text-primary-100 mb-8">
          Join leading companies using data to drive sustainability. Start exploring our database today.
        </p>
        <RouterLink to="/companies" class="btn bg-white text-primary-700 hover:bg-gray-100 btn-lg inline-flex items-center gap-2">
          <Database :size="20" />
          <span>Explore 50+ Companies</span>
          <ArrowRight :size="20" />
        </RouterLink>
      </div>
    </section>
  </div>
</template>

<style scoped>
.btn-lg {
  @apply px-8 py-4 text-lg;
}
</style>
