<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Camera } from 'lucide-vue-next'
import Avatar from '@/components/ui/Avatar.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const form = ref({
  displayName: '',
  username: '',
  bio: '',
  website: '',
  location: '',
})

const isSaving = ref(false)
const successMessage = ref('')

const handleSubmit = async () => {
  isSaving.value = true
  try {
    // TODO: Call API to update profile
    await new Promise(resolve => setTimeout(resolve, 1000))
    successMessage.value = 'Профиль успешно обновлён'
    setTimeout(() => successMessage.value = '', 3000)
  } finally {
    isSaving.value = false
  }
}

const handleAvatarUpload = () => {
  // TODO: Implement avatar upload
}

onMounted(() => {
  if (authStore.user) {
    form.value.displayName = authStore.user.displayName
    form.value.username = authStore.user.username
    form.value.bio = authStore.user.bio || ''
  }
})
</script>

<template>
  <div>
    <h2 class="text-xl font-bold text-text-primary dark:text-white mb-6">
      Профиль
    </h2>
    
    <!-- Success message -->
    <div 
      v-if="successMessage"
      class="mb-6 p-4 bg-success/10 text-success rounded-lg"
    >
      {{ successMessage }}
    </div>
    
    <form @submit.prevent="handleSubmit" class="space-y-6">
      <!-- Avatar -->
      <div>
        <label class="block text-sm font-medium text-text-primary dark:text-white mb-2">
          Аватар
        </label>
        <div class="flex items-center gap-4">
          <div class="relative">
            <Avatar 
              :src="authStore.user?.avatarUrl" 
              :alt="authStore.user?.displayName"
              :size="80"
            />
            <button
              type="button"
              @click="handleAvatarUpload"
              class="absolute bottom-0 right-0 p-1.5 bg-primary text-white rounded-full hover:bg-primary-600 transition-colors"
            >
              <Camera class="w-4 h-4" />
            </button>
          </div>
          <div class="text-sm text-text-tertiary">
            <p>JPG, PNG или GIF</p>
            <p>Максимум 2 MB</p>
          </div>
        </div>
      </div>
      
      <!-- Display name -->
      <Input
        v-model="form.displayName"
        label="Отображаемое имя"
        placeholder="Ваше имя"
        required
      />
      
      <!-- Username -->
      <Input
        v-model="form.username"
        label="Имя пользователя"
        placeholder="username"
        hint="Только латинские буквы, цифры и подчёркивания"
        required
      />
      
      <!-- Bio -->
      <div>
        <label class="block text-sm font-medium text-text-primary dark:text-white mb-1.5">
          О себе
        </label>
        <textarea
          v-model="form.bio"
          placeholder="Расскажите о себе..."
          rows="4"
          class="w-full px-4 py-2.5 bg-white dark:bg-dark-tertiary border border-border dark:border-dark-tertiary rounded-lg text-text-primary dark:text-white placeholder:text-text-tertiary resize-none"
        />
        <p class="mt-1.5 text-sm text-text-tertiary">
          {{ form.bio.length }}/300 символов
        </p>
      </div>
      
      <!-- Website -->
      <Input
        v-model="form.website"
        label="Веб-сайт"
        type="url"
        placeholder="https://example.com"
      />
      
      <!-- Location -->
      <Input
        v-model="form.location"
        label="Местоположение"
        placeholder="Москва, Россия"
      />
      
      <!-- Submit -->
      <div class="flex justify-end">
        <Button type="submit" :loading="isSaving">
          Сохранить изменения
        </Button>
      </div>
    </form>
  </div>
</template>

