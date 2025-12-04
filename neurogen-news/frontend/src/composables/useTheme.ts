import { computed, onMounted, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import type { Theme } from '@/types'

export function useTheme() {
  const authStore = useAuthStore()
  
  const theme = computed(() => authStore.theme)
  
  const isDark = computed(() => {
    if (theme.value === 'dark') return true
    if (theme.value === 'light') return false
    // System preference
    return window.matchMedia('(prefers-color-scheme: dark)').matches
  })
  
  const setTheme = (newTheme: Theme) => {
    authStore.setTheme(newTheme)
  }
  
  const toggleTheme = () => {
    if (theme.value === 'light') {
      setTheme('dark')
    } else if (theme.value === 'dark') {
      setTheme('system')
    } else {
      setTheme('light')
    }
  }
  
  // Listen for system theme changes
  onMounted(() => {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    
    const handleChange = () => {
      if (theme.value === 'system') {
        // Re-apply theme to update class
        authStore.setTheme('system')
      }
    }
    
    mediaQuery.addEventListener('change', handleChange)
    
    // Cleanup on unmount is handled by Vue's reactivity system
  })
  
  // Apply theme class to document
  watch(isDark, (dark) => {
    if (dark) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }, { immediate: true })
  
  return {
    theme,
    isDark,
    setTheme,
    toggleTheme,
  }
}
