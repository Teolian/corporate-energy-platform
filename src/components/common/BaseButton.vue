<script setup lang="ts">
import { computed } from 'vue'

type ButtonVariant = 'primary' | 'secondary' | 'ghost' | 'chip'
type ButtonSize = 'sm' | 'md' | 'lg'

interface Props {
  variant?: ButtonVariant
  size?: ButtonSize
  disabled?: boolean
  active?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  disabled: false,
  active: false
})

const buttonClass = computed(() => {
  const base = 'inline-flex items-center justify-center font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed'

  // Variant styles
  const variants = {
    primary: 'bg-primary-600 text-white hover:bg-primary-700 focus:ring-primary-500',
    secondary: 'bg-gray-200 text-gray-900 hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-100 dark:hover:bg-gray-600 focus:ring-gray-500',
    ghost: 'text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 focus:ring-gray-500',
    chip: props.active
      ? 'bg-primary-500 text-white shadow-sm'
      : 'bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600'
  }

  // Size styles
  const sizes = {
    sm: props.variant === 'chip' ? 'px-3 py-1.5 text-sm rounded-full' : 'px-3 py-1.5 text-sm rounded-lg',
    md: props.variant === 'chip' ? 'px-4 py-2 text-base rounded-full' : 'px-4 py-2 text-base rounded-lg',
    lg: props.variant === 'chip' ? 'px-6 py-3 text-lg rounded-full' : 'px-6 py-3 text-lg rounded-lg'
  }

  return [base, variants[props.variant], sizes[props.size]].join(' ')
})
</script>

<template>
  <button
    :class="buttonClass"
    :disabled="disabled"
  >
    <slot />
  </button>
</template>
