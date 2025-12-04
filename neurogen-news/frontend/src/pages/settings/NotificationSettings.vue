<script setup lang="ts">
import { ref } from 'vue'
import Button from '@/components/ui/Button.vue'

const settings = ref({
  email: {
    newFollower: true,
    newComment: true,
    commentReply: true,
    mentions: true,
    articlePublished: false,
    newsletter: true,
    marketing: false,
  },
  push: {
    newFollower: true,
    newComment: true,
    commentReply: true,
    mentions: true,
  },
})

const isSaving = ref(false)

const handleSave = async () => {
  isSaving.value = true
  try {
    // TODO: Call API
    await new Promise(resolve => setTimeout(resolve, 1000))
  } finally {
    isSaving.value = false
  }
}

const toggleGroups = [
  {
    title: 'Email уведомления',
    key: 'email',
    items: [
      { key: 'newFollower', label: 'Новые подписчики' },
      { key: 'newComment', label: 'Комментарии к статьям' },
      { key: 'commentReply', label: 'Ответы на комментарии' },
      { key: 'mentions', label: 'Упоминания' },
      { key: 'articlePublished', label: 'Публикация статей подписок' },
      { key: 'newsletter', label: 'Еженедельная рассылка' },
      { key: 'marketing', label: 'Новости и акции' },
    ],
  },
  {
    title: 'Push уведомления',
    key: 'push',
    items: [
      { key: 'newFollower', label: 'Новые подписчики' },
      { key: 'newComment', label: 'Комментарии к статьям' },
      { key: 'commentReply', label: 'Ответы на комментарии' },
      { key: 'mentions', label: 'Упоминания' },
    ],
  },
]
</script>

<template>
  <div>
    <h2 class="text-xl font-bold text-text-primary dark:text-white mb-6">
      Уведомления
    </h2>
    
    <div class="space-y-8">
      <section v-for="group in toggleGroups" :key="group.key">
        <h3 class="text-lg font-medium text-text-primary dark:text-white mb-4">
          {{ group.title }}
        </h3>
        
        <div class="space-y-3">
          <label 
            v-for="item in group.items" 
            :key="item.key"
            class="flex items-center justify-between p-3 bg-background-secondary dark:bg-dark-tertiary rounded-lg cursor-pointer hover:bg-background-tertiary dark:hover:bg-dark transition-colors"
          >
            <span class="text-text-primary dark:text-white">
              {{ item.label }}
            </span>
            <input
              v-model="settings[group.key as keyof typeof settings][item.key as keyof typeof settings.email]"
              type="checkbox"
              class="w-5 h-5 rounded text-primary focus:ring-primary cursor-pointer"
            />
          </label>
        </div>
      </section>
      
      <div class="flex justify-end pt-4">
        <Button @click="handleSave" :loading="isSaving">
          Сохранить настройки
        </Button>
      </div>
    </div>
  </div>
</template>

