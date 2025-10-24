import { ref, watch, onMounted } from 'vue'

const STORAGE_KEY = 'jp-energy-dark-mode'

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
        localStorage.setItem(STORAGE_KEY, 'enabled')
      } else {
        document.documentElement.classList.remove('dark')
        localStorage.setItem(STORAGE_KEY, 'disabled')
      }
    },
    { immediate: true }
  )

  // Initialize from localStorage or system preference
  onMounted(() => {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored !== null) {
      isDark.value = stored === 'enabled'
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
