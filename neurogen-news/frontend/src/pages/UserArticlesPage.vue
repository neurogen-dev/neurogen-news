<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { ChevronLeft } from 'lucide-vue-next'
import FeedList from '@/components/feed/FeedList.vue'
import type { ArticleCard } from '@/types'

const route = useRoute()
const username = route.params.username as string

const articles = ref<ArticleCard[]>([])
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
      Статьи @{{ username }}
    </h1>
    
    <FeedList 
      :articles="articles"
      :is-loading="isLoading"
      :has-more="false"
    />
  </div>
</template>

