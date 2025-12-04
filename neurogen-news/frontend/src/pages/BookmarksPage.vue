<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Bookmark, FolderPlus, Trash2, MoreHorizontal } from 'lucide-vue-next'
import FeedList from '@/components/feed/FeedList.vue'
import Button from '@/components/ui/Button.vue'
import type { ArticleCard } from '@/types'

const bookmarks = ref<ArticleCard[]>([])
const isLoading = ref(true)
const selectedFolder = ref<string | null>(null)

const folders = ref([
  { id: 'all', name: 'Все закладки', count: 0 },
  { id: 'read-later', name: 'Прочитать позже', count: 0 },
  { id: 'favorites', name: 'Избранное', count: 0 },
])

onMounted(async () => {
  // Simulate loading
  await new Promise(resolve => setTimeout(resolve, 500))
  isLoading.value = false
})
</script>

<template>
  <div class="min-h-screen">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-text-primary dark:text-white flex items-center gap-2">
          <Bookmark class="w-6 h-6" />
          Закладки
        </h1>
        <p class="text-text-secondary mt-1">
          Сохранённые статьи для быстрого доступа
        </p>
      </div>
      
      <Button variant="secondary">
        <FolderPlus class="w-4 h-4 mr-1" />
        Новая папка
      </Button>
    </div>
    
    <div class="flex gap-6">
      <!-- Folders sidebar -->
      <aside class="w-64 shrink-0 hidden lg:block">
        <nav class="space-y-1">
          <button
            v-for="folder in folders"
            :key="folder.id"
            @click="selectedFolder = folder.id === 'all' ? null : folder.id"
            class="w-full flex items-center justify-between px-4 py-3 rounded-lg transition-colors"
            :class="[
              (selectedFolder === null && folder.id === 'all') || selectedFolder === folder.id
                ? 'bg-primary/10 text-primary'
                : 'text-text-secondary hover:text-text-primary hover:bg-background-secondary dark:hover:bg-dark-tertiary'
            ]"
          >
            <span>{{ folder.name }}</span>
            <span class="text-sm text-text-tertiary">{{ folder.count }}</span>
          </button>
        </nav>
      </aside>
      
      <!-- Bookmarks list -->
      <main class="flex-1">
        <FeedList 
          :articles="bookmarks" 
          :is-loading="isLoading"
          :has-more="false"
        />
        
        <!-- Empty state -->
        <div 
          v-if="!isLoading && bookmarks.length === 0"
          class="text-center py-12 bg-white dark:bg-dark-secondary rounded-xl border border-border dark:border-dark-tertiary"
        >
          <Bookmark class="w-16 h-16 text-text-tertiary mx-auto mb-4" />
          <h2 class="text-xl font-medium text-text-primary dark:text-white mb-2">
            Нет закладок
          </h2>
          <p class="text-text-secondary mb-6">
            Сохраняйте интересные статьи, чтобы вернуться к ним позже
          </p>
          <Button as="RouterLink" to="/">
            Перейти к статьям
          </Button>
        </div>
      </main>
    </div>
  </div>
</template>

