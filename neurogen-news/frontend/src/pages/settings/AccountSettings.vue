<script setup lang="ts">
import { ref } from 'vue'
import { AlertTriangle } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const emailForm = ref({
  email: authStore.user?.email || '',
  password: '',
})

const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const isUpdatingEmail = ref(false)
const isUpdatingPassword = ref(false)
const showDeleteConfirm = ref(false)

const handleEmailUpdate = async () => {
  isUpdatingEmail.value = true
  try {
    // TODO: Call API
    await new Promise(resolve => setTimeout(resolve, 1000))
  } finally {
    isUpdatingEmail.value = false
  }
}

const handlePasswordUpdate = async () => {
  isUpdatingPassword.value = true
  try {
    // TODO: Call API
    await new Promise(resolve => setTimeout(resolve, 1000))
    passwordForm.value = { currentPassword: '', newPassword: '', confirmPassword: '' }
  } finally {
    isUpdatingPassword.value = false
  }
}

const handleDeleteAccount = async () => {
  // TODO: Implement account deletion
}
</script>

<template>
  <div class="space-y-8">
    <h2 class="text-xl font-bold text-text-primary dark:text-white">
      Аккаунт
    </h2>
    
    <!-- Email -->
    <section>
      <h3 class="text-lg font-medium text-text-primary dark:text-white mb-4">
        Email
      </h3>
      <form @submit.prevent="handleEmailUpdate" class="space-y-4">
        <Input
          v-model="emailForm.email"
          label="Email адрес"
          type="email"
          required
        />
        <Input
          v-model="emailForm.password"
          label="Текущий пароль"
          type="password"
          placeholder="Для подтверждения"
          required
        />
        <Button type="submit" :loading="isUpdatingEmail">
          Обновить email
        </Button>
      </form>
    </section>
    
    <hr class="border-border dark:border-dark-tertiary" />
    
    <!-- Password -->
    <section>
      <h3 class="text-lg font-medium text-text-primary dark:text-white mb-4">
        Пароль
      </h3>
      <form @submit.prevent="handlePasswordUpdate" class="space-y-4">
        <Input
          v-model="passwordForm.currentPassword"
          label="Текущий пароль"
          type="password"
          required
        />
        <Input
          v-model="passwordForm.newPassword"
          label="Новый пароль"
          type="password"
          hint="Минимум 8 символов"
          required
        />
        <Input
          v-model="passwordForm.confirmPassword"
          label="Подтвердите пароль"
          type="password"
          required
        />
        <Button type="submit" :loading="isUpdatingPassword">
          Изменить пароль
        </Button>
      </form>
    </section>
    
    <hr class="border-border dark:border-dark-tertiary" />
    
    <!-- Danger zone -->
    <section>
      <h3 class="text-lg font-medium text-error mb-4">
        Опасная зона
      </h3>
      <div class="p-4 bg-error/5 border border-error/20 rounded-lg">
        <div class="flex items-start gap-3">
          <AlertTriangle class="w-5 h-5 text-error shrink-0 mt-0.5" />
          <div>
            <h4 class="font-medium text-text-primary dark:text-white mb-1">
              Удалить аккаунт
            </h4>
            <p class="text-sm text-text-secondary mb-3">
              Это действие нельзя отменить. Все ваши данные будут удалены навсегда.
            </p>
            <Button 
              variant="danger" 
              size="sm"
              @click="showDeleteConfirm = true"
            >
              Удалить аккаунт
            </Button>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

