import { ref, watch, onMounted } from 'vue'
import { STORAGE_KEYS, DARK_MODE_VALUES } from '@/constants/storage'

export function useDarkMode() {
  const isDark = ref(false)

  const toggleDarkMode = () => {
    isDark.value = !isDark.value
  }

  const setDarkMode = (value: boolean) => {
    isDark.value = value
  }

  // Apply dark mode class to document
  watch(
    isDark,
    (dark) => {
      if (dark) {
        document.documentElement.classList.add('dark')
        localStorage.setItem(STORAGE_KEYS.DARK_MODE, DARK_MODE_VALUES.ENABLED)
      } else {
        document.documentElement.classList.remove('dark')
        localStorage.setItem(STORAGE_KEYS.DARK_MODE, DARK_MODE_VALUES.DISABLED)
      }
    },
    { immediate: true }
  )

  // Initialize from localStorage or system preference
  onMounted(() => {
    const stored = localStorage.getItem(STORAGE_KEYS.DARK_MODE)
    if (stored !== null) {
      isDark.value = stored === DARK_MODE_VALUES.ENABLED
    } else {
      // Use system preference
      isDark.value = window.matchMedia('(prefers-color-scheme: dark)').matches
    }
  })

  return {
    isDark,
    toggleDarkMode,
    setDarkMode
  }
}
