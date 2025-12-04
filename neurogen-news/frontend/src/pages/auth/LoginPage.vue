<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute, RouterLink } from 'vue-router'
import { Mail, Lock, Eye, EyeOff, Sparkles } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const form = ref({
  email: '',
  password: '',
  remember: false,
})

const errors = ref<Record<string, string>>({})
const isLoading = ref(false)
const generalError = ref('')

const redirectPath = computed(() => 
  (route.query.redirect as string) || '/'
)

const validate = (): boolean => {
  errors.value = {}
  
  if (!form.value.email) {
    errors.value.email = 'Введите email'
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.value.email)) {
    errors.value.email = 'Введите корректный email'
  }
  
  if (!form.value.password) {
    errors.value.password = 'Введите пароль'
  }
  
  return Object.keys(errors.value).length === 0
}

const handleSubmit = async () => {
  if (!validate()) return
  
  isLoading.value = true
  generalError.value = ''
  
  try {
    const success = await authStore.login({
      email: form.value.email,
      password: form.value.password,
      remember: form.value.remember,
    })
    
    if (success) {
      router.push(redirectPath.value)
    } else {
      generalError.value = 'Неверный email или пароль'
    }
  } catch {
    generalError.value = 'Произошла ошибка. Попробуйте позже.'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="min-h-[80vh] flex items-center justify-center py-12 px-4">
    <div class="w-full max-w-md">
      <!-- Logo & Header -->
      <div class="text-center mb-8 animate-slide-up">
        <RouterLink to="/" class="inline-flex items-center gap-2.5 text-2xl font-bold group">
          <span class="text-3xl transition-transform duration-300 group-hover:scale-110 glow">🧠</span>
          <span>
            <span class="text-gradient font-display">Neurogen</span>
            <span class="text-gray-400 font-normal">.News</span>
          </span>
        </RouterLink>
        <h1 class="mt-6 text-2xl font-bold text-gray-900 dark:text-white flex items-center justify-center gap-2">
          Добро пожаловать!
          <Sparkles class="w-5 h-5 text-primary animate-pulse" />
        </h1>
        <p class="mt-2 text-gray-500 dark:text-gray-400">
          Войдите в аккаунт, чтобы продолжить
        </p>
      </div>
      
      <!-- Form Card -->
      <div class="velvet-panel p-8 animate-slide-up" style="animation-delay: 0.1s;">
        <!-- OAuth buttons -->
        <div class="space-y-3 mb-6">
          <button class="w-full flex items-center justify-center gap-3 px-4 py-3 velvet-button-glass rounded-xl font-medium transition-all duration-200 hover:-translate-y-0.5">
            <svg class="w-5 h-5" viewBox="0 0 24 24">
              <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
              <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
              <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
              <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
            </svg>
            Войти через Google
          </button>
          
          <button class="w-full flex items-center justify-center gap-3 px-4 py-3 velvet-button-glass rounded-xl font-medium transition-all duration-200 hover:-translate-y-0.5">
            <svg class="w-5 h-5 text-[#0077FF]" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2C6.477 2 2 6.477 2 12c0 4.991 3.657 9.128 8.438 9.879V14.89h-2.54V12h2.54V9.797c0-2.506 1.492-3.89 3.777-3.89 1.094 0 2.238.195 2.238.195v2.46h-1.26c-1.243 0-1.63.771-1.63 1.562V12h2.773l-.443 2.89h-2.33v6.989C18.343 21.129 22 16.99 22 12c0-5.523-4.477-10-10-10z"/>
            </svg>
            Войти через VK
          </button>
          
          <button class="w-full flex items-center justify-center gap-3 px-4 py-3 velvet-button-glass rounded-xl font-medium transition-all duration-200 hover:-translate-y-0.5">
            <svg class="w-5 h-5 text-[#26A5E4]" viewBox="0 0 24 24" fill="currentColor">
              <path d="M11.944 0A12 12 0 0 0 0 12a12 12 0 0 0 12 12 12 12 0 0 0 12-12A12 12 0 0 0 12 0a12 12 0 0 0-.056 0zm4.962 7.224c.1-.002.321.023.465.14a.506.506 0 0 1 .171.325c.016.093.036.306.02.472-.18 1.898-.962 6.502-1.36 8.627-.168.9-.499 1.201-.82 1.23-.696.065-1.225-.46-1.9-.902-1.056-.693-1.653-1.124-2.678-1.8-1.185-.78-.417-1.21.258-1.91.177-.184 3.247-2.977 3.307-3.23.007-.032.014-.15-.056-.212s-.174-.041-.249-.024c-.106.024-1.793 1.14-5.061 3.345-.48.33-.913.49-1.302.48-.428-.008-1.252-.241-1.865-.44-.752-.245-1.349-.374-1.297-.789.027-.216.325-.437.893-.663 3.498-1.524 5.83-2.529 6.998-3.014 3.332-1.386 4.025-1.627 4.476-1.635z"/>
            </svg>
            Войти через Telegram
          </button>
        </div>
        
        <!-- Divider -->
        <div class="relative my-6">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full h-px bg-gradient-to-r from-transparent via-gray-200 dark:via-gray-700 to-transparent" />
          </div>
          <div class="relative flex justify-center">
            <span class="px-4 text-sm text-gray-400 bg-white dark:bg-slate-800/80 backdrop-blur-sm rounded-full">
              или
            </span>
          </div>
        </div>
        
        <!-- Email/Password form -->
        <form @submit.prevent="handleSubmit" class="space-y-5">
          <!-- General error -->
          <div 
            v-if="generalError" 
            class="p-4 bg-error/10 border border-error/20 text-error text-sm rounded-xl flex items-center gap-3"
          >
            <svg class="w-5 h-5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="8" x2="12" y2="12"/>
              <line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
            {{ generalError }}
          </div>
          
          <Input
            v-model="form.email"
            type="email"
            label="Email"
            placeholder="your@email.com"
            :error="errors.email"
            :icon="Mail"
            required
          />
          
          <Input
            v-model="form.password"
            type="password"
            label="Пароль"
            placeholder="••••••••"
            :error="errors.password"
            :icon="Lock"
            required
          />
          
          <div class="flex items-center justify-between">
            <label class="velvet-checkbox">
              <input v-model="form.remember" type="checkbox" />
              <span class="velvet-checkbox-box">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                  <polyline points="20 6 9 17 4 12"/>
                </svg>
              </span>
              <span class="text-sm text-gray-600 dark:text-gray-400">Запомнить меня</span>
            </label>
            
            <RouterLink 
              to="/forgot-password"
              class="text-sm text-primary hover:text-primary-600 transition-colors"
            >
              Забыли пароль?
            </RouterLink>
          </div>
          
          <Button 
            type="submit" 
            variant="primary"
            :loading="isLoading"
            class="w-full"
            size="lg"
          >
            Войти
          </Button>
        </form>
      </div>
      
      <!-- Register link -->
      <p class="mt-8 text-center text-gray-500 dark:text-gray-400 animate-slide-up" style="animation-delay: 0.2s;">
        Нет аккаунта?
        <RouterLink 
          :to="{ name: 'register', query: route.query }"
          class="text-primary hover:text-primary-600 font-medium transition-colors"
        >
          Зарегистрироваться
        </RouterLink>
      </p>
    </div>
  </div>
</template>
