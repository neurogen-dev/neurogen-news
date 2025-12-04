import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { CurrentUser, Theme, LoginRequest, RegisterRequest, AuthResponse } from '@/types'
import { apiClient } from '@/api/client'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<CurrentUser | null>(null)
  const accessToken = ref<string | null>(localStorage.getItem('accessToken'))
  const refreshToken = ref<string | null>(localStorage.getItem('refreshToken'))
  const theme = ref<Theme>((localStorage.getItem('theme') as Theme) || 'system')
  const hasSeenStartHere = ref(localStorage.getItem('hasSeenStartHere') === 'true')
  
  // Getters
  const isLoggedIn = computed(() => !!user.value)
  const isAdmin = computed(() => user.value?.role === 'ADMIN')
  const isModerator = computed(() => ['ADMIN', 'MODERATOR', 'EDITOR'].includes(user.value?.role || ''))
  const isAuthor = computed(() => ['ADMIN', 'MODERATOR', 'EDITOR', 'AUTHOR'].includes(user.value?.role || ''))
  
  // Actions
  async function login(credentials: LoginRequest): Promise<boolean> {
    try {
      const response = await apiClient.post<AuthResponse>('/auth/login', credentials)
      setAuth(response.data)
      return true
    } catch {
      return false
    }
  }
  
  async function register(data: RegisterRequest): Promise<boolean> {
    try {
      const response = await apiClient.post<AuthResponse>('/auth/register', data)
      setAuth(response.data)
      return true
    } catch {
      return false
    }
  }
  
  async function logout() {
    try {
      await apiClient.post('/auth/logout')
    } finally {
      clearAuth()
    }
  }
  
  async function refreshAccessToken(): Promise<boolean> {
    if (!refreshToken.value) return false
    
    try {
      const response = await apiClient.post<AuthResponse>('/auth/refresh', {
        refreshToken: refreshToken.value
      })
      setAuth(response.data)
      return true
    } catch {
      clearAuth()
      return false
    }
  }
  
  async function fetchCurrentUser() {
    if (!accessToken.value) return
    
    try {
      const response = await apiClient.get<CurrentUser>('/users/me')
      user.value = response.data
    } catch {
      clearAuth()
    }
  }
  
  function setAuth(data: AuthResponse) {
    user.value = data.user
    accessToken.value = data.accessToken
    refreshToken.value = data.refreshToken
    
    localStorage.setItem('accessToken', data.accessToken)
    localStorage.setItem('refreshToken', data.refreshToken)
  }
  
  function clearAuth() {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
  }
  
  function setTheme(newTheme: Theme) {
    theme.value = newTheme
    localStorage.setItem('theme', newTheme)
    
    // Apply theme to document
    if (newTheme === 'dark' || (newTheme === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }
  
  function markStartHereSeen() {
    hasSeenStartHere.value = true
    localStorage.setItem('hasSeenStartHere', 'true')
  }
  
  // Initialize theme on store creation
  setTheme(theme.value)
  
  return {
    // State
    user,
    accessToken,
    refreshToken,
    theme,
    hasSeenStartHere,
    
    // Getters
    isLoggedIn,
    isAdmin,
    isModerator,
    isAuthor,
    
    // Actions
    login,
    register,
    logout,
    refreshAccessToken,
    fetchCurrentUser,
    setTheme,
    markStartHereSeen,
  }
})
