<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute, RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const form = ref({
  email: '',
  username: '',
  password: '',
  confirmPassword: '',
  agreeToTerms: false,
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
  
  if (!form.value.username) {
    errors.value.username = 'Введите имя пользователя'
  } else if (!/^[a-zA-Z0-9_]{3,20}$/.test(form.value.username)) {
    errors.value.username = 'Имя должно содержать 3-20 символов (буквы, цифры, _)'
  }
  
  if (!form.value.password) {
    errors.value.password = 'Введите пароль'
  } else if (form.value.password.length < 8) {
    errors.value.password = 'Пароль должен содержать минимум 8 символов'
  }
  
  if (form.value.password !== form.value.confirmPassword) {
    errors.value.confirmPassword = 'Пароли не совпадают'
  }
  
  if (!form.value.agreeToTerms) {
    errors.value.terms = 'Необходимо принять условия'
  }
  
  return Object.keys(errors.value).length === 0
}

const handleSubmit = async () => {
  if (!validate()) return
  
  isLoading.value = true
  generalError.value = ''
  
  try {
    const success = await authStore.register({
      email: form.value.email,
      username: form.value.username,
      password: form.value.password,
    })
    
    if (success) {
      router.push(redirectPath.value)
    } else {
      generalError.value = 'Не удалось создать аккаунт. Возможно, email или имя уже заняты.'
    }
  } catch {
    generalError.value = 'Произошла ошибка. Попробуйте позже.'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="min-h-[80vh] flex items-center justify-center py-8">
    <div class="w-full max-w-md">
      <!-- Logo -->
      <div class="text-center mb-8">
        <RouterLink to="/" class="inline-flex items-center gap-2 text-2xl font-bold">
          <span class="text-3xl">🧠</span>
          <span>
            <span class="text-primary">Neurogen</span><span class="text-text-secondary">.News</span>
          </span>
        </RouterLink>
        <h1 class="mt-4 text-2xl font-bold text-text-primary dark:text-white">
          Создать аккаунт
        </h1>
        <p class="mt-2 text-text-secondary">
          Присоединяйтесь к сообществу
        </p>
      </div>
      
      <!-- Form -->
      <div class="bg-white dark:bg-dark-secondary rounded-2xl border border-border dark:border-dark-tertiary p-6">
        <!-- OAuth buttons -->
        <div class="space-y-3 mb-6">
          <Button variant="secondary" class="w-full">
            <svg class="w-5 h-5 mr-2" viewBox="0 0 24 24">
              <path fill="currentColor" d="M12.545,10.239v3.821h5.445c-0.712,2.315-2.647,3.972-5.445,3.972c-3.332,0-6.033-2.701-6.033-6.032s2.701-6.032,6.033-6.032c1.498,0,2.866,0.549,3.921,1.453l2.814-2.814C17.503,2.988,15.139,2,12.545,2C7.021,2,2.543,6.477,2.543,12s4.478,10,10.002,10c8.396,0,10.249-7.85,9.426-11.748L12.545,10.239z"/>
            </svg>
            Регистрация через Google
          </Button>
          
          <Button variant="secondary" class="w-full">
            <svg class="w-5 h-5 mr-2" viewBox="0 0 24 24">
              <path fill="currentColor" d="M12 2C6.477 2 2 6.477 2 12c0 4.991 3.657 9.128 8.438 9.879V14.89h-2.54V12h2.54V9.797c0-2.506 1.492-3.89 3.777-3.89 1.094 0 2.238.195 2.238.195v2.46h-1.26c-1.243 0-1.63.771-1.63 1.562V12h2.773l-.443 2.89h-2.33v6.989C18.343 21.129 22 16.99 22 12c0-5.523-4.477-10-10-10z"/>
            </svg>
            Регистрация через VK
          </Button>
        </div>
        
        <div class="relative my-6">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-border dark:border-dark-tertiary" />
          </div>
          <div class="relative flex justify-center text-sm">
            <span class="px-2 bg-white dark:bg-dark-secondary text-text-tertiary">
              или
            </span>
          </div>
        </div>
        
        <!-- Registration form -->
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <!-- General error -->
          <div 
            v-if="generalError" 
            class="p-3 bg-error/10 text-error text-sm rounded-lg"
          >
            {{ generalError }}
          </div>
          
          <Input
            v-model="form.email"
            type="email"
            label="Email"
            placeholder="your@email.com"
            :error="errors.email"
            required
          />
          
          <Input
            v-model="form.username"
            type="text"
            label="Имя пользователя"
            placeholder="username"
            :error="errors.username"
            hint="Только буквы, цифры и подчёркивания"
            required
          />
          
          <Input
            v-model="form.password"
            type="password"
            label="Пароль"
            placeholder="••••••••"
            :error="errors.password"
            hint="Минимум 8 символов"
            required
          />
          
          <Input
            v-model="form.confirmPassword"
            type="password"
            label="Подтвердите пароль"
            placeholder="••••••••"
            :error="errors.confirmPassword"
            required
          />
          
          <div>
            <label class="flex items-start gap-2 text-sm text-text-secondary cursor-pointer">
              <input 
                v-model="form.agreeToTerms"
                type="checkbox" 
                class="mt-1 w-4 h-4 rounded border-border dark:border-dark-tertiary text-primary focus:ring-primary"
              />
              <span>
                Я принимаю 
                <RouterLink to="/terms" class="text-primary hover:underline">условия использования</RouterLink>
                и 
                <RouterLink to="/privacy" class="text-primary hover:underline">политику конфиденциальности</RouterLink>
              </span>
            </label>
            <p v-if="errors.terms" class="mt-1 text-sm text-error">
              {{ errors.terms }}
            </p>
          </div>
          
          <Button 
            type="submit" 
            :loading="isLoading"
            class="w-full"
          >
            Создать аккаунт
          </Button>
        </form>
      </div>
      
      <!-- Login link -->
      <p class="mt-6 text-center text-text-secondary">
        Уже есть аккаунт?
        <RouterLink 
          :to="{ name: 'login', query: route.query }"
          class="text-primary hover:underline"
        >
          Войти
        </RouterLink>
      </p>
    </div>
  </div>
</template>

