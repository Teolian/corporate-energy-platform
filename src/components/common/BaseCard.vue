<script setup lang="ts">
interface Props {
  padding?: 'none' | 'sm' | 'md' | 'lg'
  hover?: boolean
}

withDefaults(defineProps<Props>(), {
  padding: 'md',
  hover: false
})
</script>

<template>
  <div
    :class="[
      'bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700',
      hover && 'hover:shadow-md transition-shadow',
      {
        'p-3': padding === 'sm',
        'p-4': padding === 'md',
        'p-6': padding === 'lg',
        'p-0': padding === 'none'
      }
    ]"
  >
    <slot />
  </div>
</template>
