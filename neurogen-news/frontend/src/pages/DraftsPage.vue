<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { FileText, PenSquare, Trash2, MoreHorizontal, Clock } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Badge from '@/components/ui/Badge.vue'
import { formatRelativeTime } from '@/utils/formatters'
import type { Draft } from '@/types'

const drafts = ref<Draft[]>([])
const isLoading = ref(true)

onMounted(async () => {
  // Simulate loading
  await new Promise(resolve => setTimeout(resolve, 500))
  isLoading.value = false
})

const deleteDraft = (id: string) => {
  if (confirm('Удалить черновик?')) {
    drafts.value = drafts.value.filter(d => d.id !== id)
  }
}
</script>

<template>
  <div class="min-h-screen">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-text-primary dark:text-white flex items-center gap-2">
          <FileText class="w-6 h-6" />
          Черновики
        </h1>
        <p class="text-text-secondary mt-1">
          Незавершённые публикации
        </p>
      </div>
      
      <Button as="RouterLink" to="/editor/new">
        <PenSquare class="w-4 h-4 mr-1" />
        Новая статья
      </Button>
    </div>
    
    <!-- Loading -->
    <div v-if="isLoading" class="space-y-4">
      <div 
        v-for="i in 3" 
        :key="i" 
        class="animate-pulse bg-white dark:bg-dark-secondary rounded-xl p-5"
      >
        <div class="h-6 bg-background-tertiary dark:bg-dark-tertiary rounded w-3/4 mb-3" />
        <div class="h-4 bg-background-tertiary dark:bg-dark-tertiary rounded w-1/2" />
      </div>
    </div>
    
    <!-- Drafts list -->
    <div v-else-if="drafts.length > 0" class="space-y-4">
      <article
        v-for="draft in drafts"
        :key="draft.id"
        class="bg-white dark:bg-dark-secondary rounded-xl border border-border dark:border-dark-tertiary p-5 hover:shadow-card-hover transition-shadow"
      >
        <div class="flex items-start justify-between gap-4">
          <div class="flex-1 min-w-0">
            <RouterLink 
              :to="`/editor/${draft.id}`"
              class="text-lg font-bold text-text-primary dark:text-white hover:text-primary transition-colors line-clamp-1"
            >
              {{ draft.title || 'Без названия' }}
            </RouterLink>
            
            <div class="flex items-center gap-3 mt-2 text-sm text-text-tertiary">
              <span class="flex items-center gap-1">
                <Clock class="w-4 h-4" />
                {{ formatRelativeTime(draft.updatedAt) }}
              </span>
              <Badge v-if="draft.level" :variant="draft.level">
                {{ draft.level }}
              </Badge>
            </div>
          </div>
          
          <div class="flex items-center gap-2">
            <Button 
              as="RouterLink" 
              :to="`/editor/${draft.id}`"
              variant="secondary" 
              size="sm"
            >
              Редактировать
            </Button>
            <button
              @click="deleteDraft(draft.id)"
              class="p-2 text-text-tertiary hover:text-error transition-colors"
            >
              <Trash2 class="w-4 h-4" />
            </button>
          </div>
        </div>
      </article>
    </div>
    
    <!-- Empty state -->
    <div 
      v-else
      class="text-center py-12 bg-white dark:bg-dark-secondary rounded-xl border border-border dark:border-dark-tertiary"
    >
      <FileText class="w-16 h-16 text-text-tertiary mx-auto mb-4" />
      <h2 class="text-xl font-medium text-text-primary dark:text-white mb-2">
        Нет черновиков
      </h2>
      <p class="text-text-secondary mb-6">
        Начните писать новую статью
      </p>
      <Button as="RouterLink" to="/editor/new">
        <PenSquare class="w-4 h-4 mr-1" />
        Создать статью
      </Button>
    </div>
  </div>
</template>

