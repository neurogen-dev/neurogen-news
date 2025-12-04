<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { ChevronLeft, MessageCircle } from 'lucide-vue-next'
import type { Comment } from '@/types'

const route = useRoute()
const username = route.params.username as string

const comments = ref<Comment[]>([])
const isLoading = ref(true)

onMounted(async () => {
  await new Promise(resolve => setTimeout(resolve, 500))
  isLoading.value = false
})
</script>

<template>
  <div class="min-h-screen">
    <RouterLink 
      :to="`/@${username}`"
      class="inline-flex items-center gap-1 text-text-secondary hover:text-primary mb-6"
    >
      <ChevronLeft class="w-4 h-4" />
      Профиль @{{ username }}
    </RouterLink>
    
    <h1 class="text-2xl font-bold text-text-primary dark:text-white mb-6">
      Комментарии @{{ username }}
    </h1>
    
    <div 
      v-if="!isLoading && comments.length === 0"
      class="text-center py-12 bg-white dark:bg-dark-secondary rounded-xl"
    >
      <MessageCircle class="w-12 h-12 text-text-tertiary mx-auto mb-3" />
      <p class="text-text-secondary">Нет комментариев</p>
    </div>
  </div>
</template>

