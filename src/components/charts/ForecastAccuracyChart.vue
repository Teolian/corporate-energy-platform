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

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend)

interface DataPoint {
  month: string
  predicted: number
  actual?: number
  accuracy?: number
}

interface Props {
  data: DataPoint[]
}

const props = defineProps<Props>()

const chartData = computed(() => {
  const labels = props.data.map((d) => d.month)
  const predicted = props.data.map((d) => d.predicted)
  const actual = props.data.map((d) => d.actual || null)

  return {
    labels,
    datasets: [
      {
        label: 'Predicted',
        data: predicted,
        borderColor: 'rgb(147, 51, 234)',
        backgroundColor: 'rgba(147, 51, 234, 0.1)',
        borderWidth: 2,
        pointRadius: 4,
        pointHoverRadius: 6,
        tension: 0.3
      },
      {
        label: 'Actual',
        data: actual,
        borderColor: 'rgb(16, 185, 129)',
        backgroundColor: 'rgba(16, 185, 129, 0.1)',
        borderWidth: 3,
        pointRadius: 4,
        pointHoverRadius: 6,
        tension: 0.3,
        spanGaps: true
      }
    ]
  }
})

const chartOptions: ChartOptions<'line'> = {
  responsive: true,
  maintainAspectRatio: false,
  interaction: {
    mode: 'index',
    intersect: false
  },
  plugins: {
    legend: {
      display: true,
      position: 'top',
      labels: {
        usePointStyle: true,
        padding: 15
      }
    },
    tooltip: {
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      padding: 12,
      cornerRadius: 8,
      callbacks: {
        label: (context) => {
          const label = context.dataset.label || ''
          const value = context.parsed.y?.toFixed(1) || 'N/A'
          return `${label}: ${value} MWh`
        },
        afterBody: (tooltipItems) => {
          const index = tooltipItems[0].dataIndex
          const dataPoint = props.data[index]
          if (dataPoint.accuracy) {
            return [`Accuracy: ${dataPoint.accuracy.toFixed(1)}%`]
          }
          return []
        }
      }
    }
  },
  scales: {
    x: {
      grid: {
        display: false
      }
    },
    y: {
      beginAtZero: false,
      grid: {
        color: 'rgba(156, 163, 175, 0.1)'
      },
      ticks: {
        callback: (value) => `${value} MWh`
      }
    }
  }
}
</script>

<template>
  <div class="w-full h-full">
    <Line :data="chartData" :options="chartOptions" />
  </div>
</template>
