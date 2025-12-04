<script setup lang="ts">
import { ref } from 'vue'
import Button from '@/components/ui/Button.vue'

const settings = ref({
  defaultSort: 'popular',
  showNSFW: false,
  autoplayVideo: true,
  compactMode: false,
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
</script>

<template>
  <div>
    <h2 class="text-xl font-bold text-text-primary dark:text-white mb-6">
      Настройки ленты
    </h2>
    
    <div class="space-y-6">
      <!-- Default sort -->
      <div>
        <label class="block text-sm font-medium text-text-primary dark:text-white mb-2">
          Сортировка по умолчанию
        </label>
        <select
          v-model="settings.defaultSort"
          class="w-full px-4 py-2.5 bg-white dark:bg-dark-tertiary border border-border dark:border-dark-tertiary rounded-lg text-text-primary dark:text-white"
        >
          <option value="popular">Популярное</option>
          <option value="new">Свежее</option>
          <option value="hot">Горячее</option>
        </select>
      </div>
      
      <!-- Toggles -->
      <div class="space-y-3">
        <label class="flex items-center justify-between p-3 bg-background-secondary dark:bg-dark-tertiary rounded-lg cursor-pointer">
          <div>
            <span class="block text-text-primary dark:text-white">
              Показывать NSFW контент
            </span>
            <span class="text-sm text-text-tertiary">
              Контент 18+ будет отображаться в ленте
            </span>
          </div>
          <input
            v-model="settings.showNSFW"
            type="checkbox"
            class="w-5 h-5 rounded text-primary focus:ring-primary cursor-pointer"
          />
        </label>
        
        <label class="flex items-center justify-between p-3 bg-background-secondary dark:bg-dark-tertiary rounded-lg cursor-pointer">
          <div>
            <span class="block text-text-primary dark:text-white">
              Автовоспроизведение видео
            </span>
            <span class="text-sm text-text-tertiary">
              Видео будут воспроизводиться автоматически
            </span>
          </div>
          <input
            v-model="settings.autoplayVideo"
            type="checkbox"
            class="w-5 h-5 rounded text-primary focus:ring-primary cursor-pointer"
          />
        </label>
        
        <label class="flex items-center justify-between p-3 bg-background-secondary dark:bg-dark-tertiary rounded-lg cursor-pointer">
          <div>
            <span class="block text-text-primary dark:text-white">
              Компактный режим
            </span>
            <span class="text-sm text-text-tertiary">
              Уменьшенные карточки статей
            </span>
          </div>
          <input
            v-model="settings.compactMode"
            type="checkbox"
            class="w-5 h-5 rounded text-primary focus:ring-primary cursor-pointer"
          />
        </label>
      </div>
      
      <div class="flex justify-end pt-4">
        <Button @click="handleSave" :loading="isSaving">
          Сохранить настройки
        </Button>
      </div>
    </div>
  </div>
</template>

