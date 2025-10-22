<script setup lang="ts">
interface Props {
  variant?: 'primary' | 'green'
  removable?: boolean
}

interface Emits {
  (e: 'remove'): void
}

withDefaults(defineProps<Props>(), {
  variant: 'primary',
  removable: true
})

const emit = defineEmits<Emits>()
</script>

<template>
  <div
    :class="[
      'inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-sm font-medium',
      variant === 'primary'
        ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-700 dark:text-primary-300'
        : 'bg-energy-green/10 dark:bg-energy-green/20 text-energy-green dark:text-energy-green'
    ]"
  >
    <slot />
    <button
      v-if="removable"
      @click="emit('remove')"
      :class="[
        'rounded-full p-0.5 transition-colors',
        variant === 'primary'
          ? 'hover:bg-primary-100 dark:hover:bg-primary-800/50'
          : 'hover:bg-energy-green/20 dark:hover:bg-energy-green/30'
      ]"
    >
      <slot name="remove-icon">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </slot>
    </button>
  </div>
</template>
