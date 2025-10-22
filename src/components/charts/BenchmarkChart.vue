<script setup lang="ts">
import { computed } from 'vue'
import { Bar } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
  type ChartOptions
} from 'chart.js'

// Register Chart.js components
ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend)

interface CompanyScore {
  name: string
  score: number
  isYou?: boolean
}

interface Props {
  companies: CompanyScore[]
}

const props = defineProps<Props>()

const chartData = computed(() => ({
  labels: props.companies.map((c) => c.name),
  datasets: [
    {
      label: 'Efficiency Score',
      data: props.companies.map((c) => c.score),
      backgroundColor: props.companies.map((c) =>
        c.isYou ? 'rgba(14, 165, 233, 0.8)' : 'rgba(156, 163, 175, 0.6)'
      ),
      borderColor: props.companies.map((c) =>
        c.isYou ? 'rgb(14, 165, 233)' : 'rgb(156, 163, 175)'
      ),
      borderWidth: 2,
      borderRadius: 6
    }
  ]
}))

const chartOptions: ChartOptions<'bar'> = {
  responsive: true,
  maintainAspectRatio: false,
  indexAxis: 'y',
  plugins: {
    legend: {
      display: false
    },
    tooltip: {
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      padding: 12,
      cornerRadius: 8,
      callbacks: {
        label: (context) => {
          const score = context.parsed.x?.toFixed(1) || '0'
          return `Efficiency: ${score}%`
        }
      }
    }
  },
  scales: {
    x: {
      beginAtZero: true,
      max: 100,
      grid: {
        color: 'rgba(156, 163, 175, 0.1)'
      },
      ticks: {
        font: {
          size: 11
        },
        callback: (value) => `${value}%`
      }
    },
    y: {
      grid: {
        display: false
      },
      ticks: {
        font: {
          size: 12,
          weight: 'bold'
        }
      }
    }
  }
}
</script>

<template>
  <div class="w-full h-full">
    <Bar :data="chartData" :options="chartOptions" />
  </div>
</template>
