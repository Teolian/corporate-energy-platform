<script setup lang="ts">
import { computed } from 'vue'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  type ChartOptions
} from 'chart.js'
import type { ChartDataPoint } from '@/types/demand'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend)

interface Props {
  title: string
  data: ChartDataPoint[]
}

const props = defineProps<Props>()

const chartData = computed(() => ({
  labels: props.data.map(d => d.time),
  datasets: [
    {
      label: 'Actual Demand',
      data: props.data.map(d => d.actual),
      borderColor: 'rgb(14, 165, 233)',
      backgroundColor: 'rgba(14, 165, 233, 0.1)',
      borderWidth: 2,
      tension: 0.4
    },
    {
      label: 'Forecast',
      data: props.data.map(d => d.forecast || null),
      borderColor: 'rgb(147, 51, 234)',
      borderWidth: 2,
      borderDash: [5, 5],
      tension: 0.4
    }
  ]
}))

const chartOptions: ChartOptions<'line'> = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: true, position: 'top' },
    tooltip: {
      callbacks: {
        label: (context) => `${context.dataset.label}: ${context.parsed.y?.toFixed(0)} MW`
      }
    }
  },
  scales: {
    y: {
      beginAtZero: false,
      ticks: { callback: (value) => `${value} MW` }
    }
  }
}
</script>

<template>
  <div class="h-80">
    <h3 class="text-lg font-semibold mb-4">{{ title }}</h3>
    <Line :data="chartData" :options="chartOptions" />
  </div>
</template>
