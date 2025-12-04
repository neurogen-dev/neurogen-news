<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import { ArrowLeft, Mail, CheckCircle } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'

const email = ref('')
const error = ref('')
const isLoading = ref(false)
const isSuccess = ref(false)

const validate = (): boolean => {
  error.value = ''
  
  if (!email.value) {
    error.value = 'Введите email'
    return false
  }
  
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value)) {
    error.value = 'Введите корректный email'
    return false
  }
  
  return true
}

const handleSubmit = async () => {
  if (!validate()) return
  
  isLoading.value = true
  
  try {
    // TODO: Call API to send reset email
    await new Promise(resolve => setTimeout(resolve, 1500))
    isSuccess.value = true
  } catch {
    error.value = 'Произошла ошибка. Попробуйте позже.'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="min-h-[80vh] flex items-center justify-center">
    <div class="w-full max-w-md">
      <!-- Logo -->
      <div class="text-center mb-8">
        <RouterLink to="/" class="inline-flex items-center gap-2 text-2xl font-bold">
          <span class="text-3xl">🧠</span>
          <span>
            <span class="text-primary">Neurogen</span><span class="text-text-secondary">.News</span>
          </span>
        </RouterLink>
      </div>
      
      <!-- Form -->
      <div class="bg-white dark:bg-dark-secondary rounded-2xl border border-border dark:border-dark-tertiary p-6">
        <!-- Success state -->
        <template v-if="isSuccess">
          <div class="text-center py-4">
            <div class="w-16 h-16 bg-success/10 rounded-full flex items-center justify-center mx-auto mb-4">
              <CheckCircle class="w-8 h-8 text-success" />
            </div>
            <h2 class="text-xl font-bold text-text-primary dark:text-white mb-2">
              Письмо отправлено!
            </h2>
            <p class="text-text-secondary mb-6">
              Мы отправили инструкции по восстановлению пароля на 
              <span class="font-medium text-text-primary dark:text-white">{{ email }}</span>
            </p>
            <p class="text-sm text-text-tertiary mb-6">
              Не получили письмо? Проверьте папку «Спам» или 
              <button @click="isSuccess = false" class="text-primary hover:underline">
                попробуйте снова
              </button>
            </p>
            <Button as="RouterLink" to="/login" variant="secondary" class="w-full">
              <ArrowLeft class="w-4 h-4 mr-2" />
              Вернуться ко входу
            </Button>
          </div>
        </template>
        
        <!-- Form state -->
        <template v-else>
          <div class="text-center mb-6">
            <div class="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center mx-auto mb-4">
              <Mail class="w-8 h-8 text-primary" />
            </div>
            <h2 class="text-xl font-bold text-text-primary dark:text-white mb-2">
              Забыли пароль?
            </h2>
            <p class="text-text-secondary">
              Введите email, который вы использовали при регистрации. 
              Мы отправим ссылку для сброса пароля.
            </p>
          </div>
          
          <form @submit.prevent="handleSubmit" class="space-y-4">
            <Input
              v-model="email"
              type="email"
              label="Email"
              placeholder="your@email.com"
              :error="error"
              required
            />
            
            <Button 
              type="submit" 
              :loading="isLoading"
              class="w-full"
            >
              Отправить ссылку
            </Button>
          </form>
          
          <div class="mt-6 text-center">
            <RouterLink 
              to="/login"
              class="inline-flex items-center gap-1 text-sm text-text-secondary hover:text-primary transition-colors"
            >
              <ArrowLeft class="w-4 h-4" />
              Вернуться ко входу
            </RouterLink>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

